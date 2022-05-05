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
	"fmt"

	v6 "github.com/miquido/conduit-connector-elasticsearch/internal/elasticsearch/v6"
	v7 "github.com/miquido/conduit-connector-elasticsearch/internal/elasticsearch/v7"
	v8 "github.com/miquido/conduit-connector-elasticsearch/internal/elasticsearch/v8"
)

type Version = string

const (
	Version6 Version = "6"
	Version7 Version = "7"
	Version8 Version = "8"
)

// NewClient creates new Elasticsearch client which supports given server version.
// Returns error when provided version is unsupported or client initialization failed.
func NewClient(version Version, config interface{}) (Client, error) {
	switch version {
	case Version6:
		return v6.NewClient(config)

	case Version7:
		return v7.NewClient(config)

	case Version8:
		return v8.NewClient(config)

	default:
		return nil, fmt.Errorf("unsupported version: %s", version)
	}
}
