package dbr

import (
	"reflect"
)

// todo for maybe later the sync.Pool code
//var wfPool = &sync.Pool{
//	New: func() interface{} {
//		return &whereFragment{}
//	},
//}
//
//// Get returns a buffer from the pool.
//func wfGet() *whereFragment {
//	return wfPool.Get().(*whereFragment)
//}
//
//// Put returns a buffer to the pool.
//// The buffer is reset before it is put back into circulation.
//func wfPut(wfs WhereFragments) {
//		wf.Condition = ""
//		wf.Values = nil
//		wf.EqualityMap = nil
//		wfPool.Put(wf)
//}

// Eq is a map column -> value pairs which must be matched in a query
type Eq map[string]interface{}

type whereFragment struct {
	Condition   string
	Values      []interface{}
	EqualityMap map[string]interface{}
}

// WhereFragments provides a list where clauses
type WhereFragments []*whereFragment

// ConditionArg
type ConditionArg func(*whereFragment)

// ConditionRaw adds a condition and checks values if they implement driver.Valuer.
func ConditionRaw(raw string, values ...interface{}) ConditionArg {
	if err := argsValuer(&values); err != nil {
		panic(err) // todo remove panic
	}
	return func(wf *whereFragment) {
		wf.Condition = raw
		wf.Values = values
	}
}

func ConditionMap(eq Eq) ConditionArg {
	return func(wf *whereFragment) {
		// todo add argsValuer
		wf.EqualityMap = eq
	}
}

func newWhereFragments(wargs ...ConditionArg) WhereFragments {
	ret := make(WhereFragments, len(wargs))
	for i, warg := range wargs {
		ret[i] = new(whereFragment)
		warg(ret[i])
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
