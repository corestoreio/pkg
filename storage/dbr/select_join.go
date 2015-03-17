package dbr

type (
	joinOn struct {
		whereSqlOrMap interface{}
		args          []interface{}
	}
	joinFragment struct {
		joinType     string
		table        string
		columns      []string
		onConditions []joinOn // slice is joined via AND
	}
)

// JoinTable is a helper func which transforms variadic arguments into a slice
func JoinTable(tableAlias ...string) []string {
	return tableAlias
}

func JoinOn(w interface{}, a ...interface{}) joinOn {
	return joinOn{
		whereSqlOrMap: w,
		args:          a,
	}
}

func (b *SelectBuilder) join(j string, t, c []string, on ...joinOn) *SelectBuilder {
	b.JoinFragments = append(b.JoinFragments, &joinFragment{
		joinType:     j,
		table:        quoteAs(t...),
		columns:      c,
		onConditions: on,
	})
	return b
}

// Join creates a join construct with the onConditions glued together with AND
func (b *SelectBuilder) Join(table, columns []string, onConditions ...joinOn) *SelectBuilder {
	return b.join("INNER", table, columns, onConditions...)
}

// LeftJoin creates a join construct with the onConditions glued together with AND
func (b *SelectBuilder) LeftJoin(table, columns []string, onConditions ...joinOn) *SelectBuilder {
	return b.join("LEFT", table, columns, onConditions...)
}

// LeftJoin creates a join construct with the onConditions glued together with AND
func (b *SelectBuilder) RightJoin(table, columns []string, onConditions ...joinOn) *SelectBuilder {
	return b.join("RIGHT", table, columns, onConditions...)
}
