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

	"github.com/corestoreio/pkg/util/csjwt"
)

func init() {
	_ = csjwt.NewSigningMethodHS256()
	_ = csjwt.NewSigningMethodHS384()
	_ = csjwt.NewSigningMethodHS512()
}

var hmacTestData = []struct {
	name        string
	tokenString []byte
	method      csjwt.Signer
	claims      map[string]interface{}
	valid       bool
}{
	{
		"web sample",
		[]byte("eyJ0eXAiOiJKV1QiLA0KICJhbGciOiJIUzI1NiJ9.eyJpc3MiOiJqb2UiLA0KICJleHAiOjEzMDA4MTkzODAsDQogImh0dHA6Ly9leGFtcGxlLmNvbS9pc19yb290Ijp0cnVlfQ.dBjftJeZ4CVP-mB92K27uhbUJU1p1r_wW1gFWFOEjXk"),
		csjwt.NewSigningMethodHS256(),
		map[string]interface{}{"iss": "joe", "exp": 1300819380, "http://example.com/is_root": true},
		true,
	},
	{
		"HS384",
		[]byte("eyJhbGciOiJIUzM4NCIsInR5cCI6IkpXVCJ9.eyJleHAiOjEuMzAwODE5MzhlKzA5LCJodHRwOi8vZXhhbXBsZS5jb20vaXNfcm9vdCI6dHJ1ZSwiaXNzIjoiam9lIn0.KWZEuOD5lbBxZ34g7F-SlVLAQ_r5KApWNWlZIIMyQVz5Zs58a7XdNzj5_0EcNoOy"),
		csjwt.NewSigningMethodHS384(),
		map[string]interface{}{"iss": "joe", "exp": 1300819380, "http://example.com/is_root": true},
		true,
	},
	{
		"HS512",
		[]byte("eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJleHAiOjEuMzAwODE5MzhlKzA5LCJodHRwOi8vZXhhbXBsZS5jb20vaXNfcm9vdCI6dHJ1ZSwiaXNzIjoiam9lIn0.CN7YijRX6Aw1n2jyI2Id1w90ja-DEMYiWixhYCyHnrZ1VfJRaFQz1bEbjjA5Fn4CLYaUG432dEYmSbS4Saokmw"),
		csjwt.NewSigningMethodHS512(),
		map[string]interface{}{"iss": "joe", "exp": 1300819380, "http://example.com/is_root": true},
		true,
	},
	{
		"web sample: invalid",
		[]byte("eyJ0eXAiOiJKV1QiLA0KICJhbGciOiJIUzI1NiJ9.eyJpc3MiOiJqb2UiLA0KICJleHAiOjEzMDA4MTkzODAsDQogImh0dHA6Ly9leGFtcGxlLmNvbS9pc19yb290Ijp0cnVlfQ.dBjftJeZ4CVP-mB92K27uhbUJU1p1r_wW1gFWFOEjXo"),
		csjwt.NewSigningMethodHS256(),
		map[string]interface{}{"iss": "joe", "exp": 1300819380, "http://example.com/is_root": true},
		false,
	},
}

func TestHMACVerify(t *testing.T) {

	for _, data := range hmacTestData {
		signing, signature, err := csjwt.SplitForVerify(data.tokenString)
		if err != nil {
			t.Fatal(err, "\n", string(data.tokenString))
		}

		err = data.method.Verify(signing, signature, csjwt.WithPassword(hmacTestKey))
		if data.valid && err != nil {
			t.Errorf("[%v] Method %s Error while verifying key: %v", data.name, data.method, err)
		}
		if !data.valid && err == nil {
			t.Errorf("[%v] Invalid key passed validation", data.name)
		}
	}
}

func TestHMACSign(t *testing.T) {

	for _, data := range hmacTestData {
		if data.valid {
			signing, signature, err := csjwt.SplitForVerify(data.tokenString)
			if err != nil {
				t.Fatal(err, "\n", string(data.tokenString))
			}

			sig, err := data.method.Sign(signing, csjwt.WithPassword(hmacTestKey))
			if err != nil {
				t.Errorf("[%v] Error signing token: %v", data.name, err)
			}
			if !bytes.Equal(sig, signature) {
				t.Errorf("[%v] Incorrect signature.\nwas:\n%v\nexpecting:\n%v", data.name, string(sig), string(signature))
			}
		}
	}
}

func BenchmarkHS256Signing(b *testing.B) {
	benchmarkSigning(b, csjwt.NewSigningMethodHS256(), csjwt.WithPassword(hmacTestKey))
}

func BenchmarkHS384Signing(b *testing.B) {
	benchmarkSigning(b, csjwt.NewSigningMethodHS384(), csjwt.WithPassword(hmacTestKey))
}

func BenchmarkHS512Signing(b *testing.B) {
	benchmarkSigning(b, csjwt.NewSigningMethodHS512(), csjwt.WithPassword(hmacTestKey))
}

