// Copyright 2015-present, Cyrill @ Schumacher.fm and the CoreStore contributors
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

package dml

import (
	"bytes"

	"github.com/corestoreio/errors"
	"github.com/corestoreio/log"
)

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
	// DB can be either a *sql.DB (connection pool), a *sql.Conn (a single
	// dedicated database session) or a *sql.Tx (an in-progress database
	// transaction).
	DB QueryPreparer

	// Type bitwise flag containing the type of the SHOW statement.
	Type uint
	// LikeCondition supports only one argument. Either LIKE or WHERE can be
	// set.
	LikeCondition  Arguments
	WhereFragments Conditions
}

// NewShow creates a new Truman SHOW.
func NewShow() *Show { return &Show{} }

// Show creates a new Show statement with a random connection from the pool.
func (c *ConnPool) Show() *Show {
	id := c.makeUniqueID()
	l := c.Log
	if l != nil {
		l = c.Log.With(log.String("ConnPool", "Show"), log.String("id", id))
	}
	return &Show{
		BuilderBase: BuilderBase{
			builderCommon: builderCommon{
				id:  id,
				Log: l,
			},
		},
		DB: c.DB,
	}
}

// Show creates a new Show statement bound to a single connection.
func (c *Conn) Show() *Show {
	id := c.makeUniqueID()
	l := c.Log
	if l != nil {
		l = c.Log.With(log.String("Conn", "Show"), log.String("id", id))
	}
	return &Show{
		BuilderBase: BuilderBase{
			builderCommon: builderCommon{
				id:  id,
				Log: l,
			},
		},
		DB: c.DB,
	}
}

// Show creates a new Show query bound to a transaction.
func (tx *Tx) Show() *Show {
	id := tx.makeUniqueID()
	l := tx.Log
	if l != nil {
		l = tx.Log.With(log.String("Tx", "Show"), log.String("id", id))
	}
	return &Show{
		BuilderBase: BuilderBase{
			builderCommon: builderCommon{
				id:  id,
				Log: l,
			},
		},
		DB: tx.DB,
	}
}

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

// WithArgs sets the interfaced arguments for the execution with Query+. It
// internally resets previously applied arguments. This function does not
// support interpolation.
func (b *Show) WithArgs(args ...interface{}) *Show {
	b.withArgs(args)
	return b
}

// WithArguments sets the arguments for the execution with Query+. It internally
// resets previously applied arguments. This function supports interpolation.
func (b *Show) WithArguments(args Arguments) *Show {
	b.withArguments(args)
	return b
}

// Like sets the comparisons LIKE condition. Either WHERE or LIKE can be used.
// Only the first argument supported.
func (b *Show) Like(arg Arguments) *Show {
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
	return b.buildArgsAndSQL(b)
}

func (b *Show) writeBuildCache(sql []byte) {
	b.cacheSQL = sql
}

func (b *Show) readBuildCache() (sql []byte) {
	return b.cacheSQL
}

// DisableBuildCache if enabled it does not cache the SQL string as a final
// rendered byte slice. Allows you to rebuild the query with different
// statements.
func (b *Show) DisableBuildCache() *Show {
	b.IsBuildCacheDisabled = true
	return b
}

// ToSQL serialized the Show to a SQL string
// It returns the string with placeholders and a slice of query arguments
func (b *Show) toSQL(w *bytes.Buffer, placeHolders []string) (_ []string, err error) {

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

	if len(b.LikeCondition) == 1 {
		Like.write(w, b.LikeCondition...)
	} else {
		placeHolders, err = b.WhereFragments.write(w, 'w', placeHolders)
		if err != nil {
			return nil, errors.WithStack(err)
		}
	}

	return placeHolders, nil
}
