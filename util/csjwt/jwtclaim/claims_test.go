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

package jwtclaim_test

import (
	"encoding/json"
	"fmt"
	"math"
	"testing"
	"time"

	"github.com/corestoreio/errors"
	"github.com/corestoreio/pkg/util/csjwt"
	"github.com/corestoreio/pkg/util/csjwt/jwtclaim"
	"github.com/stretchr/testify/assert"
)

var _ csjwt.Claimer = (*jwtclaim.Standard)(nil)
var _ fmt.Stringer = (*jwtclaim.Standard)(nil)

var _ csjwt.Claimer = (*jwtclaim.Store)(nil)
var _ fmt.Stringer = (*jwtclaim.Store)(nil)

var _ csjwt.Claimer = (*jwtclaim.Map)(nil)
var _ fmt.Stringer = (*jwtclaim.Map)(nil)
var _ fmt.GoStringer = (*jwtclaim.Map)(nil)

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
		sc          csjwt.Claimer
		wantErrKind errors.Kind
	}{
		{&jwtclaim.Standard{}, errors.NoKind},
		{&jwtclaim.Standard{ExpiresAt: time.Now().Add(time.Second).Unix()}, errors.NoKind},
		{&jwtclaim.Standard{ExpiresAt: time.Now().Add(-time.Second).Unix()}, errors.NotValid},
		{&jwtclaim.Standard{IssuedAt: time.Now().Add(-time.Second).Unix()}, errors.NoKind},
		{&jwtclaim.Standard{IssuedAt: time.Now().Add(time.Second * 5).Unix()}, errors.NotValid},
		{&jwtclaim.Standard{NotBefore: time.Now().Add(-time.Second).Unix()}, errors.NoKind},
		{&jwtclaim.Standard{NotBefore: time.Now().Add(time.Second * 5).Unix()}, errors.NotValid},
		{
			&jwtclaim.Standard{
				ExpiresAt: time.Now().Add(-time.Second).Unix(),
				IssuedAt:  time.Now().Add(time.Second * 5).Unix(),
				NotBefore: time.Now().Add(time.Second * 5).Unix(),
			},
			errors.NotValid,
		},

		{jwtclaim.Map{}, errors.NotValid},                                               // 7
		{jwtclaim.Map{"exp": time.Now().Add(time.Second).Unix()}, errors.NoKind},        // 8
		{jwtclaim.Map{"exp": time.Now().Add(-time.Second * 2).Unix()}, errors.NotValid}, // 9
		{jwtclaim.Map{"iat": time.Now().Add(-time.Second).Unix()}, errors.NoKind},       // 10
		{jwtclaim.Map{"iat": time.Now().Add(time.Second * 5).Unix()}, errors.NotValid},
		{jwtclaim.Map{"nbf": time.Now().Add(-time.Second).Unix()}, errors.NoKind},
		{jwtclaim.Map{"nbf": time.Now().Add(time.Second * 5).Unix()}, errors.NotValid},
		{
			jwtclaim.Map{
				"exp": time.Now().Add(-time.Second).Unix(),
				"iat": time.Now().Add(time.Second * 5).Unix(),
				"nbf": time.Now().Add(time.Second * 5).Unix(),
			},
			errors.NotValid,
		},
	}
	for i, test := range tests {
		if !test.wantErrKind.Empty() {
			err := test.sc.Valid()
			assert.True(t, test.wantErrKind.Match(err), "Index %d => %s", i, err)
		} else {
			assert.NoError(t, test.sc.Valid(), "Index %d", i)
		}
	}
}

