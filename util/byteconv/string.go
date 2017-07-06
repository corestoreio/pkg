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

import "database/sql"

// ParseNullStringSQL converts a maybe nil byte slice into the appropriate valid
// SQL type.
func ParseNullStringSQL(b *sql.RawBytes) (ns sql.NullString) {
	if b == nil {
		return
	}
	b2 := *b
	if b2 == nil {
		return
	}
	ns.Valid = true
	ns.String = string(b2)
	return
}

// ParseStringSQL casts the byte slice into a string.
func ParseStringSQL(b *sql.RawBytes) string {
	b2 := *b
	if len(b2) == 0 {
		return ""
	}
	return string(b2)
}
