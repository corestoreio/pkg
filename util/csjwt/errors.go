package csjwt

import (
	"github.com/corestoreio/csfw/util/cserr"
)

// Error variables predefined
const (
	ErrSignatureInvalid  cserr.Error = `[csjwt] signature is invalid`
	ErrTokenNotInRequest cserr.Error = `[csjwt] token not present in request`
	ErrTokenMalformed    cserr.Error = `[csjwt] token is malformed`
	ErrTokenUnverifiable cserr.Error = `[csjwt] token is unverifiable`

	ErrValidationUnknownAlg       cserr.Error = `[csjwt] unknown token signing algorithm`
	ErrValidationExpired          cserr.Error = `[csjwt] token is expired`
	ErrValidationUsedBeforeIssued cserr.Error = `[csjwt] token used before issued, clock skew issue?`
	ErrValidationNotValidYet      cserr.Error = `[csjwt] token is not valid yet`
	ErrValidationAudience         cserr.Error = `[csjwt] token is not valid for current audience`
	ErrValidationIssuer           cserr.Error = `[csjwt] token issue validation failed`
	ErrValidationJTI              cserr.Error = `[csjwt] token JTI validation failed`
	ErrValidationClaimsInvalid    cserr.Error = `[csjwt] token claims validation failed`
)

// Private errors no need to make them public
const (
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
