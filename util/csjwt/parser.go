package csjwt

import (
	"bytes"
	"encoding/json"
	"unicode"

	"net/http"

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
	// the token has been stored.
	FormInputName string
	// Methods for verifying and signing a token
	Methods []Signer

	// JSONer interface to pass in a custom JSON parser.
	// Can be nil
	JSONer
}

// NewVerification creates new verification parser with the default signing
// method HS256, if availableSigners slice argument is empty.
// Nil arguments are forbidden.
func NewVerification(availableSigners ...Signer) *Verification {
	if len(availableSigners) == 0 {
		availableSigners = []Signer{NewSigningMethodHS256()}
	}
	return &Verification{
		FormInputName: HTTPFormInputName,
		Methods:       availableSigners,
	}
}

// ParseWithMap validate, and return a token. Uses the MapClaims type.
// keyFunc will receive the parsed token and should return the key for validating.
func (vf *Verification) ParseWithMap(rawToken []byte, keyFunc Keyfunc) (Token, error) {
	return vf.ParseWithClaim(rawToken, keyFunc, &MapClaims{})
}

// ParseWithClaims same as Parse() but lets you specify your own Claimer.
// Claimer must be a pointer.
func (vf *Verification) ParseWithClaim(rawToken []byte, keyFunc Keyfunc, claims Claimer) (Token, error) {
	pos, valid := dotPositions(rawToken)
	if !valid {
		return Token{}, errTokenInvalidSegmentCounts
	}

	token := Token{
		Raw:    rawToken,
		Claims: claims,
	}

	js := vf.JSONer
	if js == nil {
		js = JSONDecoder{}
	}

	// parse Header
	if headerBytes, err := DecodeSegment(token.Raw[:pos[0]]); err != nil {
		if startsWithBearer(token.Raw) {
			return token, errTokenShouldNotContainBearer
		}
		return token, cserr.NewMultiErr(ErrTokenMalformed, err)
	} else if err := js.Unmarshal(headerBytes, &token.Header); err != nil {
		return token, err
	}

	// parse Claims
	if claimBytes, err := DecodeSegment(token.Raw[pos[0]+1 : pos[1]]); err != nil {
		return token, cserr.NewMultiErr(ErrTokenMalformed, err)
	} else {
		if err := js.Unmarshal(claimBytes, token.Claims); err != nil {
			return token, err
		}
	}

	// Validate Claims
	if err := token.Claims.Valid(); err != nil {
		return token, cserr.NewMultiErr(ErrValidationClaimsInvalid, err)
	}

	// Lookup key
	if keyFunc == nil {
		return token, errMissingKeyFunc
	}
	key, err := keyFunc(token)
	if err != nil {
		return token, cserr.NewMultiErr(ErrTokenUnverifiable, err)
	}

	// Lookup signature method
	method, err := vf.getMethod(&token)
	if err != nil {
		return token, err
	}

	// Perform validation
	token.Signature = token.Raw[pos[1]+1:]
	if err := method.Verify(token.Raw[:pos[1]], token.Signature, key); err != nil {
		return token, cserr.NewMultiErr(ErrSignatureInvalid, err)
	}

	token.Valid = true
	return token, nil
}

func (vf *Verification) getMethod(token *Token) (Signer, error) {

	if len(vf.Methods) == 0 {
		return nil, errors.New("[csjwt] No methods supplied to the Verfication Method slice")
	}

	alg := token.Alg()
	if alg == "" {
		return nil, errors.Errorf("[csjwt] Cannot find alg entry in token header: %#v", token.Header)
	}

	for _, m := range vf.Methods {
		if m.Alg() == alg {
			return m, nil
		}
	}
	return nil, errors.Errorf("[csjwt] Algorithm %q not found in method list %q", alg, methods(vf.Methods))
}

// ParseFromRequest same as ParseFromRequest but allows to add a custer Claimer.
// Claimer must be a pointer.
func (vf *Verification) ParseFromRequest(req *http.Request, keyFunc Keyfunc, claims Claimer) (Token, error) {
	// Look for an Authorization header
	if ah := req.Header.Get(HTTPHeaderAuthorization); ah != "" {
		// Should be a bearer token
		auth := []byte(ah)
		if startsWithBearer(auth) {
			return vf.ParseWithClaim(auth[7:], keyFunc, claims)
		}
	}

	// Look for "access_token" parameter
	_ = req.ParseMultipartForm(10e6) // ignore errors
	if tokStr := req.Form.Get(vf.FormInputName); tokStr != "" {
		return vf.ParseWithClaim([]byte(tokStr), keyFunc, claims)
	}

	return Token{}, ErrTokenNotInRequest
}

// JSONer interface to pass in a custom JSON parser.
// Can be nil in the Parser type.
type JSONer interface {
	Unmarshal(data []byte, v interface{}) error
}

// JSONDecoder default JSON decoder
type JSONDecoder struct {
	UseJSONNumber bool
}

func (jp JSONDecoder) Unmarshal(data []byte, v interface{}) error {
	dec := json.NewDecoder(bytes.NewReader(data))
	if jp.UseJSONNumber {
		dec.UseNumber()
	}
	if err := dec.Decode(v); err != nil {
		return cserr.NewMultiErr(ErrTokenMalformed, err)
	}
	return nil
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
