package csjwt

import (
	"bytes"
	"net/http"
	"unicode"

	"github.com/corestoreio/csfw/util/cserr"
	"github.com/juju/errors"
)

// HTTPHeaderAuthorization identifies the bearer token in this header key
const HTTPHeaderAuthorization = `Authorization`

// HTTPFormInputName default name for the HTML form field name
const HTTPFormInputName = `access_token`

// Verification allows to parse and verify a token with custom options.
type Verification struct {
	// FormInputName defines the name of the HTML form input type in which
	// the token has been stored. If empty, the form the gets ignored.
	FormInputName string
	// CookieName defines the name of the cookie where the token has been
	// stored. If empty, cookie parsing gets ignored.
	CookieName string
	// Methods for verifying and signing a token
	Methods SignerSlice

	// Decoder interface to pass in a custom decoder parser.
	// Can be nil, falls back to JSON
	Decoder
}

// NewVerification creates new verification parser with the default signing
// method HS256, if availableSigners slice argument is empty.
// Nil arguments are forbidden.
func NewVerification(availableSigners ...Signer) *Verification {
	if len(availableSigners) == 0 {
		availableSigners = SignerSlice{NewSigningMethodHS256()}
	}
	return &Verification{
		Methods: availableSigners,
		Decoder: JSONDecode{},
	}
}

// Parse parses a rawToken into the template token and returns the fully parsed and
// verified Token, or an error. You must make sure to set the correct expected
// headers and claims in the template Token. The Header and Claims field in the
// template token must be a pointer.
func (vf *Verification) Parse(template Token, rawToken []byte, keyFunc Keyfunc) (Token, error) {
	pos, valid := dotPositions(rawToken)
	if !valid {
		return Token{}, errTokenInvalidSegmentCounts
	}

	if template.Header == nil || template.Claims == nil {
		return template, errTokenBaseNil
	}

	dec := vf.Decoder
	if dec == nil {
		dec = JSONDecode{}
	}

	template.Raw = rawToken

	if startsWithBearer(template.Raw) {
		return template, errTokenShouldNotContainBearer
	}

	// parse Header
	if err := dec.Unmarshal(template.Raw[:pos[0]], template.Header); err != nil {
		return template, cserr.NewMultiErr(ErrTokenMalformed, err)
	}

	// parse Claims
	if err := dec.Unmarshal(template.Raw[pos[0]+1:pos[1]], template.Claims); err != nil {
		return template, cserr.NewMultiErr(ErrTokenMalformed, err)
	}

	// validate Claims
	if err := template.Claims.Valid(); err != nil {
		return template, cserr.NewMultiErr(ErrValidationClaimsInvalid, err)
	}

	// Lookup key
	if keyFunc == nil {
		return template, errMissingKeyFunc
	}
	key, err := keyFunc(template)
	if err != nil {
		return template, cserr.NewMultiErr(ErrTokenUnverifiable, err)
	}

	// Lookup signature method
	method, err := vf.getMethod(&template)
	if err != nil {
		return template, err
	}

	// Perform validation
	template.Signature = template.Raw[pos[1]+1:]
	if err := method.Verify(template.Raw[:pos[1]], template.Signature, key); err != nil {
		return template, cserr.NewMultiErr(ErrSignatureInvalid, err)
	}

	template.Valid = true
	return template, nil
}

func (vf *Verification) getMethod(t *Token) (Signer, error) {

	if len(vf.Methods) == 0 {
		return nil, errors.New("[csjwt] No methods supplied to the Verfication Method slice")
	}

	alg := t.Alg()
	if alg == "" {
		return nil, errors.Errorf("[csjwt] Cannot find alg entry in token header: %#v", t.Header)
	}

	for _, m := range vf.Methods {
		if m.Alg() == alg {
			return m, nil
		}
	}
	return nil, errors.Errorf("[csjwt] Algorithm %q not found in method list %q", alg, vf.Methods)
}

// ParseFromRequest same as Parse but extracts the token from a request.
// First it searches for the token bearer in the header HTTPHeaderAuthorization.
// If not found the request POST form gets parsed and the FormInputName gets
// used to lookup the token value.
func (vf *Verification) ParseFromRequest(template Token, keyFunc Keyfunc, req *http.Request) (Token, error) {
	// Look for an Authorization header
	if ah := req.Header.Get(HTTPHeaderAuthorization); ah != "" {
		// Should be a bearer token
		auth := []byte(ah)
		if startsWithBearer(auth) {
			return vf.Parse(template, auth[7:], keyFunc)
		}
	}

	if vf.CookieName != "" {
		tk, err := vf.parseCookie(template, keyFunc, req)
		if err != nil && err != http.ErrNoCookie {
			return Token{}, errors.Mask(err)
		}
		if tk.Valid {
			return tk, nil
		}
		// try next, the form
	}

	if vf.FormInputName != "" {
		return vf.parseForm(template, keyFunc, req)
	}

	return Token{}, errors.Mask(ErrTokenNotInRequest)
}

func (vf *Verification) parseCookie(template Token, keyFunc Keyfunc, req *http.Request) (Token, error) {
	keks, err := req.Cookie(vf.CookieName)
	if keks != nil && keks.Value != "" {
		return vf.Parse(template, []byte(keks.Value), keyFunc)
	}
	return Token{}, err
}

func (vf *Verification) parseForm(template Token, keyFunc Keyfunc, req *http.Request) (Token, error) {
	_ = req.ParseMultipartForm(10e6) // ignore errors
	if tokStr := req.Form.Get(vf.FormInputName); tokStr != "" {
		return vf.Parse(template, []byte(tokStr), keyFunc)
	}
	return Token{}, errors.Mask(ErrTokenNotInRequest)
}

// SplitForVerify splits the token into two parts: the payload and the signature.
// An error gets returned if the number of dots don't match with the JWT standard.
func SplitForVerify(rawToken []byte) (signingString, signature []byte, err error) {
	pos, valid := dotPositions(rawToken)
	if !valid {
		return nil, nil, errTokenInvalidSegmentCounts
	}
	return rawToken[:pos[1]], rawToken[pos[1]+1:], nil
}

// dotPositions returns the position of the dots within the token slice
// and if the amount of dots are valid for a JWT.
func dotPositions(t []byte) (pos [2]int, valid bool) {
	const aDot byte = '.'
	c := 0
	for i, b := range t {
		if b == aDot {
			if c < 2 {
				pos[c] = i
			}
			c++
		}
	}
	if c == 2 {
		valid = true
	}
	return
}

// length of the string "bearer "
const prefixBearerLen = 7

var prefixBearer = []byte(`bearer `)

// startsWithBearer checks if token starts with bearer
func startsWithBearer(token []byte) bool {
	if len(token) <= prefixBearerLen {
		return false
	}
	var havePrefix [prefixBearerLen]byte
	copy(havePrefix[:], token[0:prefixBearerLen])
	for i, b := range havePrefix {
		havePrefix[i] = byte(unicode.ToLower(rune(b)))
	}
	return bytes.Equal(havePrefix[:], prefixBearer)
}
