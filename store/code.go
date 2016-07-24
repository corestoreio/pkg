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

package store

import (
	"github.com/corestoreio/csfw/util"
	"github.com/corestoreio/csfw/util/errors"
)

// StoreCodeMaxLen defines the overall maximum length a store code can have.
const StoreCodeMaxLen = 32

// CodeIsValid checks if a store code is valid. Returns an ErrStoreCodeEmpty
// or an ErrStoreCodeInvalid if the first letter is not a-zA-Z and followed by
// a-zA-Z0-9_ or store code length is greater than 32 characters.
// Error behaviour: NotValid
func CodeIsValid(c string) error {
	// maybe we can weaken that to allow emoji 8-)
	if c == "" || len(c) > StoreCodeMaxLen {
		return errors.NewNotValidf(errStoreCodeInvalid, c)
	}
	c1 := c[0]
	if false == ((c1 >= 'a' && c1 <= 'z') || (c1 >= 'A' && c1 <= 'Z')) {
		return errors.NewNotValidf(errStoreCodeInvalid, c)
	}
	if false == util.StrIsAlNum(c) {
		return errors.NewNotValidf(errStoreCodeInvalid, c)
	}
	return nil
}
