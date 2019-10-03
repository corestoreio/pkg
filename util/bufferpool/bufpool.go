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

package bufferpool

import (
	"bytes"
	"sync"
)

// TODO use github.com/intel-go/bytebuf or check if it can be merged upstream into go-src/master
// TODO: https://github.com/thejerf/gomempool
// TODO: https://github.com/valyala/bytebufferpool/blob/master/pool.go => self calibrating buffer pool

var bufferPool = New(256) // estimated *cough* average size

// Get returns a buffer from the pool.
func Get() *bytes.Buffer {
	return bufferPool.Get()
}

// Put returns a buffer to the pool.
// The buffer is reset before it is put back into circulation.
func Put(buf *bytes.Buffer) {
	bufferPool.Put(buf)
}

// tank implements a sync.Pool for bytes.Buffer
type tank struct {
	p *sync.Pool
}

// Get returns type safe a buffer
func (t tank) Get() *bytes.Buffer {
	return t.p.Get().(*bytes.Buffer)
}

// Put empties the buffer and returns it back to the pool.
//
//		bp := New(320)
//		buf := bp.Get()
//		defer bp.Put(buf)
//		// your code
//		return buf.String()
//
// If you use Bytes() function to return bytes make sure you copy the data
// away otherwise your returned byte slice will be empty.
// For using String() no copying is required.
func (t tank) Put(buf *bytes.Buffer) {
	// @see https://go-review.googlesource.com/c/go/+/136116/4/src/fmt/print.go
	// Proper usage of a sync.Pool requires each entry to have approximately
	// the same memory cost. To obtain this property when the stored type
	// contains a variably-sized buffer, we add a hard limit on the maximum buffer
	// to place back in the pool.
	//
	// See https://golang.org/issue/23199
	const maxSize = 1 << 16 // 64KiB
	if buf.Cap() > maxSize {
		return
	}

	buf.Reset()
	t.p.Put(buf)
}

// New instantiates a new bytes.Buffer pool with a custom
// pre-allocated buffer size.
func New(size int) tank {
	return tank{
		p: &sync.Pool{
			New: func() interface{} {
				return bytes.NewBuffer(make([]byte, 0, size))
			},
		},
	}
}
