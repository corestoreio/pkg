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

// ParseNullBool same as strconv.ParseBool but has no allocations.
func ParseNullBool(b []byte) (val sql.NullBool, err error) {
	if b == nil {
		return
	}
	val.Bool, err = ParseBool(b)
	val.Valid = err == nil
	return
}

var bools = map[string]bool{
	"1": true, "t": true, "T": true, "true": true, "TRUE": true, "True": true,
	"0": false, "f": false, "F": false, "false": false, "FALSE": false, "False": false,
}

// ParseBool same as strconv.ParseBool but faster and no allocations.
func ParseBool(b []byte) (bool, error) {
	// The only difference between using stdlib or the map access is, that
	// stdlib does one allocation where the string(b) map access has no
	// overhead, but both have the same speed.
	if UseStdLib {
		return strconv.ParseBool(string(b))
	}

	if t, ok := bools[string(b)]; ok { // compiler optimizes the byte to string conversion
		return t, nil
	}
	return false, syntaxError("ParseBool", string(b))
}
