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

package utils

import "strings"

// StrIsAlNum returns true if an alpha numeric string consists of characters a-zA-Z0-9_
func StrIsAlNum(s string) bool {
	c := 0
	for _, r := range s {
		switch {
		case '0' <= r && r <= '9':
			c++
			break
		case 'a' <= r && r <= 'z':
			c++
			break
		case 'A' <= r && r <= 'Z':
			c++
			break
		case r == '_':
			c++
			break
		}
	}
	return len(s) == c
}

// StrContains checks if a string contains one of the multiple provided
// substrings.
func StrContains(s string, substrs ...string) bool {
	for _, subs := range substrs {
		if strings.Contains(s, subs) {
			return true
		}
	}
	return false
}

// StrStartsWith checks if a string starts with one of the multiple provided
// substrings.
func StrStartsWith(s string, substrs ...string) bool {
	for _, subs := range substrs {
		if strings.Index(s, subs) == 0 {
			return true
		}
	}
	return false
}
