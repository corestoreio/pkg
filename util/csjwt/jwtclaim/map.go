package jwtclaim

import (
	"encoding/json"
	"github.com/corestoreio/csfw/util/conv"
	"github.com/corestoreio/csfw/util/cserr"
	"github.com/juju/errors"
	"time"
)

// Map default type for the Claim field in a token. Slowest but
// most flexible type. For speed, use a custom struct type with
// embedding StandardClaims and ffjson generated en-/decoder.
type Map map[string]interface{}

// VerifyAudience compares the aud claim against cmp.
// If required is false, this method will return true if the value matches or is unset
func (m Map) VerifyAudience(cmp string, req bool) bool {
	aud := conv.ToByte(m["aud"])
	return verifyConstantTime(aud, []byte(cmp), req)
}

// Compares the exp claim against cmp.
// If required is false, this method will return true if the value matches or is unset
func (m Map) VerifyExpiresAt(cmp int64, req bool) bool {
	exp := conv.ToInt64(m["exp"])
	return verifyExp(exp, cmp, req)
}

// Compares the iat claim against cmp.
// If required is false, this method will return true if the value matches or is unset
func (m Map) VerifyIssuedAt(cmp int64, req bool) bool {
	iat := conv.ToInt64(m["iat"])
	return verifyIat(iat, cmp, req)
}

// Compares the iss claim against cmp.
// If required is false, this method will return true if the value matches or is unset
func (m Map) VerifyIssuer(cmp string, req bool) bool {
	iss := conv.ToByte(m["iss"])
	return verifyConstantTime(iss, []byte(cmp), req)
}

// Compares the nbf claim against cmp.
// If required is false, this method will return true if the value matches or is unset
func (m Map) VerifyNotBefore(cmp int64, req bool) bool {
	nbf := conv.ToInt64(m["nbf"])
	return verifyNbf(nbf, cmp, req)
}

// Validates time based claims "exp, iat, nbf". There is no accounting for
// clock skew. As well, if any of the above claims are not in the token, it
// will still be considered a valid claim.
func (m Map) Valid() error {
	var vErr *cserr.MultiErr
	now := TimeFunc().Unix()

	if len(m) == 0 {
		return ErrValidationClaimsInvalid
	}

	if m.VerifyExpiresAt(now, false) == false {
		vErr = vErr.AppendErrors(ErrValidationExpired)
	}

	if m.VerifyIssuedAt(now, false) == false {
		vErr = vErr.AppendErrors(ErrValidationUsedBeforeIssued)
	}

	if m.VerifyNotBefore(now, false) == false {
		vErr = vErr.AppendErrors(ErrValidationNotValidYet)
	}

	if vErr.HasErrors() {
		return vErr
	}
	return nil
}

func (m Map) Set(key string, value interface{}) error {
	m[key] = value
	return nil
}

func (m Map) Get(key string) (value interface{}, err error) {
	return m[key], nil
}

func (m Map) Keys() []string {
	keys := make([]string, len(m))
	i := 0
	for k := range m {
		keys[i] = k
		i++
	}
	return keys
}

// Expires duration when a token expires.
func (m Map) Expires() (exp time.Duration) {
	if cexp, ok := m["exp"]; ok {
		fexp := conv.ToFloat64(cexp)
		if fexp > 0.001 {
			tm := time.Unix(int64(fexp), 0)
			if remainer := tm.Sub(time.Now()); remainer > 0 {
				exp = remainer
			}
		}
	}
	return
}

// String human readable output via JSON, slow.
func (m Map) String() string {
	b, err := json.Marshal(m)
	if err != nil {
		return errors.Errorf("[jwtclaim] Map.String(): json.Marshal Error: %s", err).Error()
	}
	return string(b)
}
