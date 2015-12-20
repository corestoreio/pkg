// Copyright 2015, Cyrill @ Schumacher.fm and the CoreStore contributors
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
	"testing"

	"bytes"

	"github.com/corestoreio/csfw/config"
	"github.com/corestoreio/csfw/config/scope"
	"github.com/corestoreio/csfw/net/ctxjwt"
	"github.com/stretchr/testify/assert"
)

func TestPasswordFromConfig(t *testing.T) {

	cfg := config.NewMockGetter(
		config.WithMockValues(config.MockPV{
			scope.StrDefault.FQPathInt64(ctxjwt.PathJWTPassword): `Rump3lst!lzch3n`,
		}),
	)

	jm, err := ctxjwt.NewService(
		ctxjwt.WithPasswordFromConfig(cfg),
	)
	assert.NoError(t, err)

	theToken, _, err := jm.GenerateToken(nil)
	assert.NoError(t, err)
	assert.NotEmpty(t, theToken)

}

func TestWithRSAReaderFail(t *testing.T) {

	jm, err := ctxjwt.NewService(
		ctxjwt.WithRSA(bytes.NewReader([]byte(`invalid pem data`))),
	)
	assert.Nil(t, jm)
	assert.Equal(t, "Private Key from io.Reader no found", err.Error())

}
