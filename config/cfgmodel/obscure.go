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

// ErrMissingEncryptor gets returned if you have forgotten to set the
// Encryptor interface on the Obscure struct.
var ErrMissingEncryptor = errors.New("cfgmodel: Missing Encryptor")

// Encryptor functions needed for encryption and decryption of
// string values. For example implements M1 and M2 encryption key functions.
type Encryptor interface {
	Encrypt([]byte) ([]byte, error)
	Decrypt([]byte) ([]byte, error)
}

// NoopEncryptor does nothing and only used for testing.
// Noop = No operator
type NoopEncryptor struct {
	// EncErr allows to return an error while encrypting
	EncErr error
	// DecErr allows to return an error while decryption
	DecErr error
}

func (ne NoopEncryptor) Encrypt(s []byte) ([]byte, error) {
	return s, ne.EncErr
}

func (ne NoopEncryptor) Decrypt(s []byte) ([]byte, error) {
	return s, ne.DecErr
}

// WithEncryptor sets the functions for reading and writing encrypted data
// to the configuration service. May return nil.
func WithEncryptor(e Encryptor) Option {
	return func(b *optionBox) error {
		if b.Obscure == nil {
			return nil
		}
		b.Obscure.Encryptor = e
		return nil
	}
}

// Obscure backend model for handling sensible values
type Obscure struct {
	Byte
	Encryptor
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

// Get returns an encrypted value decrypted. Panics if Encryptor interface is nil.
func (p Obscure) Get(sg config.ScopedGetter) ([]byte, scope.Hash, error) {
	if p.Encryptor == nil {
		return nil, 0, ErrMissingEncryptor
	}
	s, h, err := p.Byte.Get(sg)
	if err != nil {
		return nil, h, errors.Wrap(err, "[cfgmodel] Obscure.Byte.Get")
	}
	s2, err := p.Decrypt(s)
	return s2, h, errors.Wrap(err, "[cfgmodel] Obscure.Get.Decrypt")
}

// Write writes a raw value encrypted. Panics if Encryptor interface is nil.
func (p Obscure) Write(w config.Writer, v []byte, s scope.Scope, scopeID int64) (err error) {
	if p.Encryptor == nil {
		return ErrMissingEncryptor
	}
	v, err = p.Encrypt(v)
	if err != nil {
		return err
	}
	return p.Byte.Write(w, v, s, scopeID)
}
