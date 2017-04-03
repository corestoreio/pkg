package dbr

import (
	"database/sql"
	"strconv"

	"github.com/corestoreio/csfw/util/bufferpool"
	"github.com/corestoreio/errors"
	"github.com/corestoreio/log"
)

// Update contains the clauses for an UPDATE statement
type Update struct {
	Log log.Logger
	DB  struct {
		Preparer
		Execer
	}

	RawFullSQL   string
	RawArguments Arguments

	Table      alias
	SetClauses struct {
		Columns []string
		Arguments
	}
	WhereFragments []*whereFragment
	OrderBys       []string
	LimitCount     uint64
	LimitValid     bool
	OffsetCount    uint64
	OffsetValid    bool

	// Listeners allows to dispatch certain functions in different
	// situations.
	Listeners UpdateListeners
	// PropagationStopped set to true if you would like to interrupt the
	// listener chain. Once set to true all sub sequent calls of the next
	// listeners will be suppressed.
	PropagationStopped bool
	// propagationStoppedAt position in the slice where the stopped propagation
	// has been requested. for every new iteration the propagation must stop at
	// this position.
	propagationStoppedAt int
	// previousError any error occurred during construction the SQL statement
	previousError error
}

// NewUpdate creates a new object with a black hole logger.
func NewUpdate(table ...string) *Update {
	return &Update{
		Log:   log.BlackHole{},
		Table: MakeAlias(table...),
	}
}

// Update creates a new Update for the given table
func (c *Connection) Update(table ...string) *Update {
	u := &Update{
		Log:   c.Log,
		Table: MakeAlias(table...),
	}
	u.DB.Execer = c.DB
	return u
}

// UpdateBySQL creates a new Update for the given SQL string and arguments
func (c *Connection) UpdateBySQL(sql string, args ...Argument) *Update {
	u := &Update{
		Log:          c.Log,
		RawFullSQL:   sql,
		RawArguments: args,
	}
	u.DB.Execer = c.DB
	return u
}

// Update creates a new Update for the given table bound to a transaction
func (tx *Tx) Update(table ...string) *Update {
	u := &Update{
		Log:   tx.Logger,
		Table: MakeAlias(table...),
	}
	u.DB.Execer = tx.Tx
	return u
}

// UpdateBySQL creates a new Update for the given SQL string and arguments bound
// to a transaction
func (tx *Tx) UpdateBySQL(sql string, args ...Argument) *Update {
	u := &Update{
		Log:          tx.Logger,
		RawFullSQL:   sql,
		RawArguments: args,
	}
	u.DB.Execer = tx.Tx
	return u
}

// Set appends a column/value pair for the statement
func (b *Update) Set(column string, value Argument) *Update {
	if b.previousError != nil {
		return b
	}
	b.SetClauses.Columns = append(b.SetClauses.Columns, column)
	b.SetClauses.Arguments = append(b.SetClauses.Arguments, value)
	return b
}

// SetMap appends the elements of the map at column/value pairs for the
// statement.
func (b *Update) SetMap(clauses map[string]Argument) *Update {
	if b.previousError != nil {
		return b
	}
	for col, val := range clauses {
		b.Set(col, val)
	}
	return b
}

// Where appends a WHERE clause to the statement
func (b *Update) Where(args ...ConditionArg) *Update {
	if b.previousError != nil {
		return b
	}
	b.WhereFragments = append(b.WhereFragments, newWhereFragments(args...)...)
	return b
}

// OrderBy appends a column to ORDER the statement by
func (b *Update) OrderBy(ord string) *Update {
	b.OrderBys = append(b.OrderBys, ord)
	return b
}

// OrderDir appends a column to ORDER the statement by with a given direction
func (b *Update) OrderDir(ord string, isAsc bool) *Update {
	if isAsc {
		b.OrderBys = append(b.OrderBys, ord+" ASC")
	} else {
		b.OrderBys = append(b.OrderBys, ord+" DESC")
	}
	return b
}

// Limit sets a limit for the statement; overrides any existing LIMIT
func (b *Update) Limit(limit uint64) *Update {
	b.LimitCount = limit
	b.LimitValid = true
	return b
}

// Offset sets an offset for the statement; overrides any existing OFFSET
func (b *Update) Offset(offset uint64) *Update {
	b.OffsetCount = offset
	b.OffsetValid = true
	return b
}

