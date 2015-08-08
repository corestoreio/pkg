package dbr

import "fmt"

type queryBuilder interface {
	ToSql() (string, []interface{})
	EventReceiver
}

func makeSql(b queryBuilder) string {
	sRaw, vals := b.ToSql()
	sql, err := Preprocess(sRaw, vals)
	if err != nil {
		b.EventErrKv("dbr.makeSql.string", err, kvs{"sql": sRaw, "args": fmt.Sprint(vals)})
	}
	return sql
}

// String returns a string representing a preprocessed, interpolated, query.
func (b *DeleteBuilder) String() string {
	return makeSql(b)
}

// String returns a string representing a preprocessed, interpolated, query.
func (b *InsertBuilder) String() string {
	return makeSql(b)
}

// String returns a string representing a preprocessed, interpolated, query.
func (b *SelectBuilder) String() string {
	return makeSql(b)
}

// String returns a string representing a preprocessed, interpolated, query.
func (b *UpdateBuilder) String() string {
	return makeSql(b)
}
