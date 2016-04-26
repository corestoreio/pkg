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

package csmath_test

import (
	"testing"

	"github.com/corestoreio/csfw/util/csmath"
	"github.com/stretchr/testify/assert"
)

func TestRound(t *testing.T) {
	t.Parallel()

	tests := []struct {
		val, roundOn float64
		places       int
		want         float64
	}{
		{1234.56, .5, 1, 1234.6},
		{123.445, .5, 2, 123.45},
	}

	for _, test := range tests {
		have := csmath.Round(test.val, test.roundOn, test.places)
		assert.EqualValues(t, test.want, have, "%v", test)
	}
}
