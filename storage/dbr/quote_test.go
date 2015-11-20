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
		{[]string{"a"}, "`a`"},
		{[]string{"a", "b"}, "`a` AS `b`"},
		{[]string{"`c`"}, "`c`"},
		{[]string{"d.e"}, "`d`.`e`"},
		{[]string{"`d`.`e`"}, "`d`.`e`"},
		{[]string{"f", "g", "h"}, "`f` AS `g_h`"},
		{[]string{"f", "g", "h`h"}, "`f` AS `g_hh`"},
	}
	for i, test := range tests {
		assert.Exactly(t, test.want, Quoter.QuoteAs(test.have...), "Index %d", i)
	}
}

// BenchmarkQuoteAs-4	 3000000	       417 ns/op	      48 B/op	       2 allocs/op
func BenchmarkQuoteAs(b *testing.B) {
	want := "`e`.`entity_id` AS `ee`"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if have := Quoter.QuoteAs("e.entity_id", "ee"); have != want {
			b.Errorf("Have %s\nWant %s\n", have, want)
		}
	}
}
