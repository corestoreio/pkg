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

package byteconv

import (
	"database/sql"
	"testing"

	"github.com/corestoreio/pkg/util/assert"
)

func TestParseNullBool(t *testing.T) {
	UseStdLib = false
	runner := func(have string, want sql.NullBool, wantErr bool) func(*testing.T) {
		return func(t *testing.T) {
			b := sql.RawBytes(have)
			if have == "NULL" {
				b = nil
			}
			bv, ok, err := ParseBool(b)
			if wantErr {
				assert.Error(t, err, "%q", have)
				return
			}
			assert.NoError(t, err, "%s %q", t.Name(), have)
			assert.Exactly(t, want.Valid, ok)
			assert.Exactly(t, want.Bool, bv, t.Name())
		}
	}
	t.Run("NULL is false and invalid", runner("NULL", sql.NullBool{}, false))
	t.Run("empty is false and invalid", runner("", sql.NullBool{}, true))
	t.Run(" is false and invalid", runner("", sql.NullBool{}, true))
	t.Run("£ is false and invalid", runner("£", sql.NullBool{}, true))
	t.Run("0 is false and valid", runner("0", sql.NullBool{Valid: true}, false))
	t.Run("1 is true and valid", runner("1", sql.NullBool{Valid: true, Bool: true}, false))
	t.Run("10 is false and invalid", runner("10", sql.NullBool{}, true))
	t.Run("01 is false and invalid", runner("01", sql.NullBool{}, true))
	t.Run("t is true and valid", runner("t", sql.NullBool{Valid: true, Bool: true}, false))
	t.Run("true is true and valid", runner("true", sql.NullBool{Valid: true, Bool: true}, false))
	t.Run("TRUE is true and valid", runner("TRUE", sql.NullBool{Valid: true, Bool: true}, false))
	t.Run("f is false and valid", runner("f", sql.NullBool{Valid: true, Bool: false}, false))
	t.Run("false is false and valid", runner("false", sql.NullBool{Valid: true, Bool: false}, false))
	t.Run("FALSE is false and valid", runner("FALSE", sql.NullBool{Valid: true, Bool: false}, false))
}

var benchmarkParseBool bool

// BenchmarkParseBool/no-std-map-4         	50000000	        29.7 ns/op	       0 B/op	       0 allocs/op
// BenchmarkParseBool/with-stdlib-4        	50000000	        30.2 ns/op	       4 B/op	       1 allocs/op
func BenchmarkParseBool(b *testing.B) {
	var err error
	tr := true
	true := []byte(`True`)
	b.Run("no-std-map", func(b *testing.B) {
		UseStdLib = false
		for i := 0; i < b.N; i++ {
			benchmarkParseBool, _, err = ParseBool(true)
		}
		if err != nil {
			b.Fatal(err)
		}
	})
	b.Run("with-stdlib", func(b *testing.B) {
		UseStdLib = tr
		for i := 0; i < b.N; i++ {
			benchmarkParseBool, _, err = ParseBool(true)
		}
		if err != nil {
			b.Fatal(err)
		}
	})
}
