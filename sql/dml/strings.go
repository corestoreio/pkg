// Copyright 2015-present, Cyrill @ Schumacher.fm and the CoreStore contributors
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

package dml

import (
	"database/sql/driver"
	"encoding/csv"
	"regexp"
	"strings"

	"github.com/corestoreio/errors"
)

// Strings contains several functions for data moving. You can use package
// slices.String for further modifications of this slice type.
type Strings []string

// quoteEscapeRegex is the regex to match escaped characters in a string.
var quoteEscapeRegex = regexp.MustCompile(`([^\\]([\\]{2})*)\\"`)

// Scan satisfies the sql.Scanner interface for String.
func (l *Strings) Scan(src interface{}) error {
	var str string
	switch t := src.(type) {
	case []byte:
		str = string(t)
	case string:
		str = t
	default:
		return errors.NewNotValidf("[slices] String.Scan Unknown type or not yet implemented: %#v", src)
	}

	// change quote escapes for csv parser
	str = quoteEscapeRegex.ReplaceAllString(str, `$1""`)
	str = strings.Replace(str, `\\`, `\`, -1)

	// remove braces
	str = str[1 : len(str)-1]

	// bail if only one
	if len(str) == 0 {
		*l = Strings([]string{})
		return nil
	}

	// parse with csv reader
	cr := csv.NewReader(strings.NewReader(str))
	slice, err := cr.Read()
	if err != nil {
		return errors.NewNotValidf("[slices] String.Scan CSV read error: %s", err)
	}

	*l = Strings(slice)

	return nil
}

// Value satisfies the driver.Valuer interface for String.
func (l Strings) Value() (driver.Value, error) {
	v := make([]string, len(l))
	for i, s := range l {
		v[i] = `"` + strings.Replace(strings.Replace(s, `\`, `\\\`, -1), `"`, `\"`, -1) + `"`
	}
	return "{" + strings.Join(v, ",") + "}", nil
}
