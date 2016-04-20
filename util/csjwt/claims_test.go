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

package csjwt_test

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/corestoreio/csfw/util/csjwt"
	"github.com/corestoreio/csfw/util/csjwt/jwtclaim"
	"github.com/stretchr/testify/assert"
)

var _ csjwt.Header = (*csjwt.Head)(nil)
var _ fmt.Stringer = (*csjwt.Head)(nil)
var _ fmt.GoStringer = (*csjwt.Head)(nil)

type claimMock struct {
	validErr error
	setErr   error
	getErr   error
	keys     []string
}

func (c claimMock) Valid() error                            { return c.validErr }
func (c claimMock) Expires() time.Duration                  { return 0 }
func (c claimMock) Set(key string, value interface{}) error { return c.setErr }
func (c claimMock) Get(key string) (interface{}, error) {
	return nil, c.getErr
}
func (c claimMock) Keys() []string { return []string{"k1"} }

func TestNewHeadStringer(t *testing.T) {
	t.Parallel()
	var h csjwt.Header
	h = csjwt.NewHead("Quantum")
	assert.Exactly(t, "csjwt.NewHead(\"Quantum\")", fmt.Sprintf("%s", h))
	assert.Exactly(t, "csjwt.NewHead(\"Quantum\")", fmt.Sprintf("%v", h))
	assert.Exactly(t, "csjwt.NewHead(\"Quantum\")", fmt.Sprintf("%#v", h))
}

func TestNewHead(t *testing.T) {
	t.Parallel()
	var h csjwt.Header
	h = csjwt.NewHead("X")
	assert.Exactly(t, "X", h.Alg())
	assert.Exactly(t, csjwt.ContentTypeJWT, h.Typ())
}

func TestHeadSetGet(t *testing.T) {
	t.Parallel()
	var h csjwt.Header
	h = csjwt.NewHead("X")

	assert.NoError(t, h.Set(jwtclaim.HeaderAlg, "Y"))
	g, err := h.Get(jwtclaim.HeaderAlg)
	assert.NoError(t, err)
	assert.Exactly(t, "Y", g)

	assert.NoError(t, h.Set(jwtclaim.HeaderTyp, "JWE"))
	g, err = h.Get(jwtclaim.HeaderTyp)
	assert.NoError(t, err)
	assert.Exactly(t, "JWE", g)

	assert.EqualError(t, h.Set("x", "y"), "[csjwt] Header \"x\" not yet supported. Please switch to type jwtclaim.HeadSegments.")
	g, err = h.Get("x")
	assert.EqualError(t, err, "[csjwt] Header \"x\" not yet supported. Please switch to type jwtclaim.HeadSegments.")
	assert.Empty(t, g)
}

func TestMergeClaims(t *testing.T) {
	t.Parallel()
	tests := []struct {
		dst               csjwt.Token
		srcs              csjwt.Claimer
		wantSigningString string
		wantErr           error
	}{
		{csjwt.NewToken(nil), nil, `eyJ0eXAiOiJKV1QifQo.bnVsbAo`, nil},
		{csjwt.NewToken(jwtclaim.Map{}), claimMock{getErr: errors.New("claimMerge get error")}, ``, errors.New("[csjwt] Cannot get Key \"k1\" from Claim index 0. Error: claimMerge get error")},
		{csjwt.NewToken(jwtclaim.Map{"k1": "v1"}), jwtclaim.Map{"k2": 2}, `eyJ0eXAiOiJKV1QifQo.eyJrMSI6InYxIiwiazIiOjJ9Cg`, nil},
		{csjwt.NewToken(jwtclaim.NewStore()), jwtclaim.Map{"k2": 2}, ``, errors.New("[csjwt] Cannot set Key \"k2\" with value `2'. Claim index 0. Error: [jwtclaim] Claim \"k2\" not supported.")},
		{csjwt.NewToken(&jwtclaim.Standard{}), &jwtclaim.Store{
			Standard: &jwtclaim.Standard{},
			UserID:   "Gopher",
		}, ``, errors.New("[csjwt] Cannot set Key \"store\" with value `'. Claim index 0. Error: [jwtclaim] Claim \"store\" not supported.")},
	}
	for i, test := range tests {
		haveErr := csjwt.MergeClaims(test.dst.Claims, test.srcs)
		if test.wantErr != nil {
			assert.EqualError(t, haveErr, test.wantErr.Error(), "Index %d", i)
			continue
		}
		if haveErr != nil {
			t.Fatalf("Index %d => %s", i, haveErr)
		}

		buf, err := test.dst.SigningString()
		if err != nil {
			t.Fatalf("Index %d => %s", i, err)
		}
		assert.Exactly(t, test.wantSigningString, buf.String(), "Index %d", i)
	}
}
