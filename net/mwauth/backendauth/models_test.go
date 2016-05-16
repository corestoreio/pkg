package backendauth_test

import (
	"testing"

	"github.com/corestoreio/csfw/net/mwauth/backendauth"
	"net"
)

func TestIPRange_In(t *testing.T) {
	tests := []struct {
		from   string
		to     string
		testIP string
		want   bool
	}{
		{"0.0.0.0", "255.255.255.255", "128.128.128.128", true},
		{"0.0.0.0", "128.128.128.128", "255.255.255.255", false},
		{"74.50.146.0", "74.50.146.4", "74.50.146.0", true},
		{"74.50.146.0", "74.50.146.4", "74.50.146.4", true},
		{"74.50.146.0", "74.50.146.4", "74.50.146.5", false},
		{"2002:0db8:85a3:0000:0000:8a2e:0370:7334", "74.50.153.4", "74.50.153.2", false},
		{"2002:0db8:85a3:0000:0000:8a2e:0370:7334", "2002:0db8:85a3:0000:0000:8a2e:0370:8334", "2002:0db8:85a3:0000:0000:8a2e:0370:7334", true},
		{"2002:0db8:85a3:0000:0000:8a2e:0370:7334", "2002:0db8:85a3:0000:0000:8a2e:0370:8334", "2002:0db8:85a3:0000:0000:8a2e:0370:7350", true},
		{"2002:0db8:85a3:0000:0000:8a2e:0370:7334", "2002:0db8:85a3:0000:0000:8a2e:0370:8334", "2002:0db8:85a3:0000:0000:8a2e:0370:8334", true},
		{"2002:0db8:85a3:0000:0000:8a2e:0370:7334", "2002:0db8:85a3:0000:0000:8a2e:0370:8334", "2002:0db8:85a3:0000:0000:8a2e:0370:8335", false},
		{"::ffff:183.0.2.128", "::ffff:183.0.2.250", "::ffff:183.0.2.127", false},
		{"::ffff:183.0.2.128", "::ffff:183.0.2.250", "::ffff:183.0.2.128", true},
		{"::ffff:183.0.2.128", "::ffff:183.0.2.250", "::ffff:183.0.2.129", true},
		{"::ffff:183.0.2.128", "::ffff:183.0.2.250", "::ffff:183.0.2.250", true},
		{"::ffff:183.0.2.128", "::ffff:183.0.2.250", "::ffff:183.0.2.251", false},
		{"::ffff:183.0.2.128", "::ffff:183.0.2.250", "183.0.2.130", true},
		{"183.0.2.128", "183.0.2.250", "::ffff:183.0.2.130", true},
		{"unparseable", "183.0.2.250", "::ffff:183.0.2.130", false},
	}
	for _, test := range tests {
		if have, want := backendauth.NewIPRange(test.from, test.to).InStr(test.testIP), test.want; have != want {
			t.Errorf("Assertion (have: %t want: %t) failed on range %s-%s with test %s", have, want, test.from, test.to, test.testIP)
		}
	}
}

var benchmarkIPRange bool

// BenchmarkIPRangeV6-4   	100000000	        19.7 ns/op	       0 B/op	       0 allocs/op
func BenchmarkIPRangeV6(b *testing.B) {
	ir := backendauth.NewIPRange("2002:0db8:85a3:0000:0000:8a2e:0370:7334", "2002:0db8:85a3:0000:0000:8a2e:0370:8334")
	check := net.ParseIP("2002:0db8:85a3:0000:0000:8a2e:0370:7350")
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		benchmarkIPRange = ir.In(check)
	}
	if !benchmarkIPRange {
		b.Fatal("benchmarkIPRange must be true")
	}
}

// BenchmarkIPRangeV4-4   	100000000	        20.0 ns/op	       0 B/op	       0 allocs/op
func BenchmarkIPRangeV4(b *testing.B) {
	ir := backendauth.NewIPRange("74.50.146.0", "74.50.146.4")
	check := net.ParseIP("74.50.146.3")
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		benchmarkIPRange = ir.In(check)
	}
	if !benchmarkIPRange {
		b.Fatal("benchmarkIPRange must be true")
	}
}
