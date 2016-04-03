package csjwt

import (
	"crypto/subtle"

	"github.com/corestoreio/csfw/util/conv"
)

// Claimer for a type to be a Claims object
type Claimer interface {
	// Valid method that determines if the token is invalid for any supported reason.
	// Returns nil on success
	Valid() error
}

// StandardClaims represents a structured version of Claims Section, as
// referenced at https://tools.ietf.org/html/rfc7519#section-4.1
type StandardClaims struct {
	// Audience claim identifies the recipients that the JWT is
	// intended for.  Each principal intended to process the JWT MUST
	// identify itself with a value in the audience claim.  If the principal
	// processing the claim does not identify itself with a value in the
	// "aud" claim when this claim is present, then the JWT MUST be
	// rejected.  In the general case, the "aud" value is an array of case-
	// sensitive strings, each containing a StringOrURI value.  In the
	// special case when the JWT has one audience, the "aud" value MAY be a
	// single case-sensitive string containing a StringOrURI value.  The
	// interpretation of audience values is generally application specific.
	// Use of this claim is OPTIONAL.
	Audience string `json:"aud,omitempty"`
	// ExpiresAt claim identifies the expiration time on
	// or after which the JWT MUST NOT be accepted for processing.  The
	// processing of the "exp" claim requires that the current date/time
	// MUST be before the expiration date/time listed in the "exp" claim.
	// Implementers MAY provide for some small leeway, usually no more than
	// a few minutes, to account for clock skew.  Its value MUST be a number
	// containing a NumericDate value.  Use of this claim is OPTIONAL.
	ExpiresAt int64 `json:"exp,omitempty"`
	// ID claim provides a unique identifier for the JWT.
	// The identifier value MUST be assigned in a manner that ensures that
	// there is a negligible probability that the same value will be
	// accidentally assigned to a different data object; if the application
	// uses multiple issuers, collisions MUST be prevented among values
	// produced by different issuers as well.  The "jti" claim can be used
	// to prevent the JWT from being replayed.  The "jti" value is a case-
	// sensitive string.  Use of this claim is OPTIONAL.
	ID string `json:"jti,omitempty"`
	// IssuedAt claim identifies the time at which the JWT was
	// issued.  This claim can be used to determine the age of the JWT.  Its
	// value MUST be a number containing a NumericDate value.  Use of this
	// claim is OPTIONAL.
	IssuedAt int64 `json:"iat,omitempty"`
	// Issuer claim identifies the principal that issued the
	// JWT.  The processing of this claim is generally application specific.
	// The "iss" value is a case-sensitive string containing a StringOrURI
	// value.  Use of this claim is OPTIONAL.
	Issuer string `json:"iss,omitempty"`
	// NotBefore claim identifies the time before which the JWT
	// MUST NOT be accepted for processing.  The processing of the "nbf"
	// claim requires that the current date/time MUST be after or equal to
	// the not-before date/time listed in the "nbf" claim.  Implementers MAY
	// provide for some small leeway, usually no more than a few minutes, to
	// account for clock skew.  Its value MUST be a number containing a
	// NumericDate value.  Use of this claim is OPTIONAL.
	NotBefore int64 `json:"nbf,omitempty"`
	// Subject claim identifies the principal that is the
	// subject of the JWT.  The claims in a JWT are normally statements
	// about the subject.  The subject value MUST either be scoped to be
	// locally unique in the context of the issuer or be globally unique.
	// The processing of this claim is generally application specific.  The
	// "sub" value is a case-sensitive string containing a StringOrURI
	// value.  Use of this claim is OPTIONAL.
	Subject string `json:"sub,omitempty"`
}

// Valid validates time based claims "exp, iat, nbf".
// There is no accounting for clock skew.
// As well, if any of the above claims are not in the token, it will still
// be considered a valid claim.
func (c StandardClaims) Valid() error {
	vErr := new(ValidationError)
	now := TimeFunc().Unix()

	// The claims below are optional, by default, so if they are set to the
	// default value in Go, let's not fail the verification for them.
	if c.VerifyExpiresAt(now, false) == false {
		vErr.err = "Token is expired"
		vErr.Errors |= ValidationErrorExpired
	}

	if c.VerifyIssuedAt(now, false) == false {
		vErr.err = "Token used before issued, clock skew issue?"
		vErr.Errors |= ValidationErrorIssuedAt
	}

	if c.VerifyNotBefore(now, false) == false {
		vErr.err = "Token is not valid yet"
		vErr.Errors |= ValidationErrorNotValidYet
	}

	if vErr.valid() {
		return nil
	}

	return vErr
}

