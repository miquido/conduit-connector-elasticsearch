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
