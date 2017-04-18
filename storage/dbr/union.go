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
	"strings"

	"github.com/corestoreio/csfw/util/bufferpool"
	"github.com/corestoreio/errors"
)

// Union represents a UNION SQL statement. UNION is used to combine the result
// from multiple SELECT statements into a single result set.
type Union struct {
	Selects  []*Select
	OrderBys []string
	IsAll    bool
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
	for i, s := range u.Selects {
		s.AddColumnsExprAlias(strconv.Itoa(i), "_preserve_result_set")
	}
	u.OrderBys = append([]string{"`_preserve_result_set`"}, u.OrderBys...)
	return u
}

// OrderBy appends a column or an expression to ORDER the statement ascending.
// MySQL will order the result set in a temporary table, which is slow.
// https://dev.mysql.com/doc/relnotes/mysql/5.7/en/news-5-7-3.html
func (u *Union) OrderBy(ord ...string) *Union {
	u.OrderBys = append(u.OrderBys, ord...)
	return u
}

// OrderByDesc appends a column or an expression to ORDER the statement
// descending. MySQL will order the result set in a temporary table, which is
// slow.
func (u *Union) OrderByDesc(ord ...string) *Union {
	u.OrderBys = orderByDesc(u.OrderBys, ord)
	return u
}

// ToSQL renders the UNION into a string and returns its arguments. This
// function is idempotent.
func (u *Union) ToSQL() (string, Arguments, error) {
	var w = bufferpool.Get()
	defer bufferpool.Put(w)

	args := make(Arguments, 0, len(u.Selects))
	for i, s := range u.Selects {

		if i > 0 {
			sqlWriteUnionAll(w, u.IsAll)
		}
		w.WriteRune('(')
		sArgs, err := s.toSQL(w)
		if err != nil {
			return "", nil, errors.Wrapf(err, "[dbr] Union.ToSQL at Select index %d", i)
		}
		w.WriteRune(')')
		args = append(args, sArgs...)
	}
	sqlWriteOrderBy(w, u.OrderBys, true)
	return w.String(), args, nil
}

// UnionTemplate builds multiple select statements joined by UNION and all based
// on a common template.
type UnionTemplate struct {
	Select        *Select
	oldNew        [][]string
	repls         []*strings.Replacer
	stmtCount     int
	OrderBys      []string
	IsAll         bool
	previousError error
}

// NewUnionTemplate creates a new UNION generator from a provided SELECT
// template.
func NewUnionTemplate(selectTemplate *Select) *UnionTemplate {
	return &UnionTemplate{
		Select: selectTemplate,
	}
}

// All returns all rows. The default behavior for UNION is that duplicate rows
// are removed from the result. Enabling ALL returns all rows.
func (ut *UnionTemplate) All() *UnionTemplate {
	ut.IsAll = true
	return ut
}

// PreserveResultSet enables the correct ordering of the result set from the
// Select statements. UNION by default produces an unordered set of rows. To
// cause rows in a UNION result to consist of the sets of rows retrieved by each
// SELECT one after the other, select an additional column in each SELECT to use
// as a sort column and add an ORDER BY following the last SELECT.
func (ut *UnionTemplate) PreserveResultSet() *UnionTemplate {
	// this API is different than compared to the Union.PreserveResultSet()
	// because here we can guarantee idempotent calls to ToSQL.
	ut.Select.AddColumnsExprAlias("{preserveResultSet}", "_preserve_result_set")
	ut.OrderBys = append([]string{"`_preserve_result_set`"}, ut.OrderBys...)
	for i := 0; i < ut.stmtCount; i++ {
		ut.oldNew[i] = append(ut.oldNew[i], "{preserveResultSet}", strconv.Itoa(i))
	}
	return ut
}

// OrderBy appends a column or an expression to ORDER the statement ascending.
// MySQL will order the result set in a temporary table, which is slow.
// https://dev.mysql.com/doc/relnotes/mysql/5.7/en/news-5-7-3.html
func (ut *UnionTemplate) OrderBy(ord ...string) *UnionTemplate {
	ut.OrderBys = append(ut.OrderBys, ord...)
	return ut
}

// OrderByDesc appends a column or an expression to ORDER the statement
// descending. MySQL will order the result set in a temporary table, which is
// slow.
func (ut *UnionTemplate) OrderByDesc(ord ...string) *UnionTemplate {
	ut.OrderBys = orderByDesc(ut.OrderBys, ord)
	return ut
}

// StringReplace replaces the `key` with one of the `values`. Each value defines
// a generated SELECT query. Repeating calls of StringReplace must provide the
// same amount of `values` as the first call. This function is just a simple
// string replacement. Make sure that your key does not match other parts of the
// SQL query.
func (ut *UnionTemplate) StringReplace(key string, values ...string) *UnionTemplate {
	if ut.stmtCount == 0 {
		ut.stmtCount = len(values)
		ut.oldNew = make([][]string, ut.stmtCount)
		ut.repls = make([]*strings.Replacer, ut.stmtCount)
	}
	if len(values) != ut.stmtCount {
		ut.previousError = errors.NewNotValidf("[dbr] UnionTemplate.StringReplace: Argument count for values too short. Have %d Want %d", len(values), ut.stmtCount)
		return ut
	}
	for i := 0; i < ut.stmtCount; i++ {
		ut.oldNew[i] = append(ut.oldNew[i], key, values[i])
	}
	return ut
}

// MultiplyArguments repeats the `args` variable n-times to match the number of
// generated SELECT queries in the final UNION statement. It should be called
// after all calls to `StringReplace` have been made.
func (ut *UnionTemplate) MultiplyArguments(args ...Argument) Arguments {
	ret := make(Arguments, len(args)*ut.stmtCount)
	lArgs := len(args)
	for i := 0; i < ut.stmtCount; i++ {
		copy(ret[i*lArgs:], args)
	}
	return ret
}

// Preprocess wrapper function which calls ToSQL and then Preprocess.
func (ut *UnionTemplate) Preprocess() (string, error) {
	sStr, args, err := ut.ToSQL()
	if err != nil {
		return "", errors.Wrap(err, "[dbr] UnionTemplate.Preprocess.ToSQL")
	}
	sStr, err = Preprocess(sStr, args...)
	return sStr, errors.Wrap(err, "[dbr] UnionTemplate.Preprocess.Preprocess")
}

// ToSQL generates the SQL string and its arguments. Calls to this function are
// idempotent.
func (ut *UnionTemplate) ToSQL() (string, Arguments, error) {
	if ut.previousError != nil {
		return "", nil, ut.previousError
	}

	w := bufferpool.Get()
	tplArgs, err := ut.Select.toSQL(w)
	selStr := w.String()
	bufferpool.Put(w)
	if err != nil {
		return "", nil, errors.Wrap(err, "[dbr] UnionTpl.ToSQL: toSQL template")
	}

	wu := bufferpool.Get()
	defer bufferpool.Put(wu)

	for i := 0; i < ut.stmtCount; i++ {
		repl := ut.repls[i]
		if repl == nil {
			repl = strings.NewReplacer(ut.oldNew[i]...)
			ut.repls[i] = repl
		}
		if i > 0 {
			sqlWriteUnionAll(wu, ut.IsAll)
		}
		wu.WriteRune('(')
		repl.WriteString(wu, selStr)
		wu.WriteRune(')')
	}
	sqlWriteOrderBy(wu, ut.OrderBys, true)
	return wu.String(), ut.MultiplyArguments(tplArgs...), nil
}
