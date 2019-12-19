package urlvalues_test

import (
	"net/url"
	"testing"

	"github.com/corestoreio/pkg/sql/ddl"
	"github.com/corestoreio/pkg/sql/dml"
	"github.com/corestoreio/pkg/sql/urlvalues"
	"github.com/corestoreio/pkg/util/assert"
)

// BenchmarkNewFilter_Filters-4   	 3324812	       325 ns/op	     249 B/op	       2 allocs/op
func BenchmarkNewFilter_Filters(b *testing.B) {
	values, err := url.ParseQuery("name__exclude=Mike&name__exclude=Peter")
	assert.NoError(b, err)
	cond := make(dml.Conditions, 0, 100)
	tbl := ddl.NewTable("url_values_models",
		&ddl.Column{Field: "id", Pos: 1},
		&ddl.Column{Field: "name", Pos: 2},
		&ddl.Column{Field: "city", Pos: 3},
	)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		uv := urlvalues.Values(values)
		uf := urlvalues.NewFilter(uv)

		cond, _ = uf.Filters(tbl, cond)
	}
}

func TestNewFilter_AllAllowed(t *testing.T) {
	const query = "SELECT `id`, `name`, `city` FROM `url_values_models`"

	tbl := ddl.NewTable("url_values_models",
		&ddl.Column{Field: "id", Pos: 1},
		&ddl.Column{Field: "name", Pos: 2},
		&ddl.Column{Field: "city", Pos: 3},
	)

	tests := []struct {
		queryString, wantQuery string
	}{
		{"id__gt=1", query + " WHERE (`id` > '1')"},
		{"name__gte=Michael", query + " WHERE (`name` >= 'Michael')"},
		{"id__lt=10", query + " WHERE (`id` < '10')"},
		{"name__lte=Peter", query + " WHERE (`name` <= 'Peter')"},
		{"name__exclude=Peter", query + " WHERE (`name` != 'Peter')"},
		{"name__exclude=Mike&name__exclude=Peter", query + " WHERE (`name` NOT IN ('Mike','Peter'))"},
		{"name=Mike", query + " WHERE (`name` = 'Mike')"},
		{"name__eq=Mike", query + " WHERE (`name` = 'Mike')"},
		{"name__eq=Mi%27ke", query + " WHERE (`name` = 'Mi\\'ke')"},
		{"name__ieq=mik_", query + " WHERE (`name` LIKE 'mik_')"},
		{"name__ineq=mik_", query + " WHERE (`name` NOT LIKE 'mik_')"},
		{"name__ineq=mik%25", query + " WHERE (`name` NOT LIKE 'mik%')"},
		{"name__include=Peter&name__include=Mike", query + " WHERE (`name` IN ('Peter','Mike'))"},
		{"name=Mike&name=Peter", query + " WHERE (`name` IN ('Mike','Peter'))"},
		{"name[]=Mike&name[]=Peter", query + " WHERE (`name` IN ('Mike','Peter'))"},
		{"name%5B%5D=Mike&name%5B%5D=Peter", query + " WHERE (`name` IN ('Mike','Peter'))"},
		{"invalid_field=1", query},
		{"name__gt=1&name__lt=2", query + " WHERE (`name` > '1') AND (`name` < '2')"},
		{"name__gt=1&name__lt=2", query + " WHERE (`name` > '1') AND (`name` < '2')"},
		{"name__bw=33&name__bw=44", query + " WHERE (`name` BETWEEN '33' AND '44')"},
		{"name__bw=55", query + " WHERE (`name` = '55')"},
		{"name__nbw=33&name__nbw=44", query + " WHERE (`name` NOT BETWEEN '33' AND '44')"},
		{"name__nbw=55", query + " WHERE (`name` != '55')"},
		// sorting
		{"id__lt=10&id__sort=asc", query + " WHERE (`id` < '10') ORDER BY `id` ASC"},
		{"id__lt=10&id__sort=desc", query + " WHERE (`id` < '10') ORDER BY `id` DESC"},
		{"id__lt=10&id__sort=", query + " WHERE (`id` < '10') ORDER BY `id`"},
		{"id__lt=10&id__sort=", query + " WHERE (`id` < '10') ORDER BY `id`"},
		{"name__include=Peter&name__include=Mike&id__sort=desc&name__sort=asc", query + " WHERE (`name` IN ('Peter','Mike')) ORDER BY `id` DESC, `name` ASC"},
	}
	for _, test := range tests {
		t.Run(test.queryString, func(t *testing.T) {
			values, err := url.ParseQuery(test.queryString)
			assert.NoError(t, err)

			uv := urlvalues.Values(values)
			uf := urlvalues.NewFilter(uv)
			uf.Deterministic = true

			cond, sortOrder := uf.Filters(tbl, nil)
			dmlSelect := dml.NewSelect(tbl.Columns.FieldNames()...).From(tbl.Name).Where(cond...)

			sqlStr, _, err := dmlSelect.OrderBy(sortOrder...).ToSQL()
			assert.NoError(t, err)
			assert.Exactly(t, test.wantQuery, sqlStr, "%q", sqlStr)
		})
	}
}