func BenchmarkHS256SigningFast(b *testing.B) {
	hf, err := csjwt.NewSigningMethodHS256Fast(csjwt.WithPassword(hmacTestKey))
	if err != nil {
		b.Fatal(err)
	}
	benchmarkSigning(b, hf, csjwt.Key{})
}

func BenchmarkHS384SigningFast(b *testing.B) {
	hf, err := csjwt.NewSigningMethodHS384Fast(csjwt.WithPassword(hmacTestKey))
	if err != nil {
		b.Fatal(err)
	}
	benchmarkSigning(b, hf, csjwt.Key{})
}

func BenchmarkHS512SigningFast(b *testing.B) {
	hf, err := csjwt.NewSigningMethodHS512Fast(csjwt.WithPassword(hmacTestKey))
	if err != nil {
		b.Fatal(err)
	}
	benchmarkSigning(b, hf, csjwt.Key{})
}
func BenchmarkBlake2b256SigningFast(b *testing.B) {
	hf, err := csjwt.NewSigningMethodBlake2b256(csjwt.WithPassword(hmacTestKey))
	if err != nil {
		b.Fatal(err)
	}
	benchmarkSigning(b, hf, csjwt.Key{})
}
func BenchmarkBlake2b512SigningFast(b *testing.B) {
	hf, err := csjwt.NewSigningMethodBlake2b512(csjwt.WithPassword(hmacTestKey))
	if err != nil {
		b.Fatal(err)
	}
	benchmarkSigning(b, hf, csjwt.Key{})
}

func BenchmarkHS256Verify(b *testing.B) {
	signing, signature, err := csjwt.SplitForVerify(hmacTestData[0].tokenString) // HS256 token
	if err != nil {
		b.Fatal(err)
	}
	benchmarkMethodVerify(b, csjwt.NewSigningMethodHS256(), signing, signature, csjwt.WithPassword(hmacTestKey))
}
func BenchmarkHS384Verify(b *testing.B) {
	signing, signature, err := csjwt.SplitForVerify(hmacTestData[1].tokenString) // HS384 token
	if err != nil {
		b.Fatal(err)
	}
	benchmarkMethodVerify(b, csjwt.NewSigningMethodHS384(), signing, signature, csjwt.WithPassword(hmacTestKey))
}
func BenchmarkHS512Verify(b *testing.B) {
	signing, signature, err := csjwt.SplitForVerify(hmacTestData[2].tokenString) // HS512 token
	if err != nil {
		b.Fatal(err)
	}
	benchmarkMethodVerify(b, csjwt.NewSigningMethodHS512(), signing, signature, csjwt.WithPassword(hmacTestKey))
}

func BenchmarkHS256VerifyFast(b *testing.B) {
	signing, signature, err := csjwt.SplitForVerify(hmacTestData[0].tokenString) // HS256 token
	if err != nil {
		b.Fatal(err)
	}
	hf, err := csjwt.NewSigningMethodHS256Fast(csjwt.WithPassword(hmacTestKey))
	if err != nil {
		b.Fatal(err)
	}
	benchmarkMethodVerify(b, hf, signing, signature, csjwt.Key{})
}
func BenchmarkHS384VerifyFast(b *testing.B) {
	signing, signature, err := csjwt.SplitForVerify(hmacTestData[1].tokenString) // HS384 token
	if err != nil {
		b.Fatal(err)
	}
	hf, err := csjwt.NewSigningMethodHS384Fast(csjwt.WithPassword(hmacTestKey))
	if err != nil {
		b.Fatal(err)
	}
	benchmarkMethodVerify(b, hf, signing, signature, csjwt.Key{})
}
func BenchmarkHS512VerifyFast(b *testing.B) {
	signing, signature, err := csjwt.SplitForVerify(hmacTestData[2].tokenString) // HS512 token
	if err != nil {
		b.Fatal(err)
	}
	hf, err := csjwt.NewSigningMethodHS512Fast(csjwt.WithPassword(hmacTestKey))
	if err != nil {
		b.Fatal(err)
	}
	benchmarkMethodVerify(b, hf, signing, signature, csjwt.Key{})
}
func BenchmarkBlake2b256VerifyFast(b *testing.B) {
	signing, signature, err := csjwt.SplitForVerify(hmacFastTestData[4].tokenString)
	if err != nil {
		b.Fatal(err)
	}
	hf, err := csjwt.NewSigningMethodBlake2b256(csjwt.WithPassword(hmacTestKey))
	if err != nil {
		b.Fatal(err)
	}
	benchmarkMethodVerify(b, hf, signing, signature, csjwt.Key{})
}
func BenchmarkBlake2b512VerifyFast(b *testing.B) {
	signing, signature, err := csjwt.SplitForVerify(hmacFastTestData[5].tokenString)
	if err != nil {
		b.Fatal(err)
	}
	hf, err := csjwt.NewSigningMethodBlake2b512(csjwt.WithPassword(hmacTestKey))
	if err != nil {
		b.Fatal(err)
	}
	benchmarkMethodVerify(b, hf, signing, signature, csjwt.Key{})
}
