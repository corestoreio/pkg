// Copyright 2015-2016, Cyrill @ Schumacher.fm and the CoreStore contributors
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
	"encoding/base64"
	"encoding/json"

	"encoding/gob"

	"github.com/corestoreio/csfw/util/errors"
)

// Deserializer provides an interface for providing custom deserializers.
// Also known as unserialize ;-)
type Deserializer interface {
	Deserialize(src []byte, dst interface{}) error
}

// Serializer provides an interface for providing custom serializers.
type Serializer interface {
	Serialize(src interface{}) ([]byte, error)
}

// JSONEncoding default JSON de- & serializer with base64 support
type JSONEncoding struct{}

// Deserialize decodes a value using encoding/json.
func (jp JSONEncoding) Deserialize(src []byte, dst interface{}) error {
	dec, err := DecodeSegment(src)
	if err != nil {
		return errors.Wrap(err, "[csjwt] JSONEncoding.Deserialize.DecodeSegment")
	}
	return errors.Wrap(json.Unmarshal(dec, dst), "[csjwt] JSONEncoding.Deserialize.Unmarshal")
}

// Serialize encodes a value using encoding/json.
func (jp JSONEncoding) Serialize(src interface{}) ([]byte, error) {
	buf := bufPool.Get()
	defer bufPool.Put(buf)
	if err := json.NewEncoder(buf).Encode(src); err != nil {
		return nil, errors.Wrap(err, "[csjwt] JSONEncoding.Serialize.Encode")
	}
	return EncodeSegment(buf.Bytes()), nil
}

// GobEncoding encodes JWT values using encoding/gob. This is the simplest
// encoder and can handle complex types via gob.Register.
// TODO(CS): Add gob priming to avoid storing type information in the token.
type GobEncoding struct{}

// Serialize encodes a value using gob.
func (e GobEncoding) Serialize(src interface{}) ([]byte, error) {
	buf := bufPool.Get()
	defer bufPool.Put(buf)
	// todo figure out how to use one instance of NewEncoder instead of creating each time a new one
	enc := gob.NewEncoder(buf)
	if err := enc.Encode(src); err != nil {
		return nil, errors.Wrap(err, "[csjwt] GobEncoding.Serialize.Encode")
	}
	return EncodeSegment(buf.Bytes()), nil
}

// Deserialize decodes a value using gob.
func (e GobEncoding) Deserialize(src []byte, dst interface{}) error {
	srcDec, err := DecodeSegment(src)
	if err != nil {
		return errors.Wrap(err, "[csjwt] JSONEncoding.Deserialize.DecodeSegment")
	}

	buf := bufPool.Get()
	defer bufPool.Put(buf)
	if _, err := buf.Write(srcDec); err != nil {
		return errors.Wrap(err, "[csjwt] GobEncoding.Deserialize.Write")
	}
	// todo figure out how to use one instance of NewDecoder instead of creating each time a new one
	dec := gob.NewDecoder(buf)
	if err := dec.Decode(dst); err != nil {
		return errors.Wrap(err, "[csjwt] GobEncoding.Deserialize.Decode")
	}
	return nil
}

// EncodeSegment encodes JWT specific base64url encoding with padding stripped.
// Returns a new byte slice.
func EncodeSegment(seg []byte) []byte {
	dbuf := make([]byte, base64.RawURLEncoding.EncodedLen(len(seg)))
	base64.RawURLEncoding.Encode(dbuf, seg)
	return dbuf
}

// DecodeSegment decodes JWT specific base64url encoding with padding stripped.
// Returns a new byte slice. Error behaviour: NotValid.
func DecodeSegment(seg []byte) ([]byte, error) {
	dbuf := make([]byte, base64.RawURLEncoding.DecodedLen(len(seg)))
	n, err := base64.RawURLEncoding.Decode(dbuf, seg)
	return dbuf[:n], errors.NewNotValid(err, "[csjwt] DecodeSegment")
}
