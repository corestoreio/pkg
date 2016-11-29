package dbr

import "github.com/corestoreio/csfw/util/errors"

type queryBuilder interface {
	ToSql() (string, []interface{}, error)
}

func makeSql(b queryBuilder) (string, error) {
	// todo add maybe logging in debug mode
	sRaw, vals, err := b.ToSql()
	if err != nil {
		return "", errors.Wrap(err, "[dbr] makeSql.tosql")
	}
	sql, err := Preprocess(sRaw, vals)
	if err != nil {
		return "", errors.Wrap(err, "[dbr] makeSql.string")
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
