package dbr

import "fmt"

type queryBuilder interface {
	ToSql() (string, []interface{}, error)
	EventReceiver
}

func makeSql(b queryBuilder) (string, error) {
	sRaw, vals, err := b.ToSql()
	if err != nil {
		return "", b.EventErrKv("dbr.makeSql.tosql", err, nil)
	}
	sql, err := Preprocess(sRaw, vals)
	if err != nil {
		return "", b.EventErrKv("dbr.makeSql.string", err, kvs{"sql": sRaw, "args": fmt.Sprint(vals)})
	}
	return sql, nil
}

// String returns a string representing a preprocessed, interpolated, query.
func (b *DeleteBuilder) String() (string, error) {
	return makeSql(b)
}

// String returns a string representing a preprocessed, interpolated, query.
func (b *InsertBuilder) String() (string, error) {
	return makeSql(b)
}

// String returns a string representing a preprocessed, interpolated, query.
func (b *SelectBuilder) String() (string, error) {
	return makeSql(b)
}

// String returns a string representing a preprocessed, interpolated, query.
func (b *UpdateBuilder) String() (string, error) {
	return makeSql(b)
}
