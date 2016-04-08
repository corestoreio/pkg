package csjwt

import (
	"github.com/corestoreio/csfw/util/cserr"
)

// Error variables predefined
const (
	ErrSignatureInvalid  cserr.Error = `[csjwt] signature is invalid`
	ErrNoTokenInRequest  cserr.Error = `[csjwt] no token present in request`
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
)
