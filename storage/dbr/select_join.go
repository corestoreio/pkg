package dbr

type (
	joinOn struct {
		whereSqlOrMap interface{}
		args          []interface{}
	}
	joinFragment struct {
		joinType, table string
		columns         []string
		onConditions    []joinOn // slice is joined via AND
	}
)

func JoinOn(w interface{}, a ...interface{}) joinOn {
	return joinOn{
		whereSqlOrMap: w,
		args:          a,
	}
}

func (b *SelectBuilder) join(j, t string, c []string, on ...joinOn) *SelectBuilder {
	b.JoinFragments = append(b.JoinFragments, &joinFragment{
		joinType:     j,
		table:        t,
		columns:      c,
		onConditions: on,
	})
	return b
}

// Join creates a join construct with the onConditions glued together with AND
func (b *SelectBuilder) Join(table string, columns []string, onConditions ...joinOn) *SelectBuilder {
	return b.join("INNER", table, columns, onConditions...)
}

// LeftJoin creates a join construct with the onConditions glued together with AND
func (b *SelectBuilder) LeftJoin(table string, columns []string, onConditions ...joinOn) *SelectBuilder {
	return b.join("LEFT", table, columns, onConditions...)
}

// LeftJoin creates a join construct with the onConditions glued together with AND
func (b *SelectBuilder) RightJoin(table string, columns []string, onConditions ...joinOn) *SelectBuilder {
	return b.join("RIGHT", table, columns, onConditions...)
}
