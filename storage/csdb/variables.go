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
	"github.com/corestoreio/csfw/storage/dbr"
	"github.com/corestoreio/errors"
)

// Variables contains multiple MySQL configuration variables. Not threadsafe.
type Variables struct {
	Data        map[string]string
	Show        *dbr.Show
	name, value string
}

// NewVariables creates a new variable collection. If the argument names gets
// passed, the SQL query will load the all variables matching the names.
// Empty argument loads all variables.
func NewVariables(names ...string) *Variables {
	vs := &Variables{
		Data: make(map[string]string),
		Show: dbr.NewShow().Variable().Interpolate(),
	}
	vs.Show.UseBuildCache = true
	if len(names) > 1 {
		vs.Show.Where(dbr.Column("Variable_name", dbr.In.Str(names...)))
	} else if len(names) == 1 {
		vs.Show.Where(dbr.Column("Variable_name", dbr.Like.Str(names...)))
	}
	return vs
}

// ToSQL implements dbr.QueryBuilder interface to assemble a SQL string and its
// arguments for query execution.
func (vs *Variables) ToSQL() (string, []interface{}, error) {
	return vs.Show.ToSQL()
}

// RowScan implements dbr.Scanner interface and scans a single row from the
// database query result.
func (vs *Variables) RowScan(idx int64, _ []string, scan func(...interface{}) error) error {
	if err := errors.WithStack(scan(&vs.name, &vs.value)); err != nil {
		return err
	}
	vs.Data[vs.name] = vs.value
	return nil
}

func isValidVarName(name string, allowPercent bool) error {
	if name == "" {
		return nil
	}
	for _, r := range name {
		var ok bool
		switch {
		case '0' <= r && r <= '9':
			ok = true
		case 'a' <= r && r <= 'z', 'A' <= r && r <= 'Z':
			ok = true
		case r == '_': // % can be bypassed with underscore ;-)
			ok = true
		case r == '%' && allowPercent:
			ok = true
		}
		if !ok {
			return errors.NewNotValidf("[csdb] Invalid character %q in variable name %q", string(r), name)
		}
	}
	return nil
}
