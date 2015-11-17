package dbr

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
	"time"

	"github.com/corestoreio/csfw/utils/bufferpool"
)

type expr struct {
	Sql    string
	Values []interface{}
}

// Expr is a SQL fragment with placeholders, and a slice of args to replace them with
func Expr(sql string, values ...interface{}) *expr {
	return &expr{Sql: sql, Values: values}
}

// UpdateBuilder contains the clauses for an UPDATE statement
type UpdateBuilder struct {
	*Session
	runner

	RawFullSql   string
	RawArguments []interface{}

	Table          alias
	SetClauses     []*setClause
	WhereFragments []*whereFragment
	OrderBys       []string
	LimitCount     uint64
	LimitValid     bool
	OffsetCount    uint64
	OffsetValid    bool
}

var _ queryBuilder = (*UpdateBuilder)(nil)

type setClause struct {
	column string
	value  interface{}
}

// Update creates a new UpdateBuilder for the given table
func (sess *Session) Update(table ...string) *UpdateBuilder {
	return &UpdateBuilder{
		Session: sess,
		runner:  sess.cxn.DB,
		Table:   newAlias(table...),
	}
}

// UpdateBySql creates a new UpdateBuilder for the given SQL string and arguments
func (sess *Session) UpdateBySql(sql string, args ...interface{}) *UpdateBuilder {
	if err := argsValuer(&args); err != nil {
		sess.EventErrKv("dbr.insertbuilder.values", err, kvs{"args": fmt.Sprint(args)})
	}
	return &UpdateBuilder{
		Session:      sess,
		runner:       sess.cxn.DB,
		RawFullSql:   sql,
		RawArguments: args,
	}
}

// Update creates a new UpdateBuilder for the given table bound to a transaction
func (tx *Tx) Update(table ...string) *UpdateBuilder {
	return &UpdateBuilder{
		Session: tx.Session,
		runner:  tx.Tx,
		Table:   newAlias(table...),
	}
}

// UpdateBySql creates a new UpdateBuilder for the given SQL string and arguments bound to a transaction
func (tx *Tx) UpdateBySql(sql string, args ...interface{}) *UpdateBuilder {
	if err := argsValuer(&args); err != nil {
		tx.EventErrKv("dbr.insertbuilder.values", err, kvs{"args": fmt.Sprint(args)})
	}
	return &UpdateBuilder{
		Session:      tx.Session,
		runner:       tx.Tx,
		RawFullSql:   sql,
		RawArguments: args,
	}
}

// Set appends a column/value pair for the statement
func (b *UpdateBuilder) Set(column string, value interface{}) *UpdateBuilder {
	if dbVal, ok := value.(driver.Valuer); ok {
		if val, err := dbVal.Value(); err == nil {
			value = val
		} else {
			panic(err)
		}
	}
	b.SetClauses = append(b.SetClauses, &setClause{column: column, value: value})
	return b
}

// SetMap appends the elements of the map as column/value pairs for the statement
func (b *UpdateBuilder) SetMap(clauses map[string]interface{}) *UpdateBuilder {
	for col, val := range clauses {
		b = b.Set(col, val)
	}
	return b
}

// Where appends a WHERE clause to the statement
func (b *UpdateBuilder) Where(args ...ConditionArg) *UpdateBuilder {
	b.WhereFragments = append(b.WhereFragments, newWhereFragments(args...)...)
	return b
}

// OrderBy appends a column to ORDER the statement by
func (b *UpdateBuilder) OrderBy(ord string) *UpdateBuilder {
	b.OrderBys = append(b.OrderBys, ord)
	return b
}

// OrderDir appends a column to ORDER the statement by with a given direction
func (b *UpdateBuilder) OrderDir(ord string, isAsc bool) *UpdateBuilder {
	if isAsc {
		b.OrderBys = append(b.OrderBys, ord+" ASC")
	} else {
		b.OrderBys = append(b.OrderBys, ord+" DESC")
	}
	return b
}

// Limit sets a limit for the statement; overrides any existing LIMIT
func (b *UpdateBuilder) Limit(limit uint64) *UpdateBuilder {
	b.LimitCount = limit
	b.LimitValid = true
	return b
}

// Offset sets an offset for the statement; overrides any existing OFFSET
func (b *UpdateBuilder) Offset(offset uint64) *UpdateBuilder {
	b.OffsetCount = offset
	b.OffsetValid = true
	return b
}

// ToSql serialized the UpdateBuilder to a SQL string
// It returns the string with placeholders and a slice of query arguments
func (b *UpdateBuilder) ToSql() (string, []interface{}, error) {
	if b.RawFullSql != "" {
		return b.RawFullSql, b.RawArguments, nil
	}

	if len(b.Table.Expression) == 0 {
		return "", nil, ErrMissingTable
	}
	if len(b.SetClauses) == 0 {
		return "", nil, ErrMissingSet
	}

	var sql = bufferpool.Get()
	defer bufferpool.Put(sql)

	var args []interface{}

	sql.WriteString("UPDATE ")
	sql.WriteString(b.Table.QuoteAs())
	sql.WriteString(" SET ")

	// Build SET clause SQL with placeholders and add values to args
	for i, c := range b.SetClauses {
		if i > 0 {
			sql.WriteString(", ")
		}
		Quoter.writeQuotedColumn(c.column, sql)
		if e, ok := c.value.(*expr); ok {
			sql.WriteString(" = ")
			sql.WriteString(e.Sql)
			args = append(args, e.Values...)
		} else {
			sql.WriteString(" = ?")
			args = append(args, c.value)
		}
	}

	// Write WHERE clause if we have any fragments
	if len(b.WhereFragments) > 0 {
		sql.WriteString(" WHERE ")
		writeWhereFragmentsToSql(b.WhereFragments, sql, &args)
	}

	// Ordering and limiting
	if len(b.OrderBys) > 0 {
		sql.WriteString(" ORDER BY ")
		for i, s := range b.OrderBys {
			if i > 0 {
				sql.WriteString(", ")
			}
			sql.WriteString(s)
		}
	}

	if b.LimitValid {
		sql.WriteString(" LIMIT ")
		fmt.Fprint(sql, b.LimitCount)
	}

	if b.OffsetValid {
		sql.WriteString(" OFFSET ")
		fmt.Fprint(sql, b.OffsetCount)
	}

	return sql.String(), args, nil
}

// Exec executes the statement represented by the UpdateBuilder
// It returns the raw database/sql Result and an error if there was one
func (b *UpdateBuilder) Exec() (sql.Result, error) {
	sql, args, err := b.ToSql()
	if err != nil {
		return nil, b.EventErrKv("dbr.update.exec.tosql", err, nil)
	}

	fullSql, err := Preprocess(sql, args)
	if err != nil {
		return nil, b.EventErrKv("dbr.update.exec.interpolate", err, kvs{"sql": fullSql})
	}

	// Start the timer:
	startTime := time.Now()
	defer func() { b.TimingKv("dbr.update", time.Since(startTime).Nanoseconds(), kvs{"sql": fullSql}) }()

	result, err := b.runner.Exec(fullSql)
	if err != nil {
		return result, b.EventErrKv("dbr.update.exec.exec", err, kvs{"sql": fullSql})
	}

	return result, nil
}
