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

package cfgmodel

import (
	"github.com/corestoreio/csfw/config"
	"github.com/corestoreio/csfw/store/scope"
	"github.com/corestoreio/csfw/util/errors"
)

// EncodeFunc encodes the value v into a byte slice.
type EncodeFunc func(v interface{}) (data []byte, _ error)

// DecodeFunc decodes the data into the pointer vPtr.
type DecodeFunc func(data []byte, vPtr interface{}) error

// WithEncoder sets the functions for encoding and decoding data to the
// configuration service. Tip: You can directly use json.Marshal, json.Unmarshal
// or xml.Marshal, xml.Unmarshal.
func WithEncoder(e EncodeFunc, d DecodeFunc) Option {
	return func(b *optionBox) error {
		if b.Encode == nil {
			return nil
		}
		b.Encode.EncodeFunc = e
		b.Encode.DecodeFunc = d
		return nil
	}
}

// Encode backend model for handling json, xml, toml, csv and many other formats
// which needs encoding and decoding.
type Encode struct {
	Byte
	EncodeFunc
	DecodeFunc
}

// NewEncode creates a new Encode with validation checks when writing values.
func NewEncode(path string, opts ...Option) Encode {
	ret := Encode{
		Byte: NewByte(path),
	}
	(&ret).Option(opts...)
	return ret
}

// Option sets the options and returns the last set previous option
func (p *Encode) Option(opts ...Option) error {
	ob := &optionBox{
		baseValue: &p.baseValue,
		Encode:    p,
	}
	for _, o := range opts {
		if err := o(ob); err != nil {
			return errors.Wrap(err, "[cfgmodel] Encode.Option")
		}
	}
	p = ob.Encode
	p.baseValue = *ob.baseValue
	return nil
}

// Get uses the pointer argument vPtr to decode the data into vPtr. It panics
// when the Encoder interface is nil. It does not check if vPtr has been passed
// as a pointer.
func (p Encode) Get(sg config.Scoped, vPtr interface{}) error {
	s, err := p.Byte.Get(sg)
	if err != nil {
		return errors.Wrap(err, "[cfgmodel] Encode.Byte.Get")
	}
	return errors.Wrap(p.DecodeFunc(s, vPtr), "[cfgmodel] Encode.Get.Decode")
}

// Write writes a raw value encrypted. Panics if Encryptor interface is nil.
func (p Encode) Write(w config.Writer, v interface{}, h scope.TypeID) error {
	raw, err := p.EncodeFunc(v)
	if err != nil {
		return errors.Wrap(err, "[cfgmodel] Encode.Write.Encode")
	}
	return errors.Wrap(p.Byte.Write(w, raw, h), "[cfgmodel] Encode.Write.Write")
}
