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

import "github.com/corestoreio/csfw/util/errors"

const (
	errServiceUnsupportedScope         = "[jwt] Service does not support this: %s. Only default or website scope are allowed."
	errTokenParseNotValidOrBlackListed = "[jwt] Token not valid or black listed"
	errScopedConfigNotValid            = `[jwt] ScopedConfig %s is invalid.`
	errUnknownSigningMethod            = "[jwt] Unknown signing method - Have: %q Want: %q"
	errUnknownSigningMethodOptions     = "[jwt] Unknown signing method - Have: %q Want: ES, HS or RS"
	errKeyEmpty                        = "[jwt] Provided key argument is empty"

	// ErrTokenBlacklisted returned by the middleware if the token can be found
	// within the black list.
	errTokenBlacklisted = "[jwt] Token has been black listed"

	errStoreNotFound = "[jwt] Store not found in token claim"
)

var (
	errBlacklistEmptyKID = errors.NewEmptyf("[jwt] Cannot add token to blacklist because JTI / key ID is empty.")
)