// VerifyAudience compares the aud claim against cmp.
// If required is false, this method will return true if the value matches or is unset
func (c *StandardClaims) VerifyAudience(cmp string, req bool) bool {
	return verifyConstantTime([]byte(c.Audience), []byte(cmp), req)
}

// VerifyExpiresAt compares the exp claim against cmp.
// If required is false, this method will return true if the value matches or is unset
func (c *StandardClaims) VerifyExpiresAt(cmp int64, req bool) bool {
	return verifyExp(c.ExpiresAt, cmp, req)
}

// VerifyIssuedAt compares the iat claim against cmp.
// If required is false, this method will return true if the value matches or is unset
func (c *StandardClaims) VerifyIssuedAt(cmp int64, req bool) bool {
	return verifyIat(c.IssuedAt, cmp, req)
}

// VerifyIssuer compares the iss claim against cmp.
// If required is false, this method will return true if the value matches or is unset
func (c *StandardClaims) VerifyIssuer(cmp string, req bool) bool {
	return verifyConstantTime([]byte(c.Issuer), []byte(cmp), req)
}

// VerifyNotBefore compares the nbf claim against cmp.
// If required is false, this method will return true if the value matches or is unset
func (c *StandardClaims) VerifyNotBefore(cmp int64, req bool) bool {
	return verifyNbf(c.NotBefore, cmp, req)
}

// MapClaims default type for the Claim field in a token. You can provide
type MapClaims map[string]interface{}

// VerifyAudience compares the aud claim against cmp.
// If required is false, this method will return true if the value matches or is unset
func (m MapClaims) VerifyAudience(cmp string, req bool) bool {
	aud := conv.ToByte(m["aud"])
	return verifyConstantTime(aud, []byte(cmp), req)
}

// Compares the exp claim against cmp.
// If required is false, this method will return true if the value matches or is unset
func (m MapClaims) VerifyExpiresAt(cmp int64, req bool) bool {
	exp, _ := m["exp"].(float64)
	return verifyExp(int64(exp), cmp, req)
}

// Compares the iat claim against cmp.
// If required is false, this method will return true if the value matches or is unset
func (m MapClaims) VerifyIssuedAt(cmp int64, req bool) bool {
	iat, _ := m["iat"].(float64)
	return verifyIat(int64(iat), cmp, req)
}

// Compares the iss claim against cmp.
// If required is false, this method will return true if the value matches or is unset
func (m MapClaims) VerifyIssuer(cmp string, req bool) bool {
	iss := conv.ToByte(m["iss"])
	return verifyConstantTime(iss, []byte(cmp), req)
}

// Compares the nbf claim against cmp.
// If required is false, this method will return true if the value matches or is unset
func (m MapClaims) VerifyNotBefore(cmp int64, req bool) bool {
	nbf, _ := m["nbf"].(float64)
	return verifyNbf(int64(nbf), cmp, req)
}

// Validates time based claims "exp, iat, nbf".
// There is no accounting for clock skew.
// As well, if any of the above claims are not in the token, it will still
// be considered a valid claim.
func (m MapClaims) Valid() error {
	vErr := new(ValidationError)
	now := TimeFunc().Unix()

	if m.VerifyExpiresAt(now, false) == false {
		vErr.err = "Token is expired"
		vErr.Errors |= ValidationErrorExpired
	}

	if m.VerifyIssuedAt(now, false) == false {
		vErr.err = "Token used before issued, clock skew issue?"
		vErr.Errors |= ValidationErrorIssuedAt
	}

	if m.VerifyNotBefore(now, false) == false {
		vErr.err = "Token is not valid yet"
		vErr.Errors |= ValidationErrorNotValidYet
	}

	if vErr.valid() {
		return nil
	}

	return vErr
}

func verifyConstantTime(aud, cmp []byte, required bool) bool {
	if len(aud) == 0 {
		return !required
	}
	return subtle.ConstantTimeCompare(aud, cmp) == 1
}

func verifyExp(exp int64, now int64, required bool) bool {
	if exp == 0 {
		return !required
	}
	return now <= exp
}

func verifyIat(iat int64, now int64, required bool) bool {
	if iat == 0 {
		return !required
	}
	return now >= iat
}

func verifyNbf(nbf int64, now int64, required bool) bool {
	if nbf == 0 {
		return !required
	}
	return now >= nbf
}
