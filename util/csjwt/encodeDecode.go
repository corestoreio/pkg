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

	"github.com/juju/errors"
)

// Decoder interface to pass in a custom decoding type.
type Decoder interface {
	// Unmarshal transforms a string into the arbitrary type v
	Unmarshal(data []byte, v interface{}) error
}

// Encoder interface to encode an arbitrary type into a byte slice
type Encoder interface {
	Marshal(v interface{}) ([]byte, error)
}

// JSONDecode default JSON decoder with base64 support
type JSONDecode struct{}

func (jp JSONDecode) Unmarshal(data []byte, v interface{}) error {
	dec, err := DecodeSegment(data)
	if err != nil {
		return errors.Mask(err)
	}
	return json.Unmarshal(dec, v)
}

// JSONEncode default JSON encoder with bas64 support
type JSONEncode struct{}

func (jp JSONEncode) Marshal(v interface{}) ([]byte, error) {
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
