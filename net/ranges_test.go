// Copyright 2015-present, Cyrill @ Schumacher.fm and the CoreStore contributors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package net_test

import (
	"net"
	"testing"

	csnet "github.com/corestoreio/pkg/net"
	"github.com/corestoreio/pkg/util/assert"
)

func TestIPRanges_In(t *testing.T) {
	irs := csnet.IPRanges{
		csnet.MakeIPRange("74.50.146.0", "74.50.146.4"),
		csnet.MakeIPRange("::ffff:183.0.2.128", "::ffff:183.0.2.250"),
	}
	if have, want := irs.InStr("74.50.146.2"), true; want != have {
		t.Errorf("Have %t Want %t", have, want)
	}
	if have, want := irs.InStr("::ffff:183.0.2.250"), true; want != have {
		t.Errorf("Have %t Want %t", have, want)
	}
	if have, want := irs.InStr("::ffff:183.0.2.252"), false; want != have {
		t.Errorf("Have %t Want %t", have, want)
	}
}

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
		if have, want := csnet.MakeIPRange(test.from, test.to).InStr(test.testIP), test.want; have != want {
			t.Errorf("Assertion (have: %t want: %t) failed on range %s-%s with test %s", have, want, test.from, test.to, test.testIP)
		}
	}
}

func TestPrivateIPRanges(t *testing.T) {
	tests := []struct {
		ip   net.IP
		want bool
	}{
		{net.ParseIP("74.50.146.4"), false},
		{net.ParseIP("100.64.1.0"), true},
		{net.ParseIP("192.168.1.3"), true},
		{nil, false},
	}
	for _, test := range tests {
		if have, want := csnet.PrivateIPRanges.In(test.ip), test.want; have != want {
			t.Errorf("Have %t Want %t => IP %s", have, want, test.ip)
		}
	}
}

var benchmarkIPRange bool

// BenchmarkIPRangeV6-4   	10000000	       185 ns/op	      16 B/op	       1 allocs/op
func BenchmarkIPRangeV6(b *testing.B) {
	ir := csnet.MakeIPRange("2002:0db8:85a3:0000:0000:8a2e:0370:7334", "2002:0db8:85a3:0000:0000:8a2e:0370:8334")

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		benchmarkIPRange = ir.InStr("2002:0db8:85a3:0000:0000:8a2e:0370:7350")
	}
	if !benchmarkIPRange {
		b.Fatal("benchmarkIPRange must be true")
	}
}

// BenchmarkIPRangeV4-4   	20000000	       101 ns/op	      16 B/op	       1 allocs/op
func BenchmarkIPRangeV4(b *testing.B) {
	ir := csnet.MakeIPRange("74.50.146.0", "74.50.146.4")
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		benchmarkIPRange = ir.InStr("74.50.146.3")
	}
	if !benchmarkIPRange {
		b.Fatal("benchmarkIPRange must be true")
	}
}

func TestMakeIPRanges_Imbalanced(t *testing.T) {
	if ipr := csnet.MakeIPRanges("Imbalanced"); ipr != nil {
		t.Errorf("Should create a nil slice but got %#v", ipr)
	}
}

func TestMakeIPRanges(t *testing.T) {
	ipr := csnet.MakeIPRanges(
		"10.0.0.0", "10.255.255.255",
		"100.64.0.0", "100.127.255.255",
		"172.16.0.0", "172.31.255.255",
		"192.0.0.0", "192.0.0.255",
		"192.168.0.0", "192.168.255.255",
		"198.18.0.0", "198.19.255.255",
	)
	assert.Exactly(t, csnet.PrivateIPRanges, ipr)
}

func TestIPRanges_Strings(t *testing.T) {
	assert.Exactly(t,
		[]string{"10.0.0.0", "10.255.255.255", "100.64.0.0", "100.127.255.255", "172.16.0.0", "172.31.255.255", "192.0.0.0", "192.0.0.255", "192.168.0.0", "192.168.255.255", "198.18.0.0", "198.19.255.255"},
		csnet.PrivateIPRanges.Strings())
}
