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

package log

import (
	"bytes"
	"io"
	"sync"
)

var _ io.Writer = (*MutexBuffer)(nil)

// MutexBuffer allows concurrent and parallel writes to a buffer. Mostly used
// during testing when the logger should be able to accept multiple writes.
type MutexBuffer struct {
	mu  sync.Mutex
	buf bytes.Buffer
}

// Write writes to a buffer with an acquired lock
func (pl *MutexBuffer) Write(p []byte) (n int, err error) {
	pl.mu.Lock()
	defer pl.mu.Unlock()
	return pl.buf.Write(p)
}

// String reads from the buffer and returns a string.
func (pl *MutexBuffer) String() string {
	pl.mu.Lock()
	defer pl.mu.Unlock()
	return pl.buf.String()
}

// Bytes reads from the buffer and returns the bytes
func (pl *MutexBuffer) Bytes() []byte {
	pl.mu.Lock()
	defer pl.mu.Unlock()
	return pl.buf.Bytes()
}

// Reset truncates the buffer to zero length
func (pl *MutexBuffer) Reset() {
	pl.mu.Lock()
	pl.buf.Reset()
	pl.mu.Unlock()
}
