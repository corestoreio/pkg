// Copyright 2015-present, Cyrill @ Schumacher.fm and the CoreStore contributors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package csjwt

import (
	"encoding"
	"fmt"
	"time"

	"github.com/corestoreio/errors"
	"github.com/corestoreio/log"
	"github.com/corestoreio/pkg/util/bufferpool"
	"github.com/corestoreio/pkg/util/conv"
)

// ContentTypeJWT defines the content type of a token. At the moment only JWT is
// supported. JWE may be added in the future JSON Web Encryption (JWE).
// https://tools.ietf.org/html/rfc7519
const ContentTypeJWT = `JWT`

// TimeFunc provides the current time when parsing token to validate "exp" claim
// (expiration time). You can override it to use another time value.  This is
// useful for testing or if your server uses a different time zone than your
// tokens.
var TimeFunc = time.Now

// Token represents a JWT Token.  Different fields will be used depending on
// whether you're creating or parsing/verifying a token.
type Token struct {
	Raw       []byte  // The raw token.  Populated when you Parse a token
	Header    Header  // The first segment of the token
	Claims    Claimer // The second segment of the token
	Signature []byte  // The third segment of the token.  Populated when you Parse a token
	Valid     bool    // Is the token valid?  Populated when you Parse/Verify a token
	Serializer
}

// NewToken creates a new Token and presets the header to typ = JWT. A new token
// has not yet an assigned algorithm. The underlying default template header
// consists of a two field struct for the minimum requirements. If you need more
// header fields consider using a map or the jwtclaim.HeadSegments type. Default
// header from function NewHead().
func NewToken(c Claimer) *Token {
	return &Token{
		Header: NewHead(),
		Claims: c,
	}
}

// Alg returns the assigned algorithm to this token. Can return an empty string.
func (t *Token) Alg() string {
	if t.Header == nil {
		return ""
	}
	h, _ := t.Header.Get(headerAlg)
	return conv.ToString(h)
}

// SignedString gets the complete, signed token. Sets the header alg to the
// provided Signer.Alg() value. Returns a byte slice, save for further
// processing. This functions allows to sign a token with different signing
// methods.
func (t *Token) SignedString(method Signer, key Key) ([]byte, error) {
	if err := t.Header.Set(headerAlg, method.Alg()); err != nil {
		return nil, errors.WithStack(err)
	}

	buf, err := t.SigningString(make([]byte, 0, 512))
	if err != nil {
		return nil, errors.WithStack(err)
	}
	sig, err := method.Sign(buf, key)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	buf = append(buf, '.')
	buf = append(buf, sig...)
	return buf, nil
}

func (t *Token) marshal(dst interface{}) (data []byte, err error) {
	if t.Serializer != nil {
		data, err = t.Serializer.Serialize(dst)
		if err != nil {
			err = errors.WithStack(err)
		}
		return EncodeSegment(data), err
	}

	// the order of the cases is important as JSON can be embedded but main type
	// has e.g. TextMarshaler.
	switch dt := dst.(type) {
	case encoding.BinaryMarshaler:
		data, err = dt.MarshalBinary()
		if err != nil {
			err = errors.WithStack(err)
		}
	case encoding.TextMarshaler:
		data, err = dt.MarshalText()
		if err != nil {
			err = errors.WithStack(err)
		}
	case interface{ Marshal() ([]byte, error) }:
		data, err = dt.Marshal()
		if err != nil {
			err = errors.WithStack(err)
		}
	case interface{ MarshalJSON() ([]byte, error) }:
		data, err = dt.MarshalJSON()
		if err != nil {
			err = errors.WithStack(err)
		}
	default:
		enc := jsonEncoding{}
		data, err = enc.Serialize(dst)
		if err != nil {
			err = errors.WithStack(err)
		}
	}
	return EncodeSegment(data), err
}

// SigningString generates the signing string. This is the most expensive part
// of the whole deal.  Unless you need this for something special, just go
// straight for the SignedString. Returns a buffer which can be used for further
// modifications.
// SigningString supports custom binary, text, json, protobuf encoding.
func (t *Token) SigningString(buf []byte) ([]byte, error) {
	data, err := t.marshal(t.Header)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	buf = append(buf, data...)
	buf = append(buf, '.')

	data, err = t.marshal(t.Claims)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	buf = append(buf, data...)
	return buf, nil
}

// MarshalLog marshals the token into an unsigned log.Field. It uses the function
// SigningString().
func (t *Token) MarshalLog(kv log.KeyValuer) error {
	buf := bufferpool.Get()
	defer bufferpool.Put(buf)
	if tk, err := t.SigningString(buf.Bytes()); err != nil {
		kv.AddString("token_error", fmt.Sprintf("%+v", err))
	} else {
		kv.AddString("token", string(tk))
	}
	return nil
}
