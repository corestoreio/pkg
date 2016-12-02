package dbr

// JoinFragments defines multiple join conditions.
type JoinFragments []*joinFragment

type (
	joinFragment struct {
		// left, right, inner, middle, upper, lower ...
		JoinType string
		// table name/alias which should be joined
		Table alias
		// contains all column names from the joined table
		Columns []string
		// join on condition
		OnConditions []*whereFragment // slice is joined via AND
	}
)

// JoinTable is a helper func which transforms variadic arguments into a slice
func JoinTable(tableAlias ...string) []string {
	return tableAlias
}

// JoinColumns is a helper func which transforms variadic arguments into a slice
func JoinColumns(columns ...string) []string {
	return columns
}

func (b *Select) join(j string, t, c []string, on ...ConditionArg) *Select {
	b.JoinFragments = append(b.JoinFragments, &joinFragment{
		JoinType:     j,
		Table:        MakeAlias(t...),
		Columns:      c,
		OnConditions: newWhereFragments(on...),
	})
	return b
}

// Join creates a join construct with the onConditions glued together with AND
func (b *Select) Join(table, columns []string, onConditions ...ConditionArg) *Select {
	return b.join("INNER", table, columns, onConditions...)
}

// LeftJoin creates a join construct with the onConditions glued together with AND
func (b *Select) LeftJoin(table, columns []string, onConditions ...ConditionArg) *Select {
	return b.join("LEFT", table, columns, onConditions...)
}

// LeftJoin creates a join construct with the onConditions glued together with AND
func (b *Select) RightJoin(table, columns []string, onConditions ...ConditionArg) *Select {
	return b.join("RIGHT", table, columns, onConditions...)
}
