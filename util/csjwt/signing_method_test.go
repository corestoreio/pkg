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

package csjwt

import (
	"fmt"
	"testing"

	"github.com/corestoreio/errors"
	"github.com/stretchr/testify/assert"
)

var _ fmt.Stringer = (*SignerSlice)(nil)

func TestMethodsSlice(t *testing.T) {

	var ms SignerSlice = []Signer{NewSigningMethodRS256(), NewSigningMethodPS256()}
	assert.Exactly(t, `RS256, PS256`, ms.String())
	assert.True(t, ms.Contains("PS256"))
	assert.False(t, ms.Contains("XS256"))

	ms = []Signer{NewSigningMethodRS256()}
	assert.Exactly(t, `RS256`, ms.String())
}

func TestMustSigningMethodFactory(t *testing.T) {

	defer func() {
		if r := recover(); r != nil {
			err, ok := r.(error)
			if !ok {
				t.Fatalf("Missing error interface: %#v", r)
			}
			assert.True(t, errors.NotSupported.Match(err), "Error: %s", err)
		} else {
			t.Fatal("Missing a panic!")
		}
	}()
	_ = MustSigningMethodFactory("rot13")
}

func TestSigningMethodFactory(t *testing.T) {

	tests := []struct {
		alg string
	}{
		{ES256},
		{ES384},
		{ES512},
		{HS256},
		{HS384},
		{HS512},
		{PS256},
		{PS384},
		{PS512},
		{RS256},
		{RS384},
		{RS512},
	}
	for _, test := range tests {
		assert.NotNil(t, MustSigningMethodFactory(test.alg), "Index %s", test.alg)
	}
}
