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
	"net/http"

	"github.com/corestoreio/pkg/store/scope"
	"github.com/corestoreio/pkg/util"
	"github.com/corestoreio/errors"
)

// CodeFieldName defines the filed name where store code has been saved. Used in
// Cookies and JSON Web Tokens (JWT) to identify an active store besides from
// the default loaded store.
const CodeFieldName = `store`

// CodeURLFieldName name of the GET parameter to set a new store in a current
// website/group context/request.
const CodeURLFieldName = `___store`

// CodeProcessor gets used in the middleware WithRunMode() to extract a
// store code from a Request and modify the response; for example setting
// cookies to persists the selected store.
type CodeProcessor interface {
	// FromRequest returns the valid non-empty store code. Returns an empty
	// store code on all other cases.
	FromRequest(runMode scope.TypeID, req *http.Request) (code string)
	// ProcessDenied gets called in the middleware WithRunMode whenever a store
	// ID isn't allowed to proceed. The variable newStoreID reflects the denied
	// store ID. The ResponseWriter and Request variables can be used for
	// additional information writing and extracting. The error Handler  will
	// always be called.
	ProcessDenied(runMode scope.TypeID, oldStoreID, newStoreID int64, w http.ResponseWriter, r *http.Request)
	// ProcessAllowed enables to adjust the ResponseWriter based on the new
	// store ID. The variable newStoreID contains the new ID, which can also be
	// 0. The code is guaranteed to be not empty, a valid store code, and always
	// points to an existing active store. The ResponseWriter and Request
	// variables can be used for additional information writing and extracting.
	// The next Handler in the chain will after this function be called.
	ProcessAllowed(runMode scope.TypeID, oldStoreID, newStoreID int64, newStoreCode string, w http.ResponseWriter, r *http.Request)
}

// CodeMaxLen defines the overall maximum length a store code can have.
const CodeMaxLen = 32

// CodeIsValid checks if a store code is valid. Returns an ErrStoreCodeEmpty
// or an ErrStoreCodeInvalid if the first letter is not a-zA-Z and followed by
// a-zA-Z0-9_ or store code length is greater than 32 characters.
// Error behaviour: NotValid
func CodeIsValid(c string) error {
	// maybe we can weaken that to allow emoji 8-)
	if c == "" || len(c) > CodeMaxLen {
		return errors.NewNotValidf(errStoreCodeInvalid, c)
	}
	c1 := c[0]
	if !((c1 >= 'a' && c1 <= 'z') || (c1 >= 'A' && c1 <= 'Z')) {
		return errors.NewNotValidf(errStoreCodeInvalid, c)
	}
	if !util.StrIsAlNum(c) {
		return errors.NewNotValidf(errStoreCodeInvalid, c)
	}
	return nil
}