// ToSQL serialized the Update to a SQL string
// It returns the string with placeholders and a slice of query arguments
func (b *Update) ToSQL() (string, Arguments, error) {
	if b.previousError != nil {
		return "", nil, errors.Wrap(b.previousError, "[dbr] Update.ToSQL")
	}

	if err := b.Listeners.dispatch(OnBeforeToSQL, b); err != nil {
		return "", nil, errors.Wrap(err, "[dbr] Update.Listeners.dispatch")
	}

	if b.RawFullSQL != "" {
		return b.RawFullSQL, b.RawArguments, nil
	}

	if len(b.Table.Expression) == 0 {
		return "", nil, errors.NewEmptyf("[dbr] Update: Table at empty")
	}
	if len(b.SetClauses.Columns) == 0 {
		return "", nil, errors.NewEmptyf("[dbr] Update: SetClauses are empty")
	}

	var buf = bufferpool.Get()
	defer bufferpool.Put(buf)

	var args = make(Arguments, 0, len(b.SetClauses.Columns)+len(b.WhereFragments))

	buf.WriteString("UPDATE ")
	buf.WriteString(b.Table.QuoteAs())
	buf.WriteString(" SET ")

	// Build SET clause SQL with placeholders and add values to args
	for i, c := range b.SetClauses.Columns {
		if i > 0 {
			buf.WriteString(", ")
		}
		Quoter.writeQuotedColumn(buf, c)
		// TODO fix expr ?
		arg := b.SetClauses.Arguments[i]
		if e, ok := arg.(*expr); ok {
			buf.WriteString(" = ")
			buf.WriteString(e.SQL)
			args = append(args, e.Arguments...)
		} else {
			buf.WriteString(" = ?")
			args = append(args, arg)
		}
	}

	// Write WHERE clause if we have any fragments
	if len(b.WhereFragments) > 0 {
		buf.WriteString(" WHERE ")
		writeWhereFragmentsToSQL(b.WhereFragments, buf, &args)
	}

	// Ordering and limiting
	if len(b.OrderBys) > 0 {
		buf.WriteString(" ORDER BY ")
		for i, s := range b.OrderBys {
			if i > 0 {
				buf.WriteString(", ")
			}
			buf.WriteString(s)
		}
	}

	if b.LimitValid {
		buf.WriteString(" LIMIT ")
		buf.WriteString(strconv.FormatUint(b.LimitCount, 10))
	}

	if b.OffsetValid {
		buf.WriteString(" OFFSET ")
		buf.WriteString(strconv.FormatUint(b.OffsetCount, 10))
	}

	return buf.String(), args, nil
}

// Exec executes the statement represented by the Update object. It returns the
// raw database/sql Result and an error if there was one.
func (b *Update) Exec() (sql.Result, error) {
	rawSQL, args, err := b.ToSQL()
	if err != nil {
		return nil, errors.Wrap(err, "[dbr] Update.Exec.ToSQL")
	}

	fullSQL, err := Preprocess(rawSQL, args...)
	if err != nil {
		return nil, errors.Wrap(err, "[dbr] Update.Exec.Preprocess")
	}

	if b.Log != nil && b.Log.IsInfo() {
		defer log.WhenDone(b.Log).Info("dbr.Update.Exec.Timing", log.String("sql", fullSQL))
	}

	result, err := b.DB.Exec(fullSQL)
	if err != nil {
		return result, errors.Wrap(err, "[dbr] Update.Exec.Exec")
	}

	return result, nil
}

// Prepare creates a new prepared statement represented by the Update object. It
// returns the raw database/sql Stmt and an error if there was one.
func (b *Update) Prepare() (*sql.Stmt, error) {
	rawSQL, _, err := b.ToSQL() // TODO create a ToSQL version without any arguments
	if err != nil {
		return nil, errors.Wrap(err, "[dbr] Update.Prepare.ToSQL")
	}

	if b.Log != nil && b.Log.IsInfo() {
		defer log.WhenDone(b.Log).Info("dbr.Update.Prepare.Timing", log.String("sql", rawSQL))
	}

	stmt, err := b.DB.Prepare(rawSQL)
	return stmt, errors.Wrap(err, "[dbr] Update.Prepare.Prepare")
}

// UpdateMulti TODO creates one update statement for multiple records. Uses a
// transaction and a prepared statement.
type UpdateMulti struct {
	Parent Update
	// Records
	Records struct {
		Columns []string
		Objects []RecordGenerater
	}
}

// Columns sets the columns which gets used for each record. Take care to avoid
// duplicate names.
func (b *UpdateMulti) Columns(cols ...string) *UpdateMulti {
	b.Records.Columns = append(b.Records.Columns, cols...)
	return b
}

// Record pulls in values to match Columns from the record. Think about a vector on how to use this.
func (b *UpdateMulti) Record(recs ...RecordGenerater) *UpdateMulti {
	b.Records.Objects = append(b.Records.Objects, recs...)
	return b
}

func (b *UpdateMulti) Exec() ([]sql.Result, error) {
	// TODO imlement
	return nil, nil
}
