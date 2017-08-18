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
	"bytes"
	"io/ioutil"
	"testing"

	"github.com/corestoreio/csfw/util/csjwt"
	"github.com/corestoreio/errors"
	"github.com/stretchr/testify/assert"
)

var hmacFastTestData []struct {
	name        string
	tokenString []byte
	method      csjwt.Signer
	claims      map[string]interface{}
	valid       bool
}

// Sample data from http://tools.ietf.org/html/draft-jones-json-web-signature-04#appendix-A.1
var hmacTestKey []byte

func init() {
	var err error
	hmacTestKey, err = ioutil.ReadFile("test/hmacTestKey")
	if err != nil {
		panic(err)
	}

	hf256, err := csjwt.NewSigningMethodHS256Fast(csjwt.WithPassword(hmacTestKey))
	if err != nil {
		panic(err)
	}

	hf384, err := csjwt.NewSigningMethodHS384Fast(csjwt.WithPassword(hmacTestKey))
	if err != nil {
		panic(err)
	}

	hf512, err := csjwt.NewSigningMethodHS512Fast(csjwt.WithPassword(hmacTestKey))
	if err != nil {
		panic(err)
	}

	blk256, err := csjwt.NewSigningMethodBlake2b256(csjwt.WithPassword(hmacTestKey))
	if err != nil {
		panic(err)
	}
	blk512, err := csjwt.NewSigningMethodBlake2b512(csjwt.WithPassword(hmacTestKey))
	if err != nil {
		panic(err)
	}

	hmacFastTestData = []struct {
		name        string
		tokenString []byte
		method      csjwt.Signer
		claims      map[string]interface{}
		valid       bool
	}{
		{ // 0
			"web sample",
			[]byte("eyJ0eXAiOiJKV1QiLA0KICJhbGciOiJIUzI1NiJ9.eyJpc3MiOiJqb2UiLA0KICJleHAiOjEzMDA4MTkzODAsDQogImh0dHA6Ly9leGFtcGxlLmNvbS9pc19yb290Ijp0cnVlfQ.dBjftJeZ4CVP-mB92K27uhbUJU1p1r_wW1gFWFOEjXk"),
			hf256,
			map[string]interface{}{"iss": "joe", "exp": 1300819380, "http://example.com/is_root": true},
			true,
		},
		{ // 1
			"HS384",
			[]byte("eyJhbGciOiJIUzM4NCIsInR5cCI6IkpXVCJ9.eyJleHAiOjEuMzAwODE5MzhlKzA5LCJodHRwOi8vZXhhbXBsZS5jb20vaXNfcm9vdCI6dHJ1ZSwiaXNzIjoiam9lIn0.KWZEuOD5lbBxZ34g7F-SlVLAQ_r5KApWNWlZIIMyQVz5Zs58a7XdNzj5_0EcNoOy"),
			hf384,
			map[string]interface{}{"iss": "joe", "exp": 1300819380, "http://example.com/is_root": true},
			true,
		},
		{ // 2
			"HS512",
			[]byte("eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJleHAiOjEuMzAwODE5MzhlKzA5LCJodHRwOi8vZXhhbXBsZS5jb20vaXNfcm9vdCI6dHJ1ZSwiaXNzIjoiam9lIn0.CN7YijRX6Aw1n2jyI2Id1w90ja-DEMYiWixhYCyHnrZ1VfJRaFQz1bEbjjA5Fn4CLYaUG432dEYmSbS4Saokmw"),
			hf512,
			map[string]interface{}{"iss": "joe", "exp": 1300819380, "http://example.com/is_root": true},
			true,
		},
		{ // 3
			"web sample: invalid",
			[]byte("eyJ0eXAiOiJKV1QiLA0KICJhbGciOiJIUzI1NiJ9.eyJpc3MiOiJqb2UiLA0KICJleHAiOjEzMDA4MTkzODAsDQogImh0dHA6Ly9leGFtcGxlLmNvbS9pc19yb290Ijp0cnVlfQ.dBjftJeZ4CVP-mB92K27uhbUJU1p1r_wW1gFWFOEjXo"),
			hf256,
			map[string]interface{}{"iss": "joe", "exp": 1300819380, "http://example.com/is_root": true},
			false,
		},
		{ // 4
			"web sample blake2 256",
			[]byte("eyJ0eXAiOiJKV1QiLA0KICJhbGciOiJIUzI1NiJ9.eyJpc3MiOiJqb2UiLA0KICJleHAiOjEzMDA4MTkzODAsDQogImh0dHA6Ly9leGFtcGxlLmNvbS9pc19yb290Ijp0cnVlfQ.p_126D070LelL7jHk-r05gvmTLONX0Om7SVg7YufDtY"),
			blk256,
			map[string]interface{}{"iss": "joe", "exp": 1300819380, "http://example.com/is_root": true},
			true,
		},
		{ // 5
			"web sample blake2 512",
			[]byte(`eyJhbGciOiJibGsyYjUxMiIsInR5cCI6IkpXVCJ9Cg.eyJleHAiOjEzMDA4MTkzODAsImh0dHA6Ly9leGFtcGxlLmNvbS9pc19yb290Ijp0cnVlLCJpc3MiOiJqb2UifQo.q_Hx6kyt3ErQ3FuGGc-5jMBeJOsLS2C6spJ0BEo4kCy0wtT-nRc8W5too4otF2nij-urC2NZ7seZuVhlaPl7vQ`),
			blk512,
			map[string]interface{}{"iss": "joe", "exp": 1300819380, "http://example.com/is_root": true},
			true,
		},
	}
}

