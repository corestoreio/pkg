// Copyright 2015-2017, Cyrill @ Schumacher.fm and the CoreStore contributors
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

package dbr

import (
	"bytes"
	"strconv"
)

const argBytesCap = 8

var bTextNull = []byte(`NULL`)

// argBytes multi dimensional and reusable data structure to encode primitive
// types into byte slices and later write those byte slices somewhere into. The
// reset function resets the internal allocated byte slices to reuse. Using
// different arguments but with the same memory size will not allocate new
// memory.
// 1. dimension: number of unrepeated place holders. e.g.: WHERE colA IN (?)
// 2. dimension number of args for each place holder, after repeat, e.g.: WHERE colA IN (?,?,?,?)
// 3. buffer to write into the argument, always a Go primitive type
type argBytes [][][]byte

func makeArgBytes() argBytes {
	ab := make(argBytes, argBytesCap)
	for i := range ab {
		ab[i] = [][]byte{make([]byte, 0, 8)}
	}
	ab = ab[:0]
	return ab
}

func (ab argBytes) DebugBytes() string {
	var w bytes.Buffer
	for i := range ab {
		if i > 0 {
			w.WriteByte(' ')
		}
		bufiLen := len(ab[i])
		for j := range ab[i] {
			if j == 0 {
				w.WriteString(strconv.Itoa(i))
				w.WriteString(":{")
			}
			if j > 0 {
				w.WriteByte(',')
			}
			w.Write(ab[i][j])
			if j == bufiLen-1 {
				w.WriteByte('}')
			}
		}
	}
	return w.String()
}

func (ab argBytes) growOrNewContainer(sliceCount int) [][]byte {
	n := 1 // grow for one index because we always append only one argument, but that arg might contain n-values
	if l := len(ab); l+n <= cap(ab) {
		ab = ab[:l+n] // grow
		buf := ab[l]  // get the new slice header
		//ab = ab[:l]   // shrink, because of a later append.
		if cap(buf) < sliceCount {
			buf = make([][]byte, sliceCount)
		}
		return buf[:sliceCount]
	}
	return make([][]byte, sliceCount)
}

func (ab argBytes) reset() argBytes {
	for i := range ab {
		for j := range ab[i] {
			ab[i][j] = ab[i][j][:0]
		}
		ab[i] = ab[i][:0]
	}
	return ab[:0]
}

func (ab argBytes) appendNull() argBytes {
	const nullLength = 4
	c := ab.growOrNewContainer(1)
	if cap(c[0]) >= nullLength {
		c[0] = c[0][:nullLength]
		copy(c[0], bTextNull)
	} else {
		c[0] = append(c[0], bTextNull...)
	}
	return append(ab, c)
}

func (ab argBytes) appendInt64(i int64) argBytes {
	c := ab.growOrNewContainer(1)
	c[0] = strconv.AppendInt(c[0], i, 10)
	return append(ab, c)
}

func (ab argBytes) appendInt64s(args ...int64) argBytes {
	c := ab.growOrNewContainer(len(args))
	for i, a := range args {
		c[i] = strconv.AppendInt(c[i], a, 10)
	}
	return append(ab, c)
}
