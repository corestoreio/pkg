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

// Obscure backend model for handling sensible values
type Obscure struct {
	Byte
	Encryptor
}

// NewObscure creates a new Obscure type.  It will panic while calling later
// Get()/Write() when the Encryptor has not been set.
func NewObscure(path string, opts ...Option) Obscure {
	return Obscure{
		Byte: NewByte(path, opts...),
	}
}

// Get returns an encrypted value decrypted. Panics if Encryptor interface is
// nil.
func (p Obscure) Get(sg config.Scoped) ([]byte, scope.Hash, error) {
	s, h, err := p.Byte.Get(sg)
	if err != nil {
		return nil, h, errors.Wrap(err, "[cfgmodel] Obscure.Byte.Get")
	}
	s2, err := p.Decrypt(s)
	return s2, h, errors.Wrap(err, "[cfgmodel] Obscure.Get.Decrypt")
}

// Write writes a raw value encrypted. Panics if Encryptor interface is nil.
func (p Obscure) Write(w config.Writer, v []byte, s scope.Scope, scopeID int64) (err error) {
	v, err = p.Encrypt(v)
	if err != nil {
		return err
	}
	return p.Byte.Write(w, v, s, scopeID)
}
