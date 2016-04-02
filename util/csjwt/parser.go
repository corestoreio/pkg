package csjwt

import (
	"bytes"
	"encoding/json"
	"fmt"
)

// Parser allows to parse a token with custom options.
type Parser struct {
	ValidMethods  []string // If populated, only these methods will be considered valid
	UseJSONNumber bool     // Use JSON Number format in JSON decoder
}

// Parse validate, and return a token.
// keyFunc will receive the parsed token and should return the key for validating.
// If everything is kosher, err will be nil
func (p Parser) Parse(rawToken []byte, keyFunc Keyfunc) (Token, error) {
	return p.ParseWithClaims(rawToken, keyFunc, &MapClaims{})
}

func (p Parser) ParseWithClaims(rawToken []byte, keyFunc Keyfunc, claims Claimer) (Token, error) {
	pos, valid := dotPositions(rawToken)
	if !valid {
		return Token{}, &ValidationError{err: "token contains an invalid number of segments", Errors: ValidationErrorMalformed}
	}

	token := Token{
		Raw:    rawToken,
		Claims: claims,
	}

	// parse Header
	if headerBytes, err := DecodeSegment(token.Raw[:pos[0]]); err != nil {
		if startsWithBearer(token.Raw) {
			return token, &ValidationError{err: "tokenstring should not contain 'bearer '", Errors: ValidationErrorMalformed}
		}
		return token, &ValidationError{err: err.Error(), Errors: ValidationErrorMalformed}
	} else if err := json.Unmarshal(headerBytes, &token.Header); err != nil {
		return token, &ValidationError{err: err.Error(), Errors: ValidationErrorMalformed}
	}

	// parse Claims
	if claimBytes, err := DecodeSegment(token.Raw[pos[0]+1 : pos[1]]); err != nil {
		return token, &ValidationError{err: err.Error(), Errors: ValidationErrorMalformed}
	} else {
		dec := json.NewDecoder(bytes.NewReader(claimBytes))
		if p.UseJSONNumber {
			dec.UseNumber()
		}
		if err := dec.Decode(token.Claims); err != nil {
			return token, &ValidationError{err: err.Error(), Errors: ValidationErrorMalformed}
		}
	}

	// Lookup signature method
	if method, ok := token.Header["alg"].(string); ok {
		if token.Method = GetSigningMethod(method); token.Method == nil {
			return token, &ValidationError{err: "signing method (alg) is unavailable.", Errors: ValidationErrorUnverifiable}
		}
	} else {
		return token, &ValidationError{err: "signing method (alg) is unspecified.", Errors: ValidationErrorUnverifiable}
	}

	// Verify signing method is in the required set
	if p.ValidMethods != nil {
		var signingMethodValid = false
		var alg = token.Method.Alg()
		for _, m := range p.ValidMethods {
			if m == alg {
				signingMethodValid = true
				break
			}
		}
		if !signingMethodValid {
			// signing method is not in the listed set
			return token, &ValidationError{err: fmt.Sprintf("signing method %v is invalid", alg), Errors: ValidationErrorSignatureInvalid}
		}
	}

	// Lookup key
	if keyFunc == nil {
		// keyFunc was not provided.  short circuiting validation
		return token, &ValidationError{err: "no Keyfunc was provided.", Errors: ValidationErrorUnverifiable}
	}
	key, err := keyFunc(token)
	if err != nil {
		// keyFunc returned an error
		return token, &ValidationError{err: err.Error(), Errors: ValidationErrorUnverifiable}
	}

	vErr := &ValidationError{}

	// Validate Claims
	if err := token.Claims.Valid(); err != nil {

		// If the Claims Valid returned an error, check if it is a validation error,
		// If it was another error type, create a ValidationError with a generic ClaimsInvalid flag set
		if e, ok := err.(*ValidationError); !ok {
			vErr = &ValidationError{err: err.Error(), Errors: ValidationErrorClaimsInvalid}
		} else {
			vErr = e
		}
	}

	// Perform validation
	token.Signature = token.Raw[pos[1]+1:]
	if err = token.Method.Verify(token.Raw[:pos[1]], token.Signature, key); err != nil {
		vErr.err = err.Error()
		vErr.Errors |= ValidationErrorSignatureInvalid
	}

	if vErr.valid() {
		token.Valid = true
		return token, nil
	}

	return token, vErr
}

// SplitForVerify splits the token into two parts: the payload and the signature.
// An error gets returned if the number of dots don't match with the JWT standard.
func SplitForVerify(rawToken []byte) (signingString, signature []byte, err error) {
	pos, valid := dotPositions(rawToken)
	if !valid {
		return nil, nil, &ValidationError{err: "token contains an invalid number of segments", Errors: ValidationErrorMalformed}
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

var prefixBearer = []byte(`bearer `)

func startsWithBearer(token []byte) bool {
	if len(token) <= len(prefixBearer) {
		return false
	}
	havePrefix := token[0:len(prefixBearer)]
	havePrefix = bytes.ToLower(havePrefix)
	return bytes.Equal(havePrefix, prefixBearer)
}
