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

package net_test

import (
	"testing"

	"github.com/corestoreio/pkg/net"
	"github.com/corestoreio/pkg/util/assert"
	"github.com/stretchr/testify/require"
)

func TestShiftPath(t *testing.T) {
	tests := []struct{ have, head, tail string }{
		{"/", "", "/"},
		{"/contact", "contact", "/"},
		{"/contact/post", "contact", "/post"},
		{"./contact/post", "contact", "/post"},
		{"../contact/post", "contact", "/post"},
		{"../../contact/post", "contact", "/post"},
		{".../../contact/post", "contact", "/post"},
		{"contact/post", "contact", "/post"},
		{"/catalog/product/view/id/123", "catalog", "/product/view/id/123"},
	}

	for i, test := range tests {
		hh, ht := net.ShiftPath(test.have)
		assert.Exactly(t, test.head, hh, "Head Index %d", i)
		require.Exactly(t, test.tail, ht, "Tail Index %d", i)
	}
}

var benchmarkShiftPath string

// BenchmarkShiftPath-4   	10000000	       217 ns/op	      96 B/op	       3 allocs/op
func BenchmarkShiftPath(b *testing.B) {
	for i := 0; i < b.N; i++ {
		benchmarkShiftPath, _ = net.ShiftPath("/catalog/product/view/id/123")
	}
}
