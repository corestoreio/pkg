// Copyright 2015-present, Cyrill @ Schumacher.fm and the CoreStore contributors
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

package bufferpool_test

import (
	"bytes"
	"sync/atomic"
	"testing"

	"github.com/corestoreio/pkg/sync/bgwork"
	"github.com/corestoreio/pkg/util/assert"
	"github.com/corestoreio/pkg/util/bufferpool"
)

func TestReader(t *testing.T) {
	data := []byte(`// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.`)

	counter := new(int32)

	bgwork.Wait(100, func(_ int) {
		r := bufferpool.GetReader(data)
		defer bufferpool.PutReader(r)
		var buf bytes.Buffer
		r.WriteTo(&buf)
		if !bytes.Equal(data, buf.Bytes()) {
			// *testing.T can't be used here
			panic("Buffer has no equality to data")
		}
		atomic.AddInt32(counter, 1)
	})
	assert.Exactly(t, atomic.LoadInt32(counter), int32(100))
}