func TestClaimsGetSet(t *testing.T) {
	tests := []struct {
		sc            csjwt.Claimer
		key           string
		val           interface{}
		wantSetErrBhf errors.Kind
		wantGetErrBhf errors.Kind
	}{
		{&jwtclaim.Standard{}, jwtclaim.KeyAudience, '', errors.NotValid, errors.NoKind},
		{&jwtclaim.Standard{}, jwtclaim.KeyAudience, "Go", errors.NoKind, errors.NoKind},
		{&jwtclaim.Standard{}, jwtclaim.KeyExpiresAt, time.Now().Unix(), errors.NoKind, errors.NoKind},
		{&jwtclaim.Standard{}, "Not Supported", time.Now().Unix(), errors.NotSupported, errors.NotSupported},

		{jwtclaim.Map{}, jwtclaim.KeyAudience, "Go", errors.NoKind, errors.NoKind},
		{jwtclaim.Map{}, jwtclaim.KeyExpiresAt, time.Now().Unix(), errors.NoKind, errors.NoKind},
		{jwtclaim.Map{}, "Not Supported", math.Pi, errors.NoKind, errors.NoKind},
		{&jwtclaim.Store{}, jwtclaim.KeyStore, "xde", errors.NoKind, errors.NoKind},
	}
	for i, test := range tests {

		haveSetErr := test.sc.Set(test.key, test.val)
		if !test.wantSetErrBhf.Empty() {
			assert.True(t, test.wantSetErrBhf.Match(haveSetErr), "Index %d => %s", i, haveSetErr)
		} else {
			assert.NoError(t, haveSetErr, "Index %d", i)
		}

		haveVal, haveGetErr := test.sc.Get(test.key)
		if !test.wantGetErrBhf.Empty() {
			assert.True(t, test.wantGetErrBhf.Match(haveGetErr), "Index %d => %s", i, haveGetErr)
			continue
		} else {
			assert.NoError(t, haveGetErr, "Index %d", i)
		}

		if test.wantSetErrBhf.Empty() {
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

func TestClaimsExpiresSkew(t *testing.T) {
	tm := time.Now()
	tests := []struct {
		sc      csjwt.Claimer
		skew    time.Duration
		isValid bool
	}{
		{&jwtclaim.Standard{ExpiresAt: tm.Add(time.Second * 2).Unix()}, time.Second * 1, true},
		{&jwtclaim.Standard{ExpiresAt: tm.Add(time.Second * 2).Unix()}, -time.Second * 3, false},
		{&jwtclaim.Standard{ExpiresAt: tm.Unix() - 1}, 0, false},
		{&jwtclaim.Standard{ExpiresAt: tm.Unix() - 1}, time.Second * 1, true},

		{jwtclaim.Map{"exp": tm.Add(time.Second * 2).Unix()}, time.Second * 1, true},
		{jwtclaim.Map{"exp": tm.Add(time.Second * 2).Unix()}, -time.Second * 3, false},
		{jwtclaim.Map{"exp": tm.Unix() - 1}, 0, false},
		{jwtclaim.Map{"exp": tm.Unix() - 1}, time.Second * 1, true},
	}
	for i, test := range tests {
		if err := test.sc.Set(jwtclaim.KeyTimeSkew, test.skew); err != nil {
			t.Fatalf("%+v", err)
		}
		err := test.sc.Valid()
		assert.Exactly(t, !test.isValid, errors.NotValid.Match(err), "Index %d => %s", i, err)
	}
}

func TestClaimsNotBeforeSkew(t *testing.T) {
	tm := time.Now()
	tests := []struct {
		sc      csjwt.Claimer
		skew    time.Duration
		isValid bool
	}{
		{&jwtclaim.Standard{NotBefore: tm.Add(time.Second * 2).Unix()}, time.Second * 1, false},
		{&jwtclaim.Standard{NotBefore: tm.Add(time.Second * 2).Unix()}, time.Second * 3, true},
		{&jwtclaim.Standard{NotBefore: tm.Unix() - 1}, 0, true},
		{&jwtclaim.Standard{NotBefore: tm.Unix() - 1}, time.Second * 1, true},

		{jwtclaim.Map{"nbf": tm.Add(time.Second * 2).Unix()}, time.Second * 1, false},
		{jwtclaim.Map{"nbf": tm.Add(time.Second * 2).Unix()}, time.Second * 3, true},
		{jwtclaim.Map{"nbf": tm.Unix() - 1}, 0, true},
		{jwtclaim.Map{"nbf": tm.Unix() - 1}, time.Second * 1, true},
		{jwtclaim.Map{"nbf": tm.Unix() - 1}, -time.Second * 2, false},
	}
	for i, test := range tests {
		if err := test.sc.Set(jwtclaim.KeyTimeSkew, test.skew); err != nil {
			t.Fatalf("%+v", err)
		}
		err := test.sc.Valid()
		assert.Exactly(t, !test.isValid, errors.NotValid.Match(err), "Index %d => %s", i, err)
	}
}

func TestMap_String(t *testing.T) {
	m := jwtclaim.Map{
		"k1": "v1",
		"k2": 3.14159,
		"k3": false,
	}
	assert.Exactly(t, "{\"k1\":\"v1\",\"k2\":3.14159,\"k3\":false}", m.String())
}

func TestMap_String_error(t *testing.T) {
	m := jwtclaim.Map{
		"k1": "v1",
		"k2": 3.14159,
		"k3": make(chan int),
	}
	assert.Exactly(t, "[jwtclaim] Map.String(): json.Marshal Error: json: unsupported type: chan int", m.String())
}

func TestStandard_String(t *testing.T) {
	s := &jwtclaim.Standard{
		Issuer:    "Corestore",
		ExpiresAt: 4711,
	}
	assert.Exactly(t, "{\"exp\":4711,\"iss\":\"Corestore\"}", s.String())
}

func TestStore_String(t *testing.T) {
	s := jwtclaim.NewStore()
	s.Audience = "Gopher"
	s.ID = "1"
	s.Store = "nz"
	s.UserID = "23642736"
	assert.Exactly(t, "{\"aud\":\"Gopher\",\"jti\":\"1\",\"store\":\"nz\",\"userid\":\"23642736\"}", s.String())
}
