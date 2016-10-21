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

	"github.com/corestoreio/csfw/config/cfgmock"
	"github.com/corestoreio/csfw/config/cfgmodel"
	"github.com/corestoreio/csfw/config/cfgpath"
	"github.com/corestoreio/csfw/store/scope"
	"github.com/stretchr/testify/assert"
)

var _ cfgmodel.Encrypter = (*rot13)(nil)
var _ cfgmodel.Decrypter = (*rot13)(nil)
var _ cfgmodel.Encrypter = (*cfgmodel.EncryptFunc)(nil)
var _ cfgmodel.Decrypter = (*cfgmodel.DecryptFunc)(nil)

// rot13 represents the most powerful encryption algorithm in the world ;-)
// Apply it two times to get the doubled security of encryption.
type rot13 struct{}

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

func TestObscure(t *testing.T) {

	var wantPlain = []byte(`H3llo G0phers`)
	var wantCiphered = []byte(`U3yyb T0curef`)
	const cfgPath = "aa/bb/cc"

	b := cfgmodel.NewObscure(
		cfgPath,
		cfgmodel.WithCSVComma(''), // trick it only for testing.
		cfgmodel.WithEncrypter(rot13{}),
		cfgmodel.WithDecrypter(rot13{}),
		cfgmodel.WithScopeStore(),
	)
	wantPath := cfgpath.MustNewByParts(cfgPath).String() // Default Scope

	haveSL, haveErr := b.Get(cfgmock.NewService(
		cfgmock.PathValue{
			wantPath: wantCiphered,
		}).NewScoped(34, 4))
	if haveErr != nil {
		t.Fatal(haveErr)
	}
	assert.Exactly(t, wantPlain, haveSL)

	mw := new(cfgmock.Write)
	b.Write(mw, wantPlain, scope.Store.Pack(12))
	assert.Exactly(t, wantCiphered, mw.ArgValue)
	assert.Exactly(t, "stores/12/aa/bb/cc", mw.ArgPath)
}