func TestNewFilter_SomeAllowed(t *testing.T) {
	const query = "SELECT `id`, `name` FROM `url_values_models`"

	tbl := ddl.NewTable("url_values_models",
		&ddl.Column{Field: "id", Pos: 1},
		&ddl.Column{Field: "name", Pos: 2},
	)

	tests := []struct {
		queryString, wantQuery string
	}{
		{"id__gt=1", query + " WHERE (`id` > '1')"},
		{"name__gte=Michael", query + " WHERE (`name` >= 'Michael')"},
		{"id__lt=10", query},
		{"name__lte=Peter", query},
		{"name__gt=1&name__lt=2", query},
		{"name__gt=1&name__lt=2", query},
		{"name__gt=1&name__lt=2&city_sort=asc", query}, // city not allowed and should now be shown
	}
	for _, test := range tests {
		t.Run(test.queryString, func(t *testing.T) {
			values, err := url.ParseQuery(test.queryString)
			assert.NoError(t, err)

			uv := urlvalues.Values(values)
			uf := urlvalues.NewFilter(uv)
			uf.Deterministic = true
			uf.Allow("id__gt")
			uf.Allow("name__gte")

			cond, sortOrder := uf.Filters(tbl, nil)
			sqlStr, _, err := dml.NewSelect(tbl.Columns.FieldNames()...).From(tbl.Name).Where(cond...).OrderBy(sortOrder...).ToSQL()
			assert.NoError(t, err)
			assert.Exactly(t, test.wantQuery, sqlStr, "%q", sqlStr)
		})
	}
}

func TestNewPager(t *testing.T) {
	const query = "SELECT `id`, `name` FROM `url_values_models`"

	tbl := ddl.NewTable("url_values_models",
		&ddl.Column{Field: "id", Pos: 1},
		&ddl.Column{Field: "name", Pos: 2},
	)

	tests := []struct {
		queryString, wantQuery string
	}{
		{"limit=10", query + " LIMIT 0,10"},
		{"limit=-10", query + " LIMIT 0,100"},
		{"page=-1", query + " LIMIT 9900,100"},
		{"page=1", query + " LIMIT 0,100"},
		{"page=2", query + " LIMIT 100,100"},
		{"page=3", query + " LIMIT 200,100"},
		{"page=10", query + " LIMIT 900,100"},
		{"page=1&limit=20", query + " LIMIT 0,20"},
		{"page=2&limit=20", query + " LIMIT 20,20"},
	}
	for _, test := range tests {
		t.Run(test.queryString, func(t *testing.T) {
			values, err := url.ParseQuery(test.queryString)
			assert.NoError(t, err)

			uv := urlvalues.Values(values)
			p := urlvalues.NewPager(uv)
			p.MaxLimit = 100
			p.MaxOffset = 100
			var cond dml.Conditions

			a := dml.NewSelect(tbl.Columns.FieldNames()...).From(tbl.Name).Where(cond...).WithArgs()
			a, err = p.Pagination(a)
			assert.NoError(t, err)

			sqlStr, _, err := a.ToSQL()
			assert.NoError(t, err)

			assert.Exactly(t, test.wantQuery, sqlStr, "%q", sqlStr)
		})
	}
}
