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

package csdb_test

import (
	"strings"
	"testing"

	"github.com/corestoreio/csfw/storage/csdb"
	"github.com/corestoreio/errors"
	"github.com/stretchr/testify/assert"
)

func TestIsValidIdentifier(t *testing.T) {
	t.Parallel()
	const errDummy = errors.Error("Dummy")
	tests := []struct {
		have string
		want error
	}{
		{"$catalog_product_3ntity", nil},
		{"`catalog_product_3ntity", errDummy},
		{"", errDummy},
		{strings.Repeat("a", 64), errDummy},
	}
	for i, test := range tests {
		haveErr := csdb.IsValidIdentifier(test.have)
		if test.want != nil {
			assert.True(t, errors.IsNotValid(haveErr), "Index %d", i)
		} else {
			assert.NoError(t, haveErr, "Index %d", i)
		}
	}
}

var benchmarkIsValidIdentifier error

func BenchmarkIsValidIdentifier(b *testing.B) {
	const id = `$catalog_product_3ntity`
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		benchmarkIsValidIdentifier = csdb.IsValidIdentifier(id)
	}
	if benchmarkIsValidIdentifier != nil {
		b.Fatalf("%+v", benchmarkIsValidIdentifier)
	}
}
