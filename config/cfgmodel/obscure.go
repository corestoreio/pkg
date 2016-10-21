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

// todo implement compatible type for both interfaces like in the example ExampleNewGCM_encrypt() and ExampleNewGCM_decrypt of the std lib cipher package.

// Encrypter defines a function which encrypts the plaintext input data and
// returns the encrypted data. Or may return an error.
type Encrypter interface {
	Encrypt([]byte) ([]byte, error)
}

// Decrypter defines a function which decrypts the input data and returns the
// plaintext data. Or may return an error.
type Decrypter interface {
	Decrypt([]byte) ([]byte, error)
}

// EncodeFunc defines a wrapper type to match interface Encoder
type EncryptFunc func(s []byte) ([]byte, error)

func (ef EncryptFunc) Encrypt(s []byte) ([]byte, error) {
	return ef(s)
}

// DecryptFunc defines a wrapper type to match interface Decoder
type DecryptFunc func(v interface{}) (data []byte, _ error)

func (df DecryptFunc) Decrypt(s []byte) ([]byte, error) {
	return df(s)
}

// WithEncrypter sets the encryption function. Convenient helper function.
func WithEncrypter(e Encrypter) Option {
	return func(b *optionBox) error {
		if b.Obscure == nil {
			return nil
		}
		b.Obscure.Encrypter = e
		return nil
	}
}

// WithDecrypter set the decryption function. Convenient helper function.
func WithDecrypter(d Decrypter) Option {
	return func(b *optionBox) error {
		if b.Obscure == nil {
			return nil
		}
		b.Obscure.Decrypter = d
		return nil
	}
}

// Obscure backend model for handling sensible values
type Obscure struct {
	Byte
	Encrypter
	Decrypter
}

// NewObscure creates a new Obscure with validation checks when writing values.
func NewObscure(path string, opts ...Option) Obscure {
	ret := Obscure{
		Byte: NewByte(path),
	}
	(&ret).Option(opts...)
	return ret
}

// Option sets the options and returns the last set previous option
func (p *Obscure) Option(opts ...Option) error {
	ob := &optionBox{
		baseValue: &p.baseValue,
		Obscure:   p,
	}
	for _, o := range opts {
		if err := o(ob); err != nil {
			return errors.Wrap(err, "[cfgmodel] Obscure.Option")
		}
	}
	p = ob.Obscure
	p.baseValue = *ob.baseValue
	return nil
}

// Get returns an encrypted value decrypted. Panics if Encryptor interface is
// nil.
func (p Obscure) Get(sg config.Scoped) ([]byte, error) {
	s, err := p.Byte.Get(sg)
	if err != nil {
		return nil, errors.Wrap(err, "[cfgmodel] Obscure.Byte.Get")
	}
	s2, err := p.Decrypt(s)
	return s2, errors.Wrap(err, "[cfgmodel] Obscure.Get.Decrypt")
}

// Write writes a raw value encrypted. Panics if Encryptor interface is nil.
func (p Obscure) Write(w config.Writer, v []byte, h scope.TypeID) (err error) {
	v, err = p.Encrypt(v)
	if err != nil {
		return err
	}
	return p.Byte.Write(w, v, h)
}
