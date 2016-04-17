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

package csjwt

import (
	"testing"

	"errors"

	"github.com/stretchr/testify/assert"
)

func TestVerificationGetMethod(t *testing.T) {
	tests := []struct {
		vf           *Verification
		token        *Token
		wantSigner   Signer
		wantErr      error
		haveLastUsed uint32
	}{
		{
			&Verification{},
			nil,
			nil,
			errors.New(`[csjwt] No methods supplied to the Verfication Method slice`),
			0,
		},
		{
			NewVerification(),
			&Token{},
			nil,
			errors.New(`[csjwt] Cannot find alg entry in token header: <nil>`),
			0,
		},
		{
			NewVerification(NewSigningMethodHS512()),
			&Token{
				Header: NewHead("RS4"),
			},
			nil,
			errors.New(`[csjwt] Algorithm "RS4" not found in method list "HS512"`),
			0,
		},
		{
			NewVerification(NewSigningMethodPS256(), NewSigningMethodRS512(), NewSigningMethodHS512()),
			&Token{
				Header: NewHead(HS512),
			},
			NewSigningMethodHS512(),
			nil,
			2,
		},
	}
	for i, test := range tests {
		haveSigner, haveErr := test.vf.getMethod(test.token)
		if test.wantErr != nil {
			assert.EqualError(t, haveErr, test.wantErr.Error(), "Index %d", i)
			assert.Nil(t, haveSigner, "Index %d", i)
			continue
		}
		assert.NoError(t, haveErr, "Index %d", i)
		assert.Exactly(t, test.wantSigner, haveSigner, "Index %d", i)
	}
}

// BenchmarkVerificationGetMethod-4	50000000	        37.8 ns/op	       0 B/op	       0 allocs/op
func BenchmarkVerificationGetMethod(b *testing.B) {

	vf := NewVerification(NewSigningMethodPS256(), NewSigningMethodRS384(), NewSigningMethodHS512(), NewSigningMethodES256(), NewSigningMethodHS256())

	tokens := [2]*Token{
		{
			Header: NewHead(HS256),
		},
		{
			Header: NewHead(RS384),
		},
	}
	wantAlg := [2]string{
		"HS256",
		"RS384",
	}

	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		var i int
		for pb.Next() {
			m, err := vf.getMethod(tokens[i%2])
			if err != nil {
				b.Fatal(err)
			}
			if have, want := m.Alg(), wantAlg[i%2]; have != want {
				b.Fatalf("Have %s Want %s", have, want)
			}
			i++
		}
	})
}
