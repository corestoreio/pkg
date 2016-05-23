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

	"github.com/corestoreio/csfw/net/request"
	"github.com/stretchr/testify/assert"
)

func TestGetRealIP(t *testing.T) {

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
