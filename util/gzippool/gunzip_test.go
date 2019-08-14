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

package gzippool_test

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"sync/atomic"
	"testing"

	"github.com/corestoreio/pkg/sync/bgwork"
	"github.com/corestoreio/pkg/util/assert"
	"github.com/corestoreio/pkg/util/gzippool"
)

func TestGetReader(t *testing.T) {
	data := []byte(`The quick brown fox jumps over the lazy farmer.`)

	counter := new(int32)

	bgwork.Wait(100, func(_ int) {
		var buf bytes.Buffer

		zw := gzippool.GetWriter(&buf)

		_, err := zw.Write(data)
		if err != nil {
			panic(err)
		}

		gzippool.PutWriter(zw)

		zr := gzippool.GetReader(&buf)
		defer gzippool.PutReader(zr)

		unzipped, err := ioutil.ReadAll(zr)
		if err != nil {
			panic(err)
		}
		if want, have := string(data), string(unzipped); want != have {
			panic(fmt.Sprintf("%q != %q", want, have))
		}
		atomic.AddInt32(counter, 1)
	})
	assert.Exactly(t, atomic.LoadInt32(counter), int32(100))
}
