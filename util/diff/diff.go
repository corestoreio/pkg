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

package diff

import (
	"github.com/pmezard/go-difflib/difflib"
)

// Unified compares two strings and output the differences.
func Unified(a, b string) (string, error) {
	diff := difflib.UnifiedDiff{
		A:        difflib.SplitLines(a),
		B:        difflib.SplitLines(b),
		FromFile: "Original",
		ToFile:   "Current",
		Context:  1,
	}
	return difflib.GetUnifiedDiffString(diff)
}

// MustUnified is the same as Unified but panics on error.
func MustUnified(a, b string) string {
	d, err := Unified(a, b)
	if err != nil {
		panic(err)
	}
	return d
}
