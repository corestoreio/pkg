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
	"github.com/corestoreio/errors"
	"github.com/corestoreio/log"
)

const (
	showGlobal uint = 1 << iota
	showSession
	showVariables
	showMasterStatus
)

// Show represents the SHOW syntax
type Show struct {
	Log log.Logger // Log optional logger
	// DB gets required once the Load*() functions will be used.
	DB QueryPreparer

	Type           uint
	LikeCondition  Argument
	WhereFragments WhereFragments
	IsInterpolate  bool // See Interpolate()
	// UseBuildCache if `true` the final build query including place holders
	// will be cached in a private field. Each time a call to function ToSQL
	// happens, the arguments will be re-evaluated and returned or interpolated.
	UseBuildCache bool
	cacheSQL      []byte
	cacheArgs     Arguments // like a buffer, gets reused
}

// NewShow creates a new Truman SHOW.
func NewShow() *Show { return &Show{} }

// WithDB sets the database query object.
func (b *Show) WithDB(db QueryPreparer) *Show {
	b.DB = db
	return b
}

// Global displays with a GLOBAL modifier, the statement displays global system
// variable values. These are the values used to initialize the corresponding
// session variables for new connections to MySQL. If a variable has no global
// value, no value is displayed.
func (b *Show) Global() *Show {
	b.Type = b.Type | showGlobal
	return b
}

// Session displays with a SESSION modifier, the statement displays the system
// variable values that are in effect for the current connection. If a variable
// has no session value, the global value is displayed. LOCAL is a synonym for
// SESSION.
func (b *Show) Session() *Show {
	b.Type = b.Type | showSession
	return b
}

// Variable shows the values of MySQL|MariaDB system variables (“Server System
// Variables”). This statement does not require any privilege. It requires only
// the ability to connect to the server.
func (b *Show) Variable() *Show {
	b.Type = b.Type | showVariables
	return b
}

func (b *Show) MasterStatus() *Show {
	b.Type = b.Type | showMasterStatus
	return b
}

// Where appends a WHERE clause to the statement for the given string and args
// or map of column/value pairs. Either WHERE or LIKE can be used.
func (b *Show) Where(c ...ConditionArg) *Show {
	b.WhereFragments = b.WhereFragments.append(c...)
	return b
}

// Like sets the comparisons LIKE condition. Either WHERE or LIKE can be used.
func (b *Show) Like(arg Argument) *Show {
	b.LikeCondition = arg
	return b
}

// Interpolate if set stringyfies the arguments into the SQL string and returns
// pre-processed SQL command when calling the function ToSQL. Not suitable for
// prepared statements. ToSQLs second argument `Arguments` will then be nil.
func (b *Show) Interpolate() *Show {
	b.IsInterpolate = true
	return b
}

// ToSQL converts the select statement into a string and returns its arguments.
func (b *Show) ToSQL() (string, []interface{}, error) {
	return toSQL(b, b.IsInterpolate, _isNotPrepared)
}

// argumentCapacity returns the total possible guessed size of a new Arguments
// slice. Use as the cap parameter in a call to `make`.
func (b *Show) argumentCapacity() int {
	return len(b.WhereFragments)
}

func (b *Show) writeBuildCache(sql []byte) {
	b.cacheSQL = sql
}

func (b *Show) readBuildCache() (sql []byte, _ Arguments, err error) {
	if b.cacheSQL == nil {
		return nil, nil, nil
	}
	b.cacheArgs, err = b.appendArgs(b.cacheArgs[:0])
	return b.cacheSQL, b.cacheArgs, err
}

func (b *Show) hasBuildCache() bool {
	return b.UseBuildCache
}

// ToSQL serialized the Show to a SQL string
// It returns the string with placeholders and a slice of query arguments
func (b *Show) toSQL(w queryWriter) error {

	w.WriteString("SHOW ")

	if b.Type&showSession != 0 {
		w.WriteString("SESSION ")
	}

	switch {
	case b.Type&showVariables != 0:
		w.WriteString("VARIABLES")
	case b.Type&showMasterStatus != 0:
		w.WriteString("MASTER STATUS")
	}

	if b.LikeCondition != nil {
		b.LikeCondition = b.LikeCondition.applyOperator(Like)
		_ = writeOperator(w, true, b.LikeCondition)
	} else if err := b.WhereFragments.write(w, 'w'); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

// ToSQL serialized the Show to a SQL string
// It returns the string with placeholders and a slice of query arguments
func (b *Show) appendArgs(args Arguments) (_ Arguments, err error) {

	if cap(args) == 0 {
		args = make(Arguments, 0, b.argumentCapacity())
	}

	if b.LikeCondition != nil {
		args = append(args, b.LikeCondition)
	} else if args, _, err = b.WhereFragments.appendArgs(args, 'w'); err != nil {
		return nil, errors.WithStack(err)
	}

	return args, nil
}
