package csjwt

import (
	"bytes"
	"encoding/json"
	"unicode"

	"github.com/corestoreio/csfw/util/cserr"
	"github.com/juju/errors"
)

// Parser allows to parse a token with custom options.
type Parser struct {
	// ValidMethods if populated, only these methods will be considered valid
	ValidMethods []string
	// UseJSONNumber format in JSON decoder
	UseJSONNumber bool
	// JSONer interface to pass in a custom JSON parser.
	// Can be nil
	JSONer
}

// JSONer interface to pass in a custom JSON parser.
// Can be nil in the Parser type.
type JSONer interface {
	Unmarshal(data []byte, v interface{}) error
}

type jsonParser struct {
	useJSONNumber bool
}

func (jp jsonParser) Unmarshal(data []byte, v interface{}) error {
	dec := json.NewDecoder(bytes.NewReader(data))
	if jp.useJSONNumber {
		dec.UseNumber()
	}
	if err := dec.Decode(v); err != nil {
		return cserr.NewMultiErr(ErrTokenMalformed, err)
	}
	return nil
}

// Parse validate, and return a token.
// keyFunc will receive the parsed token and should return the key for validating.
// If everything is kosher, err will be nil
func (p Parser) Parse(rawToken []byte, keyFunc Keyfunc) (Token, error) {
	return p.ParseWithClaims(rawToken, keyFunc, &MapClaims{})
}

// ParseWithClaims same as Parse() but lets you specify your own Claimer.
// Claimer must be a pointer.
func (p Parser) ParseWithClaims(rawToken []byte, keyFunc Keyfunc, claims Claimer) (Token, error) {
	pos, valid := dotPositions(rawToken)
	if !valid {
		return Token{}, errTokenInvalidSegmentCounts
	}

	token := Token{
		Raw:    rawToken,
		Claims: claims,
	}

	if p.JSONer == nil {
		p.JSONer = jsonParser{
			useJSONNumber: p.UseJSONNumber,
		}
	}

	// parse Header
	if headerBytes, err := DecodeSegment(token.Raw[:pos[0]]); err != nil {
		if startsWithBearer(token.Raw) {
			return token, errTokenShouldNotContainBearer
		}
		return token, cserr.NewMultiErr(ErrTokenMalformed, err)
	} else if err := p.JSONer.Unmarshal(headerBytes, &token.Header); err != nil {
		return token, err
	}

	// parse Claims
	if claimBytes, err := DecodeSegment(token.Raw[pos[0]+1 : pos[1]]); err != nil {
		return token, cserr.NewMultiErr(ErrTokenMalformed, err)
	} else {
		if err := p.JSONer.Unmarshal(claimBytes, token.Claims); err != nil {
			return token, err
		}
	}

	// Lookup signature method
	if err := token.updateMethod(); err != nil {
		return token, err
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
			return token, errors.Errorf("Token signing method %s is invalid", alg)
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

	// Perform validation
	token.Signature = token.Raw[pos[1]+1:]
	if err := token.Method.Verify(token.Raw[:pos[1]], token.Signature, key); err != nil {
		return token, cserr.NewMultiErr(ErrSignatureInvalid, err)
	}

	token.Valid = true
	return token, nil
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
