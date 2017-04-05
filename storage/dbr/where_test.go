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
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsValidIdentifier(t *testing.T) {
	tests := []struct {
		have string
		want int8
	}{
		{"", 1},
		{"a", 0},
		{"a.", 1},
		{"a.b", 0},
		{".b", 1},
		{"", 2},
		{"Goooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooopher", 1},
		{"Gooooooooooooooooooooooooooooooooooooooooooooooooooooooooooopher", 0},
		{"Gooooooooooooooooooooooooooooooooooooooooooooooooooooooooooopher.Gooooooooooooooooooooooooooooooooooooooooooooooooooooooooooopher", 0},
		{"Gooooooooooooooooooooooooooooooooooooooooooooooooooooooooooopher.Goooooooooooooooooooooooooooooooooooooooooooooooooooooooooo0opher", 1},
		{"Goooooooooooooooooooooooooooooooooooooooooooooooooooooooooopher.Gooooooooooooooooooooooooooooooooooooooooooooooooooooooo0oph€r", 2},
	}
	for i, test := range tests {
		assert.Exactly(t, test.want, isValidIdentifier(test.have), "Index %d with %q", i, test.have)
	}
}

var benchmarkIsValidIdentifier int8

// BenchmarkIsValidIdentifier-4   	20000000	        92.0 ns/op	       0 B/op	       0 allocs/op
func BenchmarkIsValidIdentifier(b *testing.B) {
	for i := 0; i < b.N; i++ {
		benchmarkIsValidIdentifier = isValidIdentifier(`store_owner.catalog_product_entity_varchar`)
	}
	if benchmarkIsValidIdentifier != 0 {
		b.Fatalf("Should be zero but got %d", benchmarkIsValidIdentifier)
	}
}
