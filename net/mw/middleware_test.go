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

package mw_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/corestoreio/pkg/net/mw"
	"github.com/corestoreio/pkg/util/assert"
)

type key uint

const ctxKey key = 0

func fromContext(ctx context.Context) (string, bool) {
	value, ok := ctx.Value(ctxKey).(string)
	return value, ok
}

func serveHTTPContext(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	time.Sleep(time.Millisecond) // wait for other goroutines
	val, _ := fromContext(ctx)
	if _, ok := ctx.Deadline(); ok {
		val += " with deadline"
	}
	if ctx.Err() == context.Canceled {
		val += " canceled"
	}
	_, _ = w.Write([]byte(val))
}

type closeNotifyWriter struct {
	*httptest.ResponseRecorder
}

func (w *closeNotifyWriter) CloseNotify() <-chan bool {
	// return an already "closed" notifier
	notify := make(chan bool, 1)
	notify <- true
	return notify
}

func TestWithCloseHandler(t *testing.T) {
	ctx := context.WithValue(context.Background(), ctxKey, "gopher life")
	mws := mw.MiddlewareSlice{mw.WithCloseNotify()}
	finalCH := mws.Chain(http.HandlerFunc(serveHTTPContext))

	w := &closeNotifyWriter{httptest.NewRecorder()}
	r, err := http.NewRequest("GET", "http://corestore.io/catalog/product/id/3452", nil)
	if err != nil {
		t.Fatal(err)
	}

	finalCH.ServeHTTP(w, r.WithContext(ctx))
	assert.Equal(t, "gopher life canceled", w.Body.String())
}

func TestWithTimeoutHandler(t *testing.T) {
	ctx := context.WithValue(context.Background(), ctxKey, "gopher life")
	finalCH := mw.ChainFunc(serveHTTPContext, mw.WithTimeout(time.Second))

	w := httptest.NewRecorder()
	r, err := http.NewRequest("GET", "http://corestore.io/catalog/product/id/3452", nil)
	if err != nil {
		t.Fatal(err)
	}
	finalCH.ServeHTTP(w, r.WithContext(ctx))
	assert.Equal(t, "gopher life with deadline", w.Body.String())
}

func TestWithHeader(t *testing.T) {
	finalCH := mw.Chain(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`Confirmed landing on drone ship.`))
	}), mw.WithHeader("X-CoreStore-CartID", "0815"))

	w := httptest.NewRecorder()
	r, err := http.NewRequest("GET", "http://corestore.io/catalog/product/id/3452", nil)
	if err != nil {
		t.Fatal(err)
	}
	finalCH.ServeHTTP(w, r)
	if have, want := w.HeaderMap.Get("X-CoreStore-CartID"), "0815"; want != have {
		t.Errorf("Want: %q Have: %q\nHeader: %#v", want, have, w.HeaderMap)
	}
}

func TestWithXHTTPMethodOverrideForm(t *testing.T) {
	var mws mw.MiddlewareSlice
	mws = mws.Append(mw.WithXHTTPMethodOverride(mw.SetMethodOverrideFormKey("_mykey")))

	finalCH := mw.ChainFunc(func(w http.ResponseWriter, r *http.Request) {
		if have, want := r.Method, "HEAD"; want != have {
			t.Errorf("Want: %q Have: %q", want, have)
		}
	}, mws...)

	w := httptest.NewRecorder()
	r, err := http.NewRequest("GET", "http://corestore.io/catalog/product/id/3452", nil)
	if err != nil {
		t.Fatal(err)
	}
	r.Form = url.Values{
		"_mykey": []string{"HEAD"},
	}

	finalCH.ServeHTTP(w, r)
}

func TestWithXHTTPMethodOverrideHeader(t *testing.T) {
	finalCH := mw.ChainFunc(func(w http.ResponseWriter, r *http.Request) {
		if have, want := r.Method, "OPTIONS"; want != have {
			t.Errorf("Want: %q Have: %q", want, have)
		}
	}, mw.WithXHTTPMethodOverride())

	w := httptest.NewRecorder()
	r, err := http.NewRequest("GET", "http://corestore.io/catalog/product/id/3452", nil)
	if err != nil {
		t.Fatal(err)
	}
	r.Header.Set(mw.MethodOverrideHeader, "OPTIONS")
	finalCH.ServeHTTP(w, r)
}

func TestWithXHTTPMethodOverrideNone(t *testing.T) {
	finalCH := mw.ChainFunc(func(w http.ResponseWriter, r *http.Request) {
		if have, want := r.Method, "PATCH"; want != have {
			t.Errorf("Want: %q Have: %q", want, have)
		}
	}, mw.WithXHTTPMethodOverride())

	w := httptest.NewRecorder()
	r, err := http.NewRequest("PATCH", "http://corestore.io/catalog/product/id/3452", nil)
	if err != nil {
		t.Fatal(err)
	}
	finalCH.ServeHTTP(w, r)
}
