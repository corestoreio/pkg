package csjwt

import (
	"bytes"
	"time"

	"github.com/corestoreio/csfw/storage/text"
	"github.com/corestoreio/csfw/util/conv"
	"github.com/juju/errors"
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
	Raw       text.Chars // The raw token.  Populated when you Parse a token
	Header    Header     // The first segment of the token
	Claims    Claimer    // The second segment of the token
	Signature text.Chars // The third segment of the token.  Populated when you Parse a token
	Valid     bool       // Is the token valid?  Populated when you Parse/Verify a token
	Encoder
}

// NewToken creates a new Token and presets the header to typ = JWT.
// A new token has not yet an assigned algorithm. The underlying default
// template header consists of a two field struct for the minimum
// requirements. If you need more header fields consider using a map or
// the jwtclaim.HeadSegments type.
func NewToken(c Claimer) Token {
	return Token{
		Header: NewHead(),
		Claims: c,
	}
}

// Alg returns the assigned algorithm to this token.
// Can return an empty string.
func (t Token) Alg() string {
	if t.Header == nil {
		return ""
	}
	h, _ := t.Header.Get(headerAlg)
	return conv.ToString(h)
}

// SignedString gets the complete, signed token.
// Sets the header alg to the provided Signer.Alg() value.
// Returns a byte slice, save for further processing.
// This functions allows to sign a token with different signing methods.
func (t Token) SignedString(method Signer, key Key) (text.Chars, error) {

	if err := t.Header.Set(headerAlg, method.Alg()); err != nil {
		return nil, errors.Mask(err)
	}

	buf, err := t.SigningString()
	if err != nil {
		return nil, errors.Mask(err)
	}
	sig, err := method.Sign(buf.Bytes(), key)
	if err != nil {
		return nil, errors.Mask(err)
	}

	if _, err := buf.WriteRune('.'); err != nil {
		return nil, errors.Mask(err)
	}
	if _, err := buf.Write(sig); err != nil {
		return nil, errors.Mask(err)
	}
	return buf.Bytes(), nil
}

// SigningString generates the signing string.  This is the
// most expensive part of the whole deal.  Unless you
// need this for something special, just go straight for
// the SignedString.
// Returns a buffer which can be used for further modifications.
func (t Token) SigningString() (buf bytes.Buffer, err error) {

	enc := t.Encoder
	if enc == nil {
		enc = JSONEncode{}
	}

	var j []byte
	j, err = enc.Marshal(t.Header)
	if err != nil {
		err = errors.Mask(err)
		return
	}
	if _, err = buf.Write(j); err != nil {
		err = errors.Mask(err)
		return
	}
	if _, err = buf.WriteRune('.'); err != nil {
		err = errors.Mask(err)
		return
	}
	j, err = enc.Marshal(t.Claims)
	if err != nil {
		err = errors.Mask(err)
		return
	}
	if _, err = buf.Write(j); err != nil {
		err = errors.Mask(err)
		return
	}
	return
}
