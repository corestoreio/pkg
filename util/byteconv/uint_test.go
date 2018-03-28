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

	"github.com/stretchr/testify/assert"
)

func TestParseUint(t *testing.T) {

	tests := []struct {
		raw     sql.RawBytes
		base    int
		bitSize int
		want    uint64
		wantErr string
	}{
		{sql.RawBytes{0x31, 0x36}, 10, 64, 16, ""},
		{sql.RawBytes{0x31, 0x36}, 10, 8, 16, ""},
		{sql.RawBytes{0x34, 0x38, 0x31}, 10, 64, 481, ""},
		{sql.RawBytes{0x34, 0x38, 0x31}, 10, 8, 0, `strconv.ParseUint: parsing "481": value out of range`},
		{[]byte(`1`), 10, 64, 1, ""},
		{[]byte(`0`), 10, 64, 0, ""},
		{[]byte(``), 10, 64, 0, "strconv.ParseUint: parsing \"\": invalid syntax"},
		{[]byte{}, 10, 64, 0, "strconv.ParseUint: parsing \"\": invalid syntax"},
		{[]byte(`-1`), 10, 64, 0, "strconv.ParseUint: parsing \"-1\": invalid syntax"},
	}

	for i, test := range tests {
		have, ok, err := ParseUint(test.raw, test.base, test.bitSize)
		if test.wantErr != "" {
			assert.EqualError(t, err, test.wantErr, "Index %d", i)
			continue
		} else if err != nil {
			t.Fatalf("Index %d: %s", i, err)
		}
		if !ok {
			t.Fatalf("must be true Index %d: Have:%d Want:%d", i, have, test.want)
		}
		if have != test.want {
			t.Fatalf("Index %d: Have:%d Want:%d", i, have, test.want)
		}
	}
}

var benchmarkParseUint uint64

// BenchmarkParseUint-4   	50000000	        32.1 ns/op	       0 B/op	       0 allocs/op
func BenchmarkParseUint(b *testing.B) {
	data := []byte(`123456789`)
	const want uint64 = 123456789
	for i := 0; i < b.N; i++ {
		var err error
		benchmarkParseUint, _, err = ParseUint(data, 10, 64)
		if err != nil {
			b.Fatal(err)
		}
		if benchmarkParseUint != want {
			b.Fatalf("Have %d Want %d", benchmarkParseUint, want)
		}
	}
}
