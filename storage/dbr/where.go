package dbr

import "strings"

// Eq is a map Expression -> value pairs which must be matched in a query.
// Joined at AND statements to the WHERE clause. Implements ConditionArg
// interface. Eq = EqualityMap.
type Eq map[string]Argument

func (eq Eq) appendConditions(wfs *WhereFragments) {
	for c, arg := range eq {
		if arg == nil {
			arg = ArgNull()
		}
		*wfs = append(*wfs, &whereFragment{
			Condition: c,
			Arguments: Arguments{arg},
		})
	}
}

type whereFragment struct {
	// Condition can contain either a column name in the form of table.column or
	// just column. Or Condition can contain an expression. Whenever a condition
	// is not a valid identifier we treat it as an expression.
	Condition string
	Arguments Arguments
}

func (wf *whereFragment) appendConditions(wfs *WhereFragments) {
	*wfs = append(*wfs, wf)
}

// WhereFragments provides a list where clauses
type WhereFragments []*whereFragment

// ConditionArg used at argument in Where()
type ConditionArg interface {
	appendConditions(*WhereFragments)
}

// Condition adds a condition and checks values if they implement driver.Valuer.
func Condition(rawStatementOrColumnName string, arg ...Argument) ConditionArg {
	return &whereFragment{
		Condition: rawStatementOrColumnName,
		Arguments: arg,
	}
}

func appendConditions(wf *WhereFragments, wargs ...ConditionArg) {
	for _, warg := range wargs {
		warg.appendConditions(wf)
	}
}

// Invariant: only called when len(fragments) > 0
func writeWhereFragmentsToSQL(fragments WhereFragments, w queryWriter, args *Arguments) {
	anyConditions := false
	for _, f := range fragments {

		if anyConditions {
			_, _ = w.WriteString(" AND (")
		} else {
			_, _ = w.WriteRune('(')
			anyConditions = true
		}

		addArg := false
		if isValidIdentifier(f.Condition) > 0 {
			_, _ = w.WriteString(f.Condition)
			addArg = true
		} else {
			Quoter.quoteAs(w, f.Condition)
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

// maxIdentifierLength see http://dev.mysql.com/doc/refman/5.7/en/identifiers.html
const maxIdentifierLength = 64

// IsValidIdentifier checks the permissible syntax for identifiers. Certain
// objects within MySQL, including database, table, index, column, alias, view,
// stored procedure, partition, tablespace, and other object names are known as
// identifiers. ASCII: [0-9,a-z,A-Z$_] (basic Latin letters, digits 0-9, dollar,
// underscore) Max length 63 characters. Returns errors.NotValid
//
// http://dev.mysql.com/doc/refman/5.7/en/identifiers.html
func isValidIdentifier(objectName string) int8 {

	qualifier := "z" // just a dummy value, can be optimized later
	if i := strings.IndexByte(objectName, '.'); i >= 0 {
		qualifier = objectName[:i]
		objectName = objectName[i+1:]
	}

	for _, name := range [2]string{qualifier, objectName} {
		if len(name) > maxIdentifierLength || name == "" {
			return 1 //errors.NewNotValidf("[csdb] Incorrect identifier. Too long or empty: %q", name)
		}

		for i := 0; i < len(name); i++ {
			if !mapAlNum(name[i]) {
				return 2 // errors.NewNotValidf("[csdb] Invalid character in name %q", name)
			}
		}
	}

	return 0
}

func mapAlNum(r byte) bool {
	var ok bool
	switch {
	case '0' <= r && r <= '9':
		ok = true
	case 'a' <= r && r <= 'z', 'A' <= r && r <= 'Z':
		ok = true
	case r == '$', r == '_':
		ok = true
	}
	return ok
}
