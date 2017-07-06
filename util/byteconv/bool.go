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
	"strconv"
)

// ParseNullBoolSQL parses RawBytes into an integer and on success returns a
// valid bool. The Bool is only true if the underlying value compares equally to
// one.
func ParseNullBoolSQL(b *sql.RawBytes) (val sql.NullBool) {
	b2 := *b
	if len(b2) != 1 {
		return
	}
	i, _ := ParseInt(b2)
	val.Bool = i == 1
	val.Valid = true
	return
}

// ParseBoolSQL parses the underlying bytes into an integer and returns only
// then true if the integer is equal to one.
func ParseBoolSQL(b *sql.RawBytes) bool {
	b2 := *b
	if len(b2) != 1 {
		return false
	}
	i, _ := ParseInt(b2)
	return i == 1
}

func ParseBool(b []byte) (bool, error) {
	if UseStdLib {
		return strconv.ParseBool(string(b))
	}
	return false, nil
}
