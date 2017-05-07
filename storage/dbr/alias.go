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

import "github.com/corestoreio/errors"

type alias struct {
	// Select used in cases where a sub-select is required.
	Select *Select
	// IsExpression if true the field `Name` will be treated as an expression and
	// won't get quoted when generating the SQL.
	IsExpression bool
	// Name can be any kind of SQL expression or a valid identifier. It gets
	// quoted when `IsExpression` is false.
	Name string
	// Alias must be a valid identifier allowed for alias usage.
	Alias string
}

// MakeAlias creates a new name with an optional alias. Supports two arguments.
// 1. a qualifier name and 2. an alias.
func MakeAlias(nameAlias ...string) alias {
	a := alias{
		Name: nameAlias[0],
	}
	if len(nameAlias) > 1 {
		a.Alias = nameAlias[1]
	}
	return a
}

// MakeAliasExpr creates a new expression with an optional alias. Supports two
// arguments. 1. an expression and 2. an alias.
func MakeAliasExpr(expressionAlias ...string) alias {
	a := alias{
		IsExpression: true,
		Name:         expressionAlias[0],
	}
	if len(expressionAlias) > 1 {
		a.Alias = expressionAlias[1]
	}
	return a
}

// String returns the correct stringyfied statement.
func (a alias) String() string {
	if a.IsExpression {
		return Quoter.exprAlias(a.Name, a.Alias)
	}
	return a.QuoteAs()
}

// QuoteAs always quuotes the name and the alias
func (a alias) QuoteAs() string {
	return Quoter.QuoteAs(a.Name, a.Alias)
}

// FquoteAs writes the quoted table and its maybe alias into w.
func (a alias) FquoteAs(w queryWriter) (Arguments, error) {
	if a.Select != nil {
		w.WriteByte('(')
		args, err := a.Select.toSQL(w)
		w.WriteByte(')')
		w.WriteString(" AS ")
		Quoter.quote(w, a.Alias)
		return args, errors.Wrap(err, "[dbr] FquoteAs.SubSelect")
	}

	if a.IsExpression {
		Quoter.FquoteExprAs(w, a.Name, a.Alias)
	} else {
		Quoter.FquoteAs(w, a.Name, a.Alias)
	}

	return nil, nil
}

// TODO(CyS) if we need to distinguish between table name and the column or even need
// a sub select in the column list, then we can implement type aliases and replace
// all []string with type aliases. This costs some allocs but for modifying queries
// in dispatched events, it's getting easier ...
type aliases []alias

func (as aliases) fQuoteAs(w queryWriter, args Arguments) (Arguments, error) {
	for i, a := range as {
		if i > 0 {
			w.WriteString(", ")
		}
		args2, err := a.FquoteAs(w)
		if err != nil {
			return nil, errors.Wrapf(err, "[dbr] aliases.fQuoteAs")
		}
		if args2 != nil {
			args = append(args, args2...)
		}
	}
	return args, nil
}

func appendColumns(as aliases, columns []string) aliases {
	if len(as) == 0 {
		as = make(aliases, 0, len(columns))
	}
	for _, c := range columns {
		as = append(as, alias{Name: c})
	}
	return as
}

// columns must be balanced slice. i=column name, i+1=alias name
func appendColumnsAliases(as aliases, columns []string, isExpression bool) aliases {
	if len(as) == 0 {
		as = make(aliases, 0, len(columns)/2)
	}
	for i := 0; i < len(columns); i = i + 2 {
		as = append(as, alias{Name: columns[i], Alias: columns[i+1], IsExpression: isExpression})
	}
	return as
}
