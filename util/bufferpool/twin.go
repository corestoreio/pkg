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

package bufferpool

import (
	"bytes"
	"fmt"
	"io"
	"sync"
)

// TODO use github.com/intel-go/bytebuf or check if it can be merged upstream into go-src/master

var twinBufferPool = NewTwin(1024) // estimated *cough* average size

// TwinBuffer contains two buffers.
type TwinBuffer struct {
	First  *bytes.Buffer
	Second *bytes.Buffer
}

// String prints the buffers content for debug purposes.
func (tw *TwinBuffer) String() string {
	return fmt.Sprintf("%q\n%q", tw.First.String(), tw.Second.String())
}

// Write writes to both buffers one after another.
func (tw *TwinBuffer) Write(p []byte) (n int, err error) {
	n, err = tw.First.Write(p)
	if err != nil {
		return
	}
	if n != len(p) {
		return 0, io.ErrShortWrite
	}
	n, err = tw.Second.Write(p)
	return n, err
}

// CopyFirstToSecond resets the second buffer and copies the content from the
// first buffer to the second buffer. First buffer gets eventually reset.
func (tw *TwinBuffer) CopyFirstToSecond() (n int64, err error) {
	tw.Second.Reset()
	n, err = tw.First.WriteTo(tw.Second)
	tw.First.Reset()
	return
}

// CopySecondToFirst resets the first buffer and copies the content from the
// second buffer to the first buffer. Second buffer gets eventually reset.
func (tw *TwinBuffer) CopySecondToFirst() (n int64, err error) {
	tw.First.Reset()
	n, err = tw.Second.WriteTo(tw.First)
	tw.Second.Reset()
	return
}

// Reset resets both buffers.
func (tw *TwinBuffer) Reset() {
	tw.First.Reset()
	tw.Second.Reset()
}

// GetTwin returns a buffer containing two buffers, `First` and `Second`, from the pool.
func GetTwin() *TwinBuffer {
	return twinBufferPool.Get()
}

// PutTwin returns a twin buffer to the pool. The buffers get reset before they
// are put back into circulation.
func PutTwin(buf *TwinBuffer) {
	// @see https://go-review.googlesource.com/c/go/+/136116/4/src/fmt/print.go
	// Proper usage of a sync.Pool requires each entry to have approximately
	// the same memory cost. To obtain this property when the stored type
	// contains a variably-sized buffer, we add a hard limit on the maximum buffer
	// to place back in the pool.
	//
	// See https://golang.org/issue/23199
	const maxSize = 1 << 16 // 64KiB
	if buf.First.Cap() > maxSize || buf.Second.Cap() > maxSize {
		return
	}
	twinBufferPool.Put(buf)
}

// PutTwinCallBack same as PutTwin but executes fn after buf has been returned
// into the pool.
//		buf := twinBuf.Get()
//		defer twinBuf.PutCallBack(buf, wg.Done)
func PutTwinCallBack(buf *TwinBuffer, fn func()) {
	const maxSize = 1 << 16 // 64KiB
	if buf.First.Cap() > maxSize || buf.Second.Cap() > maxSize {
		fn()
		return
	}
	twinBufferPool.PutCallBack(buf, fn)
}

// twinTank implements a sync.Pool for TwinBuffer
type twinTank struct {
	p *sync.Pool
}

// Get returns type safe a buffer
func (t twinTank) Get() *TwinBuffer {
	return t.p.Get().(*TwinBuffer)
}

// Put empties the buffer and returns it back to the pool.
//
//		bp := NewTwin(512)
//		buf := bp.Get()
//		defer bp.Put(buf)
//		// your code
//		return buf.String()
//
// If you use Bytes() function to return bytes make sure you copy the data
// away otherwise your returned byte slice will be empty.
// For using String() no copying is required.
func (t twinTank) Put(buf *TwinBuffer) {
	buf.First.Reset()
	buf.Second.Reset()
	t.p.Put(buf)
}

// PutCallBack same as Put but executes fn after buf has been returned into the
// pool. Good use case when you might have multiple defers in your code.
func (t twinTank) PutCallBack(buf *TwinBuffer, fn func()) {
	buf.First.Reset()
	buf.Second.Reset()
	t.p.Put(buf)
	fn()
}

// NewTwin instantiates a new TwinBuffer pool with a custom pre-allocated
// buffer size. The fields `First` and `Second` will have the same size.
func NewTwin(size int) twinTank {
	return twinTank{
		p: &sync.Pool{
			New: func() interface{} {
				return &TwinBuffer{
					First:  bytes.NewBuffer(make([]byte, 0, size)),
					Second: bytes.NewBuffer(make([]byte, 0, size)),
				}
			},
		},
	}
}
