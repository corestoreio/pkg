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
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

var _ fmt.Stringer = (*SignerSlice)(nil)

func TestMethodsSlice(t *testing.T) {
	t.Parallel()
	var ms SignerSlice = []Signer{NewSigningMethodRS256(), NewSigningMethodPS256()}
	assert.Exactly(t, `RS256, PS256`, ms.String())
	assert.True(t, ms.Contains("PS256"))
	assert.False(t, ms.Contains("XS256"))

	ms = []Signer{NewSigningMethodRS256()}
	assert.Exactly(t, `RS256`, ms.String())
}

func TestMustNewSigningMethodByAlg(t *testing.T) {
	t.Parallel()
	defer func() {
		if r := recover(); r != nil {
			assert.EqualError(t, r.(error), "[csjwt] Unknown signing algorithm \"rot13\"")
		} else {
			t.Fatal("Missing a panic!")
		}
	}()
	_ = MustNewSigningMethodByAlg("rot13")
}

func TestNewSigningMethodByAlg(t *testing.T) {
	t.Parallel()
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
		assert.NotNil(t, MustNewSigningMethodByAlg(test.alg), "Index %s", test.alg)
	}
}
