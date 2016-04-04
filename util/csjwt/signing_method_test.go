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

func TestAlgorithm(t *testing.T) {
	tests := []struct {
		a    csjwt.Algorithm
		s    string
		want string
	}{
		{0, "Algorithm(0)", "Algorithm(0)"},
		{0, "Algorithm(-10)", "Algorithm(0)"},
		{0, "Algorithm(x'x)", "Algorithm(0)"},
		{csjwt.ES256, "ES256", "ES256"},
		{csjwt.HS256, "HS256", "HS256"},
		{csjwt.RS512, "RS512", "RS512"},
		{4711, "Algorithm(4711)", "Algorithm(4711)"},
		{0, "Algorithm(47232876482736486723411)", "Algorithm(0)"},
	}
	for i, test := range tests {
		assert.Exactly(t, test.want, test.a.String(), "Index %d", i)
		assert.Exactly(t, test.a, csjwt.ToAlgorithm(test.s), "Index %d", i)
	}
}
