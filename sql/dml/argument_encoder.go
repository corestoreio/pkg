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

package dml

import (
	"bytes"
	"strconv"
	"time"

	"github.com/corestoreio/pkg/storage/null"
)

const argBytesCap = 8

var bTextNullUC = []byte(sqlStrNullUC)

// argEncoded multi dimensional and reusable data structure to encode primitive
// types into byte slices and later write those byte slices somewhere into. The
// reset function resets the internal allocated byte slices to reuse. Using
// different arguments but with the same memory size will not allocate new
// memory.
// 1. dimension: number of unrepeated place holders. e.g.: WHERE colA IN (?)
// 2. dimension number of args for each place holder, after repeatPlaceHolders, e.g.: WHERE colA IN (?,?,?,?)
// 3. buffer to write into the argument, always a Go primitive type
//
// Reusing an already initialized type for different arguments improves
// performance a lot and allocates maybe a little. See benchmarks.
type argEncoded [][][]byte

func makeArgBytes() argEncoded {
	ae := make(argEncoded, argBytesCap)
	for i := range ae {
		ae[i] = [][]byte{make([]byte, 0, 16)}
	}
	ae = ae[:0]
	return ae
}

func (ae argEncoded) DebugBytes() string {
	var w bytes.Buffer
	for i := range ae {
		if i > 0 {
			w.WriteByte(' ')
		}
		bufiLen := len(ae[i])
		for j := range ae[i] {
			if j == 0 {
				w.WriteString(strconv.Itoa(i))
				w.WriteString(":{")
			}
			if j > 0 {
				w.WriteByte(',')
			}
			w.Write(ae[i][j])
			if j == bufiLen-1 {
				w.WriteByte('}')
			}
		}
	}
	return w.String()
}

func (ae argEncoded) growOrNewContainer(sliceCount int) [][]byte {
	n := 1 // grow for one index because we always append only one argument, but that arg might contain n-values
	if l := len(ae); l+n <= cap(ae) {
		ae = ae[:l+n] // grow
		buf := ae[l]  // get the new slice
		//ae = ae[:l]   // shrink, because of a later append.
		if cap(buf) < sliceCount {
			buf = make([][]byte, sliceCount)
		}
		return buf[:sliceCount]
	}
	return make([][]byte, sliceCount)
}

func (ae argEncoded) reset() argEncoded {
	for i := range ae {
		for j := range ae[i] {
			ae[i][j] = ae[i][j][:0]
		}
		ae[i] = ae[i][:0]
	}
	return ae[:0]
}

func (ae argEncoded) appendNull() argEncoded {
	const nullLength = 4
	c := ae.growOrNewContainer(1)
	if cap(c[0]) >= nullLength {
		c[0] = c[0][:nullLength]
		copy(c[0], bTextNullUC)
	} else {
		c[0] = append(c[0], bTextNullUC...)
	}
	return append(ae, c)
}

func (ae argEncoded) appendInt(i int) argEncoded {
	c := ae.growOrNewContainer(1)
	c[0] = strconv.AppendInt(c[0], int64(i), 10)
	return append(ae, c)
}

func (ae argEncoded) appendInts(args ...int) argEncoded {
	c := ae.growOrNewContainer(len(args))
	for i, a := range args {
		c[i] = strconv.AppendInt(c[i], int64(a), 10)
	}
	return append(ae, c)
}

func (ae argEncoded) appendInt64(i int64) argEncoded {
	c := ae.growOrNewContainer(1)
	c[0] = strconv.AppendInt(c[0], i, 10)
	return append(ae, c)
}

func (ae argEncoded) appendInt64s(args ...int64) argEncoded {
	c := ae.growOrNewContainer(len(args))
	for i, a := range args {
		c[i] = strconv.AppendInt(c[i], a, 10)
	}
	return append(ae, c)
}

func (ae argEncoded) appendUint64(i uint64) argEncoded {
	c := ae.growOrNewContainer(1)
	c[0] = strconv.AppendUint(c[0], i, 10)
	return append(ae, c)
}

