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
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/corestoreio/cspkg/net"
	"github.com/corestoreio/cspkg/net/signed"
	"github.com/corestoreio/errors"
	"github.com/stretchr/testify/assert"
)

var _ signed.HeaderParseWriter = (*signed.ContentSignature)(nil)

func TestSignature_Write(t *testing.T) {

	w := httptest.NewRecorder()
	sig := signed.NewContentSignature("myKeyID", "hmac-sha1")
	sig.Write(w, []byte(`Hello Gophers`))

	const wantSig = `keyId="myKeyID",algorithm="hmac-sha1",signature="48656c6c6f20476f7068657273"`
	if have, want := w.Header().Get(net.ContentSignature), wantSig; have != want {
		t.Errorf("Have: %v Want: %v", have, want)
	}
}

// 3000000	       568 ns/op	     160 B/op	       4 allocs/op
func BenchmarkSignature_Write(b *testing.B) {
	const wantSig = `keyId="myKeyID",algorithm="hmac-sha1",signature="48656c6c6f20476f7068657273"`

	sig := signed.NewContentSignature("myKeyID", "hmac-sha1")
	sig.HeaderName = "Content-S1gnatur3"

	w := httptest.NewRecorder()
	s := []byte(`Hello Gophers`)
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		sig.Write(w, s)
	}
	if have, want := w.Header().Get("Content-S1gnatur3"), wantSig; have != want {
		b.Errorf("Have: %v Want: %v", have, want)
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
		0: {
			newReqHeader(`keyId="myKeyID",algorithm="hmac-sha1",signature="48656c6c6f20476f7068657273"`),
			"hmac-sha1",
			"myKeyID",
			"hmac-sha1",
			[]byte(`Hello Gophers`),
			nil,
		},
		1: {
			newReqHeader(`   keyId="myKeyID"	,  algorithm="hmac-sha1"	,		signature="48656c6c6f20476f7068657273"	`),
			"hmac-sha1",
			"myKeyID",
			"hmac-sha1",
			[]byte(`Hello Gophers`),
			nil,
		},
		2: {
			newReqHeader(`   k3y1d="myKeyID"	,  alg0r1thm="hmac-sha1"	,		s1gnatur3="48656c6c6f20476f7068657273"	`),
			"hmac-sha1",
			"myKeyID",
			"hmac-sha1",
			nil,
			errors.IsNotValid,
		},
		3: { // algorithm key too short
			newReqHeader(`   keyId="myKeyID"	,  alg0r1thm="hmac-sha1"	,		s1gnar3="48656c6c6f20476f7068657273"	`),
			"hmac-sha1",
			"myKeyID",
			"hmac-sha1",
			nil,
			errors.IsNotValid,
		},
		4: { // signature key too short
			newReqHeader(`   keyId="myKeyID"	,  algorithm="hmac-sha1"	,		s1gnar3="48656c6c6f20476f7068657273"	`),
			"hmac-sha1",
			"myKeyID",
			"hmac-sha1",
			nil,
			errors.IsNotValid,
		},
		5: { // algorithm key too short
			newReqHeader(`   keyId="myKeyID"	,  alg0r1m="hmac-sha1"	,		s1gnatur3="48656c6c6f20476f7068657273"	`),
			"hmac-sha1",
			"myKeyID",
			"hmac-sha1",
			nil,
			errors.IsNotValid,
		},
		6: {
			newReqHeader(`keyId="",algorithm="hmac-sha1",signature="48656c6c6f20476f7068657273"`),
			"hmac-sha1",
			"",
			"hmac-sha1",
			nil,
			errors.IsNotValid,
		},
		7: {
			newReqHeader(`keyId="",algorithm="none",signature="48656c6c6f20476f7068657273"`),
			"hmac-sha1",
			"",
			"hmac-sha1",
			nil,
			errors.IsNotValid,
		},
		8: {
			newReqHeader(``),
			"hmac-sha1",
			"",
			"hmac-sha1",
			nil,
			errors.IsNotFound,
		},
		9: {
			newReqHeader(`keyId="",algorithm="",signature="48656c6c6f20476f7068657273"`),
			"hmac-sha1",
			"",
			"",
			nil,
			errors.IsNotValid,
		},
		10: {
			newReqHeader(`keyId=,algorithm="",signature="48656c6c6f20476f7068657273"`),
			"hmac-sha1",
			"",
			"",
			nil,
			errors.IsNotValid,
		},
		11: {
			newReqHeader(`keyId="",algorithm="asdasd",signature=""`),
			"hmac-sha1",
			"",
			"",
			nil,
			errors.IsNotValid,
		},
		12: {
			newReqHeader(`keyId="",algori,thm="asdasd",signature=""`),
			"hmac-sha1",
			"",
			"",
			nil,
			errors.IsNotValid,
		},
		13: {
			newReqHeader(`keyId="",algorithm="asdasd",signature="asdas,dsad"`),
			"hmac-sha1",
			"",
			"",
			nil,
			errors.IsNotValid,
		},
		14: {
			newReqHeader(`k="",a="",s=""`),
			"hmac-sha1",
			"",
			"",
			nil,
			errors.IsNotValid,
		},
		15: {
			newReqHeader(`keyId="",algorithm="asdasd",`),
			"hmac-sha1",
			"",
			"",
			nil,
			errors.IsNotValid,
		},
		16: {
			newReqHeader(`signature="asdasd",`),
			"hmac-sha1",
			"",
			"",
			nil,
			errors.IsNotValid,
		},
		17: {
			newReqHeader(`keyId="myKeyID",algorithm="",signature="48656c6c6f20476f7068657273"`),
			"hmac-sha1",
			"myKeyID",
			"hmac-sha1",
			nil,
			errors.IsNotValid,
		},
	}
	for i, test := range tests {
		sig := signed.NewContentSignature(test.wantKeyID, test.haveAlgorithm)
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

// 1000000	      1798 ns/op	     464 B/op	       5 allocs/op
func BenchmarkSignature_Parse(b *testing.B) {

	req := httptest.NewRequest("GET", "http://corestore.io", nil)
	req.Header.Set("Content-S1gnatur3", `keyId="myKeyID",algorithm="hmac-sha1",signature="48656c6c6f20476f7068657273"`)

	sig := &signed.ContentSignature{
		KeyID: "myKeyID",
		ContentHMAC: signed.ContentHMAC{
			Algorithm:  "hmac-sha1",
			HeaderName: "Content-S1gnatur3",
		},
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
