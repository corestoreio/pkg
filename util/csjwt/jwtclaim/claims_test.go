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

package jwtclaim_test

import (
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"testing"
	"time"

	"github.com/corestoreio/csfw/util/csjwt"
	"github.com/corestoreio/csfw/util/csjwt/jwtclaim"
	"github.com/stretchr/testify/assert"
)

var _ csjwt.Claimer = (*jwtclaim.Standard)(nil)
var _ csjwt.Claimer = (*jwtclaim.Map)(nil)

func TestStandardClaimsParseJSON(t *testing.T) {

	sc := jwtclaim.Standard{
		Audience:  `LOTR`,
		ExpiresAt: time.Now().Add(time.Hour).Unix(),

		IssuedAt:  time.Now().Unix(),
		Issuer:    `Gandalf`,
		NotBefore: time.Now().Unix(),
		Subject:   `Test Subject`,
	}
	rawJSON, err := json.Marshal(sc)
	if err != nil {
		t.Fatal(err)
	}
	assert.Len(t, rawJSON, 102, "JSON: %s", rawJSON)

	scNew := jwtclaim.Standard{}
	if err := json.Unmarshal(rawJSON, &scNew); err != nil {
		t.Fatal(err)
	}
	assert.Exactly(t, sc, scNew)
	assert.NoError(t, scNew.Valid())
}

func TestClaimsValid(t *testing.T) {
	tests := []struct {
		sc        csjwt.Claimer
		wantValid error
	}{
		{&jwtclaim.Standard{}, jwtclaim.ErrValidationClaimsInvalid},
		{&jwtclaim.Standard{ExpiresAt: time.Now().Add(time.Second).Unix()}, nil},
		{&jwtclaim.Standard{ExpiresAt: time.Now().Add(-time.Second).Unix()}, jwtclaim.ErrValidationExpired},
		{&jwtclaim.Standard{IssuedAt: time.Now().Add(-time.Second).Unix()}, nil},
		{&jwtclaim.Standard{IssuedAt: time.Now().Add(time.Second * 5).Unix()}, jwtclaim.ErrValidationUsedBeforeIssued},
		{&jwtclaim.Standard{NotBefore: time.Now().Add(-time.Second).Unix()}, nil},
		{&jwtclaim.Standard{NotBefore: time.Now().Add(time.Second * 5).Unix()}, jwtclaim.ErrValidationNotValidYet},
		{
			&jwtclaim.Standard{
				ExpiresAt: time.Now().Add(-time.Second).Unix(),
				IssuedAt:  time.Now().Add(time.Second * 5).Unix(),
				NotBefore: time.Now().Add(time.Second * 5).Unix(),
			},
			fmt.Errorf("%s\n%s\n%s", jwtclaim.ErrValidationExpired, jwtclaim.ErrValidationUsedBeforeIssued, jwtclaim.ErrValidationNotValidYet),
		},

		{jwtclaim.Map{}, jwtclaim.ErrValidationClaimsInvalid},                                         // 7
		{jwtclaim.Map{"exp": time.Now().Add(time.Second).Unix()}, nil},                                // 8
		{jwtclaim.Map{"exp": time.Now().Add(-time.Second * 2).Unix()}, jwtclaim.ErrValidationExpired}, // 9
		{jwtclaim.Map{"iat": time.Now().Add(-time.Second).Unix()}, nil},                               // 10
		{jwtclaim.Map{"iat": time.Now().Add(time.Second * 5).Unix()}, jwtclaim.ErrValidationUsedBeforeIssued},
		{jwtclaim.Map{"nbf": time.Now().Add(-time.Second).Unix()}, nil},
		{jwtclaim.Map{"nbf": time.Now().Add(time.Second * 5).Unix()}, jwtclaim.ErrValidationNotValidYet},
		{
			jwtclaim.Map{
				"exp": time.Now().Add(-time.Second).Unix(),
				"iat": time.Now().Add(time.Second * 5).Unix(),
				"nbf": time.Now().Add(time.Second * 5).Unix(),
			},
			fmt.Errorf("%s\n%s\n%s", jwtclaim.ErrValidationExpired, jwtclaim.ErrValidationUsedBeforeIssued, jwtclaim.ErrValidationNotValidYet),
		},
	}
	for i, test := range tests {
		if test.wantValid != nil {
			assert.EqualError(t, test.sc.Valid(), test.wantValid.Error(), "Index %d", i)
		} else {
			assert.NoError(t, test.sc.Valid(), "Index %d", i)
		}
	}
}

func TestClaimsGetSet(t *testing.T) {
	tests := []struct {
		sc         csjwt.Claimer
		key        string
		val        interface{}
		wantSetErr error
		wantGetErr error
	}{
		{&jwtclaim.Standard{}, jwtclaim.ClaimAudience, '', errors.New("Unable to cast 63743 to string"), nil},
		{&jwtclaim.Standard{}, jwtclaim.ClaimAudience, "Go", nil, nil},
		{&jwtclaim.Standard{}, jwtclaim.ClaimExpiresAt, time.Now().Unix(), nil, nil},
		{&jwtclaim.Standard{}, "Not Supported", time.Now().Unix(), errors.New("Claim \"Not Supported\" not supported. Please see constants Claim*."), errors.New("Claim \"Not Supported\" not supported. Please see constants Claim*.")},

		{jwtclaim.Map{}, jwtclaim.ClaimAudience, "Go", nil, nil},
		{jwtclaim.Map{}, jwtclaim.ClaimExpiresAt, time.Now().Unix(), nil, nil},
		{jwtclaim.Map{}, "Not Supported", math.Pi, nil, nil},
	}
	for i, test := range tests {

		haveSetErr := test.sc.Set(test.key, test.val)
		if test.wantSetErr != nil {
			assert.EqualError(t, haveSetErr, test.wantSetErr.Error(), "Index %d", i)
		} else {
			assert.NoError(t, haveSetErr, "Index %d", i)
		}

		haveVal, haveGetErr := test.sc.Get(test.key)
		if test.wantGetErr != nil {
			assert.EqualError(t, haveGetErr, test.wantGetErr.Error(), "Index %d", i)
			continue
		} else {
			assert.NoError(t, haveGetErr, "Index %d", i)
		}

		if test.wantSetErr == nil {
			assert.Exactly(t, test.val, haveVal, "Index %d", i)
		}
	}
}

func TestClaimsExpires(t *testing.T) {
	tm := time.Now()
	tests := []struct {
		sc      csjwt.Claimer
		wantExp time.Duration
	}{
		{&jwtclaim.Standard{ExpiresAt: tm.Add(time.Second * 2).Unix()}, time.Second * 1},
		{&jwtclaim.Standard{ExpiresAt: tm.Add(time.Second * 5).Unix()}, time.Second * 4},
		{&jwtclaim.Standard{ExpiresAt: -123123}, time.Duration(0)},
		{&jwtclaim.Standard{}, time.Duration(0)},

		{jwtclaim.Map{"exp": tm.Add(time.Second * 2).Unix()}, time.Second * 1},
		{jwtclaim.Map{"exp": tm.Add(time.Second * 22).Unix()}, time.Second * 21},
		{jwtclaim.Map{"exp": -123123}, time.Duration(0)},
		{jwtclaim.Map{"eXp": 23}, time.Duration(0)},
		{jwtclaim.Map{"exp": fmt.Sprintf("%d", tm.Unix()+10)}, time.Second * 9},
	}
	for i, test := range tests {
		assert.Exactly(t, int64(test.wantExp.Seconds()), int64(test.sc.Expires().Seconds()), "Index %d", i)
	}
}
