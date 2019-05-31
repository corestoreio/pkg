// Copyright 2015-present, Cyrill @ Schumacher.fm and the CoreStore contributors
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

// Error variables predefined
const (
	errSignatureInvalid         = "[csjwt] signature is invalid: %s\nToken: %#v\n"
	errTokenMalformed           = `[csjwt] token is malformed`
	errTokenUnverifiable        = `[csjwt] token is unverifiable: %s`
	errValidationClaimsInvalid  = `[csjwt] token claims validation failed`
	errVerificationMethodsEmpty = `[csjwt] No methods supplied to the Verfication Method slice`
	errAlgorithmEmpty           = `[csjwt] Cannot find alg entry in token header: %#v`
	errAlgorithmNotFound        = `[csjwt] Algorithm %q not found in method list %q`
)

// Private errors no need to make them public
const (
	errTokenBaseNil                  = `[csjwt] template token header and/or claim are nil`
	errTokenInvalidSegmentCounts     = `[csjwt] token contains an invalid number of segments`
	errMissingKeyFunc                = `[csjwt] Missing KeyFunc`
	errTokenShouldNotContainBearer   = `[csjwt] tokenstring should not contain 'bearer '`
	errKeyEmptyPassword              = "[csjwt] Empty password provided"
	errKeyMissingPassword            = "[csjwt] Missing password to decrypt private key"
	errKeyDecryptPEMBlockFailed      = "[csjwt] Failed to decrypt PEMBlock: %s"
	errKeyParsePKCS8PrivateKeyFailed = "[csjwt] Failed to parse PKCS8PrivateKey: %s"
	errKeyParseCertificateFailed     = "[csjwt] Failed to parse Certificate: %s"
	errKeyMustBePEMEncoded           = "[csjwt] invalid key: Key must be PEM encoded PKCS1 or PKCS8 private key"
	errKeyNonECDSAPublicKey          = "[csjwt] invalid key: Not a valid ECDSA public key"
	errKeyNonRSAPrivateKey           = "[csjwt] invalid key: Not a valid RSA private key"
)

// ErrECDSAVerification sadly this is missing from crypto/ecdsa compared to crypto/rsa
const errECDSAVerification = "crypto/ecdsa: verification error"

const (
	errECDSAPublicKeyEmpty     = `[csjwt] ECDSA Public Key not provided`
	errECDSAPrivateKeyEmpty    = `[csjwt] ECDSA Private Key not provided`
	errECDSAPrivateInvalidBits = `[csjwt] ECDSA Private Key has invalid curve bits`
	errECDSAHashUnavailable    = `[csjwt] ECDSA Hash unavaiable`

	errHmacPasswordEmpty    = `[csjwt] HMAC-SHA Password not provided`
	errHmacHashUnavailable  = `[csjwt] HMAC-SHA Hash unavaiable`
	errHmacSignatureInvalid = `[csjwt] HMAC-SHA Signature invalid`

	errRSAPublicKeyEmpty  = `[csjwt] RSA Public Key not provided`
	errRSAPrivateKeyEmpty = `[csjwt] RSA Private Key not provided`
	errRSAHashUnavailable = `[csjwt] RSA Hash unavaiable`
)
