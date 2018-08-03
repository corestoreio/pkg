// Copyright 2015-2016, Cyrill @ Schumacher.fm and the CoreStore contributors
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
	"sync"
	"testing"

	"github.com/corestoreio/pkg/util/assert"
	"github.com/corestoreio/pkg/util/bufferpool"
)

var twinBuf = bufferpool.NewTwin(4096)

func TestBufferPoolMultiSize(t *testing.T) {
	t.Parallel()

	const iterations = 30
	var wg sync.WaitGroup
	wg.Add(iterations)
	for i := 0; i < iterations; i++ {
		go func(wg *sync.WaitGroup) {
			buf := twinBuf.Get()
			defer twinBuf.PutCallBack(buf, wg.Done)

			assert.Exactly(t, 4096, buf.First.Cap())
			assert.Exactly(t, 4096, buf.Second.Cap())
			assert.Exactly(t, 0, buf.First.Len())
			assert.Exactly(t, 0, buf.Second.Len())

			have := []byte(`Unless required by applicable law or agreed to in writing, software`)
			n, err := buf.Write(have)
			assert.NoError(t, err)
			assert.Exactly(t, n, len(have))
			assert.Exactly(t, string(have), buf.First.String())
			assert.Exactly(t, string(have), buf.Second.String())
			assert.Exactly(t, "\"Unless required by applicable law or agreed to in writing, software\"\n\"Unless required by applicable law or agreed to in writing, software\"", buf.String())

		}(&wg)
	}
	wg.Wait()
}

func TestTwinBuffer_CopyFirstToSecond(t *testing.T) {
	t.Parallel()

	buf := bufferpool.GetTwin()
	defer bufferpool.PutTwin(buf)
	data := []byte(`// Licensed under the Apache License, Version 2.0 (the "License");`)
	_, err := buf.First.Write(data)
	assert.NoError(t, err, "First.Write should not fail")

	_, err = buf.CopyFirstToSecond()
	assert.NoError(t, err, "CopyFirstToSecond should not fail")

	assert.Exactly(t, string(data), buf.Second.String())
	assert.Exactly(t, "", buf.First.String())
}

func TestTwinBuffer_CopySecondToFirst(t *testing.T) {
	t.Parallel()

	buf := bufferpool.GetTwin()
	defer bufferpool.PutTwin(buf)
	data := []byte(`// Licensed under the Apache License, Version 2.0 (the "License");`)
	_, err := buf.Second.Write(data)
	assert.NoError(t, err, "First.Write should not fail")

	_, err = buf.CopySecondToFirst()
	assert.NoError(t, err, "CopyFirstToSecond should not fail")

	assert.Exactly(t, string(data), buf.First.String())
	assert.Exactly(t, "", buf.Second.String())
}
