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

package destination

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"sync"

	sdk "github.com/conduitio/conduit-connector-sdk"
	"github.com/miquido/conduit-connector-elasticsearch/internal/elasticsearch"
)

func NewDestination() sdk.Destination {
	return &Destination{}
}

type Destination struct {
	sdk.UnimplementedDestination

	config         Config
	client         elasticsearch.Client
	mutex          sync.Mutex
	recordsBuffer  []sdk.Record
	ackFuncsBuffer map[string]sdk.AckFunc
}

func (d *Destination) GetClient() elasticsearch.Client {
	return d.client
}

func (d *Destination) Configure(_ context.Context, cfgRaw map[string]string) (err error) {
	d.config, err = ParseConfig(cfgRaw)

	return
}

func (d *Destination) Open(ctx context.Context) (err error) {
	// Initialize Elasticsearch client
	d.client, err = elasticsearch.NewClient(d.config.Version, d.config)
	if err != nil {
		return fmt.Errorf("connection could not be established: %w", err)
	}

	// Check the connection
	if err := d.client.Ping(ctx); err != nil {
		return fmt.Errorf("connection could not be established: %w", err)
	}

	// Initializing the buffer
	d.mutex = sync.Mutex{}
	d.recordsBuffer = make([]sdk.Record, 0, d.config.BulkSize)
	d.ackFuncsBuffer = make(map[string]sdk.AckFunc, d.config.BulkSize)

	return nil
}

func (d *Destination) WriteAsync(ctx context.Context, record sdk.Record, ackFunc sdk.AckFunc) error {
	key := string(record.Key.Bytes())
	if key == "" {
		if err := ackFunc(fmt.Errorf("record's Key is required")); err != nil {
			return err
		}

		return nil
	}

	d.mutex.Lock()
	defer d.mutex.Unlock()

	d.recordsBuffer = append(d.recordsBuffer, record)
	d.ackFuncsBuffer[key] = ackFunc

	if uint64(len(d.recordsBuffer)) >= d.config.BulkSize {
		if err := d.Flush(ctx); err != nil {
			return err
		}
	}

	return nil
}

func (d *Destination) Flush(ctx context.Context) error {
	if len(d.recordsBuffer) == 0 {
		return nil
	}

	// Prepare request payload
	data := &bytes.Buffer{}

	for _, item := range d.recordsBuffer {
		key := string(item.Key.Bytes())

		switch action := item.Metadata["action"]; action {
		case "create", "created",
			"update", "updated":
			if err := d.writeUpsertOperation(key, data, item); err != nil {
				return err
			}

		case "delete", "deleted":
			if err := d.writeDeleteOperation(key, data); err != nil {
				return err
			}

		default:
			sdk.Logger(ctx).Warn().Msgf("unsupported action: %+v", action)

			continue
		}
	}

	// Send the bulk request
	response, err := d.executeBulkRequest(ctx, data)
	if err != nil {
		return err
	}

	// Ack results
	for _, item := range response.Items {
		var itemResponse BulkResponseItem

		switch {
		case item.Update != nil:
			itemResponse = *item.Update

		case item.Delete != nil:
			itemResponse = *item.Delete

		default:
			sdk.Logger(ctx).Warn().Msg("no update or delete details were found in Elasticsearch response")

			continue
		}

		if err := d.sendAckForOperation(itemResponse); err != nil {
			return err
		}
	}

	// Reset buffers
	d.recordsBuffer = d.recordsBuffer[:0]
	d.ackFuncsBuffer = make(map[string]sdk.AckFunc, d.config.BulkSize)

	return nil
}

func (d *Destination) Teardown(context.Context) error {
	return nil // No close routine needed
}

func (d *Destination) writeUpsertOperation(key string, data *bytes.Buffer, item sdk.Record) error {
	// BulkRequestActionAndMetadata
	bulkRequestUpdateAction := BulkRequestUpdateAction{
		ID:              key,
		Index:           d.config.Index,
		RetryOnConflict: 3,
	}

	if d.config.Version == elasticsearch.Version6 {
		bulkRequestUpdateAction.Type = &d.config.Type
	}

	entryMetadata, err := json.Marshal(BulkRequestActionAndMetadata{
		Update: &bulkRequestUpdateAction,
	})
	if err != nil {
		return fmt.Errorf("failed to prepare metadata with key=%s: %w", key, err)
	}

	data.Write(entryMetadata)
	data.WriteRune('\n')

	// BulkRequestOptionalSource
	sourcePayload := BulkRequestOptionalSource{
		DocAsUpsert: true,
	}

	switch itemPayload := item.Payload.(type) {
	case sdk.StructuredData:
		// Payload is potentially convertable into JSON
		itemPayloadMarshalled, err := json.Marshal(itemPayload)
		if err != nil {
			return fmt.Errorf("failed to prepare data with key=%s: %w", key, err)
		}

		sourcePayload.Doc = itemPayloadMarshalled

	default:
		// Nothing more can be done, we can trust the source to provide valid JSON
		sourcePayload.Doc = itemPayload.Bytes()
	}

	entrySource, err := json.Marshal(sourcePayload)
	if err != nil {
		return fmt.Errorf("failed to prepare data with key=%s: %w", key, err)
	}

	data.Write(entrySource)
	data.WriteRune('\n')
	return nil
}

