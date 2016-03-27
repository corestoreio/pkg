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

package ctxjwt_test

import (
	"crypto/x509"
	"path/filepath"
	"testing"

	"github.com/corestoreio/csfw/net/ctxjwt"
	"github.com/corestoreio/csfw/store/scope"
	"github.com/stretchr/testify/assert"
)

func TestOptionWithTokenID(t *testing.T) {
	t.Parallel()
	jwts, err := ctxjwt.NewService(ctxjwt.WithTokenID(scope.DefaultID, 0, true))
	assert.NoError(t, err)

	theToken, jti, err := jwts.GenerateToken(scope.DefaultID, 0, nil)
	assert.NoError(t, err)
	assert.NotEmpty(t, jti)
	assert.NotEmpty(t, theToken)
	assert.Len(t, jti, uuidLen)
}

func TestOptionWithRSAReaderFail(t *testing.T) {
	t.Parallel()
	jm, err := ctxjwt.NewService(
		ctxjwt.WithRSA(scope.DefaultID, 0, []byte(`invalid pem data`)),
	)
	assert.Nil(t, jm)
	assert.Equal(t, "Private Key from io.Reader no found", err.Error())

}

var pkFile = filepath.Join("testdata", "test_rsa")

func TestOptionWithRSAFromFileNoOrFailedPassword(t *testing.T) {
	t.Parallel()
	jm, err := ctxjwt.NewService(ctxjwt.WithRSAFromFile(scope.DefaultID, 0, pkFile))
	assert.EqualError(t, err, ctxjwt.ErrPrivateKeyNoPassword.Error())
	assert.Nil(t, jm)

	jm2, err2 := ctxjwt.NewService(ctxjwt.WithRSAFromFile(scope.DefaultID, 0, pkFile, []byte(`adfasdf`)))
	assert.EqualError(t, err2, x509.IncorrectPasswordError.Error())
	assert.Nil(t, jm2)
}

func testRsaOption(t *testing.T, opt ctxjwt.Option) {
	jm, err := ctxjwt.NewService(opt)
	assert.NoError(t, err)
	assert.NotNil(t, jm)

	theToken, _, err := jm.GenerateToken(scope.DefaultID, 0, nil)
	assert.NoError(t, err)
	assert.NotEmpty(t, theToken)

	tk, err := jm.Parse(theToken)
	assert.NoError(t, err)
	assert.NotNil(t, tk)
	assert.True(t, tk.Valid)
}

func TestOptionWithRSAFromFilePassword(t *testing.T) {
	t.Parallel()
	pw := []byte("cccamp")
	testRsaOption(t, ctxjwt.WithRSAFromFile(scope.DefaultID, 0, pkFile, pw))
}

func TestOptionWithRSAFromFileNoPassword(t *testing.T) {
	t.Parallel()
	// pkFileNP := filepath.Join(cstesting.RootPath, "net", "ctxjwt", "test_rsa_np")
	pkFileNP := filepath.Join("testdata", "test_rsa_np")
	testRsaOption(t, ctxjwt.WithRSAFromFile(scope.DefaultID, 0, pkFileNP))
}

func TestOptionWithRSAGenerator(t *testing.T) {
	t.Parallel()
	if testing.Short() {
		t.Skip("Test skipped in short mode")
	}
	testRsaOption(t, ctxjwt.WithRSAGenerator(scope.DefaultID, 0))
}

func TestOptionWithBackend(t *testing.T) {
	t.Parallel()
	t.Log("todo")
}
