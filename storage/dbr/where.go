package dbr

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
	Column    string
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

// ConditionColumn adds a column to a WHERE statement
func ConditionColumn(column string, arg Argument) ConditionArg {
	return conditionArgFunc(func() (*whereFragment, error) {
		return &whereFragment{
			Column:    column,
			Arguments: Arguments{arg},
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
func writeWhereFragmentsToSQL(fragments WhereFragments, w queryWriter, args *Arguments) {
	anyConditions := false
	for _, f := range fragments {
		if f.EqualityMap != nil {
			anyConditions = writeEqualityMapToSQL(f.EqualityMap, w, args, anyConditions)
			continue
		}

		if anyConditions {
			_, _ = w.WriteString(" AND (")
		} else {
			_, _ = w.WriteRune('(')
			anyConditions = true
		}

		addArg := false
		if f.Condition != "" {
			_, _ = w.WriteString(f.Condition)
			addArg = true
		} else {
			Quoter.writeQuotedColumn(w, f.Column)
			// a column only supports one argument. If not provided we panic with an index out of bounds error.
			arg := f.Arguments[0]
			switch arg.operator() {
			case OperatorNull:
				w.WriteString(" IS NULL")
			case OperatorNotNull:
				w.WriteString(" IS NOT NULL")
			case OperatorIn:
				w.WriteString(" IN ?")
				addArg = true
			case OperatorNotIn:
				w.WriteString(" NOT IN ?")
				addArg = true
			case OperatorLike:
				w.WriteString(" LIKE ?")
				addArg = true
			case OperatorNotLike:
				w.WriteString(" NOT LIKE ?")
				addArg = true
			case OperatorBetween:
				w.WriteString(" BETWEEN ? AND ?")
				addArg = true
			case OperatorNotBetween:
				w.WriteString(" NOT BETWEEN ? AND ?")
				addArg = true
			case OperatorGreatest:
				w.WriteString(" GREATEST ?")
				addArg = true
			case OperatorLeast:
				w.WriteString(" LEAST ?")
				addArg = true
			case OperatorEqual:
				w.WriteString(" = ?")
				addArg = true
			case OperatorNotEqual:
				w.WriteString(" != ?")
				addArg = true
			default:
				w.WriteString(" = ?")
				addArg = true
			}
		}
		_, _ = w.WriteRune(')')
		if addArg {
			*args = append(*args, f.Arguments...)
		}

	}
}

func writeEqualityMapToSQL(eq map[string]Argument, w queryWriter, args *Arguments, anyConditions bool) bool {
	for k, arg := range eq {
		if arg == nil || arg.operator() == OperatorNull {
			anyConditions = writeWhereCondition(w, k, " IS NULL", anyConditions)
			continue
		}
		if arg.operator() == OperatorNotNull {
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
