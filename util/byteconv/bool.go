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

import "strconv"

var bools = map[string]bool{
	"1": true, "t": true, "T": true, "true": true, "TRUE": true, "True": true, "yes": true, "YES": true,
	"0": false, "f": false, "F": false, "false": false, "FALSE": false, "False": false, "NULL": false, "null": false, "nil": false, "no": false, "NO": false,
}

// ParseBool same as strconv.ParseBool but faster and no allocations.
// Use err == nil to check if a bool value is valid.
func ParseBool(b []byte) (v bool, ok bool, err error) {
	// The only difference between using stdlib or the map access is, that
	// stdlib does one allocation where the string(b) map access has no
	// overhead, but both have the same speed.
	if b == nil {
		return
	}
	if UseStdLib {
		v, err = strconv.ParseBool(string(b))
		ok = err == nil
		return
	}
	if t, ok := bools[string(b)]; ok { // compiler optimizes the byte to string conversion
		return t, ok, nil
	}
	return false, false, syntaxError("ParseBool", string(b))
}
