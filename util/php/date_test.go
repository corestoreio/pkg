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

package php_test

import (
	"testing"

	"github.com/corestoreio/csfw/util/php"
	"github.com/stretchr/testify/assert"
)

func TestToGoFormat(t *testing.T) {

	t.Skip("todo")

	tests := []struct {
		phpf string
		gof  string
	}{
		{`l jS \of F Y h:i:s A`, `a b of `}, // e.g. Monday 8th of August 2005 03:12:46 PM
		{`l \t\h\e jS`, `x the DATE`},       // e.g. Wednesday the 15th
	}
	for _, test := range tests {
		assert.Exactly(t, test.gof, php.DateToGoFormat(test.phpf))
	}
}
