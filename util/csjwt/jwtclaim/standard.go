package jwtclaim

import (
	"encoding/json"
	"time"

	"github.com/corestoreio/cspkg/util/conv"
	"github.com/corestoreio/errors"
)

//go:generate ffjson $GOFILE

// Standard represents a structured version of Claims Section, as
// referenced at https://tools.ietf.org/html/rfc7519#section-4.1
// ffjson: noencoder
type Standard struct {
	// TimeSkew duration of time skew we allow between signer and verifier.
	TimeSkew time.Duration `json:"-"`
	// Audience claim identifies the recipients that the JWT is intended for.
	// Each principal intended to process the JWT MUST identify itself with a
	// value in the audience claim.  If the principal processing the claim does
	// not identify itself with a value in the "aud" claim when this claim is
	// present, then the JWT MUST be rejected.  In the general case, the "aud"
	// value is an array of case- sensitive strings, each containing a
	// StringOrURI value.  In the special case when the JWT has one audience,
	// the "aud" value MAY be a single case-sensitive string containing a
	// StringOrURI value.  The interpretation of audience values is generally
	// application specific. Use of this claim is OPTIONAL.
	Audience string `json:"aud,omitempty"`
	// ExpiresAt claim identifies the expiration time on or after which the JWT
	// MUST NOT be accepted for processing.  The processing of the "exp" claim
	// requires that the current date/time MUST be before the expiration
	// date/time listed in the "exp" claim. Implementers MAY provide for some
	// small leeway, usually no more than a few minutes, to account for clock
	// skew.  Its value MUST be a number containing a NumericDate value.  Use of
	// this claim is OPTIONAL.
	ExpiresAt int64 `json:"exp,omitempty"`
	// ID claim provides a unique identifier for the JWT. The identifier value
	// MUST be assigned in a manner that ensures that there is a negligible
	// probability that the same value will be accidentally assigned to a
	// different data object; if the application uses multiple issuers,
	// collisions MUST be prevented among values produced by different issuers
	// as well.  The "jti" claim can be used to prevent the JWT from being
	// replayed.  The "jti" value is a case- sensitive string.  Use of this
	// claim is OPTIONAL.
	ID string `json:"jti,omitempty"`
	// IssuedAt claim identifies the time at which the JWT was issued.  This
	// claim can be used to determine the age of the JWT.  Its value MUST be a
	// number containing a NumericDate value.  Use of this claim is OPTIONAL.
	IssuedAt int64 `json:"iat,omitempty"`
	// Issuer claim identifies the principal that issued the JWT.  The
	// processing of this claim is generally application specific. The "iss"
	// value is a case-sensitive string containing a StringOrURI value.  Use of
	// this claim is OPTIONAL.
	Issuer string `json:"iss,omitempty"`
	// NotBefore claim identifies the time before which the JWT MUST NOT be
	// accepted for processing.  The processing of the "nbf" claim requires that
	// the current date/time MUST be after or equal to the not-before date/time
	// listed in the "nbf" claim.  Implementers MAY provide for some small
	// leeway, usually no more than a few minutes, to account for clock skew.
	// Its value MUST be a number containing a NumericDate value.  Use of this
	// claim is OPTIONAL.
	NotBefore int64 `json:"nbf,omitempty"`
	// Subject claim identifies the principal that is the subject of the JWT.
	// The claims in a JWT are normally statements about the subject.  The
	// subject value MUST either be scoped to be locally unique in the context
	// of the issuer or be globally unique. The processing of this claim is
	// generally application specific.  The "sub" value is a case-sensitive
	// string containing a StringOrURI value.  Use of this claim is OPTIONAL.
	Subject string `json:"sub,omitempty"`
}

// Valid validates time based claims "exp, iat, nbf". There is no accounting for
// clock skew. As well, if any of the above claims are not in the token, it will
// still be considered a valid claim. Error behaviour: NotValid
func (s *Standard) Valid() error {

	now := TimeFunc().Unix()

	// The claims below are optional, by default, so if they are set to the
	// default value in Go, let's not fail the verification for them.

	switch {
	//case s.ExpiresAt == 0 && s.IssuedAt == 0 && s.NotBefore == 0:
	//	return errors.NewNotValidf(`[jwtclaim] token claims validation failed`)

	case !s.VerifyExpiresAt(now, false):
		return errors.NewNotValidf(`[jwtclaim] token is expired %s ago`, TimeFunc().Sub(time.Unix(s.ExpiresAt, 0)))

	case !s.VerifyIssuedAt(now, false):
		return errors.NewNotValidf(`[jwtclaim] token used before issued, clock skew issue? Diff %s`, time.Unix(s.IssuedAt, 0).Sub(TimeFunc()))

	case !s.VerifyNotBefore(now, false):
		return errors.NewNotValidf(`[jwtclaim] token is not valid yet. Diff %s`, time.Unix(s.NotBefore, 0).Sub(TimeFunc()))
	}

	return nil
}

