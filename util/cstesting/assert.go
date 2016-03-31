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
	"testing"
)

// EqualPointers compares pointers for equality
func EqualPointers(t *testing.T, expected, actual interface{}) bool {
	wantP := reflect.ValueOf(expected)
	haveP := reflect.ValueOf(actual)
	if wantP.Pointer() != haveP.Pointer() {
		t.Errorf("Expecting equal pointers\nWant: %p\nHave: %p", expected, actual)
		return false
	}
	return true
}

func UnEqualPointers(t *testing.T, expected, actual interface{}) bool {
	wantP := reflect.ValueOf(expected)
	haveP := reflect.ValueOf(actual)
	if wantP.Pointer() == haveP.Pointer() {
		t.Errorf("Expecting unequal pointers\nWant: %p\nHave: %p", expected, actual)
		return false
	}
	return true
}
