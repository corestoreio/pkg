// Copyright 2015, Cyrill @ Schumacher.fm and the CoreStore contributors
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

package dbr

import "strings"

// TableColumnQuote prefixes all columns with a table name/alias and puts quotes around them.
func TableColumnQuote(t string, cols ...string) []string {
	//r := make([]string, len(cols), len(cols))
	for i, c := range cols {
		if strings.Contains(c, Quote) {
			cols[i] = c
		} else {
			cols[i] = Quote + t + Quote + "." + Quote + c + Quote
		}
	}
	return cols
}

// IfNullAs returns IFNULL(t1.c1,t2.c2) AS as
func IfNullAs(t1, c1, t2, c2, as string) string {
	return "IFNULL(" + Quote + t1 + Quote + "." + Quote + c1 + Quote + ", " + Quote + t2 + Quote + "." + Quote + c2 + Quote + ") AS " + Quote + as + Quote
}
