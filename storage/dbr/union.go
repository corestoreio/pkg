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
	"strconv"

	"github.com/corestoreio/csfw/util/bufferpool"
	"github.com/corestoreio/errors"
)

// Union represents a UNION SQL statement. UNION is used to combine the result
// from multiple SELECT statements into a single result set.
type Union struct {
	Selects  []*Select
	OrderBys []string
	IsAll    bool
	// IsPreserveResultSet enables the correct ordering of the result set from
	// the Select statements. Setting this field to true will modify the Select
	// pointers.
	IsPreserveResultSet bool
}

// NewUnion creates a new Union object.
func NewUnion(selects ...*Select) *Union {
	return &Union{
		Selects: selects,
	}
}

// Append adds more Select objects to the Union object.
func (u *Union) Append(selects ...*Select) *Union {
	u.Selects = append(u.Selects, selects...)
	return u
}

// All returns all rows. The default behavior for UNION is that duplicate rows
// are removed from the result. Enabling ALL returns all rows.
func (u *Union) All() *Union {
	u.IsAll = true
	return u
}

// PreserveResultSet enables the correct ordering of the result set from the
// Select statements. UNION by default produces an unordered set of rows. To
// cause rows in a UNION result to consist of the sets of rows retrieved by each
// SELECT one after the other, select an additional column in each SELECT to use
// as a sort column and add an ORDER BY following the last SELECT.
func (u *Union) PreserveResultSet() *Union {
	u.IsPreserveResultSet = true
	return u
}

// OrderBy appends a column or an expression to ORDER the statement by
func (u *Union) OrderBy(ord ...string) *Union {
	u.OrderBys = append(u.OrderBys, ord...)
	return u
}

// OrderDir appends a column to ORDER the statement uy with a given direction
func (u *Union) OrderDir(ord string, isAsc bool) *Union {
	if isAsc {
		u.OrderBys = append(u.OrderBys, ord+" ASC")
	} else {
		u.OrderBys = append(u.OrderBys, ord+" DESC")
	}
	return u
}

// ToSQL renders the UNION into a string and returns its arguments.
func (u *Union) ToSQL() (string, Arguments, error) {
	var w = bufferpool.Get()
	defer bufferpool.Put(w)

	args := make(Arguments, 0, len(u.Selects))
	for i, s := range u.Selects {

		if i > 0 {
			w.WriteString(" UNION ")
			if u.IsAll {
				w.WriteString("ALL ")
			}
		}
		w.WriteRune('(')

		if u.IsPreserveResultSet {
			s.AddColumnsExprAlias(strconv.Itoa(i), "_preserve_result_set")
		}

		sArgs, err := s.toSQL(w)
		if err != nil {
			return "", nil, errors.Wrapf(err, "[dbr] Union.ToSQL at Select index %d", i)
		}
		w.WriteRune(')')
		args = append(args, sArgs...)
	}

	if u.IsPreserveResultSet {
		u.OrderBys = append([]string{"`_preserve_result_set`"}, u.OrderBys...)
	}
	if len(u.OrderBys) > 0 {
		w.WriteString(" ORDER BY ")
		for i, s := range u.OrderBys {
			if i > 0 {
				w.WriteString(", ")
			}
			w.WriteString(s)
		}
	}

	return w.String(), args, nil
}
