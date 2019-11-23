package urlvalues

import (
	"sort"
	"strings"

	"github.com/corestoreio/pkg/sql/ddl"
	"github.com/corestoreio/pkg/sql/dml"
)

// URLFilter is used with Query.Apply to add WHERE clauses from the URL values:
//   - ?foo=bar - Where(`"foo" = 'bar'`)
//   - ?foo=hello&foo=world - Where(`"foo" IN ('hello','world')`)
//   - ?foo__neq=bar - Where(`"foo" != 'bar'`)
//   - ?foo__exclude=bar - Where(`"foo" != 'bar'`)
//   - ?foo__gt=42 - Where(`"foo" > 42`)
//   - ?foo__gte=42 - Where(`"foo" >= 42`)
//   - ?foo__lt=42 - Where(`"foo" < 42`)
//   - ?foo__lte=42 - Where(`"foo" <= 42`)
//   - ?foo__ieq=bar - Where(`"foo" ILIKE 'bar'`)
//   - ?foo__match=bar - Where(`"foo" SIMILAR TO 'bar'`)
type Filter struct {
	Deterministic bool // if true internal map gets printed in its order, otherwise flaky tests.
	values        Values
	allowed       map[string]struct{}
}

func NewFilter(values Values) *Filter {
	return &Filter{
		values: values,
	}
}

// Values returns URL values.
func (f *Filter) Values() Values {
	return f.values
}

// Allow only columnName__operators are allowed
func (f *Filter) Allow(filters ...string) {
	if f.allowed == nil {
		f.allowed = make(map[string]struct{})
	}
	for _, filter := range filters {
		f.allowed[filter] = struct{}{}
	}
}

func (f *Filter) isAllowed(filter string) bool {
	if len(f.allowed) == 0 {
		return true
	}
	_, ok := f.allowed[filter]
	return ok
}

func (f *Filter) Filters(tbl *ddl.Table, cond dml.Conditions) dml.Conditions {
	if f == nil {
		return cond
	}
	if f.Deterministic {
		keys := make([]string, 0, len(f.values))
		for filter := range f.values {
			keys = append(keys, filter)
		}
		sort.Strings(keys)
		for _, filter := range keys {
			cond = f.iterate(tbl, cond, filter, f.values[filter])
		}
		return cond
	}
	for filter, values := range f.values {
		cond = f.iterate(tbl, cond, filter, values)
	}

	return cond
}

func (f *Filter) iterate(tbl *ddl.Table, cond dml.Conditions, filter string, values []string) dml.Conditions {
	if strings.HasSuffix(filter, "[]") {
		filter = filter[:len(filter)-2]
	}

	if !f.isAllowed(filter) {
		return cond
	}

	var operation string
	if ind := strings.Index(filter, "__"); ind != -1 {
		filter, operation = filter[:ind], filter[ind+2:]
	}

	if tbl.HasColumn(filter) {
		// TODO AND or OR
		cond = addOperator(cond, filter, operation, values)
	}
	return cond
}

func addOperator(b dml.Conditions, field, op string, values []string) dml.Conditions {
	switch op {
	case "", "eq", "include":
		b = forAllValues(b, field, values, dml.Equal, dml.In)
	case "exclude", "neq":
		b = forAllValues(b, field, values, dml.NotEqual, dml.NotIn)
	case "gt":
		b = forEachValue(b, field, values, dml.Greater)
	case "gte":
		b = forEachValue(b, field, values, dml.GreaterOrEqual)
	case "lt":
		b = forEachValue(b, field, values, dml.Less)
	case "lte":
		b = forEachValue(b, field, values, dml.LessOrEqual)
	case "ieq":
		b = forEachValue(b, field, values, dml.Like)
	case "ineq":
		b = forEachValue(b, field, values, dml.NotLike)
	case "bw":
		b = forAllValues(b, field, values, dml.Equal, dml.Between)
	case "nbw":
		b = forAllValues(b, field, values, dml.NotEqual, dml.NotBetween)
		// case "match":
		//	b = forEachValue(b, field, values, " SIMILAR TO ")
	}
	return b
}

func forEachValue(b dml.Conditions, field string, values []string, opValue dml.Op) dml.Conditions {
	for _, value := range values {
		b = append(b, dml.Column(field).Op(opValue).Str(value))
	}
	return b
}

func forAllValues(
	b dml.Conditions, field string, values []string, singleOpValue, multiOpValue dml.Op,
) dml.Conditions {
	if len(values) <= 1 {
		return forEachValue(b, field, values, singleOpValue)
	}

	return append(b, dml.Column(field).Op(multiOpValue).Strs(values...))
}
