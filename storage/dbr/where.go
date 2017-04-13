package dbr

import (
	"strings"

	"github.com/corestoreio/errors"
)

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
	Sub       struct {
		// Select adds a sub-select to the where statement. Condition must be either
		// a column name or anything else which can handle the result of a
		// sub-select.
		Select   *Select
		Operator byte
	}
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

func SubSelect(rawStatementOrColumnName string, operator byte, s *Select) ConditionArg {
	wf := &whereFragment{
		Condition: rawStatementOrColumnName,
	}
	wf.Sub.Select = s
	wf.Sub.Operator = operator
	return wf
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
func writeWhereFragmentsToSQL(fragments WhereFragments, w queryWriter, args *Arguments) error {
	anyConditions := false
	for _, f := range fragments {

		if anyConditions {
			_, _ = w.WriteString(" AND (")
		} else {
			_, _ = w.WriteRune('(')
			anyConditions = true
		}

		addArg := false

		if isValidIdentifier(f.Condition) > 0 { // must be an expression
			_, _ = w.WriteString(f.Condition)
			addArg = true
		} else {
			Quoter.quoteAs(w, f.Condition)

			if f.Sub.Select != nil {
				writeOperator(w, f.Sub.Operator, false)
				w.WriteRune('(')
				subArgs, err := f.Sub.Select.toSQL(w)
				w.WriteRune(')')
				if err != nil {
					return errors.Wrapf(err, "[dbr] writeWhereFragmentsToSQL failed SubSelect for table: %q", f.Sub.Select.FromTable.String())
				}
				*args = append(*args, subArgs...)
			} else {
				// a column only supports one argument. If not provided we panic
				// with an index out of bounds error.
				addArg = writeOperator(w, f.Arguments[0].operator(), true)
			}
		}
		_, _ = w.WriteRune(')')
		if addArg {
			*args = append(*args, f.Arguments...)
		}
	}
	return nil
}

// maxIdentifierLength see http://dev.mysql.com/doc/refman/5.7/en/identifiers.html
const maxIdentifierLength = 64

// IsValidIdentifier checks the permissible syntax for identifiers. Certain
// objects within MySQL, including database, table, index, column, alias, view,
// stored procedure, partition, tablespace, and other object names are known as
// identifiers. ASCII: [0-9,a-z,A-Z$_] (basic Latin letters, digits 0-9, dollar,
// underscore) Max length 63 characters.
//
// Returns 0 if the identifier is valid.
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
