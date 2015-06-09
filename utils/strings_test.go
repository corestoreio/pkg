// Copyright 2015, Cyrill @ Schumacher.fm and the CoreStore contributors
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

// IsAlphaNumeric returns true if a string consists of characters a-zA-Z0-9_
func TestIsAlphaNumeric(t *testing.T) {
	tests := []struct {
		have string
		want bool
	}{
		{"Hello World", false},
		{"HelloWorld", true},
		{"Hello1World", true},
		{"Hello0123456789", true},
		{"Hello0123456789€", false},
		{" Hello0123456789", false},
	}

	for _, test := range tests {
		assert.True(t, utils.IsAlphaNumeric(test.have) == test.want, "%#v", test)
	}
}

var benchIsAlphaNumeric bool

// BenchmarkIsAlphaNumeric	10000000	       132 ns/op	       0 B/op	       0 allocs/op
func BenchmarkIsAlphaNumeric(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		benchIsAlphaNumeric = utils.IsAlphaNumeric("Hello1WorldOfGophers")
	}
}
