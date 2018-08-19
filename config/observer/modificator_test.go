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

package observer

import (
	"crypto/sha256"
	"testing"

	"github.com/corestoreio/errors"
	"github.com/corestoreio/pkg/config"
	"github.com/corestoreio/pkg/util/assert"
	"github.com/corestoreio/pkg/util/hashpool"
)

func init() {
	hashpool.Register("sha256", sha256.New)
}

func errCheck(t *testing.T) func([]byte, error) []byte {
	return func(data []byte, err error) []byte {
		assert.NoError(t, err)
		return data
	}
}

func TestOperators(t *testing.T) {
	t.Parallel()

	assert.Exactly(t, []byte(`X`), errCheck(t)(trim(nil, []byte(" X\t\n"))))
	assert.Exactly(t, []byte(`HELLO`), errCheck(t)(toUpper(nil, []byte("helLo"))))
	assert.Exactly(t, []byte(`wor€d`), errCheck(t)(toLower(nil, []byte("WOr€D"))))
	assert.Exactly(t, []byte(`The Small Universe`), errCheck(t)(toTitle(nil, []byte("the small universe"))))
	assert.Exactly(t, []byte(`dGhlIHNtYWxsIHVuaXZlcnNl`), errCheck(t)(base64Encode(nil, []byte("the small universe"))))
	assert.Exactly(t, []byte(`the small universe`), errCheck(t)(base64Decode(nil, []byte("dGhlIHNtYWxsIHVuaXZlcnNl"))))

	assert.Exactly(t, []byte(`74686520736d616c6c20756e697665727365`), errCheck(t)(hexEncode(nil, []byte("the small universe"))))
	assert.Exactly(t, []byte(`the small universe`), errCheck(t)(hexDecode(nil, []byte("74686520736d616c6c20756e697665727365"))))

	p0 := config.MustNewPath("general/store_information/address")
	p1 := p0.BindStore(2)

	assert.Exactly(t, []byte("\xfa\x93\xfcRzW,p=\x1fV1\x9b\xb8n\x05\xa3\xc4\xe3A\xb7\xdbU%\xbd݇3ñ\x8f\xbf"),
		errCheck(t)(hash256(p0, []byte("The Small Universe"))))
	assert.Exactly(t, []byte("\xb46ܷ\xaeM\xa8\xaaׄ?B\\\xab\xee4u \xeb\xfe$\xe5\x1d\x9a\x973\x99\x15\xae\xef@\x06"),
		errCheck(t)(hash256(p1, []byte("The Small Universe"))))

	assert.Exactly(t, []byte("\x1f\x8b\b\x00\x00\x00\x00\x00\x00\xff\n\xc9HU\b\xceM\xcc\xc9Q\b\xcd\xcb,K-*N\x05\x04\x00\x00\xff\xff\x05\xd1r\xe9\x12\x00\x00\x00"),
		errCheck(t)(dataGzip(nil, []byte("The Small Universe"))))

	assert.Exactly(t, []byte("The Small Universe"),
		errCheck(t)(dataGunzip(nil, []byte("\x1f\x8b\b\x00\x00\x00\x00\x00\x00\xff\n\xc9HU\b\xceM\xcc\xc9Q\b\xcd\xcb,K-*N\x05\x04\x00\x00\xff\xff\x05\xd1r\xe9\x12\x00\x00\x00"))))
}

func TestMustNewModificator(t *testing.T) {
	t.Parallel()
	defer func() {
		if r := recover(); r != nil {
			if err, ok := r.(error); ok {
				assert.True(t, errors.NotSupported.Match(err))
			} else {
				t.Errorf("Panic should contain an error but got:\n%+v", r)
			}
		} else {
			t.Error("Expecting a panic but got nothing")
		}
	}()
	_ = MustNewModificator(ModificatorArg{
		Funcs: []string{"neverGonnaGiveYouUp"},
	})
}

func TestNewModificator(t *testing.T) {
	t.Parallel()

	t.Run("trim upper", func(t *testing.T) {
		ms := MustNewModificator(ModificatorArg{
			Funcs: []string{"trim", "upper"},
		})

		var p config.Path
		data := []byte(" \thello\n \t")
		have, err := ms.Observe(p, data, true)
		assert.NoError(t, err)
		assert.Exactly(t, "HELLO", string(have))
	})

	t.Run("custom operator returns error ", func(t *testing.T) {
		RegisterModificator("csx", func(*config.Path, []byte) ([]byte, error) {
			return nil, errors.New("An error")
		})

		ms := MustNewModificator(ModificatorArg{
			Funcs: []string{"trim", "csx", "upper"},
		})

		var p config.Path
		data := []byte(" \thello\n \t")
		have, err := ms.Observe(p, data, true)
		assert.Nil(t, have)
		assert.True(t, errors.Interrupted.Match(err), "%+v", err)

	})
}
