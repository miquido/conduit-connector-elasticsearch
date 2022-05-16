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
	"testing"
	"time"

	"github.com/jaswdr/faker"
	"github.com/stretchr/testify/require"
)

func TestBufferQueue_Empty(t *testing.T) {
	t.Run("Returns true when queue is empty", func(t *testing.T) {
		queue := BufferQueue{}

		require.True(t, queue.Empty())
	})

	t.Run("Returns false when queue is not empty", func(t *testing.T) {
		queue := BufferQueue{}

		queue = append(queue, &operation{})

		require.False(t, queue.Empty())
	})
}

func TestBufferQueue_Len(t *testing.T) {
	t.Run("Returns 0 when queue is empty", func(t *testing.T) {
		queue := BufferQueue{}

		require.Equal(t, 0, queue.Len())
	})

	t.Run("Returns number of elements when queue is not empty", func(t *testing.T) {
		queue := BufferQueue{}

		queue = append(queue, &operation{})
		queue = append(queue, &operation{})
		queue = append(queue, &operation{})

		require.Equal(t, 3, queue.Len())
	})
}

func TestBufferQueue_Enqueue(t *testing.T) {
	fakerInstance := faker.New()

	t.Run("Appends new item to empty queue", func(t *testing.T) {
		var (
			timeMax    = time.Now()
			timeMin    = timeMax.AddDate(10, 0, 0)
			item1Dummy = &operation{
				CreatedAt: fakerInstance.Time().TimeBetween(timeMin, timeMax),
			}
		)

		queue := BufferQueue{}

		queue.Enqueue(item1Dummy)

		require.Len(t, queue, 1)
		require.Same(t, item1Dummy, queue[0])
	})

	t.Run("Item is appended to the end of the queue", func(t *testing.T) {
		var (
			timeMax    = time.Now()
			timeMin    = timeMax.AddDate(10, 0, 0)
			item1Dummy = &operation{
				CreatedAt: fakerInstance.Time().TimeBetween(timeMin, timeMax),
			}
			item2Dummy = &operation{
				CreatedAt: fakerInstance.Time().TimeBetween(timeMin, timeMax),
			}
			item3Dummy = &operation{
				CreatedAt: fakerInstance.Time().TimeBetween(timeMin, timeMax),
			}
		)

		queue := BufferQueue{}

		queue.Enqueue(item1Dummy)
		require.Len(t, queue, 1)
		require.Same(t, item1Dummy, queue[0])

		queue.Enqueue(item2Dummy)
		require.Len(t, queue, 2)
		require.Same(t, item2Dummy, queue[1])

		queue.Enqueue(item3Dummy)
		require.Len(t, queue, 3)
		require.Same(t, item3Dummy, queue[2])
	})
}

func TestBufferQueue_Sort(t *testing.T) {
	t.Run("Empty queue can be sorted", func(t *testing.T) {
		queue := BufferQueue{}

		require.Empty(t, queue)

		queue.Sort()

		require.Empty(t, queue)
	})

	t.Run("Elements are sorted in ascending order", func(t *testing.T) {
		var (
			now        = time.Now()
			item1Dummy = &operation{
				CreatedAt: now.AddDate(1, 0, 0),
			}
			item2Dummy = &operation{
				CreatedAt: now,
			}
			item3Dummy = &operation{
				CreatedAt: now.AddDate(-1, 0, 0),
			}
			item4Dummy = &operation{
				CreatedAt: now,
			}
		)

		queue := BufferQueue{}

		queue = append(queue, item1Dummy, item2Dummy, item3Dummy, item4Dummy)

		require.Same(t, item1Dummy, queue[0])
		require.Same(t, item2Dummy, queue[1])
		require.Same(t, item3Dummy, queue[2])
		require.Same(t, item4Dummy, queue[3])

		queue.Sort()

		require.Len(t, queue, 4)
		require.Same(t, item3Dummy, queue[0])
		require.Same(t, item2Dummy, queue[1])
		require.Same(t, item4Dummy, queue[2])
		require.Same(t, item1Dummy, queue[3])
	})
}