func (ae argEncoded) appendUint64s(args ...uint64) argEncoded {
	c := ae.growOrNewContainer(len(args))
	for i, a := range args {
		c[i] = strconv.AppendUint(c[i], a, 10)
	}
	return append(ae, c)
}

func (ae argEncoded) appendFloat64(i float64) argEncoded {
	c := ae.growOrNewContainer(1)
	c[0] = strconv.AppendFloat(c[0], i, 'g', -1, 64)
	return append(ae, c)
}

func (ae argEncoded) appendFloat64s(args ...float64) argEncoded {
	c := ae.growOrNewContainer(len(args))
	for i, a := range args {
		c[i] = strconv.AppendFloat(c[i], a, 'g', -1, 64)
	}
	return append(ae, c)
}

func (ae argEncoded) appendString(s string) argEncoded {
	c := ae.growOrNewContainer(1)
	c[0] = append(c[0], []byte(s)...) // todo use copy
	return append(ae, c)
}

func (ae argEncoded) appendStrings(args ...string) argEncoded {
	c := ae.growOrNewContainer(len(args))
	for i, a := range args {
		c[i] = append(c[i], []byte(a)...) // todo use copy
	}
	return append(ae, c)
}

func (ae argEncoded) appendBool(b bool) argEncoded {
	c := ae.growOrNewContainer(1)
	var val byte = '0'
	if b {
		val = '1'
	}
	c[0] = append(c[0], val)
	return append(ae, c)
}

func (ae argEncoded) appendBools(args ...bool) argEncoded {
	c := ae.growOrNewContainer(len(args))
	for i, a := range args {
		var val byte = '0'
		if a {
			val = '1'
		}
		c[i] = append(c[i], val)
	}
	return append(ae, c)
}

func (ae argEncoded) appendTime(t time.Time) argEncoded {
	c := ae.growOrNewContainer(1)
	c[0] = t.AppendFormat(c[0], mysqlTimeFormat)
	return append(ae, c)
}

func (ae argEncoded) appendTimes(args ...time.Time) argEncoded {
	c := ae.growOrNewContainer(len(args))
	for i, a := range args {
		c[i] = a.AppendFormat(c[i], mysqlTimeFormat)
	}
	return append(ae, c)
}

func (ae argEncoded) appendNullString(args ...null.String) argEncoded {
	c := ae.growOrNewContainer(len(args))
	for i, a := range args {
		if a.Valid {
			c[i] = append(c[i], []byte(a.String)...) // todo use copy
		} else {
			c[i] = append(c[i], bTextNullUC...) // todo use copy
		}
	}
	return append(ae, c)
}

func (ae argEncoded) appendNullFloat64(args ...null.Float64) argEncoded {
	c := ae.growOrNewContainer(len(args))
	for i, a := range args {
		if a.Valid {
			c[i] = strconv.AppendFloat(c[i], a.Float64, 'g', -1, 64)
		} else {
			c[i] = append(c[i], bTextNullUC...)
		}
	}
	return append(ae, c)
}

func (ae argEncoded) appendNullInt64(args ...null.Int64) argEncoded {
	c := ae.growOrNewContainer(len(args))
	for i, a := range args {
		if a.Valid {
			c[i] = strconv.AppendInt(c[i], a.Int64, 10)
		} else {
			c[i] = append(c[i], bTextNullUC...)
		}
	}
	return append(ae, c)
}

func (ae argEncoded) appendNullBool(args ...null.Bool) argEncoded {
	c := ae.growOrNewContainer(len(args))
	for i, a := range args {
		if a.Valid {
			var val byte = '0'
			if a.Bool {
				val = '1'
			}
			c[i] = append(c[i], val)
		} else {
			c[i] = append(c[i], bTextNullUC...)
		}
	}
	return append(ae, c)
}

func (ae argEncoded) appendNullTime(args ...null.Time) argEncoded {
	c := ae.growOrNewContainer(len(args))
	for i, a := range args {
		if a.Valid {
			c[i] = a.Time.AppendFormat(c[i], mysqlTimeFormat)
		} else {
			c[i] = append(c[i], bTextNullUC...)
		}
	}
	return append(ae, c)
}
