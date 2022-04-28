package destination

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"

	sdk "github.com/conduitio/conduit-connector-sdk"
	"github.com/elastic/go-elasticsearch/v7"
)

func NewDestination() sdk.Destination {
	return &Destination{}
}

type Destination struct {
	sdk.UnimplementedDestination

	config         Config
	client         *elasticsearch.Client
	mutex          sync.Mutex
	recordsBuffer  []sdk.Record
	ackFuncsBuffer map[string]sdk.AckFunc
}

func (d *Destination) Configure(_ context.Context, cfgRaw map[string]string) (err error) {
	d.config, err = ParseConfig(cfgRaw)

	return
}

func (d *Destination) Open(_ context.Context) (err error) {
	d.client, err = elasticsearch.NewClient(elasticsearch.Config{
		Addresses: []string{d.config.Host},
		Username:  d.config.Username,
		Password:  d.config.Password,
	})
	if err != nil {
		return fmt.Errorf("connection could not be established: %w", err)
	}

	ping, err := d.client.Ping()
	if err != nil {
		return fmt.Errorf("connection could not be established: %w", err)
	}
	if ping.IsError() {
		return fmt.Errorf("connection could not be established: host ping failed: %s", ping.Status())
	}

	d.mutex = sync.Mutex{}

	// Initializing the buffer
	d.recordsBuffer = make([]sdk.Record, 0, d.config.BulkSize)
	d.ackFuncsBuffer = make(map[string]sdk.AckFunc, d.config.BulkSize)

	return nil
}

func (d *Destination) WriteAsync(ctx context.Context, record sdk.Record, ackFunc sdk.AckFunc) error {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	d.recordsBuffer = append(d.recordsBuffer, record)
	d.ackFuncsBuffer[string(record.Key.Bytes())] = ackFunc

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
	var data bytes.Buffer

	for _, item := range d.recordsBuffer {
		key := string(item.Key.Bytes())

		// actionAndMetadata
		entryMetadata, err := json.Marshal(actionAndMetadata{
			Update: struct {
				Id              string `json:"_id"`
				Index           string `json:"_index"`
				RetryOnConflict int    `json:"retry_on_conflict"`
			}{
				Id:              key,
				Index:           d.config.Index,
				RetryOnConflict: 3,
			},
		})
		if err != nil {
			return fmt.Errorf("failed to prepare metadata with key=%s: %w", key, err)
		}

		data.Write(entryMetadata)
		data.WriteRune('\n')

		// optionalSource
		entrySource, err := json.Marshal(optionalSource{
			Doc:         item.Payload.Bytes(),
			DocAsUpsert: true,
		})
		if err != nil {
			return fmt.Errorf("failed to prepare data with key=%s: %w", key, err)
		}

		data.Write(entrySource)
		data.WriteRune('\n')
	}

	fmt.Println(data.String())

	// Send the bulk request
	result, err := d.client.Bulk(bytes.NewReader(data.Bytes()), d.client.Bulk.WithContext(ctx))
	if err != nil {
		return fmt.Errorf("bulk request failure: %w", err)
	}
	if result.IsError() {
		bodyContents, err := io.ReadAll(result.Body)
		if err != nil {
			return fmt.Errorf("bulk request failure: failed to read the result: %w", err)
		}
		defer result.Body.Close()

		var errorDetails genericError
		if err := json.Unmarshal(bodyContents, &errorDetails); err != nil {
			return fmt.Errorf("bulk request failure: %s", result.Status())
		}

		return fmt.Errorf("bulk request failure: %s: %s", errorDetails.Error.Type, errorDetails.Error.Reason)
	}

	bodyContents, err := io.ReadAll(result.Body)
	if err != nil {
		return fmt.Errorf("bulk response failure: failed to read the result: %w", err)
	}
	defer result.Body.Close()

	fmt.Println(string(bodyContents))

	// Read individual errors
	var response bulkResponse
	if err := json.Unmarshal(bodyContents, &response); err != nil {
		return fmt.Errorf("bulk response failure: could not read the response: %w", err)
	}

	for _, item := range response.Items {
		ackFunc, exists := d.ackFuncsBuffer[item.Update.Id]
		if !exists {
			return fmt.Errorf("bulk response failure: could not ack item with key=%s: no ack function was registered", item.Update.Id)
		}

		if item.Update.Status == http.StatusOK {
			if err := ackFunc(nil); err != nil {
				return err
			}

			continue
		}

		if err := ackFunc(fmt.Errorf(
			"item with key=%s upsert failure: [%s] %s: %s",
			item.Update.Id,
			item.Update.Error.Type,
			item.Update.Error.Reason,
			item.Update.Error.CausedBy,
		)); err != nil {
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

type actionAndMetadata struct {
	Update struct {
		Id              string `json:"_id"`
		Index           string `json:"_index"`
		RetryOnConflict int    `json:"retry_on_conflict"`
	} `json:"update"`
}

type optionalSource struct {
	Doc         json.RawMessage `json:"doc"`
	DocAsUpsert bool            `json:"doc_as_upsert"`
}

type bulkResponse struct {
	Took   int  `json:"took"`
	Errors bool `json:"errors"`
	Items  []struct {
		Update struct {
			Index   string `json:"_index"`
			Type    string `json:"_type"`
			Id      string `json:"_id"`
			Version int    `json:"_version,omitempty"`
			Result  string `json:"result,omitempty"`
			Shards  struct {
				Total      int `json:"total"`
				Successful int `json:"successful"`
				Failed     int `json:"failed"`
			} `json:"_shards,omitempty"`
			SeqNo       int `json:"_seq_no,omitempty"`
			PrimaryTerm int `json:"_primary_term,omitempty"`
			Status      int `json:"status"`
			Error       struct {
				Type     string          `json:"type"`
				Reason   string          `json:"reason"`
				CausedBy json.RawMessage `json:"caused_by"`
			} `json:"error,omitempty"`
		} `json:"update"`
	} `json:"items"`
}

type genericError struct {
	Error struct {
		Type   string `json:"type"`
		Reason string `json:"reason"`
	} `json:"error"`
}
