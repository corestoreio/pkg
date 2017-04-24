package dbr

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestQuoteAs(t *testing.T) {
	tests := []struct {
		have []string
		want string
	}{
		0: {[]string{"a"}, "`a`"},
		1: {[]string{"a", "b"}, "`a` AS `b`"},
		2: {[]string{"a", ""}, "`a`"},
		3: {[]string{"`c`"}, "`c`"},
		4: {[]string{"d.e"}, "`d`.`e`"},
		5: {[]string{"`d`.`e`"}, "`d`.`e`"},
		6: {[]string{"f", "g", "h"}, "`f` AS `g_h`"},
		7: {[]string{"f", "g", "h`h"}, "`f` AS `g_hh`"},
	}
	for i, test := range tests {
		assert.Exactly(t, test.want, Quoter.QuoteAs(test.have...), "Index %d", i)
	}
}

// BenchmarkQuoteAs-4	 3000000	       417 ns/op	      48 B/op	       2 allocs/op
// BenchmarkQuoteAs-4   10000000	       231 ns/op	      48 B/op	       2 allocs/op
// BenchmarkQuoteAs-4    5000000	       287 ns/op	      32 B/op	       1 allocs/op
func BenchmarkQuoteAs(b *testing.B) {
	const want = "`e`.`entity_id` AS `ee`"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if have := Quoter.QuoteAs("e.entity_id", "ee"); have != want {
			b.Fatalf("Have %s\nWant %s\n", have, want)
		}
	}
}

// BenchmarkQuoteAlias-4   	20000000	        96.3 ns/op	      48 B/op	       1 allocs/op
func BenchmarkQuoteAlias(b *testing.B) {
	const want = "(e.price * a.tax * e.weee) AS `final_price`"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if have := Quoter.ExprAlias("(e.price * a.tax * e.weee)", "final_price"); have != want {
			b.Fatalf("Have %s\nWant %s\n", have, want)
		}
	}
}

// BenchmarkQuoteQuote/Worse_Case-4         	 5000000	       402 ns/op	      96 B/op	       5 allocs/op
// BenchmarkQuoteQuote/Best_Case-4          	20000000	       108 ns/op	      32 B/op	       1 allocs/op
func BenchmarkQuoteQuote(b *testing.B) {
	const want = "`databaseName`.`tableName`"

	b.ReportAllocs()
	b.ResetTimer()
	b.Run("Worse Case", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			if have := Quoter.Quote("database`Name", "table`Name"); have != want {
				b.Fatalf("Have %s\nWant %s\n", have, want)
			}
		}
	})
	b.Run("Best Case", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			if have := Quoter.Quote("databaseName", "tableName"); have != want {
				b.Fatalf("Have %s\nWant %s\n", have, want)
			}
		}
	})
}

func TestMysqlQuoter_Quote(t *testing.T) {
	assert.Exactly(t, "`tableName`", Quoter.Quote("tableName"))
	assert.Exactly(t, "`databaseName`.`tableName`", Quoter.Quote("databaseName", "tableName"))
	assert.Exactly(t, "`tableName`", Quoter.Quote("", "tableName")) // qualifier is empty
	assert.Exactly(t, "`databaseName`.`tableName`", Quoter.Quote("database`Name", "table`Name"))
}
