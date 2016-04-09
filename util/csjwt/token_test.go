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
	"testing"

	"github.com/corestoreio/csfw/util/csjwt"
	"github.com/stretchr/testify/assert"
)

func TestTokenAlg(t *testing.T) {
	t.Parallel()
	tests := []struct {
		tok     csjwt.Token
		wantAlg string
	}{
		{csjwt.NewToken(nil), ""},
		{
			csjwt.Token{
				Header: map[string]interface{}{
					"alg": 3,
				},
			},
			"",
		},
		{
			csjwt.Token{
				Header: map[string]interface{}{
					"alg": "Gopher",
				},
			},
			"Gopher",
		},
	}
	for i, test := range tests {
		assert.Exactly(t, test.wantAlg, test.tok.Alg(), "Index %d", i)
	}
}
