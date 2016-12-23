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

package csdb

import (
	"github.com/corestoreio/errors"
)

// IsValidIdentifier checks the permissible syntax for identifiers. Certain
// objects within MySQL, including database, table, index, column, alias, view,
// stored procedure, partition, tablespace, and other object names are known as
// identifiers. ASCII: [0-9,a-z,A-Z$_] (basic Latin letters, digits 0-9, dollar,
// underscore) Max length 63 characters. Returns errors.NotValid
//
// http://dev.mysql.com/doc/refman/5.7/en/identifiers.html
func IsValidIdentifier(names ...string) error {
	for _, name := range names {
		if len(name) > 63 || name == "" {
			return errors.NewNotValidf("[csdb] Incorrect identifier. Too long or empty: %q", name)
		}

		for _, r := range name {
			var ok bool
			switch {
			case '0' <= r && r <= '9':
				ok = true
			case 'a' <= r && r <= 'z', 'A' <= r && r <= 'Z':
				ok = true
			case r == '$', r == '_':
				ok = true
			}
			if !ok {
				return errors.NewNotValidf("[csdb] Invalid character %q in name %q", string(r), name)
			}
		}
	}
	return nil
}
