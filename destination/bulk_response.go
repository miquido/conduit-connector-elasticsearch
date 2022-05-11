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

import "encoding/json"

type bulkResponse struct {
	Took   int  `json:"took"`
	Errors bool `json:"errors"`
	Items  []struct {
		Update *bulkResponseItem `json:"update,omitempty"`
		Delete *bulkResponseItem `json:"delete,omitempty"`
	} `json:"items"`
}

type bulkResponseItem struct {
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
