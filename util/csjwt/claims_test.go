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

	"github.com/corestoreio/csfw/util/csjwt/jwtclaim"
	"github.com/stretchr/testify/assert"
)

var _ Header = (*Head)(nil)
var _ fmt.Stringer = (*Head)(nil)
var _ fmt.GoStringer = (*Head)(nil)

func TestNewHeadStringer(t *testing.T) {
	t.Parallel()
	var h Header
	h = NewHead("Quantum")
	assert.Exactly(t, "csjwt.NewHead(\"Quantum\")", fmt.Sprintf("%s", h))
	assert.Exactly(t, "csjwt.NewHead(\"Quantum\")", fmt.Sprintf("%v", h))
	assert.Exactly(t, "csjwt.NewHead(\"Quantum\")", fmt.Sprintf("%#v", h))
}

func TestNewHead(t *testing.T) {
	t.Parallel()
	var h Header
	h = NewHead("X")
	assert.Exactly(t, "X", h.Alg())
	assert.Exactly(t, ContentTypeJWT, h.Typ())
}

func TestHeadSetGet(t *testing.T) {
	t.Parallel()
	var h Header
	h = NewHead("X")

	assert.NoError(t, h.Set(jwtclaim.HeaderAlg, "Y"))
	g, err := h.Get(jwtclaim.HeaderAlg)
	assert.NoError(t, err)
	assert.Exactly(t, "Y", g)

	assert.NoError(t, h.Set(jwtclaim.HeaderTyp, "JWE"))
	g, err = h.Get(jwtclaim.HeaderTyp)
	assert.NoError(t, err)
	assert.Exactly(t, "JWE", g)

	assert.EqualError(t, h.Set("x", "y"), "[csjwt] Header \"x\" not yet supported. Please switch to type jwtclaim.HeadSegments.")
	g, err = h.Get("x")
	assert.EqualError(t, err, "[csjwt] Header \"x\" not yet supported. Please switch to type jwtclaim.HeadSegments.")
	assert.Empty(t, g)
}
