package csjwt

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"time"

	"github.com/corestoreio/csfw/storage/text"
)

// ContentTypeJWT defines the content type of a token. At the moment only JWT
// is supported. JWE may be added in the future JSON Web Encryption (JWE).
// https://tools.ietf.org/html/rfc7519
const ContentTypeJWT = `JWT`

// TimeFunc provides the current time when parsing token to validate "exp" claim (expiration time).
// You can override it to use another time value.  This is useful for testing or if your
// server uses a different time zone than your tokens.
var TimeFunc = time.Now

// Token represents a JWT Token.  Different fields will be used depending on
// whether you're creating or parsing/verifying a token.
type Token struct {
	Raw       text.Chars             // The raw token.  Populated when you Parse a token
	Header    map[string]interface{} // The first segment of the token
	Claims    Claimer                // The second segment of the token
	Signature text.Chars             // The third segment of the token.  Populated when you Parse a token
	Valid     bool                   // Is the token valid?  Populated when you Parse/Verify a token
}

// NewToken creates a new Token and presets the header to typ = JWT.
// A new token has not yet an assigned algorithm.
func NewToken(c Claimer) Token {
	return Token{
		Header: map[string]interface{}{
			"typ": ContentTypeJWT,
		},
		Claims: c,
	}
}

// Alg returns the assigned algorithm to this token.
// Can return an empty string.
func (t Token) Alg() string {
	algRaw, ok := t.Header["alg"]
	if !ok {
		return ""
	}
	if a, ok := algRaw.(string); ok {
		return a
	}
	return ""
}

// SignedString gets the complete, signed token.
// Sets the header alg to the provided Signer.Alg() value.
// Returns a byte slice, save for further processing.
// This functions allows to sign a token with different signing methods.
func (t Token) SignedString(method Signer, key Key) (text.Chars, error) {

	t.Header["alg"] = method.Alg()

	buf, err := t.SigningString()
	if err != nil {
		return nil, err
	}
	sig, err := method.Sign(buf.Bytes(), key)
	if err != nil {
		return nil, err
	}

	if _, err := buf.WriteRune('.'); err != nil {
		return nil, err
	}
	if _, err := buf.Write(sig); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
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

func marshalBase64(v interface{}) ([]byte, error) {
	buf := bufPool.Get()
	defer bufPool.Put(buf)
	if err := json.NewEncoder(buf).Encode(v); err != nil {
		return nil, err
	}
	return EncodeSegment(buf.Bytes()), nil
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
