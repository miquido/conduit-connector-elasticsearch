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

package elasticsearch

import (
	sdk "github.com/conduitio/conduit-connector-sdk"
)

func Specification() sdk.Specification {
	return sdk.Specification{
		Name:    "elasticsearch",
		Summary: "An Elasticsearch destination plugin for Conduit.",
		Version: "v0.1.0",
		Author:  "Miquido",
		DestinationParams: map[string]sdk.Parameter{
			"version": {
				Default:     "",
				Required:    false,
				Description: "The version of the Elasticsearch service.",
			},
			"host": {
				Default:     "",
				Required:    true,
				Description: "The Elasticsearh host and port (e.g.: http://127.0.0.1:9200).",
			},
			"username": {
				Default:     "",
				Required:    false,
				Description: "The username for HTTP Basic Authentication.",
			},
			"password": {
				Default:     "",
				Required:    false,
				Description: "The password for HTTP Basic Authentication.",
			},
			"cloudId": {
				Default:     "",
				Required:    false,
				Description: "Endpoint for the Elastic Service (https://elastic.co/cloud).",
			},
			"apiKey": {
				Default:     "",
				Required:    false,
				Description: "Base64-encoded token for authorization; if set, overrides username/password and service token.",
			},
			"serviceToken": {
				Default:     "",
				Required:    false,
				Description: "Service token for authorization; if set, overrides username/password.",
			},
			"certificateFingerprint": {
				Default:     "",
				Required:    false,
				Description: "SHA256 hex fingerprint given by Elasticsearch on first launch.",
			},
			"index": {
				Default:     "",
				Required:    true,
				Description: "The name of the index to write the data to.",
			},
			"bulkSize": {
				Default:     "1000",
				Required:    true,
				Description: "The maximum size of operations sent to Elasticsearch server.",
			},
		},
		SourceParams: map[string]sdk.Parameter{
			//
		},
	}
}
