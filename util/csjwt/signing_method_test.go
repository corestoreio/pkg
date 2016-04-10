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
