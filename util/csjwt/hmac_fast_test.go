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
	"testing"

	"io/ioutil"

	"github.com/corestoreio/csfw/util/csjwt"
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

	hf256, err := csjwt.NewHMACFast256(csjwt.WithPassword(hmacTestKey))
	if err != nil {
		panic(err)
	}

	hf384, err := csjwt.NewHMACFast384(csjwt.WithPassword(hmacTestKey))
	if err != nil {
		panic(err)
	}

	hf512, err := csjwt.NewHMACFast512(csjwt.WithPassword(hmacTestKey))
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
		{
			"web sample",
			[]byte("eyJ0eXAiOiJKV1QiLA0KICJhbGciOiJIUzI1NiJ9.eyJpc3MiOiJqb2UiLA0KICJleHAiOjEzMDA4MTkzODAsDQogImh0dHA6Ly9leGFtcGxlLmNvbS9pc19yb290Ijp0cnVlfQ.dBjftJeZ4CVP-mB92K27uhbUJU1p1r_wW1gFWFOEjXk"),
			hf256,
			map[string]interface{}{"iss": "joe", "exp": 1300819380, "http://example.com/is_root": true},
			true,
		},
		{
			"HS384",
			[]byte("eyJhbGciOiJIUzM4NCIsInR5cCI6IkpXVCJ9.eyJleHAiOjEuMzAwODE5MzhlKzA5LCJodHRwOi8vZXhhbXBsZS5jb20vaXNfcm9vdCI6dHJ1ZSwiaXNzIjoiam9lIn0.KWZEuOD5lbBxZ34g7F-SlVLAQ_r5KApWNWlZIIMyQVz5Zs58a7XdNzj5_0EcNoOy"),
			hf384,
			map[string]interface{}{"iss": "joe", "exp": 1300819380, "http://example.com/is_root": true},
			true,
		},
		{
			"HS512",
			[]byte("eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJleHAiOjEuMzAwODE5MzhlKzA5LCJodHRwOi8vZXhhbXBsZS5jb20vaXNfcm9vdCI6dHJ1ZSwiaXNzIjoiam9lIn0.CN7YijRX6Aw1n2jyI2Id1w90ja-DEMYiWixhYCyHnrZ1VfJRaFQz1bEbjjA5Fn4CLYaUG432dEYmSbS4Saokmw"),
			hf512,
			map[string]interface{}{"iss": "joe", "exp": 1300819380, "http://example.com/is_root": true},
			true,
		},
		{
			"web sample: invalid",
			[]byte("eyJ0eXAiOiJKV1QiLA0KICJhbGciOiJIUzI1NiJ9.eyJpc3MiOiJqb2UiLA0KICJleHAiOjEzMDA4MTkzODAsDQogImh0dHA6Ly9leGFtcGxlLmNvbS9pc19yb290Ijp0cnVlfQ.dBjftJeZ4CVP-mB92K27uhbUJU1p1r_wW1gFWFOEjXo"),
			hf256,
			map[string]interface{}{"iss": "joe", "exp": 1300819380, "http://example.com/is_root": true},
			false,
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
		if data.valid && err != nil {
			t.Errorf("[%v] Method %s Error while verifying key: %v", data.name, data.method, err)
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
