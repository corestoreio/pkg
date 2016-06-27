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

package ratelimit_test

import (
	"github.com/corestoreio/csfw/net"
	"github.com/corestoreio/csfw/net/ratelimit"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestVaryBy_Key(t *testing.T) {
	tests := []struct {
		varyBy func() ratelimit.VaryByer
		req    func() *http.Request
		want   string
	}{
		{
			func() ratelimit.VaryByer { return (*ratelimit.VaryBy)(nil) },
			func() *http.Request {
				return httptest.NewRequest("GET", "https://corestore.io/_1", nil)
			},
			"",
		},
		{
			func() ratelimit.VaryByer {
				return &ratelimit.VaryBy{
					RemoteAddr: true,
				}
			},
			func() *http.Request {
				r := httptest.NewRequest("GET", "https://corestore.io/_2", nil)
				r.Header.Set(net.XClusterClientIP, "123.123.22.11")
				return r
			},
			"123.123.22.11\n",
		},
		{
			func() ratelimit.VaryByer {
				return &ratelimit.VaryBy{
					RemoteAddr: true,
					Method:     true,
				}
			},
			func() *http.Request {
				r := httptest.NewRequest("PUT", "https://corestore.io/_3", nil)
				r.Header.Set(net.XClusterClientIP, "123.123.22.11")
				return r
			},
			"123.123.22.11\nput\n",
		},
		{
			func() ratelimit.VaryByer {
				return &ratelimit.VaryBy{
					RemoteAddr: true,
					Method:     true,
					Path:       true,
				}
			},
			func() *http.Request {
				r := httptest.NewRequest("PUT", "https://corestore.IO/_4Xy", nil)
				r.Header.Set(net.XClusterClientIP, "123.123.22.11")
				return r
			},
			"123.123.22.11\nput\n/_4Xy\n",
		},
		{
			func() ratelimit.VaryByer {
				return &ratelimit.VaryBy{
					RemoteAddr: true,
					Method:     true,
					Path:       true,
					Headers:    []string{"X-Key"},
				}
			},
			func() *http.Request {
				r := httptest.NewRequest("PUT", "https://corestore.IO/_5_Xy", nil)
				r.RemoteAddr = "123.123.22.11"
				r.Header.Set("X-Key", "y-Value")
				return r
			},
			"123.123.22.11\nput\ny-value\n/_5_Xy\n",
		},
		{
			func() ratelimit.VaryByer {
				return &ratelimit.VaryBy{
					RemoteAddr: true,
					Method:     true,
					Path:       true,
					Headers:    []string{"X-Key"},
					Params:     []string{"aa"},
				}
			},
			func() *http.Request {
				r := httptest.NewRequest("PUT", "https://corestore.IO/_6_Xy?aa=bb", nil)
				r.RemoteAddr = "123.123.22.11"
				r.Header.Set("X-Key", "y-valÜe")
				return r
			},
			"123.123.22.11\nput\ny-valÜe\n/_6_Xy\nbb\n",
		},
		{
			func() ratelimit.VaryByer {
				return &ratelimit.VaryBy{
					RemoteAddr: true,
					Method:     true,
					Path:       true,
					Headers:    []string{"X-Key"},
					Params:     []string{"aa"},
					Cookies:    []string{"keksKey"},
				}
			},
			func() *http.Request {
				r := httptest.NewRequest("PUT", "https://corestore.IO/_7_Xy?aa=bb", nil)
				r.RemoteAddr = "123.123.22.11"
				r.Header.Set("X-Key", "y-Value")
				r.AddCookie(&http.Cookie{Name: "keksKey", Value: "keksVal"})
				return r
			},
			"123.123.22.11\nput\ny-value\n/_7_Xy\nbb\nkeksVal\n",
		},
		{
			func() ratelimit.VaryByer {
				return &ratelimit.VaryBy{
					RemoteAddr:  true,
					Method:      true,
					Path:        true,
					Headers:     []string{"X-Key"},
					Params:      []string{"aa"},
					SafeUnicode: true,
				}
			},
			func() *http.Request {
				r := httptest.NewRequest("PUT", "https://corestore.IO/_8_Xy?aa=bb", nil)
				r.RemoteAddr = "123.123.22.11"
				r.Header.Set("X-Key", "y-valÜe")
				return r
			},
			"123.123.22.11\nput\ny-valüe\n/_8_Xy\nbb\n",
		},
	}
	for _, test := range tests {
		vb := test.varyBy()
		r := test.req()
		if have, want := vb.Key(r), test.want; have != want {
			t.Errorf("Case: %q => Have: %q Want: %q", r.URL.String(), have, want)
		}
	}
}

// todo: optimize allocs
// BenchmarkVaryBy_Key/non-unicode-4         	  500000	      3694 ns/op	     288 B/op	      21 allocs/op
// BenchmarkVaryBy_Key/full-unicode-4        	  500000	      3942 ns/op	     288 B/op	      21 allocs/op
func BenchmarkVaryBy_Key(b *testing.B) {

	vb := &ratelimit.VaryBy{
		RemoteAddr:  true,
		Method:      true,
		Path:        true,
		Headers:     []string{"X-Key"},
		Params:      []string{"aa"},
		SafeUnicode: false,
	}
	r := httptest.NewRequest("PUT", "https://corestore.IO/_8_Xy?aa=bb", nil)
	r.RemoteAddr = "123.123.22.11"
	r.Header.Set("X-Key", "y-valÜe")
	const wantNon = "123.123.22.11\nput\ny-valÜe\n/_8_Xy\nbb\n"
	const wantFull = "123.123.22.11\nput\ny-valüe\n/_8_Xy\nbb\n"

	b.Run("non-unicode", func(b *testing.B) {
		b.ReportAllocs()
		//b.ResetTimer()
		for i := 0; i < b.N; i++ {
			if have, want := vb.Key(r), wantNon; have != want {
				b.Errorf("Case Non: Have: %q Want: %q", have, want)
			}
		}
	})

	b.Run("full-unicode", func(b *testing.B) {
		b.ReportAllocs()
		//b.ResetTimer()
		vb.SafeUnicode = true
		for i := 0; i < b.N; i++ {
			if have, want := vb.Key(r), wantFull; have != want {
				b.Errorf("Case Full: Have: %q Want: %q", have, want)
			}
		}
	})
}
