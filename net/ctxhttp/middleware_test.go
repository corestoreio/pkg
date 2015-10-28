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

package ctxhttp_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"time"

	"github.com/corestoreio/csfw/net/ctxhttp"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
)

type key uint

const ctxKey key = 0

func newContext(ctx context.Context, value string) context.Context {
	return context.WithValue(ctx, ctxKey, value)
}

func fromContext(ctx context.Context) (string, bool) {
	value, ok := ctx.Value(ctxKey).(string)
	return value, ok
}

type handler struct{}

func (h handler) ServeHTTPContext(ctx context.Context, w http.ResponseWriter, r *http.Request) error {

	time.Sleep(time.Millisecond) // wait for other goroutines
	val, _ := fromContext(ctx)
	if _, ok := ctx.Deadline(); ok {
		val += " with deadline"
	}
	if ctx.Err() == context.Canceled {
		val += " canceled"
	}
	_, err := w.Write([]byte(val))
	return err
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
	finalCH := ctxhttp.Chain(&handler{}, ctxhttp.WithCloseNotify())

	w := &closeNotifyWriter{httptest.NewRecorder()}
	r, err := http.NewRequest("GET", "http://corestore.io/catalog/product/id/3452", nil)
	if err != nil {
		t.Fatal(err)
	}
	if err := finalCH.ServeHTTPContext(ctx, w, r); err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, "gopher life canceled", w.Body.String())
}

func TestWithTimeoutHandler(t *testing.T) {
	ctx := context.WithValue(context.Background(), ctxKey, "gopher life")
	finalCH := ctxhttp.Chain(&handler{}, ctxhttp.WithTimeout(time.Second))

	w := httptest.NewRecorder()
	r, err := http.NewRequest("GET", "http://corestore.io/catalog/product/id/3452", nil)
	if err != nil {
		t.Fatal(err)
	}
	if err := finalCH.ServeHTTPContext(ctx, w, r); err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, "gopher life with deadline", w.Body.String())
}
