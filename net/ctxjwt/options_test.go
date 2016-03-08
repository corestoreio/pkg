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
	"bytes"
	"testing"

	"github.com/corestoreio/csfw/config/mock"
	"github.com/corestoreio/csfw/config/model"
	"github.com/corestoreio/csfw/config/path"
	"github.com/corestoreio/csfw/net/ctxjwt"
	"github.com/stretchr/testify/assert"
)

func TestPasswordFromConfig(t *testing.T) {
	t.Parallel()
	srvSG := mock.NewService(
		mock.WithPV(mock.PathValue{
			path.MustNewByParts(ctxjwt.PathJWTHMACPassword).String(): `Rump3lst!lzch3n`,
		}),
	).NewScoped(1, 2)

	jm, err := ctxjwt.NewService(
		ctxjwt.WithPasswordFromConfig(srvSG, model.NoopEncryptor{}),
	)
	assert.NoError(t, err)

	theToken, _, err := jm.GenerateToken(nil)
	assert.NoError(t, err)
	assert.NotEmpty(t, theToken)

}

func TestWithRSAReaderFail(t *testing.T) {
	t.Parallel()
	jm, err := ctxjwt.NewService(
		ctxjwt.WithRSA(bytes.NewReader([]byte(`invalid pem data`))),
	)
	assert.Nil(t, jm)
	assert.Equal(t, "Private Key from io.Reader no found", err.Error())

}
