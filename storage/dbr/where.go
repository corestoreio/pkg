package dbr

import (
	"reflect"

	"github.com/corestoreio/errors"
)

// Eq is a map Expression -> value pairs which must be matched in a query.
// Joined as AND statements to the WHERE clause. Implements ConditionArg
// interface.
type Eq map[string]interface{}

func (eq Eq) newWhereFragment() (*whereFragment, error) {
	// todo add argsValuer
	//if err := argsValuer(&values); err != nil {
	//	panic(err)
	//}
	return &whereFragment{
		EqualityMap: eq,
	}, nil
}

type whereFragment struct {
	Condition   string
	Values      []interface{}
	EqualityMap map[string]interface{}
}

// WhereFragments provides a list where clauses
type WhereFragments []*whereFragment

// ConditionArg used as argument in Where()
type ConditionArg interface {
	newWhereFragment() (*whereFragment, error)
}

// implements ConditionArg interface ;-)
type conditionArgFunc func() (*whereFragment, error)

func (f conditionArgFunc) newWhereFragment() (*whereFragment, error) {
	return f()
}

// ConditionRaw adds a condition and checks values if they implement driver.Valuer.
func ConditionRaw(raw string, values ...interface{}) ConditionArg {
	return conditionArgFunc(func() (*whereFragment, error) {
		if err := argsValuer(&values); err != nil {
			return nil, errors.Wrapf(err, "[dbr] Raw: %q; Values %v", raw, values)
		}
		return &whereFragment{
			Condition: raw,
			Values:    values,
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
func writeWhereFragmentsToSQL(fragments WhereFragments, sql QueryWriter, args *[]interface{}) {
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
			if len(f.Values) > 0 {
				*args = append(*args, f.Values...)
			}
		} else if f.EqualityMap != nil {
			anyConditions = writeEqualityMapToSQL(f.EqualityMap, sql, args, anyConditions)
		}
	}
}

func writeEqualityMapToSQL(eq map[string]interface{}, sql QueryWriter, args *[]interface{}, anyConditions bool) bool {
	for k, v := range eq {
		if v == nil {
			anyConditions = writeWhereCondition(sql, k, " IS NULL", anyConditions)
			continue
		}

		vVal := reflect.ValueOf(v)

		if vVal.Kind() == reflect.Array || vVal.Kind() == reflect.Slice {
			vValLen := vVal.Len()
			if vValLen == 0 {
				if vVal.IsNil() {
					anyConditions = writeWhereCondition(sql, k, " IS NULL", anyConditions)
				} else {
					if anyConditions {
						_, _ = sql.WriteString(" AND (1=0)")
					} else {
						_, _ = sql.WriteString("(1=0)")
					}
				}
			} else if vValLen == 1 {
				anyConditions = writeWhereCondition(sql, k, " = ?", anyConditions)
				*args = append(*args, vVal.Index(0).Interface())
			} else {
				anyConditions = writeWhereCondition(sql, k, " IN ?", anyConditions)
				*args = append(*args, v)
			}
		} else {
			anyConditions = writeWhereCondition(sql, k, " = ?", anyConditions)
			*args = append(*args, v)
		}

	}

	return anyConditions
}

func writeWhereCondition(sql QueryWriter, k string, pred string, anyConditions bool) bool {
	if anyConditions {
		_, _ = sql.WriteString(" AND (")
	} else {
		_, _ = sql.WriteRune('(')
		anyConditions = true
	}
	Quoter.writeQuotedColumn(k, sql)
	_, _ = sql.WriteString(pred)
	_, _ = sql.WriteRune(')')

	return anyConditions
}
