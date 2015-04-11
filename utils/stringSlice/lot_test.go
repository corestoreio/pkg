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

package stringSlice_test

import (
	"testing"

	"github.com/corestoreio/csfw/utils/stringSlice"
	"github.com/stretchr/testify/assert"
)

func TestFilterContains(t *testing.T) {
	tests := []struct {
		haveSL stringSlice.Lot
		haveIN []string
		want   []string
	}{
		{
			stringSlice.Lot{
				"IFNULL(`scope_table`.`is_visible`, `additional_table`.`is_visible`) AS `is_visible`",
				"IFNULL(`scope_table`.`is_required`, `main_table`.`is_required`) AS `is_required`",
				"IFNULL(`scope_table`.`default_value`, `main_table`.`default_value`) AS `default_value`",
				"IFNULL(`scope_table`.`multiline_count`, `additional_table`.`multiline_count`) AS `multiline_count`",
			},
			[]string{"is_required", "default_value"},
			[]string{
				"IFNULL(`scope_table`.`is_visible`, `additional_table`.`is_visible`) AS `is_visible`",
				"IFNULL(`scope_table`.`multiline_count`, `additional_table`.`multiline_count`) AS `multiline_count`",
			},
		},
	}

	for _, test := range tests {
		test.haveSL.FilterContains(test.haveIN...)
		assert.Equal(t, test.want, test.haveSL.ToString())
	}
}

var benchFilterContains stringSlice.Lot
var benchFilterContainsData = []string{"is_required", "default_value"}

// BenchmarkFilterContains	 1000000	      1841 ns/op	      96 B/op	       2 allocs/op
func BenchmarkFilterContains(b *testing.B) {

	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		l := stringSlice.Lot{
			"IFNULL(`scope_table`.`is_visible`, `additional_table`.`is_visible`) AS `is_visible`",
			"IFNULL(`scope_table`.`is_required`, `main_table`.`is_required`) AS `is_required`",
			"IFNULL(`scope_table`.`default_value`, `main_table`.`default_value`) AS `default_value`",
			"IFNULL(`scope_table`.`multiline_count`, `additional_table`.`multiline_count`) AS `multiline_count`",
		}
		l.FilterContains(benchFilterContainsData...)
		benchFilterContains = l
	}
}

func TestUpdate(t *testing.T) {
	tests := []struct {
		haveSL stringSlice.Lot
		haveD  string
		haveI  int
		err    error
		want   []string
	}{
		{
			haveSL: stringSlice.Lot{
				"IFNULL(`scope_table`.`is_visible`, `additional_table`.`is_visible`) AS `is_visible`",
				"IFNULL(`scope_table`.`is_required`, `main_table`.`is_required`) AS `is_required`",
				"IFNULL(`scope_table`.`default_value`, `main_table`.`default_value`) AS `default_value`",
				"IFNULL(`scope_table`.`multiline_count`, `additional_table`.`multiline_count`) AS `multiline_count`",
			},
			haveD: "default_value",
			haveI: 1,
			err:   nil,
			want: []string{
				"IFNULL(`scope_table`.`is_visible`, `additional_table`.`is_visible`) AS `is_visible`",
				"default_value",
				"IFNULL(`scope_table`.`default_value`, `main_table`.`default_value`) AS `default_value`",
				"IFNULL(`scope_table`.`multiline_count`, `additional_table`.`multiline_count`) AS `multiline_count`",
			},
		},
		{
			haveSL: stringSlice.Lot{
				"IFNULL(`scope_table`.`is_visible`, `additional_table`.`is_visible`) AS `is_visible`",
				"IFNULL(`scope_table`.`is_required`, `main_table`.`is_required`) AS `is_required`",
				"IFNULL(`scope_table`.`default_value`, `main_table`.`default_value`) AS `default_value`",
				"IFNULL(`scope_table`.`multiline_count`, `additional_table`.`multiline_count`) AS `multiline_count`",
			},
			haveD: "default_value",
			haveI: 6,
			err:   stringSlice.ErrOutOfRange,
			want: []string{
				"IFNULL(`scope_table`.`is_visible`, `additional_table`.`is_visible`) AS `is_visible`",
				"IFNULL(`scope_table`.`is_required`, `main_table`.`is_required`) AS `is_required`",
				"IFNULL(`scope_table`.`default_value`, `main_table`.`default_value`) AS `default_value`",
				"IFNULL(`scope_table`.`multiline_count`, `additional_table`.`multiline_count`) AS `multiline_count`",
			},
		},
	}

	for _, test := range tests {
		assert.Equal(t, test.err, test.haveSL.Update(test.haveI, test.haveD))
		assert.Equal(t, test.want, test.haveSL.ToString())
	}
}

func TestLot(t *testing.T) {
	l := stringSlice.Lot{"Maybe", "GoLang", "should", "have", "generics", "but", "who", "needs", "them", "?", ";-)"}
	assert.Len(t, l, l.Len())
	assert.Equal(t, 1, l.Index("GoLang"))
	assert.Equal(t, -1, l.Index("Golang"))
	assert.True(t, l.Include("GoLang"))
	assert.False(t, l.Include("Golang"))

	l2 := stringSlice.Lot{"Maybe", "GoLang"}
	l2.Map(func(s string) string {
		return s + "2"
	})
	assert.Equal(t, []string{"Maybe2", "GoLang2"}, l2.ToString())
	l2.Append("will", "be")
	assert.Equal(t, []string{"Maybe2", "GoLang2", "will", "be"}, l2.ToString())

}

func TestDelete(t *testing.T) {
	l := stringSlice.Lot{"Maybe", "GoLang", "should"}
	assert.NoError(t, l.Delete(1))
	assert.Equal(t, []string{"Maybe", "should"}, l.ToString())
	assert.NoError(t, l.Delete(1))
	assert.Equal(t, []string{"Maybe"}, l.ToString())
	assert.EqualError(t, l.Delete(1), stringSlice.ErrOutOfRange.Error())
}

func TestFilter(t *testing.T) {
	l := stringSlice.Lot{"Maybe", "GoLang", "should"}
	assert.Equal(t, []string{"GoLang"}, l.Filter(func(s string) bool {
		return s == "GoLang"
	}).ToString())
}

func TestUnique(t *testing.T) {
	l := stringSlice.Lot{"Maybe", "GoLang", "GoLang", "GoLang", "or", "or", "RostLang", "RostLang"}
	assert.Equal(t, []string{"Maybe", "GoLang", "or", "RostLang"}, l.Unique().ToString())
}

var benchUnique stringSlice.Lot

// BenchmarkUnique	 2000000	       612 ns/op	     160 B/op	       2 allocs/op
func BenchmarkUnique(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		l := stringSlice.Lot{"Maybe", "GoLang", "GoLang", "GoLang", "or", "or", "RostLang", "RostLang"}
		l.Unique()
		benchUnique = l
	}
}
