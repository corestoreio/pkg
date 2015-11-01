// Copyright 2015, Cyrill @ Schumacher.fm and the CoreStore contributors
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

package httputils_test

import (
	"errors"
	"net/http"
	"testing"

	"crypto/tls"

	"github.com/corestoreio/csfw/config"
	"github.com/corestoreio/csfw/net/httputils"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
)

func TestIsSecure(t *testing.T) {
	tests := []struct {
		ctx          context.Context
		req          *http.Request
		wantIsSecure bool
	}{
		{
			context.Background(),
			func() *http.Request {
				r, err := http.NewRequest("GET", "/", nil)
				if err != nil {
					t.Fatal(err)
				}
				r.TLS = new(tls.ConnectionState)
				return r
			}(),
			true,
		},
		{
			config.NewContextMockReader(context.Background(), config.WithMockValues(config.MockPV{
				config.MockPathScopeDefault(httputils.PathOffloaderHeader): "X_FORWARDED_PROTO",
			})),
			func() *http.Request {
				r, err := http.NewRequest("GET", "/", nil)
				if err != nil {
					t.Fatal(err)
				}
				r.Header.Set("HTTP_X_FORWARDED_PROTO", "https")
				return r
			}(),
			true,
		},
		{
			config.NewContextMockReader(context.Background(), config.WithMockValues(config.MockPV{
				config.MockPathScopeDefault(httputils.PathOffloaderHeader): "X_FORWARDED_PROTO",
			})),
			func() *http.Request {
				r, err := http.NewRequest("GET", "/", nil)
				if err != nil {
					t.Fatal(err)
				}
				r.Header.Set("HTTP_X_FORWARDED_PROTO", "tls")
				return r
			}(),
			false,
		},
		{
			config.NewContextMockReader(context.Background(), config.WithMockValues(config.MockPV{})),
			func() *http.Request {
				r, err := http.NewRequest("GET", "/", nil)
				if err != nil {
					t.Fatal(err)
				}
				r.Header.Set("HTTP_X_FORWARDED_PROTO", "does not matter")
				return r
			}(),
			false,
		},
	}
	for _, test := range tests {
		assert.Exactly(t, test.wantIsSecure, httputils.IsSecure(test.ctx, test.req))
	}
}

func TestIsBaseUrlCorrect(t *testing.T) {

	var nr = func(urlStr string) *http.Request {
		r, err := http.NewRequest("GET", urlStr, nil)
		if err != nil {
			t.Fatal(err)
		}
		return r
	}

	tests := []struct {
		req         *http.Request
		haveBaseURL string
		wantErr     error
	}{
		{nr("http://corestore.io/"), "http://corestore.io/", nil},
		{nr("http://www.corestore.io/"), "http://corestore.io/", httputils.ErrBaseUrlDoNotMatch},
		{nr("http://corestore.io/"), "https://corestore.io/", httputils.ErrBaseUrlDoNotMatch},
		{nr("http://corestore.io/"), "http://corestore.io/subpath", httputils.ErrBaseUrlDoNotMatch},
		{nr("http://corestore.io/subpath"), "http://corestore.io/subpath", nil},
		{nr("http://corestore.io/subpath"), "://cs.io", errors.New("parse ://cs.io: missing protocol scheme")},
	}
	for i, test := range tests {
		haveErr := httputils.IsBaseUrlCorrect(test.req, test.haveBaseURL)
		if test.wantErr != nil {
			assert.EqualError(t, haveErr, test.wantErr.Error(), "Index %d", i)
		} else {
			assert.NoError(t, haveErr, "Index %d", i)
		}
	}
}
