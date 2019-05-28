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

package csjwt_test

import (
	"encoding/gob"
	"regexp"
	"testing"
	"time"

	"github.com/corestoreio/pkg/util/assert"
	"github.com/corestoreio/pkg/util/csjwt"
	"github.com/corestoreio/pkg/util/csjwt/jwtclaim"
)

func isBase64Token(str []byte) bool {
	r, _ := regexp.Compile(`^[A-Za-z0-9\-_\.]+$`)
	return r.Match(str)
}

func init() {
	gob.Register(csjwt.NewHead())
	gob.Register(jwtclaim.NewStore())
}

func TestGobEncoding(t *testing.T) {

	storeClaim := jwtclaim.NewStore()

	gobEncDec := csjwt.NewGobEncoding(csjwt.NewHead(), storeClaim)

	storeClaim.Store = "ch-en"
	storeClaim.ID = "2342-234345-234234-23435"
	storeClaim.ExpiresAt = time.Now().Add(time.Minute * 2).Unix()
	storeClaim.IssuedAt = time.Now().Unix()

	tk := csjwt.NewToken(storeClaim)
	tk.Serializer = gobEncDec

	m := csjwt.NewSigningMethodHS512()
	pw := csjwt.WithPasswordRandom()
	tkChar, err := tk.SignedString(m, pw)
	assert.NoError(t, err)

	// correct old token:
	// D_-EAQVIUzUxMgEDSldUAA.Mv-GAQP8udiHwAEYMjM0Mi0yMzQzNDUtMjM0MjM0LTIzNDM1Afy52IbQAAEFY2gtZW4A.VvwSTTBuaY7kxbbrZ44YXjFDoyhLXFEpIYjHI-mOOKzhZKpBiq3z3qqyqwyYlYxLzA914PGMmNiQCk1UdY7HZg

	// invalid new token:
	// D_-EAQVIUzUxMgEDSldUAA.eyJzdG9yZSI6ImNoLWVuIiwiZXhwIjoxNTU4OTg3NjczLCJqdGkiOiIyMzQyLTIzNDM0NS0yMzQyMzQtMjM0MzUiLCJpYXQiOjE1NTg5ODc1NTN9.JrFW-ex4BW1e4ICTGsnWRoIm-OLN_TCdMINy5q74v47G7BWXVyIFtnNzcuE6BdWaZhTcSp4J-3WTyKw5XQgIHQ

	t.Logf("gob %q", tkChar)

	if have, want := len(tkChar), 178; have != want {
		t.Errorf("Gob length tkChar mismatch: Have: %d Want: %d", have, want)
	}

	// check if it is base64 encoded
	assert.True(t, isBase64Token(tkChar), "Token is not base64 encoded! %q", tkChar)

	vrf := csjwt.NewVerification(m)
	vrf.Deserializer = gobEncDec

	newTk := csjwt.NewToken(jwtclaim.NewStore())

	assert.NoError(t, vrf.Parse(newTk, tkChar, csjwt.NewKeyFunc(m, pw)))

	haveStoreClaim := newTk.Claims.(*jwtclaim.Store)
	assert.Exactly(t, "ch-en", haveStoreClaim.Store)
	assert.Exactly(t, "2342-234345-234234-23435", haveStoreClaim.ID)
}

func BenchmarkTokenDecode(b *testing.B) {

	var testRunner = func(b *testing.B, encDec interface {
		csjwt.Serializer
		csjwt.Deserializer
	}) {
		storeClaim := jwtclaim.NewStore()

		storeClaim.Store = "ch-de"
		storeClaim.ID = "2342-987325-234234-23435"
		storeClaim.ExpiresAt = time.Now().Add(time.Minute * 2).Unix()
		storeClaim.IssuedAt = time.Now().Unix()

		tk := csjwt.NewToken(storeClaim)
		tk.Serializer = encDec

		m := csjwt.NewSigningMethodHS256()
		pw := csjwt.WithPasswordRandom()
		tkChar, err := tk.SignedString(m, pw)
		if err != nil {
			b.Fatalf("%+v", err)
		}

		vrf := csjwt.NewVerification(m)
		vrf.Deserializer = encDec

		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			newTk := csjwt.NewToken(jwtclaim.NewStore())
			if err := vrf.Parse(newTk, tkChar, csjwt.NewKeyFunc(m, pw)); err != nil {
				b.Fatalf("%+v", err)
			}
			haveStoreClaim := newTk.Claims.(*jwtclaim.Store)
			if have, want := haveStoreClaim.Store, "ch-de"; have != want {
				b.Errorf("Have: %v Want: %v", have, want)
			}
		}
	}

	b.Run("Gob_HS256", func(b *testing.B) {
		testRunner(b, csjwt.NewGobEncoding(csjwt.NewHead(), jwtclaim.NewStore()))
	})
	b.Run("Json_HS256", func(b *testing.B) {
		testRunner(b, nil) // falls back to default JSON serializer
	})
}
