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
	"context"
	"database/sql"

	"github.com/corestoreio/csfw/util/errors"
)

// Variables contains multiple MySQL configurations.
type Variables []*Variable

// Variable represents one MySQL configuration value retrieved from the database.
type Variable struct {
	Name  string
	Value string
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

// LoadOne loads a single variable identified by name for the current session.
// For now MySQL DSN must have set interpolateParams to true.
func (v *Variable) LoadOne(ctx context.Context, db QueryRower, name string) error {
	if err := isValidVarName(name, false); err != nil {
		return errors.Wrap(err, "[csdb] Variable.ShowVariable")
	}
	row := db.QueryRowContext(ctx, "SHOW SESSION VARIABLES LIKE ?", name)
	if err := row.Scan(&v.Name, &v.Value); err != nil {
		return errors.Wrap(err, "[csdb] ShowVariable")
	}
	return nil
}

// AppendFiltered appends multiple variables to the current slice. If name is
// empty, all variables will be loaded. Name argument can contain the SQL
// wildcard. For now MySQL DSN must have set interpolateParams to true.
func (vs *Variables) AppendFiltered(ctx context.Context, db Querier, name string) (err error) {
	if err := isValidVarName(name, true); err != nil {
		return errors.Wrap(err, "[csdb] Variables.isValidVarName")
	}

	var rows *sql.Rows
	if name != "" {
		rows, err = db.QueryContext(ctx, "SHOW SESSION VARIABLES LIKE ?", name)
	} else {
		rows, err = db.QueryContext(ctx, "SHOW SESSION VARIABLES")
	}
	if err != nil {
		return errors.Wrap(err, "[csdb] csdb.QueryContext")
	}

	defer rows.Close()
	for rows.Next() {
		v := new(Variable)
		if err := rows.Scan(&v.Name, &v.Value); err != nil {
			return errors.Wrap(err, "[csdb] Variables.Scan")
		}
		*vs = append(*vs, v)
	}

	return nil
}

// FindOne finds one entry in the slice. May return an empty type. Comparing
// is case sensitive.
func (vs Variables) FindOne(name string) (v Variable) {
	for _, vv := range vs {
		if name == vv.Name {
			return *vv
		}
	}
	return v
}

// Len returns the length
func (vs Variables) Len() int { return len(vs) }

// Less compares two slice values
func (vs Variables) Less(i, j int) bool { return vs[i].Name < vs[j].Name }

// Swap changes the position
func (vs Variables) Swap(i, j int) { vs[i], vs[j] = vs[j], vs[i] }
