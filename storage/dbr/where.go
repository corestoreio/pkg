package dbr

import "reflect"

// Eq is a map column -> value pairs which must be matched in a query
type Eq map[string]interface{}

type whereFragment struct {
	Condition   string
	Values      []interface{}
	EqualityMap map[string]interface{}
}

type ConditionArg func(*whereFragment)

func ConditionRaw(raw string, values ...interface{}) ConditionArg {
	if err := argsValuer(&values); err != nil {
		PkgLog.Info("dbr.insertbuilder.values", "err", err, "args", values)
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

func newWhereFragments(wargs ...ConditionArg) []*whereFragment {
	ret := make([]*whereFragment, len(wargs))
	for i, warg := range wargs {
		ret[i] = new(whereFragment)
		warg(ret[i])
	}
	return ret
}

// Invariant: only called when len(fragments) > 0
func writeWhereFragmentsToSql(fragments []*whereFragment, sql QueryWriter, args *[]interface{}) {
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
			anyConditions = writeEqualityMapToSql(f.EqualityMap, sql, args, anyConditions)
		}
	}
}

func writeEqualityMapToSql(eq map[string]interface{}, sql QueryWriter, args *[]interface{}, anyConditions bool) bool {
	for k, v := range eq {
		if v == nil {
			anyConditions = writeWhereCondition(sql, k, " IS NULL", anyConditions)
		} else {
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
