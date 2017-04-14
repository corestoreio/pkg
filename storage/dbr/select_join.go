package dbr

// JoinFragments defines multiple join conditions.
type JoinFragments []*joinFragment

type joinFragment struct {
	// JoinType can be LEFT, RIGHT, INNER, OUTER or CROSS
	JoinType string
	// Table name and alias of the table
	Table alias
	// OnConditions join on those conditions
	OnConditions WhereFragments
}

func (b *Select) join(j string, t alias, on ...ConditionArg) *Select {
	jf := &joinFragment{
		JoinType: j,
		Table:    t,
	}
	appendConditions(&jf.OnConditions, on...)
	b.JoinFragments = append(b.JoinFragments, jf)
	return b
}

// Join creates an INNER join construct. By default, the onConditions are glued
// together with AND.
func (b *Select) Join(table alias, onConditions ...ConditionArg) *Select {
	return b.join("INNER", table, onConditions...)
}

// LeftJoin creates a LEFT join construct. By default, the onConditions are
// glued together with AND.
func (b *Select) LeftJoin(table alias, onConditions ...ConditionArg) *Select {
	return b.join("LEFT", table, onConditions...)
}

// RightJoin creates a RIGHT join construct. By default, the onConditions are
// glued together with AND.
func (b *Select) RightJoin(table alias, onConditions ...ConditionArg) *Select {
	return b.join("RIGHT", table, onConditions...)
}

// OuterJoin creates an OUTER join construct. By default, the onConditions are
// glued together with AND.
func (b *Select) OuterJoin(table alias, onConditions ...ConditionArg) *Select {
	return b.join("OUTER", table, onConditions...)
}

// CrossJoin creates a CROSS join construct. By default, the onConditions are
// glued together with AND.
func (b *Select) CrossJoin(table alias, onConditions ...ConditionArg) *Select {
	return b.join("CROSS", table, onConditions...)
}
