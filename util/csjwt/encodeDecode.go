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
	"bytes"
	"encoding/base64"
	"encoding/gob"
	"encoding/json"
	"sync"

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
type gobEncoding struct {
	// TODO(CS): for higher performance remove the mutex and add a sync.Pool
	// pattern like in the transcache package for all encoders.
	mu   sync.Mutex
	pipe *bytes.Buffer
	enc  *gob.Encoder
	dec  *gob.Decoder
}

// NewGobEncoding creates a new primed gob Encoder/Decoder. Newly created
// Encoder/Decoder will Encode/Decode the passed sample structs without actually
// writing/reading from their respective Writer/Readers. This is useful for gob
// which encodes/decodes extra type information whenever it sees a new type.
// Pass sample values for primeObjects you plan on Encoding/Decoding to this
// method in order to avoid the storage overhead of encoding their type
// information for every NewEncoder/NewDecoder. Make sure you use gob.Register()
// for every type you plan to use otherwise there will be errors. Setting the
// primeObjects causes a priming of the encoder and decoder for each type. This
// function panics if the types, used for priming, can neither be encoded nor
// decoded.
func NewGobEncoding(primeObjects ...interface{}) *gobEncoding {
	pipe := new(bytes.Buffer)
	ge := &gobEncoding{
		pipe: pipe,
		enc:  gob.NewEncoder(pipe),
		dec:  gob.NewDecoder(pipe),
	}

	if len(primeObjects) > 0 {
		if err := ge.enc.Encode(primeObjects); err != nil {
			panic(err)
		}
		var testTypes []interface{}
		if err := ge.dec.Decode(&testTypes); err != nil {
			panic(err)
		}
		ge.pipe.Reset()
	}
	return ge
}

// Serialize encodes a value using gob.
func (e *gobEncoding) Serialize(src interface{}) ([]byte, error) {
	e.mu.Lock()
	defer e.mu.Unlock()
	defer e.pipe.Reset()
	if err := e.enc.Encode(src); err != nil {
		return nil, errors.Wrap(err, "[csjwt] GobEncoding.Serialize.Encode")
	}
	return EncodeSegment(e.pipe.Bytes()), nil
}

// Deserialize decodes a value using gob.
func (e *gobEncoding) Deserialize(src []byte, dst interface{}) error {
	srcDec, err := DecodeSegment(src)
	if err != nil {
		return errors.Wrap(err, "[csjwt] JSONEncoding.Deserialize.DecodeSegment")
	}

	e.mu.Lock()
	defer e.mu.Unlock()
	e.pipe.Reset()

	if _, err := e.pipe.Write(srcDec); err != nil {
		return errors.Wrap(err, "[csjwt] GobEncoding.Deserialize.Write")
	}
	if err := e.dec.Decode(dst); err != nil {
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
