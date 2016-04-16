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

package csjwt

import (
	"github.com/corestoreio/csfw/util/cserr"
)

// Error variables predefined
const (
	ErrSignatureInvalid        cserr.Error = `[csjwt] signature is invalid`
	ErrTokenNotInRequest       cserr.Error = `[csjwt] token not present in request`
	ErrTokenMalformed          cserr.Error = `[csjwt] token is malformed`
	ErrTokenUnverifiable       cserr.Error = `[csjwt] token is unverifiable`
	ErrValidationClaimsInvalid cserr.Error = `[csjwt] token claims validation failed`
)

// Private errors no need to make them public
const (
	errTokenBaseNil                cserr.Error = `[csjwt] template token header and/or claim are nil`
	errTokenInvalidSegmentCounts   cserr.Error = `[csjwt] token contains an invalid number of segments`
	errMissingKeyFunc              cserr.Error = `[csjwt] Missing KeyFunc`
	errTokenShouldNotContainBearer cserr.Error = `[csjwt] tokenstring should not contain 'bearer '`
	errKeyEmptyPassword            cserr.Error = "[csjwt] Empty password provided"
	errKeyMissingPassword          cserr.Error = "[csjwt] Missing password to decrypt private key"
	errKeyMustBePEMEncoded         cserr.Error = "[csjwt] invalid key: Key must be PEM encoded PKCS1 or PKCS8 private key"
	errKeyNonECDSAPublicKey        cserr.Error = "[csjwt] invalid key: Not a valid ECDSA public key"
	errKeyNonRSAPrivateKey         cserr.Error = "[csjwt] invalid key: Not a valid RSA private key"
)

// ErrECDSAVerification sadly this is missing from crypto/ecdsa compared to crypto/rsa
const ErrECDSAVerification cserr.Error = "crypto/ecdsa: verification error"

const (
	errECDSAPublicKeyEmpty     cserr.Error = `[csjwt] ECDSA Public Key not provided`
	errECDSAPrivateKeyEmpty    cserr.Error = `[csjwt] ECDSA Private Key not provided`
	errECDSAPrivateInvalidBits cserr.Error = `[csjwt] ECDSA Private Key has invalid curve bits`
	errECDSAHashUnavailable    cserr.Error = `[csjwt] ECDSA Hash unavaiable`

	errHmacPasswordEmpty    cserr.Error = `[csjwt] HMAC-SHA Password not provided`
	errHmacHashUnavailable  cserr.Error = `[csjwt] HMAC-SHA Hash unavaiable`
	errHmacSignatureInvalid cserr.Error = `[csjwt] HMAC-SHA Signature invalid`

	errRSAPublicKeyEmpty  cserr.Error = `[csjwt] RSA Public Key not provided`
	errRSAPrivateKeyEmpty cserr.Error = `[csjwt] RSA Private Key not provided`
	errRSAHashUnavailable cserr.Error = `[csjwt] RSA Hash unavaiable`
)