func TestHMACVerifyFast(t *testing.T) {

	for _, data := range hmacFastTestData {
		signing, signature, err := csjwt.SplitForVerify(data.tokenString)
		if err != nil {
			t.Fatal(err, "\n", string(data.tokenString))
		}

		err = data.method.Verify(signing, signature, csjwt.Key{})
		if err != nil && !errors.IsNotValid(err) {
			t.Errorf("[%v] Expecting a not valid error behaviour : %+v", data.name, err)
		}
		if data.valid && err != nil {
			t.Errorf("[%v] Method %s Error while verifying key: %+v", data.name, data.method, err)
		}
		if !data.valid && err == nil {
			t.Errorf("[%v] Invalid key passed validation", data.name)
		}
	}
}

func TestHMACSignFast(t *testing.T) {

	for _, data := range hmacFastTestData {
		if data.valid {
			signing, signature, err := csjwt.SplitForVerify(data.tokenString)
			if err != nil {
				t.Fatal(err, "\n", string(data.tokenString))
			}

			sig, err := data.method.Sign(signing, csjwt.Key{})
			if err != nil {
				t.Errorf("[%v] Error signing token: %v", data.name, err)
			}
			if !bytes.Equal(sig, signature) {
				t.Errorf("[%v] Incorrect signature.\nwas:\n%v\nexpecting:\n%v", data.name, string(sig), string(signature))
			}
		}
	}
}

func TestNewBlake2b256(t *testing.T) {
	s, err := csjwt.NewSigningMethodBlake2b256(csjwt.Key{
		Error: errors.NewAlreadyClosedf("Registration"),
	})
	assert.Nil(t, s)
	assert.True(t, errors.IsAlreadyClosed(err))

	s, err = csjwt.NewSigningMethodBlake2b256(csjwt.Key{})
	assert.Nil(t, s)
	assert.True(t, errors.IsEmpty(err))
}

func TestNewBlake2b512(t *testing.T) {
	s, err := csjwt.NewSigningMethodBlake2b512(csjwt.Key{
		Error: errors.NewAlreadyClosedf("Registration"),
	})
	assert.Nil(t, s)
	assert.True(t, errors.IsAlreadyClosed(err))

	s, err = csjwt.NewSigningMethodBlake2b512(csjwt.Key{})
	assert.Nil(t, s)
	assert.True(t, errors.IsEmpty(err))
}
