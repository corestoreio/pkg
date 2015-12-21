// Copyright © 2013-14 Steve Francia <spf@spf13.com>.
//
// Licensed under the Simple Public License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
// http://opensource.org/licenses/Simple-2.0
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

// AverageBufferSize should be adjusted to the average size of a bytes.buffer
// in your application
var AverageBufferSize int = 16

var bufferPool = &sync.Pool{
	New: func() interface{} {
		b := bytes.NewBuffer(make([]byte, AverageBufferSize))
		b.Reset()
		return b
	},
}

// Get returns a buffer from the pool.
func Get() (buf *bytes.Buffer) {
	return bufferPool.Get().(*bytes.Buffer)
}

// Put returns a buffer to the pool.
// The buffer is reset before it is put back into circulation.
func Put(buf *bytes.Buffer) {
	// println(buf.Len()) todo for some statistics
	buf.Reset()
	bufferPool.Put(buf)
}
