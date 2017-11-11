// Copyright 2015-2017, Cyrill @ Schumacher.fm and the CoreStore contributors
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

package eav

import "github.com/corestoreio/cspkg/util/bufferpool"

// DefaultScopeNames specifies the name of the scopes used in all EAV* function
// to generate scope based hierarchical fall backs.
var DefaultScopeNames = [...]string{"Store", "Group", "Website", "Default"}

// IfNull creates a nested IFNULL SQL statement when a scope based fall back
// hierarchy is required. Alias argument will be used as a prefix for the alias
// table name and as the final alias name.
func IfNull(alias, columnName, defaultVal string, scopeNames ...string) string {
	buf := bufferpool.Get()
	defer bufferpool.Put(buf)

	if len(scopeNames) == 0 {
		scopeNames = DefaultScopeNames[:]
	}

	brackets := 0
	for _, n := range scopeNames {
		buf.WriteString("IFNULL(")
		buf.WriteRune('`')
		buf.WriteString(alias)
		buf.WriteString(n)
		buf.WriteRune('`')
		buf.WriteRune('.')
		buf.WriteRune('`')
		buf.WriteString(columnName)
		buf.WriteRune('`')
		if brackets < len(scopeNames)-1 {
			buf.WriteRune(',')
		}
		brackets++
	}

	if defaultVal == "" {
		defaultVal = `''`
	}
	buf.WriteRune(',')
	buf.WriteString(defaultVal)
	for i := 0; i < brackets; i++ {
		buf.WriteRune(')')
	}
	buf.WriteString(" AS ")
	buf.WriteRune('`')
	buf.WriteString(alias)
	buf.WriteRune('`')
	return buf.String()
}
