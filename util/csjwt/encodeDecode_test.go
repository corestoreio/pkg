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

var BenchmarkencodeSegment []byte

func BenchmarkBase64(b *testing.B) {
	data := []byte(`assert.True(t, isBase64Token(tkChar), "Token is not base64 encoded! %q", tkChar)`)
	dataDec := []byte(`YXNzZXJ0LlRydWUodCwgaXNCYXNlNjRUb2tlbih0a0NoYXIpLCAiVG9rZW4gaXMgbm90IGJhc2U2NCBlbmNvZGVkISAlcSIsIHRrQ2hhcik`)

	// BenchmarkBase64/Encode-4         	20000000	       106 ns/op	     112 B/op	       1 allocs/op
	// BenchmarkBase64/Encode-4         	20000000	       105 ns/op	     112 B/op	       1 allocs/op
	b.Run("Encode", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			BenchmarkencodeSegment = csjwt.EncodeSegment(data)
		}
		if have, want := len(BenchmarkencodeSegment), len(dataDec); have != want {
			b.Fatalf("Invalid length of BenchmarkencodeSegment:\n%d => %q\n%d => %q", have, BenchmarkencodeSegment, want, dataDec)
		}
	})

	// BenchmarkBase64/Decode-4         	10000000	       160 ns/op	      80 B/op	       1 allocs/op
	// BenchmarkBase64/Decode-4         	10000000	       165 ns/op	      80 B/op	       1 allocs/op
	b.Run("Decode", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			var err error
			BenchmarkencodeSegment, err = csjwt.DecodeSegment(dataDec)
			if err != nil {
				b.Fatalf("%+v", err)
			}
		}
	})

}
