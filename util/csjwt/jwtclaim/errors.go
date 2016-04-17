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

package jwtclaim

import (
	"github.com/corestoreio/csfw/util/cserr"
)

// Error variables predefined
const (
	ErrValidationUnknownAlg       cserr.Error = `[jwtclaim] unknown token signing algorithm`
	ErrValidationExpired          cserr.Error = `[jwtclaim] token is expired`
	ErrValidationUsedBeforeIssued cserr.Error = `[jwtclaim] token used before issued, clock skew issue?`
	ErrValidationNotValidYet      cserr.Error = `[jwtclaim] token is not valid yet`
	ErrValidationAudience         cserr.Error = `[jwtclaim] token is not valid for current audience`
	ErrValidationIssuer           cserr.Error = `[jwtclaim] token issue validation failed`
	ErrValidationJTI              cserr.Error = `[jwtclaim] token JTI validation failed`
	ErrValidationClaimsInvalid    cserr.Error = `[jwtclaim] token claims validation failed`
)

const (
	errHeaderKeyNotSupported = "[jwtclaim] Header %q not yet supported. Please see constants Header*."
	errClaimKeyNotSupported  = "[jwtclaim] Claim %q not supported."
)