// Set sets a value. Key must be one of the constants Claim*. Error behaviour:
// NotSupported, NotValid
func (s *Standard) Set(key string, value interface{}) (err error) {
	switch key {
	case KeyAudience:
		s.Audience, err = conv.ToStringE(value)
		err = errors.Wrap(err, "[jwtclaim] ToString")
	case KeyExpiresAt:
		s.ExpiresAt, err = conv.ToInt64E(value)
		err = errors.Wrap(err, "[jwtclaim] ToInt64")
	case KeyID:
		s.ID, err = conv.ToStringE(value)
		err = errors.Wrap(err, "[jwtclaim] ToString")
	case KeyIssuedAt:
		s.IssuedAt, err = conv.ToInt64E(value)
		err = errors.Wrap(err, "[jwtclaim] ToInt64")
	case KeyIssuer:
		s.Issuer, err = conv.ToStringE(value)
		err = errors.Wrap(err, "[jwtclaim] ToString")
	case KeyNotBefore:
		s.NotBefore, err = conv.ToInt64E(value)
		err = errors.Wrap(err, "[jwtclaim] ToInt64")
	case KeySubject:
		s.Subject, err = conv.ToStringE(value)
		err = errors.Wrap(err, "[jwtclaim] ToString")
	case KeyTimeSkew:
		s.TimeSkew, err = conv.ToDurationE(value)
		err = errors.Wrap(err, "[jwtclaim] ToDurationE")
	default:
		return errors.NewNotSupportedf(errClaimKeyNotSupported, key)
	}
	return err
}

// Get returns a value or nil or an error. Key must be one of the constants
// Claim*. Error behaviour: NotSupported
func (s *Standard) Get(key string) (interface{}, error) {
	switch key {
	case KeyAudience:
		return s.Audience, nil
	case KeyExpiresAt:
		return s.ExpiresAt, nil
	case KeyID:
		return s.ID, nil
	case KeyIssuedAt:
		return s.IssuedAt, nil
	case KeyIssuer:
		return s.Issuer, nil
	case KeyNotBefore:
		return s.NotBefore, nil
	case KeySubject:
		return s.Subject, nil
	case KeyTimeSkew:
		return s.TimeSkew, nil
	}
	return nil, errors.NewNotSupportedf(errClaimKeyNotSupported, key)
}

// Expires duration when a token expires.
func (s *Standard) Expires() (exp time.Duration) {
	if s.ExpiresAt > 0 {
		tm := time.Unix(s.ExpiresAt, 0)
		if remainer := tm.Sub(time.Now()); remainer > 0 {
			exp = remainer
		}
	}
	return
}

// Keys returns all available keys which this type supports.
func (s *Standard) Keys() []string {
	return allKeys[:7]
}

// VerifyAudience compares the aud claim against cmp. If required is false, this
// method will return true if the value matches or is unset.
func (s *Standard) VerifyAudience(cmp string, req bool) bool {
	return verifyConstantTime([]byte(s.Audience), []byte(cmp), req)
}

// VerifyExpiresAt compares the exp claim against cmp. If required is false,
// this method will return true if the value matches or is unset.
func (s *Standard) VerifyExpiresAt(cmp int64, req bool) bool {
	return verifyExp(s.TimeSkew, s.ExpiresAt, cmp, req)
}

// VerifyIssuedAt compares the iat claim against cmp. If required is false, this
// method will return true if the value matches or is unset.
func (s *Standard) VerifyIssuedAt(cmp int64, req bool) bool {
	return verifyIat(s.IssuedAt, cmp, req)
}

// VerifyIssuer compares the iss claim against cmp. If required is false, this
// method will return true if the value matches or is unset.
func (s *Standard) VerifyIssuer(cmp string, req bool) bool {
	return verifyConstantTime([]byte(s.Issuer), []byte(cmp), req)
}

// VerifyNotBefore compares the nbf claim against cmp. If required is false,
// this method will return true if the value matches or is unset.
func (s *Standard) VerifyNotBefore(cmp int64, req bool) bool {
	return verifyNbf(s.TimeSkew, s.NotBefore, cmp, req)
}

// String human readable output via JSON, slow.
func (s *Standard) String() string {
	b, err := json.Marshal(s)
	if err != nil {
		return errors.NewFatalf("[jwtclaim] Standard.String(): json.Marshal Error: %s", err).Error()
	}
	return string(b)
}
