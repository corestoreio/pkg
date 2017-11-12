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
	"net/http"
	"net/url"
	"testing"

	"github.com/corestoreio/pkg/net/httputil"
	"github.com/corestoreio/errors"
	"github.com/stretchr/testify/assert"
)

//func TestCtxIsSecure(t *testing.T) {
//
//	woh, err := backend.Backend.WebSecureOffloaderHeader.ToPath(scope.Default, 0)
//	if err != nil {
//		t.Fatal(err)
//	}
//
//	tests := []struct {
//		ctx          context.Context
//		req          *http.Request
//		wantIsSecure bool
//	}{
//		{
//			context.Background(),
//			func() *http.Request {
//				r, err := http.NewRequest("GET", "/", nil)
//				if err != nil {
//					t.Fatal(err)
//				}
//				r.TLS = new(tls.ConnectionState)
//				return r
//			}(),
//			true,
//		},
//		{
//			cfgmock.WithContextScopedGetter(3, 1, context.Background(), cfgmock.WithPV(cfgmock.PathValue{
//				woh.String(): "X_FORWARDED_PROTO",
//			})),
//			func() *http.Request {
//				r, err := http.NewRequest("GET", "/", nil)
//				if err != nil {
//					t.Fatal(err)
//				}
//				r.Header.Set("HTTP_X_FORWARDED_PROTO", "https")
//				return r
//			}(),
//			true,
//		},
//		{
//			cfgmock.WithContextScopedGetter(1, 3, context.Background(), cfgmock.WithPV(cfgmock.PathValue{
//				woh.String(): "X_FORWARDED_PROTO",
//			})),
//			func() *http.Request {
//				r, err := http.NewRequest("GET", "/", nil)
//				if err != nil {
//					t.Fatal(err)
//				}
//				r.Header.Set("HTTP_X_FORWARDED_PROTO", "tls")
//				return r
//			}(),
//			false,
//		},
//		{
//			cfgmock.WithContextScopedGetter(3, 5, context.Background(), cfgmock.WithPV(cfgmock.PathValue{})),
//			func() *http.Request {
//				r, err := http.NewRequest("GET", "/", nil)
//				if err != nil {
//					t.Fatal(err)
//				}
//				r.Header.Set("HTTP_X_FORWARDED_PROTO", "does not matter")
//				return r
//			}(),
//			false,
//		},
//	}
//
//	secReq := httputil.NewCeckSecureRequest(backend.Backend.WebSecureOffloaderHeader)
//	for i, test := range tests {
//		assert.Exactly(t, test.wantIsSecure, secReq.CtxIs(test.ctx, test.req), "Index %d", i)
//	}
//}

func TestIsBaseUrlCorrect(t *testing.T) {

	var nr = func(urlStr string) *http.Request {
		r, err := http.NewRequest("GET", urlStr, nil)
		if err != nil {
			t.Fatal(err)
		}
		return r
	}

	var pu = func(rawURL string) *url.URL {
		u, err := url.Parse(rawURL)
		if err != nil {
			t.Fatal(err)
		}
		return u
	}

	tests := []struct {
		req         *http.Request
		haveBaseURL *url.URL
		wantErrBhf  errors.BehaviourFunc
	}{
		{nr("http://corestore.io/"), pu("http://corestore.io/"), nil},
		{nr("http://www.corestore.io/"), pu("http://corestore.io/"), errors.IsNotValid},
		{nr("http://corestore.io/"), pu("https://corestore.io/"), errors.IsNotValid},
		{nr("http://corestore.io/"), pu("http://corestore.io/subpath"), errors.IsNotValid},
		{nr("http://corestore.io/subpath"), pu("http://corestore.io/subpath"), nil},
		{nr("http://corestore.io/"), pu("http://corestore.io/"), nil},
		{nr("http://corestore.io/subpath/catalog/product/list"), pu("http://corestore.io/subpath"), nil},
	}
	for i, test := range tests {
		haveErr := httputil.IsBaseURLCorrect(test.req, test.haveBaseURL)
		if test.wantErrBhf != nil {
			assert.True(t, test.wantErrBhf(haveErr), "Index %d => %s", i, haveErr)
		} else {
			assert.NoError(t, haveErr, "Index %d", i)
		}
	}
}
