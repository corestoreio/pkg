package dbr

import "github.com/corestoreio/csfw/util/bufferpool"

// Eq is a map Expression -> value pairs which must be matched in a query.
// Joined at AND statements to the WHERE clause. Implements ConditionArg
// interface. Eq = EqualityMap.
type Eq map[string]Argument

func (eq Eq) newWhereFragment() (*whereFragment, error) {
	return &whereFragment{
		EqualityMap: eq,
	}, nil
}

type whereFragment struct {
	Condition string
	Arguments
	EqualityMap Eq
}

// WhereFragments provides a list where clauses
type WhereFragments []*whereFragment

// ConditionArg used at argument in Where()
type ConditionArg interface {
	newWhereFragment() (*whereFragment, error)
}

// implements ConditionArg interface ;-)
type conditionArgFunc func() (*whereFragment, error)

func (f conditionArgFunc) newWhereFragment() (*whereFragment, error) {
	return f()
}

// ConditionColumn TODO
func ConditionColumn(column string, arg Argument) ConditionArg {
	return conditionArgFunc(func() (*whereFragment, error) {
		buf := bufferpool.Get()
		defer bufferpool.Put(buf)

		Quoter.writeQuotedColumn(buf, column)

		var args Arguments
		switch arg.options() {
		case argOptionNull:
			buf.WriteString(" IS NULL")
		case argOptionNotNull:
			buf.WriteString(" IS NOT NULL")
		case argOptionIsIN:
			buf.WriteString(" IN ?")
			args = Arguments{arg}
		default:
			buf.WriteString(" = ?")
			args = Arguments{arg}
		}

		return &whereFragment{
			Condition: buf.String(),
			Arguments: args,
		}, nil
	})
}

// ConditionRaw adds a condition and checks values if they implement driver.Valuer.
func ConditionRaw(raw string, arg ...Argument) ConditionArg {
	return conditionArgFunc(func() (*whereFragment, error) {
		return &whereFragment{
			Condition: raw,
			Arguments: arg,
		}, nil
	})
}

func newWhereFragments(wargs ...ConditionArg) WhereFragments {
	ret := make(WhereFragments, len(wargs))
	for i, warg := range wargs {
		wf, err := warg.newWhereFragment()
		if err != nil {
			panic(err) // damn it ... TODO remove panic
		}
		ret[i] = wf
	}
	return ret
}

// Invariant: only called when len(fragments) > 0
func writeWhereFragmentsToSQL(fragments WhereFragments, sql queryWriter, args *Arguments) {
	anyConditions := false
	for _, f := range fragments {
		if f.Condition != "" {
			if anyConditions {
				_, _ = sql.WriteString(" AND (")
			} else {
				_, _ = sql.WriteRune('(')
				anyConditions = true
			}
			_, _ = sql.WriteString(f.Condition)
			_, _ = sql.WriteRune(')')
			if f.Arguments != nil {
				*args = append(*args, f.Arguments...)
			}
		} else if f.EqualityMap != nil {
			anyConditions = writeEqualityMapToSQL(f.EqualityMap, sql, args, anyConditions)
		}
	}
}

func writeEqualityMapToSQL(eq map[string]Argument, w queryWriter, args *Arguments, anyConditions bool) bool {
	for k, arg := range eq {
		if arg == nil || arg.options() == argOptionNull {
			anyConditions = writeWhereCondition(w, k, " IS NULL", anyConditions)
			continue
		}
		if arg.options() == argOptionNotNull {
			anyConditions = writeWhereCondition(w, k, " IS NOT NULL", anyConditions)
			continue
		}

		if arg.len() > 1 {
			anyConditions = writeWhereCondition(w, k, " IN ?", anyConditions)
			*args = append(*args, arg)
		} else {
			anyConditions = writeWhereCondition(w, k, " = ?", anyConditions)
			*args = append(*args, arg)
		}
	}

	return anyConditions
}

func writeWhereCondition(w queryWriter, column string, pred string, anyConditions bool) bool {
	if anyConditions {
		_, _ = w.WriteString(" AND (")
	} else {
		_, _ = w.WriteRune('(')
		anyConditions = true
	}
	Quoter.writeQuotedColumn(w, column)
	_, _ = w.WriteString(pred)
	_, _ = w.WriteRune(')')

	return anyConditions
}
