package csjwt

import (
	"errors"
)

// Error variables predefined
var (
	ErrSignatureInvalid            = errors.New("signature is invalid")
	ErrInvalidKey                  = errors.New("key is invalid or of invalid type")
	ErrHashUnavailable             = errors.New("the requested hash function is unavailable")
	ErrNoTokenInRequest            = errors.New("no token present in request")
	ErrTokenMalformed              = errors.New("token is malformed")
	ErrTokenInvalidSegmentCounts   = errors.New("token contains an invalid number of segments")
	ErrTokenShouldNotContainBearer = errors.New("tokenstring should not contain 'bearer '")
	ErrTokenUnverifiable           = errors.New("token is unverifiable")
	ErrMissingKeyFunc              = errors.New("missing KeyFunc")

	ErrValidationUnknownAlg       = errors.New("unknown token signing algorithm")
	ErrValidationExpired          = errors.New("token is expired")
	ErrValidationUsedBeforeIssued = errors.New("token used before issued, clock skew issue?")
	ErrValidationNotValidYet      = errors.New("token is not valid yet")
	ErrValidationAudience         = errors.New("token is not valid for current audience")
	ErrValidationIssuer           = errors.New("token issue validation failed")
	ErrValidationJTI              = errors.New("token JTI validation failed")
	ErrValidationClaimsInvalid    = errors.New("token claims validation failed")
)
