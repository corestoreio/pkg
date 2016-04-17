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
	"testing"

	"errors"
	"github.com/corestoreio/csfw/util/csjwt"
	"github.com/corestoreio/csfw/util/csjwt/jwtclaim"
	"github.com/stretchr/testify/assert"
	"time"
)

func TestTokenAlg(t *testing.T) {
	t.Parallel()
	tests := []struct {
		tok     csjwt.Token
		wantAlg string
	}{
		{csjwt.NewToken(nil), ""},
		{
			csjwt.Token{
				Header: jwtclaim.NewHeadSegments("3"),
			},
			"3",
		},
		{
			csjwt.Token{
				Header: jwtclaim.NewHeadSegments("Gopher"),
			},
			"Gopher",
		},
	}
	for i, test := range tests {
		assert.Exactly(t, test.wantAlg, test.tok.Alg(), "Index %d", i)
	}
}

type claimMerge struct{}

func (claimMerge) Valid() error                            { return nil }
func (claimMerge) Expires() time.Duration                  { return 0 }
func (claimMerge) Set(key string, value interface{}) error { return nil }
func (claimMerge) Get(key string) (value interface{}, err error) {
	return nil, errors.New("claimMerge get error")
}
func (claimMerge) Keys() []string { return []string{"k1"} }

func TestToken_Merge(t *testing.T) {
	t.Parallel()
	tests := []struct {
		tk                csjwt.Token
		toMerge           csjwt.Claimer
		wantSigningString string
		wantErr           error
	}{
		{csjwt.NewToken(nil), nil, `eyJ0eXAiOiJKV1QifQo.bnVsbAo`, nil},
		{csjwt.NewToken(jwtclaim.Map{}), claimMerge{}, ``, errors.New("[csjwt] Cannot get Key \"k1\". Error: claimMerge get error")},
		{csjwt.NewToken(jwtclaim.Map{"k1": "v1"}), jwtclaim.Map{"k2": 2}, `eyJ0eXAiOiJKV1QifQo.eyJrMSI6InYxIiwiazIiOjJ9Cg`, nil},
		{csjwt.NewToken(jwtclaim.NewStore()), jwtclaim.Map{"k2": 2}, ``, errors.New(`[csjwt] Cannot set Key "k2" with value 2. Error: [jwtclaim] Claim "k2" not supported.`)},
		{csjwt.NewToken(&jwtclaim.Standard{}), &jwtclaim.Store{
			Standard: &jwtclaim.Standard{},
			UserID:   "Gopher",
		}, ``, errors.New(`[csjwt] Cannot set Key "store" with value . Error: [jwtclaim] Claim "store" not supported.`)},
	}
	for i, test := range tests {
		haveErr := test.tk.Merge(test.toMerge)
		if test.wantErr != nil {
			assert.EqualError(t, haveErr, test.wantErr.Error(), "Index %d", i)
			continue
		}
		if haveErr != nil {
			t.Fatalf("Index %d => %s", i, haveErr)
		}

		buf, err := test.tk.SigningString()
		if err != nil {
			t.Fatalf("Index %d => %s", i, err)
		}
		assert.Exactly(t, test.wantSigningString, buf.String(), "Index %d", i)
	}
}
