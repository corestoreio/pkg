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

package cfgmodel_test

import (
	"testing"

	"github.com/corestoreio/csfw/config"
	"github.com/corestoreio/csfw/config/cfgmock"
	"github.com/corestoreio/csfw/config/cfgmodel"
	"github.com/corestoreio/csfw/config/cfgpath"
	"github.com/corestoreio/csfw/store/scope"
	"github.com/stretchr/testify/assert"
)

var _ cfgmodel.Encryptor = (*rot13)(nil)
var _ cfgmodel.Encryptor = (*cfgmodel.NoopEncryptor)(nil)

type rot13 struct {
}

func (rt rot13) Encrypt(s []byte) ([]byte, error) {
	var buf [1024]byte
	n := copy(buf[:], s)
	for i, b := range buf[:n] {
		switch {
		case 'a' <= b && b <= 'm', 'A' <= b && b <= 'M':
			buf[i] = b + 13
		case 'n' <= b && b <= 'z', 'N' <= b && b <= 'Z':
			buf[i] = b - 13
		}
	}
	return buf[:n], nil
}

func (rt rot13) Decrypt(s []byte) ([]byte, error) {
	return rt.Encrypt(s)
}
func TestObscureMissingEncryptor(t *testing.T) {

	m := cfgmodel.NewObscure(`aa/bb/cc`)
	val, h, err := m.Get(config.Scoped{})
	assert.Nil(t, val)
	assert.EqualError(t, err, cfgmodel.ErrMissingEncryptor.Error())
	assert.EqualError(t, m.Write(nil, nil, scope.Default, 0), cfgmodel.ErrMissingEncryptor.Error())
	assert.Exactly(t, scope.Hash(0).String(), h.String())
}

func TestObscure(t *testing.T) {

	var wantPlain = []byte(`H3llo G0phers`)
	var wantCiphered = []byte(`U3yyb T0curef`)
	const cfgPath = "aa/bb/cc"

	b := cfgmodel.NewObscure(
		cfgPath,
		cfgmodel.WithCSVComma(''), // trick it
		cfgmodel.WithEncryptor(rot13{}),
	)
	wantPath := cfgpath.MustNewByParts(cfgPath).String() // Default Scope

	haveSL, haveH, haveErr := b.Get(cfgmock.NewService(
		cfgmock.WithPV(cfgmock.PathValue{
			wantPath: wantCiphered,
		}),
	).NewScoped(34, 4))
	if haveErr != nil {
		t.Fatal(haveErr)
	}
	assert.Exactly(t, wantPlain, haveSL)
	assert.Exactly(t, scope.DefaultHash.String(), haveH.String())

	mw := new(cfgmock.Write)
	b.Write(mw, wantPlain, scope.Store, 12)
	assert.Exactly(t, wantCiphered, mw.ArgValue)
	assert.Exactly(t, "stores/12/aa/bb/cc", mw.ArgPath)
}

func TestNoopEncryptor(t *testing.T) {

	ne := cfgmodel.NoopEncryptor{}

	var a = []byte("a")
	e, err := ne.Encrypt(a)
	assert.Exactly(t, a, e)
	assert.NoError(t, err)

	var b = []byte("b")
	d, err := ne.Decrypt(b)
	assert.Exactly(t, b, d)
	assert.NoError(t, err)
}
