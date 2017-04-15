package dbr

import (
	"github.com/corestoreio/csfw/util/bufferpool"
	"github.com/corestoreio/errors"
)

type alias struct {
	Select     *Select
	Expression string
	Alias      string
}

// MakeAlias creates a new alias expression
func MakeAlias(as ...string) alias {
	a := alias{
		Expression: as[0],
	}
	if len(as) > 1 {
		a.Alias = as[1]
	}
	return a
}

func (t alias) String() string {
	return Quoter.ExprAlias(t.Expression, t.Alias)
}

func (t alias) QuoteAs() string {
	return Quoter.QuoteAs(t.Expression, t.Alias)
}

// QuoteAsWriter writes the quote table and its maybe alias into w.
func (t alias) QuoteAsWriter(w queryWriter) (Arguments, error) {
	if t.Select != nil {
		w.WriteRune('(')
		args, err := t.Select.toSQL(w)
		w.WriteRune(')')
		w.WriteString(" AS ")
		Quoter.quote(w, t.Alias)
		return args, errors.Wrap(err, "[dbr] QuoteAsWriter.SubSelect")
	}
	Quoter.quoteAs(w, t.Expression, t.Alias)
	return nil, nil
}

// DefaultScopeNames specifies the name of the scopes used in all EAV* function
// to generate scope based hierarchical fall backs.
var DefaultScopeNames = [...]string{"Store", "Group", "Website", "Default"}

// EAVIfNull creates a nested IFNULL SQL statement when a scope based fall back
// hierarchy is required. Alias argument will be used as a prefix for the alias
// table name and as the final alias name.
// TODO: Migrate into EAV package
func EAVIfNull(alias, columnName, defaultVal string, scopeNames ...string) string {
	buf := bufferpool.Get()
	defer bufferpool.Put(buf)

	if len(scopeNames) == 0 {
		scopeNames = DefaultScopeNames[:]
	}

	brackets := 0
	for _, n := range scopeNames {
		buf.WriteString("IFNULL(")
		buf.WriteRune('`')
		buf.WriteString(alias)
		buf.WriteString(n)
		buf.WriteRune('`')
		buf.WriteRune('.')
		buf.WriteRune('`')
		buf.WriteString(columnName)
		buf.WriteRune('`')
		if brackets < len(scopeNames)-1 {
			buf.WriteRune(',')
		}
		brackets++
	}

	if defaultVal == "" {
		defaultVal = `''`
	}
	buf.WriteRune(',')
	buf.WriteString(defaultVal)
	for i := 0; i < brackets; i++ {
		buf.WriteRune(')')
	}
	buf.WriteString(" AS ")
	buf.WriteRune('`')
	buf.WriteString(alias)
	buf.WriteRune('`')
	return buf.String()
}
