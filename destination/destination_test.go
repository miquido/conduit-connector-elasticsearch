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
	"context"
	"io"
	"testing"

	sdk "github.com/conduitio/conduit-connector-sdk"
	"github.com/stretchr/testify/require"
)

func TestNewDestination(t *testing.T) {
	t.Run("New Destination can be created", func(t *testing.T) {
		require.IsType(t, &Destination{}, NewDestination())
	})
}

type ClientMock struct {
}

func (c ClientMock) Ping(ctx context.Context) error {
	panic("should not be called")
}

func (c ClientMock) Bulk(ctx context.Context, reader io.Reader) (io.ReadCloser, error) {
	panic("should not be called")
}

func (c ClientMock) PrepareCreateOperation(item sdk.Record) (metadata interface{}, payload interface{}, err error) {
	panic("should not be called")
}

func (c ClientMock) PrepareUpsertOperation(key string, item sdk.Record) (metadata interface{}, payload interface{}, err error) {
	panic("should not be called")
}

func (c ClientMock) PrepareDeleteOperation(key string) (metadata interface{}, err error) {
	panic("should not be called")
}

func TestDestination_GetClient(t *testing.T) {
	clientMock := &ClientMock{}

	destination := Destination{
		client: clientMock,
	}

	require.Same(t, clientMock, destination.GetClient())
}
