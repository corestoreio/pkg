package dbr

import "fmt"

// QueryBuilder assembles a query and returns the raw SQL without parameter
// substitution and the arguments.
type QueryBuilder interface {
	ToSQL() (string, []interface{}, error)
}

func makeSQL(b QueryBuilder) string {
	sRaw, vals, err := b.ToSQL()
	if err != nil {
		return fmt.Sprintf("[dbr] ToSQL Error: %+v", err)
	}
	sql, err := Preprocess(sRaw, vals)
	if err != nil {
		return fmt.Sprintf("[dbr] Preprocess Error: %+v", err)
	}
	return sql
}

// String returns a string representing a preprocessed, interpolated, query.
// On error, the error gets printed. Fulfills interface fmt.Stringer.
func (b *Delete) String() string {
	return makeSQL(b)
}

// String returns a string representing a preprocessed, interpolated, query.
// On error, the error gets printed. Fulfills interface fmt.Stringer.
func (b *Insert) String() string {
	return makeSQL(b)
}

// String returns a string representing a preprocessed, interpolated, query.
// On error, the error gets printed. Fulfills interface fmt.Stringer.
func (b *Select) String() string {
	return makeSQL(b)
}

// String returns a string representing a preprocessed, interpolated, query.
// On error, the error gets printed. Fulfills interface fmt.Stringer.
func (b *Update) String() string {
	return makeSQL(b)
}
