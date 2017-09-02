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
	"database/sql"

	"github.com/corestoreio/csfw/storage/dbr"
	"github.com/corestoreio/errors"
)

// Variables contains multiple MySQL configuration variables. Not threadsafe.
type Variables struct {
	Convert dbr.RowConvert
	Data    map[string]string
	Show    *dbr.Show
}

// NewVariables creates a new variable collection. If the argument names gets
// passed, the SQL query will load the all variables matching the names.
// Empty argument loads all variables.
func NewVariables(names ...string) *Variables {
	vs := &Variables{
		Data: make(map[string]string),
		Show: dbr.NewShow().Variable().Interpolate(),
	}
	vs.Show.IsBuildCache = true
	if len(names) > 1 {
		vs.Show.Where(dbr.Column("Variable_name").In().Strs(names...))
	} else if len(names) == 1 {
		vs.Show.Where(dbr.Column("Variable_name").Like().Str(names[0]))
	}
	return vs
}

// ToSQL implements dbr.QueryBuilder interface to assemble a SQL string and its
// arguments for query execution.
func (vs *Variables) ToSQL() (string, []interface{}, error) {
	return vs.Show.ToSQL()
}

// RowScan implements dbr.Scanner interface and scans a single row from the
// database query result. It expects that the variable name is in column 0 and
// the variable value in column 1.
func (vs *Variables) RowScan(r *sql.Rows) error {
	if err := vs.Convert.Scan(r); err != nil {
		return err
	}
	name, err := vs.Convert.Index(0).Str()
	if err != nil {
		return errors.Wrapf(err, "[csdb] Variables.RowScan.Index.0 at Row %d\nRaw Values: %q\n", vs.Convert.Count, vs.Convert.String())
	}
	value, err := vs.Convert.Index(1).Str()
	if err != nil {
		return errors.Wrapf(err, "[csdb] Variables.RowScan.Index.1 at Row %d\nRaw Values: %q\n", vs.Convert.Count, vs.Convert.String())
	}
	vs.Data[name] = value
	return nil
}
