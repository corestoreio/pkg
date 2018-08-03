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

package slices_test

import (
	"testing"

	"github.com/corestoreio/errors"
	"github.com/corestoreio/pkg/util/assert"
	"github.com/corestoreio/pkg/util/slices"
)

func TestInt64Sort(t *testing.T) {

	is := slices.Int64{100, 10, 3, 20, 9, 30, -1}
	assert.EqualValues(t, []int64{-1, 3, 9, 10, 20, 30, 100}, is.Sort().ToInt64())
	assert.EqualValues(t, []int64{100, 30, 20, 10, 9, 3, -1}, is.Reverse().ToInt64())
}

func TestInt64Append(t *testing.T) {

	is := slices.Int64{30, -1}
	is.Append(6)
	assert.EqualValues(t, []int64{30, -1, 6}, is.ToInt64())
}

func TestInt64Update(t *testing.T) {

	is := slices.Int64{-29, 30, -1}
	assert.NoError(t, is.Update(1, 31))
	assert.EqualValues(t, 31, is[1])
	assert.True(t, errors.OutOfRange.Match(is.Update(100, 2)))
}

func TestInt64Delete(t *testing.T) {

	is := slices.Int64{-29, 30, -1}
	assert.NoError(t, is.Delete(1))
	assert.EqualValues(t, []int64{-29, -1}, is.ToInt64())
	assert.True(t, errors.OutOfRange.Match(is.Delete(100)))
}

func TestInt64Index(t *testing.T) {

	is := slices.Int64{-29, 30, -1}
	assert.EqualValues(t, 2, is.Index(-1))
	assert.EqualValues(t, -1, is.Index(123))
}

func TestInt64Contains(t *testing.T) {

	is := slices.Int64{-29, 30, -1}
	assert.True(t, is.Contains(-1))
	assert.False(t, is.Contains(-100))
}

func TestInt64Any(t *testing.T) {

	l := slices.Int64{33, 44, 55, 66}
	assert.True(t, l.Any(func(i int64) bool {
		return i == 44
	}))
	assert.False(t, l.Any(func(i int64) bool {
		return i == 77
	}))
}

func TestInt64All(t *testing.T) {

	af := func(i int64) bool {
		return (i & 1) == 0
	}
	l := slices.Int64{2, 4, 30, 22}
	assert.True(t, l.All(af))
	l.Append(11)
	assert.False(t, l.All(af))
}

func TestInt64Reduce(t *testing.T) {

	af := func(i int64) bool {
		return (i & 1) == 1
	}
	l := slices.Int64{2, 4, 30, 22}
	assert.EqualValues(t, []int64{}, l.Reduce(af).ToInt64())
	l.Append(3, 5)
	assert.EqualValues(t, []int64{3, 5}, l.Reduce(af).ToInt64())
}

func TestInt64Map(t *testing.T) {

	af := func(i int64) int64 {
		return i + 1
	}
	l := slices.Int64{2, 4, 30, 22}
	assert.EqualValues(t, []int64{3, 5, 31, 23}, l.Map(af).ToInt64())
}

func TestInt64Sum(t *testing.T) {

	l := slices.Int64{2, 4, 30, 22}
	assert.EqualValues(t, 58, l.Sum())
}

func TestInt64Unique(t *testing.T) {

	l := slices.Int64{30, 2, 4, 30, 2, 22}
	assert.EqualValues(t, []int64{30, 2, 4, 22}, l.Unique().ToInt64())
}

// 100	  23609788 ns/op	       0 B/op	       0 allocs/op old with O(n^2) with 10k elements with 2 slices
// 1000	   1696117 ns/op	  383012 B/op	     624 allocs/op without len() as 2nd param in make
// 2000	    937479 ns/op	  176523 B/op	     135 allocs/op with len() as 2nd param in make map
func BenchmarkInt64_Unique(b *testing.B) {
	const size = 10000
	input := make(slices.Int64, size)
	for i := range input {
		input[i] = int64(i)
	}
	input[size/2] = size/2 - 1
	var have slices.Int64
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		have = input.Unique()
	}
	if have, want := size-1, len(have); have != want {
		b.Errorf("Have: %v Want: %v", have, want)
	}

}
