package dbr

type (
	joinFragment struct {
		// left, right, inner, middle, upper, lower ...
		joinType string
		// table name/alias which should be joined
		table string
		// contains all column names from the joined table
		columns []string
		// if set to yes then the columns have already been added to select.columns slice
		columnsAdded bool
		// join on condition
		onConditions []*whereFragment // slice is joined via AND
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

func (b *SelectBuilder) join(j string, t, c []string, on ...ConditionArg) *SelectBuilder {
	b.JoinFragments = append(b.JoinFragments, &joinFragment{
		joinType:     j,
		table:        quoteAs(t...),
		columns:      c,
		columnsAdded: false,
		onConditions: newWhereFragments(on...),
	})
	return b
}

// Join creates a join construct with the onConditions glued together with AND
func (b *SelectBuilder) Join(table, columns []string, onConditions ...ConditionArg) *SelectBuilder {
	return b.join("INNER", table, columns, onConditions...)
}

// LeftJoin creates a join construct with the onConditions glued together with AND
func (b *SelectBuilder) LeftJoin(table, columns []string, onConditions ...ConditionArg) *SelectBuilder {
	return b.join("LEFT", table, columns, onConditions...)
}

// LeftJoin creates a join construct with the onConditions glued together with AND
func (b *SelectBuilder) RightJoin(table, columns []string, onConditions ...ConditionArg) *SelectBuilder {
	return b.join("RIGHT", table, columns, onConditions...)
}
