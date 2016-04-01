package csjwt

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"time"
)

// HTTPHeaderAuthorization identifies the bearer token in this header key
const HTTPHeaderAuthorization = `Authorization`

// HTTPFormInputName HTML form field name
const HTTPFormInputName = `access_token`

// TimeFunc provides the current time when parsing token to validate "exp" claim (expiration time).
// You can override it to use another time value.  This is useful for testing or if your
// server uses a different time zone than your tokens.
var TimeFunc = time.Now

// Token represents a JWT Token.  Different fields will be used depending on
// whether you're creating or parsing/verifying a token.
type Token struct {
	Raw       []byte                 // The raw token.  Populated when you Parse a token
	Method    Signer                 // The signing method used or to be used
	Header    map[string]interface{} // The first segment of the token
	Claims    map[string]interface{} // The second segment of the token
	Signature []byte                 // The third segment of the token.  Populated when you Parse a token
	Valid     bool                   // Is the token valid?  Populated when you Parse/Verify a token
}

// New creates a new Token. Takes a signing method
func New(method Signer) Token {
	return Token{
		Header: map[string]interface{}{
			"typ": "JWT",
			"alg": method.Alg(),
		},
		Claims: make(map[string]interface{}),
		Method: method,
	}
}

// SignedString gets the complete, signed token.
// Returns a byte slice, save for further processing.
func (t Token) SignedString(key Key) ([]byte, error) {
	sstr, err := t.SigningString()
	if err != nil {
		return nil, err
	}
	sig, err := t.Method.Sign(sstr.Bytes(), key)
	if err != nil {
		return nil, err
	}

	if _, err := sstr.WriteRune('.'); err != nil {
		return nil, err
	}
	if _, err := sstr.Write(sig); err != nil {
		return nil, err
	}
	return sstr.Bytes(), nil
}

// SigningString generates the signing string.  This is the
// most expensive part of the whole deal.  Unless you
// need this for something special, just go straight for
// the SignedString.
// Returns a buffer which can be used for further modifications.
func (t Token) SigningString() (buf bytes.Buffer, err error) {

	var j []byte
	j, err = marshalBase64(t.Header)
	if err != nil {
		return
	}
	if _, err = buf.Write(j); err != nil {
		return
	}
	if _, err = buf.WriteRune('.'); err != nil {
		return
	}
	j, err = marshalBase64(t.Claims)
	if err != nil {
		return
	}
	if _, err = buf.Write(j); err != nil {
		return
	}
	return
}

func marshalBase64(source map[string]interface{}) ([]byte, error) {
	buf := bufPool.Get()
	defer bufPool.Put(buf)
	if err := json.NewEncoder(buf).Encode(source); err != nil {
		return nil, err
	}
	return EncodeSegment(buf.Bytes()), nil
}

// Parse validates and returns a token.
// keyFunc will receive the parsed token and should return the key for validating.
// If everything is kosher, err will be nil
func Parse(tokenString []byte, keyFunc Keyfunc) (Token, error) {
	return Parser{}.Parse(tokenString, keyFunc)
}

// ParseFromRequest tries to find the token in an http.Request.
// This method will call ParseMultipartForm if there's no token in the header.
// Currently, it looks in the Authorization header as well as
// looking for an 'access_token' request parameter in req.Form.
func ParseFromRequest(req *http.Request, keyFunc Keyfunc) (token Token, err error) {

	// Look for an Authorization header
	if ah := req.Header.Get(HTTPHeaderAuthorization); ah != "" {
		// Should be a bearer token
		auth := []byte(ah)
		if startsWithBearer(auth) {
			return Parse(auth[7:], keyFunc)
		}
	}

	// Look for "access_token" parameter
	_ = req.ParseMultipartForm(10e6) // ignore errors
	if tokStr := req.Form.Get(HTTPFormInputName); tokStr != "" {
		return Parse([]byte(tokStr), keyFunc)
	}

	return Token{}, ErrNoTokenInRequest
}

// EncodeSegment encodes JWT specific base64url encoding with padding stripped.
// Returns a new byte slice.
func EncodeSegment(seg []byte) []byte {
	dbuf := make([]byte, base64.RawURLEncoding.EncodedLen(len(seg)))
	base64.RawURLEncoding.Encode(dbuf, seg)
	return dbuf
}

// DecodeSegment decodes JWT specific base64url encoding with padding stripped.
// Returns a new byte slice.
func DecodeSegment(seg []byte) ([]byte, error) {
	dbuf := make([]byte, base64.RawURLEncoding.DecodedLen(len(seg)))
	n, err := base64.RawURLEncoding.Decode(dbuf, seg)
	return dbuf[:n], err
}
