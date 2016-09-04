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

package signed_test

import (
	"encoding/hex"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/corestoreio/csfw/net"
	"github.com/corestoreio/csfw/net/signed"
	"github.com/corestoreio/csfw/util/errors"
	"github.com/stretchr/testify/assert"
)

var _ signed.HTTPWriter = (*signed.Signature)(nil)
var _ signed.HTTPParser = (*signed.Signature)(nil)

func TestSignature_Write(t *testing.T) {

	w := httptest.NewRecorder()
	sig := signed.Signature{
		EncodeFn:  hex.EncodeToString,
		KeyID:     "myKeyID",
		Algorithm: "hmac-sha1",
	}
	sig.Write(w, []byte(`Hello Gophers`))

	const wantSig = `keyId="myKeyID",algorithm="hmac-sha1",signature="48656c6c6f20476f7068657273"`
	if have, want := w.Header().Get(net.ContentSignature), wantSig; have != want {
		t.Errorf("Have: %v Want: %v", have, want)
	}
}

func TestSignature_Parse(t *testing.T) {
	var newReqHeader = func(value string) *http.Request {
		req := httptest.NewRequest("GET", "http://corestore.io", nil)
		req.Header.Set(net.ContentSignature, value)
		return req
	}
	tests := []struct {
		req           *http.Request
		haveAlgorithm string
		wantKeyID     string
		wantAlgorithm string
		wantSignature []byte
		wantErrBhf    errors.BehaviourFunc
	}{
		{
			newReqHeader(`keyId="myKeyID",algorithm="hmac-sha1",signature="48656c6c6f20476f7068657273"`),
			"hmac-sha1",
			"",
			"hmac-sha1",
			[]byte(`Hello Gophers`),
			nil,
		},
		{
			newReqHeader(`   keyId="myKeyID"	,  algorithm="hmac-sha1"	,		signature="48656c6c6f20476f7068657273"	`),
			"hmac-sha1",
			"",
			"hmac-sha1",
			[]byte(`Hello Gophers`),
			nil,
		},
		{
			newReqHeader(`   k3y1d="myKeyID"	,  alg0r1thm="hmac-sha1"	,		s1gnatur3="48656c6c6f20476f7068657273"	`),
			"hmac-sha1",
			"",
			"hmac-sha1",
			[]byte(`Hello Gophers`),
			nil,
		},
		{
			newReqHeader(`keyId="",algorithm="hmac-sha1",signature="48656c6c6f20476f7068657273"`),
			"hmac-sha1",
			"",
			"hmac-sha1",
			[]byte(`Hello Gophers`),
			nil,
		},
		{
			newReqHeader(`keyId="",algorithm="none",signature="48656c6c6f20476f7068657273"`),
			"hmac-sha1",
			"",
			"hmac-sha1",
			nil,
			errors.IsNotValid,
		},
		{
			newReqHeader(``),
			"hmac-sha1",
			"",
			"hmac-sha1",
			nil,
			errors.IsNotFound,
		},
		{
			newReqHeader(`keyId="",algorithm="",signature="48656c6c6f20476f7068657273"`),
			"hmac-sha1",
			"",
			"",
			nil,
			errors.IsNotValid,
		},
		{
			newReqHeader(`keyId=,algorithm="",signature="48656c6c6f20476f7068657273"`),
			"hmac-sha1",
			"",
			"",
			nil,
			errors.IsNotValid,
		},
		{
			newReqHeader(`keyId="",algorithm="asdasd",signature=""`),
			"hmac-sha1",
			"",
			"",
			nil,
			errors.IsNotValid,
		},
		{
			newReqHeader(`keyId="",algori,thm="asdasd",signature=""`),
			"hmac-sha1",
			"",
			"",
			nil,
			errors.IsNotValid,
		},
		{
			newReqHeader(`keyId="",algorithm="asdasd",signature="asdas,dsad"`),
			"hmac-sha1",
			"",
			"",
			nil,
			errors.IsNotValid,
		},
		{
			newReqHeader(`keyId="",algorithm="asdasd",`),
			"hmac-sha1",
			"",
			"",
			nil,
			errors.IsNotValid,
		},
		{
			newReqHeader(`signature="asdasd",`),
			"hmac-sha1",
			"",
			"",
			nil,
			errors.IsNotValid,
		},
	}
	for i, test := range tests {
		sig := &signed.Signature{
			Algorithm: test.haveAlgorithm,
			DecodeFn:  hex.DecodeString,
		}
		haveSig, haveErr := sig.Parse(test.req)
		if test.wantErrBhf != nil {
			assert.Nil(t, haveSig, "Index %d", i)
			assert.True(t, test.wantErrBhf(haveErr), "Error: %+v", haveErr)
			continue
		}
		assert.Exactly(t, test.wantKeyID, sig.KeyID, "Index %d", i)
		assert.Exactly(t, test.wantAlgorithm, sig.Algorithm, "Index %d", i)
		assert.Exactly(t, string(test.wantSignature), string(haveSig), "Index %d", i)
		assert.NoError(t, haveErr, "Index %d: %+v", i, haveErr)
	}
}

// 1000000	      2087 ns/op	     448 B/op	       4 allocs/op
func BenchmarkSignature_Parse(b *testing.B) {

	req := httptest.NewRequest("GET", "http://corestore.io", nil)
	req.Header.Set("Content-S1gnatur3", `keyId="myKeyID",algorithm="hmac-sha1",signature="48656c6c6f20476f7068657273"`)

	sig := &signed.Signature{
		Algorithm: "hmac-sha1",
		HeaderKey: "Content-S1gnatur3",
		DecodeFn:  hex.DecodeString,
	}

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		sig, err := sig.Parse(req)
		if err != nil {
			b.Fatalf("%+v", err)
		}
		if len(sig) < 3 {
			b.Fatal("Invalid length of signature")
		}
	}
}
