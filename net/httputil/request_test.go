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

package httputil_test

import (
	"net"
	"net/http"
	"testing"

	"github.com/corestoreio/csfw/net/httputil"
	"github.com/stretchr/testify/assert"
)

func TestGetRemoteAddr(t *testing.T) {
	t.Parallel()
	tests := []struct {
		r      *http.Request
		wantIP net.IP
	}{
		{func() *http.Request {
			r, err := http.NewRequest("GET", "http://gopher.go", nil)
			assert.NoError(t, err)
			r.Header.Set("X-Real-IP", "123.123.123.123")
			return r
		}(), net.ParseIP("123.123.123.123")},
		{func() *http.Request {
			r, err := http.NewRequest("GET", "http://gopher.go", nil)
			assert.NoError(t, err)
			r.Header.Set("X-Forwarded-For", "200.100.50.3")
			return r
		}(), net.ParseIP("200.100.50.3")},
		{func() *http.Request {
			r, err := http.NewRequest("GET", "http://gopher.go", nil)
			assert.NoError(t, err)
			r.RemoteAddr = "100.200.50.3"
			return r
		}(), net.ParseIP("100.200.50.3")},
		{func() *http.Request {
			r, err := http.NewRequest("GET", "http://gopher.go", nil)
			assert.NoError(t, err)
			r.RemoteAddr = "100.200.a.3"
			return r
		}(), nil},
	}

	for i, test := range tests {
		haveIP := httputil.GetRemoteAddr(test.r)
		assert.Exactly(t, test.wantIP, haveIP, "Index: %d", i)
	}
}
