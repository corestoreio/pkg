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

// Always in alphabetical order. We can add more once needed.
const (
	showBinaryLog uint = 1 << iota
	showGlobal
	showMasterStatus
	showSession
	showStatus
	showTableStatus
	showVariables
)

// Show represents the SHOW syntax
type Show struct {
	BuilderBase
	// DB gets required once the Load*() functions will be used.
	DB QueryPreparer

	// Type bitwise flag containing the type of the SHOW statement.
	Type           uint
	LikeCondition  Argument
	WhereFragments Conditions
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
// SESSION. If no modifier is present, the default is SESSION.
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

// MasterStatus provides status information about the binary log files of the
// master. It requires either the SUPER or REPLICATION CLIENT privilege.
func (b *Show) MasterStatus() *Show {
	b.Type = b.Type | showMasterStatus
	return b
}

// TableStatus works likes SHOW TABLES, but provides a lot of information about
// each non-TEMPORARY table. The LIKE clause, if present, indicates which table
// names to match. The WHERE clause can be given to select rows using more
// general conditions. This statement also displays information about views.
func (b *Show) TableStatus() *Show {
	b.Type = b.Type | showTableStatus
	return b
}

// Status provides server status information. This statement does not require
// any privilege. It requires only the ability to connect to the server.
func (b *Show) Status() *Show {
	b.Type = b.Type | showStatus
	return b
}

// BinaryLog lists the binary log files on the server.
func (b *Show) BinaryLog() *Show {
	b.Type = b.Type | showBinaryLog
	return b
}

// Where appends a WHERE clause to the statement for the given string and args
// or map of column/value pairs. Either WHERE or LIKE can be used.
func (b *Show) Where(wf ...*Condition) *Show {
	b.WhereFragments = append(b.WhereFragments, wf...)
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

// IsBuildCache if `true` the final build query including place holders will be
// cached in a private field. Each time a call to function ToSQL happens, the
// arguments will be re-evaluated and returned or interpolated.
func (b *Show) BuildCache() *Show {
	b.IsBuildCache = true
	return b
}

func (b *Show) hasBuildCache() bool {
	return b.IsBuildCache
}

// ToSQL serialized the Show to a SQL string
// It returns the string with placeholders and a slice of query arguments
func (b *Show) toSQL(w queryWriter) error {

	w.WriteString("SHOW ")

	switch {
	case b.Type&showSession != 0:
		w.WriteString("SESSION ")
	case b.Type&showGlobal != 0:
		w.WriteString("GLOBAL ")
	}

	switch {
	case b.Type&showVariables != 0:
		w.WriteString("VARIABLES")
	case b.Type&showStatus != 0:
		w.WriteString("STATUS")
	case b.Type&showMasterStatus != 0:
		w.WriteString("MASTER STATUS")
	case b.Type&showTableStatus != 0:
		w.WriteString("TABLE STATUS")
	case b.Type&showBinaryLog != 0:
		w.WriteString("BINARY LOG")
	}

	if b.LikeCondition != nil {
		Like.write(w, 1)
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
	} else if args, _, err = b.WhereFragments.appendArgs(args, appendArgsWHERE); err != nil {
		return nil, errors.WithStack(err)
	}

	return args, nil
}
