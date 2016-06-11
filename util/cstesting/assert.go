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

package cstesting

import (
	"reflect"
)

// ErrorFormater defines the function needed to print out an formatted error.
type errorFormater interface {
	Errorf(format string, args ...interface{})
}

// EqualPointers compares pointers for equality. errorFormater is *testing.T.
func EqualPointers(t errorFormater, expected, actual interface{}) bool {
	wantP := reflect.ValueOf(expected)
	haveP := reflect.ValueOf(actual)
	if wantP.Pointer() != haveP.Pointer() {
		t.Errorf("Expecting equal pointers\nWant: %p\nHave: %p", expected, actual)
		return false
	}
	return true
}
