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
	"errors"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"testing"
	"time"

	sdk "github.com/conduitio/conduit-connector-sdk"
	"github.com/jaswdr/faker"
	"github.com/miquido/conduit-connector-elasticsearch/internal"
	"github.com/stretchr/testify/require"
)

func TestNewDestination(t *testing.T) {
	t.Run("New Destination can be created", func(t *testing.T) {
		require.IsType(t, &Destination{}, NewDestination())
	})
}

func TestDestination_GetClient(t *testing.T) {
	clientMock := &clientMock{}

	destination := Destination{
		client: clientMock,
	}

	require.Same(t, clientMock, destination.GetClient())
}

func TestDestination_Flush(t *testing.T) {
	fakerInstance := faker.New()

	t.Run("Does not perform any action when queue is empty", func(t *testing.T) {
		esClientMock := clientMock{}

		destination := Destination{
			config:          Config{},
			client:          &esClientMock,
			operationsQueue: make(BufferQueue, 0),
		}

		require.NoError(t, destination.Flush(context.Background()))
		require.Len(t, esClientMock.PrepareCreateOperationCalls(), 0)
		require.Len(t, esClientMock.PrepareUpsertOperationCalls(), 0)
		require.Len(t, esClientMock.PrepareDeleteOperationCalls(), 0)
		require.Len(t, esClientMock.BulkCalls(), 0)
	})

	t.Run("Does not retry when there were no failures", func(t *testing.T) {
		var (
			operationMetadata = fakerInstance.Lorem().Sentence(6)
			operationPayload  = fakerInstance.Lorem().Sentence(6)
		)

		esClientMock := clientMock{
			PrepareCreateOperationFunc: func(item sdk.Record) (interface{}, interface{}, error) {
				return operationMetadata, operationPayload, nil
			},

			BulkFunc: func(ctx context.Context, reader io.Reader) (io.ReadCloser, error) {
				bulkRequest, err := io.ReadAll(reader)
				require.NoError(t, err)
				require.Equal(t, fmt.Sprintf("%q\n%q\n", operationMetadata, operationPayload), string(bulkRequest))

				data, err := json.Marshal(bulkResponse{
					Took:   0,
					Errors: false,
					Items: []bulkResponseItems{
						{
							Create: &bulkResponseItem{
								Status: http.StatusOK,
							},
						},
					},
				})
				require.NoError(t, err)

				return io.NopCloser(bytes.NewReader(data)), nil
			},
		}

		destination := Destination{
			config: Config{
				BulkSize: 1,
				Retries:  2,
			},
			client:          &esClientMock,
			operationsQueue: make(BufferQueue, 0),
		}

		destination.operationsQueue.Enqueue(&operation{
			CreatedAt: time.Now(),
			Record: sdk.Record{
				CreatedAt: time.Now(),
				Metadata: map[string]string{
					"action": internal.OperationInsert,
				},
				Payload: sdk.StructuredData{
					"id": fakerInstance.Int32(),
				},
			},
			AckFunc: successfulAckFunc(t),
		})

		require.NoError(t, destination.Flush(context.Background()))
		require.Len(t, esClientMock.PrepareCreateOperationCalls(), 1)
		require.Len(t, esClientMock.PrepareUpsertOperationCalls(), 0)
		require.Len(t, esClientMock.PrepareDeleteOperationCalls(), 0)
		require.Len(t, esClientMock.BulkCalls(), 1)
	})

	t.Run("Fails after failure and 2 retries", func(t *testing.T) {
		var (
			// Succeeds
			record1OperationMetadata = fakerInstance.Lorem().Sentence(6)
			record1OperationPayload  = fakerInstance.Lorem().Sentence(6)
			record1                  = sdk.Record{
				CreatedAt: time.Now().Add(-time.Hour),
				Metadata: map[string]string{
					"action": internal.OperationInsert,
				},
				Payload: sdk.StructuredData{
					"id": fakerInstance.Int32(),
				},
			}

			// Fails for the first time
			record2OperationMetadata = fakerInstance.Lorem().Sentence(6)
			record2OperationPayload  = fakerInstance.Lorem().Sentence(6)
			record2                  = sdk.Record{
				CreatedAt: time.Now(),
				Metadata: map[string]string{
					"action": internal.OperationInsert,
				},
				Payload: sdk.StructuredData{
					"id": fakerInstance.Int32(),
				},
			}

			// Fails for the second time
			record3OperationMetadata = fakerInstance.Lorem().Sentence(6)
			record3OperationPayload  = fakerInstance.Lorem().Sentence(6)
			record3                  = sdk.Record{
				CreatedAt: time.Now().Add(time.Hour),
				Metadata: map[string]string{
					"action": internal.OperationInsert,
				},
				Payload: sdk.StructuredData{
					"id": fakerInstance.Int32(),
				},
			}
		)

		bulkFuncConditionsCounter := 0

		esClientMock := clientMock{
			PrepareCreateOperationFunc: func(item sdk.Record) (interface{}, interface{}, error) {
				switch {
				case recordsAreEqual(item, record1):
					return record1OperationMetadata, record1OperationPayload, nil

				case recordsAreEqual(item, record2):
					return record2OperationMetadata, record2OperationPayload, nil

				case recordsAreEqual(item, record3):
					return record3OperationMetadata, record3OperationPayload, nil
				}

				return nil, nil, errors.New("PrepareCreateOperation: unexpected call")
			},

			BulkFunc: func(ctx context.Context, reader io.Reader) (io.ReadCloser, error) {
				if bulkFuncConditionsCounter > 2 {
					return nil, errors.New("BulkFunc: unexpected call")
				}

				bulkRequest, err := io.ReadAll(reader)
				require.NoError(t, err)

				defer func() { bulkFuncConditionsCounter++ }()

				var data []byte

				switch bulkFuncConditionsCounter {
				case 0:
					// First call: bulk request with all records
					require.Equal(t, fmt.Sprintf(
						"%q\n%q\n%q\n%q\n%q\n%q\n",
						record1OperationMetadata,
						record1OperationPayload,
						record2OperationMetadata,
						record2OperationPayload,
						record3OperationMetadata,
						record3OperationPayload,
					), string(bulkRequest))

					data, err = json.Marshal(bulkResponse{
						Took:   0,
						Errors: false,
						Items: []bulkResponseItems{
							{
								Create: &bulkResponseItem{
									Status: http.StatusOK,
								},
							},
							{
								Create: &bulkResponseItem{
									Status: http.StatusInternalServerError,
								},
							},
							{
								Create: &bulkResponseItem{
									Status: http.StatusInternalServerError,
								},
							},
						},
					})

				case 1:
					// Second call: bulk request with failed records from call 1
					require.Equal(t, fmt.Sprintf(
						"%q\n%q\n%q\n%q\n",
						record2OperationMetadata,
						record2OperationPayload,
						record3OperationMetadata,
						record3OperationPayload,
					), string(bulkRequest))

					data, err = json.Marshal(bulkResponse{
						Took:   0,
						Errors: false,
						Items: []bulkResponseItems{
							{
								Create: &bulkResponseItem{
									Status: http.StatusInternalServerError,
								},
							},
							{
								Create: &bulkResponseItem{
									Status: http.StatusOK,
								},
							},
						},
					})

				case 2:
					// Third call: bulk request with failed records from call 2
					require.Equal(t, fmt.Sprintf(
						"%q\n%q\n",
						record2OperationMetadata,
						record2OperationPayload,
					), string(bulkRequest))

					data, err = json.Marshal(bulkResponse{
						Took:   0,
						Errors: false,
						Items: []bulkResponseItems{
							{
								Create: &bulkResponseItem{
									Status: http.StatusInternalServerError,
									Error: &bulkResponseItemError{
										Type:     "foo",
										Reason:   "bar",
										CausedBy: json.RawMessage(`"baz"`),
									},
								},
							},
						},
					})
				}

				require.NoError(t, err)
				return io.NopCloser(bytes.NewReader(data)), nil
			},
		}

		destination := Destination{
			config: Config{
				BulkSize: 1,
				Retries:  2,
			},
			client:          &esClientMock,
			operationsQueue: make(BufferQueue, 0),
		}

		destination.operationsQueue.Enqueue(&operation{
			CreatedAt: record1.CreatedAt,
			Record:    record1,
			AckFunc:   successfulAckFunc(t),
		})

		destination.operationsQueue.Enqueue(&operation{
			CreatedAt: record2.CreatedAt,
			Record:    record2,
			AckFunc:   unsuccessfulAckFunc(t, fmt.Sprintf("item with key= create failure: [%s] %s: %q", "foo", "bar", "baz")),
		})

		destination.operationsQueue.Enqueue(&operation{
			CreatedAt: record3.CreatedAt,
			Record:    record3,
			AckFunc:   successfulAckFunc(t),
		})

		require.NoError(t, destination.Flush(context.Background()))
		require.Len(t, esClientMock.PrepareCreateOperationCalls(), 3+2+1)
		require.Len(t, esClientMock.PrepareUpsertOperationCalls(), 0)
		require.Len(t, esClientMock.PrepareDeleteOperationCalls(), 0)
		require.Len(t, esClientMock.BulkCalls(), 3)
	})

	t.Run("Succeeds in the second retry", func(t *testing.T) {
		var (
			// Succeeds
			record1OperationMetadata = fakerInstance.Lorem().Sentence(6)
			record1OperationPayload  = fakerInstance.Lorem().Sentence(6)
			record1                  = sdk.Record{
				CreatedAt: time.Now().Add(-time.Hour),
				Metadata: map[string]string{
					"action": internal.OperationInsert,
				},
				Payload: sdk.StructuredData{
					"id": fakerInstance.Int32(),
				},
			}

			// Fails for the first time
			record2OperationMetadata = fakerInstance.Lorem().Sentence(6)
			record2OperationPayload  = fakerInstance.Lorem().Sentence(6)
			record2                  = sdk.Record{
				CreatedAt: time.Now(),
				Metadata: map[string]string{
					"action": internal.OperationInsert,
				},
				Payload: sdk.StructuredData{
					"id": fakerInstance.Int32(),
				},
			}

			// Succeeds in the second retry
			record3OperationMetadata = fakerInstance.Lorem().Sentence(6)
			record3OperationPayload  = fakerInstance.Lorem().Sentence(6)
			record3                  = sdk.Record{
				CreatedAt: time.Now().Add(time.Hour),
				Metadata: map[string]string{
					"action": internal.OperationInsert,
				},
				Payload: sdk.StructuredData{
					"id": fakerInstance.Int32(),
				},
			}
		)

		bulkFuncConditionsCounter := 0

		esClientMock := clientMock{
			PrepareCreateOperationFunc: func(item sdk.Record) (interface{}, interface{}, error) {
				switch {
				case recordsAreEqual(item, record1):
					return record1OperationMetadata, record1OperationPayload, nil

				case recordsAreEqual(item, record2):
					return record2OperationMetadata, record2OperationPayload, nil

				case recordsAreEqual(item, record3):
					return record3OperationMetadata, record3OperationPayload, nil
				}

				return nil, nil, errors.New("PrepareCreateOperation: unexpected call")
			},

			BulkFunc: func(ctx context.Context, reader io.Reader) (io.ReadCloser, error) {
				if bulkFuncConditionsCounter > 2 {
					return nil, errors.New("BulkFunc: unexpected call")
				}

				bulkRequest, err := io.ReadAll(reader)
				require.NoError(t, err)

				defer func() { bulkFuncConditionsCounter++ }()

				var data []byte

				switch bulkFuncConditionsCounter {
				case 0:
					// First call: bulk request with all records
					require.Equal(t, fmt.Sprintf(
						"%q\n%q\n%q\n%q\n%q\n%q\n",
						record1OperationMetadata,
						record1OperationPayload,
						record2OperationMetadata,
						record2OperationPayload,
						record3OperationMetadata,
						record3OperationPayload,
					), string(bulkRequest))

					data, err = json.Marshal(bulkResponse{
						Took:   0,
						Errors: false,
						Items: []bulkResponseItems{
							{
								Create: &bulkResponseItem{
									Status: http.StatusOK,
								},
							},
							{
								Create: &bulkResponseItem{
									Status: http.StatusInternalServerError,
								},
							},
							{
								Create: &bulkResponseItem{
									Status: http.StatusInternalServerError,
								},
							},
						},
					})

				case 1:
					// Second call: bulk request with failed records from call 1
					require.Equal(t, fmt.Sprintf(
						"%q\n%q\n%q\n%q\n",
						record2OperationMetadata,
						record2OperationPayload,
						record3OperationMetadata,
						record3OperationPayload,
					), string(bulkRequest))

					data, err = json.Marshal(bulkResponse{
						Took:   0,
						Errors: false,
						Items: []bulkResponseItems{
							{
								Create: &bulkResponseItem{
									Status: http.StatusInternalServerError,
								},
							},
							{
								Create: &bulkResponseItem{
									Status: http.StatusOK,
								},
							},
						},
					})

				case 2:
					// Third call: bulk request with failed records from call 2
					require.Equal(t, fmt.Sprintf(
						"%q\n%q\n",
						record2OperationMetadata,
						record2OperationPayload,
					), string(bulkRequest))

					data, err = json.Marshal(bulkResponse{
						Took:   0,
						Errors: false,
						Items: []bulkResponseItems{
							{
								Create: &bulkResponseItem{
									Status: http.StatusOK,
								},
							},
						},
					})
				}

				require.NoError(t, err)
				return io.NopCloser(bytes.NewReader(data)), nil
			},
		}

		destination := Destination{
			config: Config{
				BulkSize: 1,
				Retries:  2,
			},
			client:          &esClientMock,
			operationsQueue: make(BufferQueue, 0),
		}

		destination.operationsQueue.Enqueue(&operation{
			CreatedAt: record1.CreatedAt,
			Record:    record1,
			AckFunc:   successfulAckFunc(t),
		})

		destination.operationsQueue.Enqueue(&operation{
			CreatedAt: record2.CreatedAt,
			Record:    record2,
			AckFunc:   successfulAckFunc(t),
		})

		destination.operationsQueue.Enqueue(&operation{
			CreatedAt: record3.CreatedAt,
			Record:    record3,
			AckFunc:   successfulAckFunc(t),
		})

		require.NoError(t, destination.Flush(context.Background()))
		require.Len(t, esClientMock.PrepareCreateOperationCalls(), 3+2+1)
		require.Len(t, esClientMock.PrepareUpsertOperationCalls(), 0)
		require.Len(t, esClientMock.PrepareDeleteOperationCalls(), 0)
		require.Len(t, esClientMock.BulkCalls(), 3)
	})
}

func recordsAreEqual(record1, record2 sdk.Record) bool {
	return reflect.DeepEqual(record1, record2)
}

func successfulAckFunc(t *testing.T) sdk.AckFunc {
	callsCounter := 0

	return func(err error) error {
		if callsCounter > 0 {
			return errors.New("AckFunc: unexpected call")
		}

		callsCounter++

		require.NoError(t, err)

		return nil
	}
}

func unsuccessfulAckFunc(t *testing.T, expectedError string) sdk.AckFunc {
	callsCounter := 0

	return func(err error) error {
		if callsCounter > 0 {
			return errors.New("AckFunc: unexpected call")
		}

		callsCounter++

		require.EqualError(t, err, expectedError)

		return nil
	}
}
