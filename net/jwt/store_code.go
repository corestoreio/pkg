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

package jwt

import (
	"github.com/corestoreio/csfw/store"
	"github.com/corestoreio/csfw/store/scope"
	"github.com/corestoreio/csfw/util/csjwt"
	"github.com/corestoreio/csfw/util/errors"
)

// ParamName use in Cookie and JWT important when the user selects a different
// store within the current website/group context. This name will be used in a
// cookie or as key value in a token to permanently save the new selected store
// code.
//
// Copied from storenet.ParamName to avoid dependency hell.
const StoreParamName = `store`

// ScopeOptionFromClaim returns a valid store code from a JSON web token or
// ErrStoreNotFound. Please make sure to add the key storenet.ParamName with the
// store code to the token claim.
func ScopeOptionFromClaim(tc csjwt.Claimer) (o scope.Option, err error) {
	err = errors.NewNotFoundf(errStoreNotFound)
	if tc == nil {
		return
	}

	raw, _ := tc.Get(StoreParamName)
	if scopeCode, ok := raw.(string); ok && scopeCode != "" {
		err = store.CodeIsValid(scopeCode)
		if err == nil {
			o, err = scope.SetByCode(scope.Store, scopeCode)
			err = errors.Wrap(err, "[jwt] scope.SetByCode")
		}
	}
	return
}
