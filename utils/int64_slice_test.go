// Copyright 2015 CoreStore Authors
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

package utils_test

import (
	"testing"

	"github.com/corestoreio/csfw/utils"
	"github.com/stretchr/testify/assert"
)

func TestInt64SliceSort(t *testing.T) {
	is := utils.Int64Slice{100, 10, 3, 20, 9, 30, -1}
	assert.EqualValues(t, []int64{-1, 3, 9, 10, 20, 30, 100}, is.Sort().ToInt64())
	assert.EqualValues(t, []int64{100, 30, 20, 10, 9, 3, -1}, is.Reverse().ToInt64())
}

func TestInt64SliceAppend(t *testing.T) {
	is := utils.Int64Slice{30, -1}
	is.Append(6)
	assert.EqualValues(t, []int64{30, -1, 6}, is.ToInt64())
}

func TestInt64SliceUpdate(t *testing.T) {
	is := utils.Int64Slice{-29, 30, -1}
	assert.NoError(t, is.Update(1, 31))
	assert.EqualValues(t, 31, is[1])
	assert.EqualError(t, utils.ErrOutOfRange, is.Update(100, 2).Error())
}

func TestInt64SliceDelete(t *testing.T) {
	is := utils.Int64Slice{-29, 30, -1}
	assert.NoError(t, is.Delete(1))
	assert.EqualValues(t, []int64{-29, -1}, is.ToInt64())
	assert.EqualError(t, utils.ErrOutOfRange, is.Delete(100).Error())
}

func TestInt64SliceIndex(t *testing.T) {
	is := utils.Int64Slice{-29, 30, -1}
	assert.EqualValues(t, 2, is.Index(-1))
	assert.EqualValues(t, -1, is.Index(123))
}

func TestInt64SliceInclude(t *testing.T) {
	is := utils.Int64Slice{-29, 30, -1}
	assert.True(t, is.Include(-1))
	assert.False(t, is.Include(-100))
}

func TestInt64SliceAny(t *testing.T) {
	l := utils.Int64Slice{33, 44, 55, 66}
	assert.True(t, l.Any(func(i int64) bool {
		return i == 44
	}))
	assert.False(t, l.Any(func(i int64) bool {
		return i == 77
	}))
}

func TestInt64SliceAll(t *testing.T) {
	af := func(i int64) bool {
		return (i & 1) == 0
	}
	l := utils.Int64Slice{2, 4, 30, 22}
	assert.True(t, l.All(af))
	l.Append(11)
	assert.False(t, l.All(af))
}

func TestInt64SliceReduce(t *testing.T) {
	af := func(i int64) bool {
		return (i & 1) == 1
	}
	l := utils.Int64Slice{2, 4, 30, 22}
	assert.EqualValues(t, []int64{}, l.Reduce(af).ToInt64())
	l.Append(3, 5)
	assert.EqualValues(t, []int64{3, 5}, l.Reduce(af).ToInt64())
}

func TestInt64SliceMap(t *testing.T) {
	af := func(i int64) int64 {
		return i + 1
	}
	l := utils.Int64Slice{2, 4, 30, 22}
	assert.EqualValues(t, []int64{3, 5, 31, 23}, l.Map(af).ToInt64())
}

func TestInt64SliceSum(t *testing.T) {
	l := utils.Int64Slice{2, 4, 30, 22}
	assert.EqualValues(t, 58, l.Sum())
}

func TestInt64SliceUnique(t *testing.T) {
	l := utils.Int64Slice{30, 2, 4, 30, 2, 22}
	assert.EqualValues(t, []int64{30, 2, 4, 22}, l.Unique().ToInt64())
}
