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

package cfgsource_test

import (
	"testing"

	"github.com/corestoreio/csfw/config/cfgsource"
	"github.com/corestoreio/csfw/util/errors"
)

func TestMustNewByString(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			err, ok := r.(error)
			if !ok {
				t.Fatal("Expecting an error")
			} else {
				if have, want := errors.IsNotValid(err), true; have != want {
					t.Fatalf("Have %t Want %t", have, want)
				}
			}
		} else {
			t.Fatal("Expecting a panic")
		}
	}()
	_ = cfgsource.MustNewByString("Panic")
}
