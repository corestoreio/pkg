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

package dbr

import (
	"unicode/utf8"

	"github.com/corestoreio/csfw/util/null"
	"github.com/corestoreio/errors"
)

type argStringNull struct {
	str null.String
}

func (a argStringNull) toIFace(args *[]interface{}) {
	if a.str.Valid {
		*args = append(*args, a.str.String)
	} else {
		*args = append(*args, nil)
	}
}

func (a argStringNull) writeTo(w queryWriter, _ int) error {
	if a.str.Valid {
		if !utf8.ValidString(a.str.String) {
			return errors.NewNotValidf("[dbr] Argument.WriteTo: StringNull is not UTF-8: %q", a.str.String)
		}
		dialect.EscapeString(w, a.str.String)
	} else {
		w.WriteString("NULL")
	}
	return nil
}

func (a argStringNull) len() int           { return 1 }
func (a argStringNull) INClause() Argument { return a }
func (a argStringNull) options() uint {
	if a.str.Valid {
		return 0
	}
	return argOptionNull
}

type argStringNulls struct {
	opts uint
	data []null.String
}

func (a argStringNulls) toIFace(args *[]interface{}) {
	for _, s := range a.data {
		if s.Valid {
			*args = append(*args, s.String)
		} else {
			*args = append(*args, nil)
		}
	}
}

func (a argStringNulls) writeTo(w queryWriter, pos int) error {
	if a.options() == 0 {
		if s := a.data[pos]; s.Valid {
			if !utf8.ValidString(s.String) {
				return errors.NewNotValidf("[dbr] Argument.WriteTo: String is not UTF-8: %q", s.String)
			}
			dialect.EscapeString(w, s.String)
			return nil
		}
		_, err := w.WriteString("NULL")
		return err
	}
	l := len(a.data) - 1
	w.WriteRune('(')
	for i, v := range a.data {
		if v.Valid {
			if !utf8.ValidString(v.String) {
				return errors.NewNotValidf("[dbr] Argument.WriteTo: StringNull is not UTF-8: %q", v.String)
			}
			dialect.EscapeString(w, v.String)
			if i < l {
				w.WriteRune(',')
			}
		} else {
			w.WriteString("NULL")
		}
	}
	_, err := w.WriteRune(')')
	return err
}

func (a argStringNulls) len() int {
	if a.options() == 0 {
		return len(a.data)
	}
	return 1
}

func (a argStringNulls) INClause() Argument {
	a.opts = argOptionIsIN
	return a
}

func (a argStringNulls) options() uint { return a.opts }

// ArgStringNull adds a nullable string or a slice of nullable strings to the
// argument list. Providing no arguments returns a NULL type. All arguments mut
// be a valid utf-8 string.
func ArgStringNull(args ...null.String) Argument {
	if len(args) == 1 {
		return argStringNull{str: args[0]}
	}
	return argStringNulls{data: args}
}
