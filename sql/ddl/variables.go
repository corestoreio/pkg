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

package ddl

import (
	"strings"

	"github.com/corestoreio/csfw/sql/dml"
	"github.com/corestoreio/errors"
)

// Variables contains multiple MySQL configuration variables. Not threadsafe.
type Variables struct {
	Data map[string]string
	Show *dml.Show
}

// NewVariables creates a new variable collection. If the argument names gets
// passed, the SQL query will load the all variables matching the names.
// Empty argument loads all variables.
func NewVariables(names ...string) *Variables {
	vs := &Variables{
		Data: make(map[string]string),
		Show: dml.NewShow().Variable().Interpolate(),
	}
	vs.Show.IsBuildCache = true
	if len(names) > 1 {
		vs.Show.Where(dml.Column("Variable_name").In().Strs(names...))
	} else if len(names) == 1 {
		vs.Show.Where(dml.Column("Variable_name").Like().Str(names[0]))
	}
	return vs
}

// EqualFold reports whether the value of key and `expected`, interpreted as
// UTF-8 strings, are equal under Unicode case-folding.
func (vs *Variables) EqualFold(key, expected string) bool {
	return strings.EqualFold(vs.Data[key], expected)
}

// Equal compares case sensitive the value of key with the `expected`.
func (vs *Variables) Equal(key, expected string) bool {
	return vs.Data[key] == expected
}

// ToSQL implements dml.QueryBuilder interface to assemble a SQL string and its
// arguments for query execution.
func (vs *Variables) ToSQL() (string, []interface{}, error) {
	return vs.Show.ToSQL()
}

// MapColumns implements dml.ColumnMapper interface and scans a single row from
// the database query result.
func (vs *Variables) MapColumns(rc *dml.ColumnMap) error {
	// TODO: how to handle to load all variables stored in the map?
	name, value := "", ""
	for rc.Next() {
		switch col := rc.Column(); col {
		case "Variable_name":
			rc.String(&name)
		case "Value":
			rc.String(&value)
		default:
			return errors.NewNotFoundf("[ddl] Column %q not found in SHOW VARIABLES", col)
		}
	}
	vs.Data[name] = value
	return errors.WithStack(rc.Err())
}
