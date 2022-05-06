// Copyright © 2022 Meroxa, Inc. and Miquido
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package v8

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"testing"
	"time"

	sdk "github.com/conduitio/conduit-connector-sdk"
	esV8 "github.com/elastic/go-elasticsearch/v8"
	"github.com/jaswdr/faker"
	"github.com/miquido/conduit-connector-elasticsearch/destination"
	"github.com/miquido/conduit-connector-elasticsearch/internal/elasticsearch"
	v8 "github.com/miquido/conduit-connector-elasticsearch/internal/elasticsearch/v8"
	"github.com/stretchr/testify/require"
)

func TestOperations(t *testing.T) {
	fakerInstance := faker.New()
	dest := destination.NewDestination().(*destination.Destination)

	cfgRaw := map[string]string{
		destination.ConfigKeyVersion:  elasticsearch.Version8,
		destination.ConfigKeyHost:     "http://127.0.0.1:9200",
		destination.ConfigKeyIndex:    "users",
		destination.ConfigKeyBulkSize: "1",
	}

	require.NoError(t, dest.Configure(context.Background(), cfgRaw))
	require.NoError(t, dest.Open(context.Background()))

	esClient := dest.GetClient().(*v8.Client).GetClient()

	require.True(t, assertIndexIsDeleted(esClient, "users"))

	t.Cleanup(func() {
		assertIndexIsDeleted(esClient, "users")

		require.NoError(t, dest.Teardown(context.Background()))
	})

	var (
		user1 = map[string]interface{}{
			"id":    float64(fakerInstance.Int32Between(100, 200)),
			"email": fakerInstance.Internet().Email(),
		}
		user2 = map[string]interface{}{
			"id":    float64(fakerInstance.Int32Between(201, 300)),
			"email": fakerInstance.Internet().Email(),
		}
	)

	ackChannel := make(chan bool, 1)
	ackFunc := func(err error) error {
		require.NoError(t, err)

		ackChannel <- true

		return nil
	}

	t.Run("records can be upserted", func(t *testing.T) {
		go func() {
			require.NoError(t, dest.WriteAsync(context.Background(), sdk.Record{
				Metadata: map[string]string{
					"action": "updated",
				},
				Payload:   sdk.StructuredData(user1),
				Key:       sdk.RawData(fmt.Sprintf("%.0f", user1["id"])),
				CreatedAt: time.Now(),
			}, ackFunc))
			require.NoError(t, dest.WriteAsync(context.Background(), sdk.Record{
				Metadata: map[string]string{
					"action": "created",
				},
				Payload:   sdk.StructuredData(user2),
				Key:       sdk.RawData(fmt.Sprintf("%.0f", user2["id"])),
				CreatedAt: time.Now(),
			}, ackFunc))
		}()

		select {
		case <-time.After(5 * time.Second):
			t.Fatalf("Timeout")

			return

		case <-ackChannel:
			time.Sleep(time.Second)
		}

		require.NoError(t, assertIndexContainsDocuments(t, esClient, []map[string]interface{}{
			user1,
			user2,
		}))
	})

	t.Run("record can be deleted", func(t *testing.T) {
		go func() {
			require.NoError(t, dest.WriteAsync(context.Background(), sdk.Record{
				Metadata: map[string]string{
					"action": "deleted",
				},
				Payload:   nil,
				Key:       sdk.RawData(fmt.Sprintf("%.0f", user1["id"])),
				CreatedAt: time.Now(),
			}, ackFunc))
		}()

		select {
		case <-time.After(5 * time.Second):
			t.Fatalf("Timeout")

			return

		case <-ackChannel:
			time.Sleep(time.Second)
		}

		require.NoError(t, assertIndexContainsDocuments(t, esClient, []map[string]interface{}{
			user2,
		}))
	})
}

func assertIndexIsDeleted(esClient *esV8.Client, index string) bool {
	res, err := esClient.Indices.Delete([]string{index}, esClient.Indices.Delete.WithIgnoreUnavailable(true))
	if err != nil || res.IsError() {
		log.Fatalf("Cannot delete index: %s", err)

		return false
	}

	return true
}

func assertIndexContainsDocuments(t *testing.T, esClient *esV8.Client, documents []map[string]interface{}) error {
	// Build the request body.
	var buf bytes.Buffer
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"match_all": map[string]interface{}{},
		},
		"sort": []map[string]interface{}{
			{
				"id": map[string]string{
					"order": "asc",
				},
			},
		},
	}

	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		return fmt.Errorf("error encoding query: %s", err)
	}

	// Search
	response, err := esClient.Search(
		esClient.Search.WithIndex("users"),
		esClient.Search.WithBody(&buf),
	)
	if err != nil {
		return fmt.Errorf("error getting response: %s", err)
	}

	defer response.Body.Close()

	if response.IsError() {
		var e map[string]interface{}

		if err := json.NewDecoder(response.Body).Decode(&e); err != nil {
			return fmt.Errorf("error parsing the response body: %s", err)
		}

		// Print the response status and error information.
		return fmt.Errorf("[%s] %s: %s",
			response.Status(),
			e["error"].(map[string]interface{})["type"],
			e["error"].(map[string]interface{})["reason"],
		)
	}

	var r map[string]interface{}

	if err := json.NewDecoder(response.Body).Decode(&r); err != nil {
		return fmt.Errorf("error parsing the response body: %s", err)
	}

	hitsMetadata := r["hits"].(map[string]interface{})
	totalMetadata := hitsMetadata["total"].(map[string]interface{})

	require.Equal(t, len(documents), int(totalMetadata["value"].(float64)))

	hits := hitsMetadata["hits"].([]interface{})

	for i, document := range documents {
		hit := hits[i].(map[string]interface{})
		source := hit["_source"].(map[string]interface{})

		require.EqualValues(t, document, source)
	}

	return nil
}
