package dbr

import (
	"bytes"
	"fmt"
)

type (
	joinOn struct {
		whereSqlOrMap interface{}
		args          []interface{}
	}
)

func JoinOn(whereSqlOrMap interface{}, args ...interface{}) joinOn {
	return joinOn{
		whereSqlOrMap: whereSqlOrMap,
		args:          args,
	}
}

func (b *SelectBuilder) join(joinType, table string, columns []string, onConditions ...joinOn) *SelectBuilder {
	var sql bytes.Buffer
	var args []interface{}
	sql.WriteString(" " + joinType + " JOIN " + table + " ON ")
	var w []*whereFragment
	for _, oc := range onConditions {
		w = append(w, newWhereFragment(oc.whereSqlOrMap, oc.args))
	}
	writeWhereFragmentsToSql(w, &sql, &args)
	b.JoinFragments = append(b.JoinFragments, sql.String())

	fmt.Printf("\n%#v\n", b.JoinFragments)

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