func (d *Destination) writeDeleteOperation(key string, data *bytes.Buffer) error {
	// BulkRequestActionAndMetadata
	bulkRequestDeleteAction := BulkRequestDeleteAction{
		ID:    key,
		Index: d.config.Index,
	}

	if d.config.Version == elasticsearch.Version6 {
		bulkRequestDeleteAction.Type = &d.config.Type
	}

	entryMetadata, err := json.Marshal(BulkRequestActionAndMetadata{
		Delete: &bulkRequestDeleteAction,
	})
	if err != nil {
		return fmt.Errorf("failed to prepare metadata with key=%s: %w", key, err)
	}

	data.Write(entryMetadata)
	data.WriteRune('\n')

	return nil
}

func (d *Destination) executeBulkRequest(ctx context.Context, data *bytes.Buffer) (BulkResponse, error) {
	// Check if there is any job to do
	if data.Len() < 1 {
		sdk.Logger(ctx).Info().Msg("no operations to execute in bulk, skipping")

		return BulkResponse{}, nil
	}

	defer data.Reset()

	// Execute the request
	responseBody, err := d.client.Bulk(ctx, bytes.NewReader(data.Bytes()))
	if err != nil {
		return BulkResponse{}, fmt.Errorf("bulk request failure: %w", err)
	}

	// Get the response
	bodyContents, err := io.ReadAll(responseBody)
	if err != nil {
		return BulkResponse{}, fmt.Errorf("bulk response failure: failed to read the result: %w", err)
	}
	defer responseBody.Close()

	// Read individual errors
	var response BulkResponse
	if err := json.Unmarshal(bodyContents, &response); err != nil {
		return BulkResponse{}, fmt.Errorf("bulk response failure: could not read the response: %w", err)
	}

	return response, nil
}

func (d *Destination) sendAckForOperation(itemResponse BulkResponseItem) error {
	ackFunc, exists := d.ackFuncsBuffer[itemResponse.ID]
	if !exists {
		return fmt.Errorf("bulk response failure: could not ack item with key=%s: no ack function was registered", itemResponse.ID)
	}

	if itemResponse.Status >= 200 && itemResponse.Status < 300 {
		if err := ackFunc(nil); err != nil {
			return err
		}

		return nil
	}

	var operationError error

	if itemResponse.Error == nil {
		operationError = fmt.Errorf(
			"item with key=%s upsert/delete failure: unknown error",
			itemResponse.ID,
		)
	} else {
		operationError = fmt.Errorf(
			"item with key=%s upsert/delete failure: [%s] %s: %s",
			itemResponse.ID,
			itemResponse.Error.Type,
			itemResponse.Error.Reason,
			itemResponse.Error.CausedBy,
		)
	}

	if err := ackFunc(operationError); err != nil {
		return err
	}

	return nil
}

type BulkRequestActionAndMetadata struct {
	Update *BulkRequestUpdateAction `json:"update,omitempty"`
	Delete *BulkRequestDeleteAction `json:"delete,omitempty"`
}

type BulkRequestUpdateAction struct {
	ID              string  `json:"_id"`
	Index           string  `json:"_index"`
	Type            *string `json:"_type"`
	RetryOnConflict int     `json:"retry_on_conflict"`
}

type BulkRequestDeleteAction struct {
	ID    string  `json:"_id"`
	Index string  `json:"_index"`
	Type  *string `json:"_type"`
}

type BulkRequestOptionalSource struct {
	Doc         json.RawMessage `json:"doc"`
	DocAsUpsert bool            `json:"doc_as_upsert"`
}

type BulkResponse struct {
	Took   int  `json:"took"`
	Errors bool `json:"errors"`
	Items  []struct {
		Update *BulkResponseItem `json:"update,omitempty"`
		Delete *BulkResponseItem `json:"delete,omitempty"`
	} `json:"items"`
}

type BulkResponseItem struct {
	Index   string `json:"_index"`
	Type    string `json:"_type"`
	ID      string `json:"_id"`
	Version int    `json:"_version,omitempty"`
	Result  string `json:"result,omitempty"`
	Shards  *struct {
		Total      int `json:"total"`
		Successful int `json:"successful"`
		Failed     int `json:"failed"`
	} `json:"_shards,omitempty"`
	SeqNo       int `json:"_seq_no,omitempty"`
	PrimaryTerm int `json:"_primary_term,omitempty"`
	Status      int `json:"status"`
	Error       *struct {
		Type     string          `json:"type"`
		Reason   string          `json:"reason"`
		CausedBy json.RawMessage `json:"caused_by"`
	} `json:"error,omitempty"`
}
