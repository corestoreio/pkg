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

package modification_test

import (
	"os"
	"testing"

	"github.com/corestoreio/errors"
	"github.com/corestoreio/pkg/config"
	"github.com/corestoreio/pkg/config/modification"
	"github.com/corestoreio/pkg/util/assert"
)

func TestNewAESGCM(t *testing.T) {
	t.Parallel()
	t.Run("encrypt and decrypt with default values", func(t *testing.T) {
		o := &modification.AESGCMOptions{}
		obEnc, err := modification.NewAESGCM(config.EventOnBeforeSet, o)
		assert.NoError(t, err)
		obDec, err := modification.NewAESGCM(config.EventOnAfterGet, o) // o contains now the random nonce value

		p := *config.MustNewPath("aa/bb/cc")
		plainText := []byte(`X-Fit Games 2018`)
		encText, err := obEnc.Observe(p, plainText, false)
		assert.NoError(t, err)

		decText, err := obDec.Observe(p, encText, true)
		assert.NoError(t, err, "%+v", err)
		assert.Exactly(t, plainText, decText)
	})

	t.Run("invalid encryption key", func(t *testing.T) {
		o := &modification.AESGCMOptions{
			Key: "a",
		}
		obEnc, err := modification.NewAESGCM(config.EventOnBeforeSet, o)
		assert.Nil(t, obEnc)
		assert.True(t, errors.NotValid.Match(err), "%+v", err)
	})

	t.Run("event EventOnAfterGet does not decrypt", func(t *testing.T) {
		o := &modification.AESGCMOptions{
			Key: "abcdefghijklmnop",
		}
		obEnc, err := modification.NewAESGCM(config.EventOnAfterGet, o)
		p := *config.MustNewPath("aa/bb/cc")
		decText, err := obEnc.Observe(p, nil, false)
		assert.NoError(t, err, "%+v", err)
		assert.Nil(t, decText)
	})

	t.Run("encrypt and decrypt with environment key and nonce", func(t *testing.T) {
		os.Setenv("AESGCM_KEY", "randomKeyERTYUIO")
		os.Setenv("AESGCM_NONCE", "randomNonce!")
		defer func() {
			os.Unsetenv("AESGCM_KEY")
			os.Unsetenv("AESGCM_NONCE")
		}()

		o := &modification.AESGCMOptions{
			KeyEnvironmentVariableName:   "AESGCM_KEY",
			NonceEnvironmentVariableName: "AESGCM_NONCE",
		}
		obEnc, err := modification.NewAESGCM(config.EventOnBeforeSet, o)
		assert.NoError(t, err, "%+v", err)
		obDec, err := modification.NewAESGCM(config.EventOnAfterGet, o) // o contains now the random nonce value

		p := *config.MustNewPath("aa/bb/cc")
		plainText := []byte(`X-Fit Games 2018`)
		encText, err := obEnc.Observe(p, plainText, false)
		assert.NoError(t, err)

		decText, err := obDec.Observe(p, encText, true)
		assert.NoError(t, err, "%+v", err)
		assert.Exactly(t, plainText, decText)
	})

}
