// Copyright 2015-2016, Cyrill @ Schumacher.fm and the CoreStore contributors
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

package request_test

import (
	"net"
	"net/http"
	"testing"

	"net/http/httptest"

	csnet "github.com/corestoreio/cspkg/net"
	"github.com/corestoreio/cspkg/net/auth"
	"github.com/corestoreio/cspkg/net/request"
	"github.com/stretchr/testify/assert"
)

func TestGetRealIP(t *testing.T) {
	t.Parallel()
	tests := []struct {
		r      *http.Request
		opt    int
		wantIP net.IP
	}{
		{func() *http.Request {
			r, _ := http.NewRequest("GET", "http://gopher.go", nil)
			r.Header.Set("X-Real-IP", "123.123.123.123")
			return r
		}(), request.IPForwardedTrust, net.ParseIP("123.123.123.123")},
		{func() *http.Request {
			r, _ := http.NewRequest("GET", "http://gopher.go", nil)
			r.Header.Set("Forwarded-For", "200.100.50.3")
			return r
		}(), request.IPForwardedTrust, net.ParseIP("200.100.50.3")},
		{func() *http.Request {
			r, _ := http.NewRequest("GET", "http://gopher.go", nil)
			r.Header.Set("X-Forwarded", "2002:0db8:85a3:0000:0000:8a2e:0370:7335")
			return r
		}(), request.IPForwardedTrust, net.ParseIP("2002:0db8:85a3:0000:0000:8a2e:0370:7335")},
		{func() *http.Request {
			r, _ := http.NewRequest("GET", "http://gopher.go", nil)
			r.Header.Set("X-Forwarded-For", "200.100.54.4, 192.168.0.100:8080")
			return r
		}(), request.IPForwardedTrust, net.ParseIP("200.100.54.4")},
		{func() *http.Request {
			r, _ := http.NewRequest("GET", "http://gopher.go", nil)
			r.Header.Set("X-Cluster-Client-Ip", "127.0.0.1:8080")
			r.RemoteAddr = "200.100.54.6:8181"
			return r
		}(), request.IPForwardedTrust, net.ParseIP("200.100.54.6")},
		{func() *http.Request {
			r, _ := http.NewRequest("GET", "http://gopher.go", nil)
			r.RemoteAddr = "100.200.50.3"
			return r
		}(), request.IPForwardedTrust, net.ParseIP("100.200.50.3")},
		{func() *http.Request {
			r, _ := http.NewRequest("GET", "http://gopher.go", nil)
			r.Header.Set("X-Forwarded-For", "127.0.0.1:8080")
			r.RemoteAddr = "2002:0db8:85a3:0000:0000:8a2e:0370:7334"
			return r
		}(), request.IPForwardedIgnore, net.ParseIP("2002:0db8:85a3:0000:0000:8a2e:0370:7334")},
		{func() *http.Request {
			r, _ := http.NewRequest("GET", "http://gopher.go", nil)
			r.RemoteAddr = "100.200.a.3"
			return r
		}(), request.IPForwardedIgnore, nil},
	}

	for i, test := range tests {
		haveIP := request.RealIP(test.r, test.opt)
		assert.Exactly(t, test.wantIP, haveIP, "Index: %d Want %s Have %s", i, test.wantIP, haveIP)
	}
}

// check if the returned function conforms with the auth package
var _ auth.TriggerFunc = request.InIPRange("192.168.0.1", "192.168.0.100")
var _ auth.TriggerFunc = request.NotInIPRange("192.168.0.1", "192.168.0.100")

func TestInIPRange(t *testing.T) {
	t.Parallel()
	rf := request.InIPRange(
		"10.0.0.0", "10.255.255.255",
		"100.64.0.0", "100.127.255.255",
	)
	r := httptest.NewRequest("GET", "/", nil)
	r.RemoteAddr = "10.2.3.4"
	assert.True(t, rf(r))
	r.RemoteAddr = "192.168.0.1"
	assert.False(t, rf(r))

}

func TestNotInIPRange(t *testing.T) {
	t.Parallel()
	rf := request.NotInIPRange(
		"10.0.0.0", "10.255.255.255",
		"100.64.0.0", "100.127.255.255",
	)
	r := httptest.NewRequest("GET", "/", nil)
	r.RemoteAddr = "10.2.3.4"
	assert.False(t, rf(r))
	r.RemoteAddr = "192.168.0.1"
	assert.True(t, rf(r))

	rf = request.NotInIPRange()
	assert.True(t, rf(r))
}

var benchmarkInIPRange bool

func BenchmarkInIPRange(b *testing.B) {
	rf := request.InIPRange("8.8.0.0", "8.8.255.255")
	r := httptest.NewRequest("GET", "/", nil)
	r.Header.Set(csnet.XRealIP, "8.8.8.8")
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		benchmarkInIPRange = rf(r)
	}
	if !benchmarkInIPRange {
		b.Fatalf("Expecting true but got %t", benchmarkInIPRange)
	}
}
