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

package phpdate_test

import (
	"testing"

	"github.com/corestoreio/csfw/util/php/phpdate"
	"github.com/stretchr/testify/assert"
)

func TestToGoFormat(t *testing.T) {
	t.Parallel()

	tests := []struct {
		phpf string
		gof  string
	}{
		{`Y-m-d H:i:s`, `2006-01-02 15:04:05`},
		{`F j, Y, g:i a`, `January 2, 2006, 3:04 pm`},                     // e.g. March 10, 2001, 5:16 pm
		{`l jS \of F Y h:i:s A`, `Monday 2S of January 2006 03:04:05 PM`}, // e.g. Monday 8th of August 2005 03:12:46 PM
		{`l \t\h\e jS`, `Monday the 2S`},                                  // e.g. Wednesday the 15th
		{`D M j G:i:s T Y`, `Mon Jan 2 15:04:05 MST 2006`},                // Sat Mar 10 17:16:18 MST 2001
	}
	for i, test := range tests {
		assert.Exactly(t, test.gof, phpdate.ToGoFormat(test.phpf), "Index %d", i)
	}
}
