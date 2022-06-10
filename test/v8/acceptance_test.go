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
	"testing"

	sdk "github.com/conduitio/conduit-connector-sdk"
	es "github.com/miquido/conduit-connector-elasticsearch"
	esDestination "github.com/miquido/conduit-connector-elasticsearch/destination"
	"github.com/miquido/conduit-connector-elasticsearch/internal/elasticsearch"
)

func TestAcceptance(t *testing.T) {
	sdk.AcceptanceTest(t, sdk.ConfigurableAcceptanceTestDriver{
		Config: sdk.ConfigurableAcceptanceTestDriverConfig{
			Connector: sdk.Connector{
				NewSpecification: es.Specification,
				NewSource:        nil,
				NewDestination:   esDestination.NewDestination,
			},

			SourceConfig: map[string]string{},

			DestinationConfig: map[string]string{
				esDestination.ConfigKeyVersion:  elasticsearch.Version8,
				esDestination.ConfigKeyHost:     "http://127.0.0.1:9200",
				esDestination.ConfigKeyIndex:    "acceptance",
				esDestination.ConfigKeyBulkSize: "1",
			},
		},
	})
}
