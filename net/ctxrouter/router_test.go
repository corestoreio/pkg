// Copyright 2013 Julien Schmidt. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be found
// in the LICENSE file.

package ctxrouter

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"net"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/corestoreio/csfw/net/ctxhttp"
	"github.com/juju/errgo"
	"golang.org/x/net/context"
)

type mockResponseWriter struct{}

func (m *mockResponseWriter) Header() (h http.Header) {
	return http.Header{}
}

func (m *mockResponseWriter) Write(p []byte) (n int, err error) {
	return len(p), nil
}

func (m *mockResponseWriter) WriteString(s string) (n int, err error) {
	return len(s), nil
}

func (m *mockResponseWriter) WriteHeader(int) {}

func TestParams(t *testing.T) {
	ps := Params{
		Param{"param1", "value1"},
		Param{"param2", "value2"},
		Param{"param3", "value3"},
	}
	for i := range ps {
		if val := ps.ByName(ps[i].Key); val != ps[i].Value {
			t.Errorf("Wrong value for %s: Got %s; Want %s", ps[i].Key, val, ps[i].Value)
		}
	}
	if val := ps.ByName("noKey"); val != "" {
		t.Errorf("Expected empty string for not found key; got: %s", val)
	}
	haveP := FromContextParams(context.Background())
	if len(haveP) != 0 || haveP == nil {
		t.Errorf("ParamsFromContext should return a non-nil slice with length zero.\nHave: %#v", haveP)
	}
}

func TestRouter(t *testing.T) {
	router := New()
	router.Use(noopMW())

	routed := false
	router.Handle("GET", "/user/:name", func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
		routed = true
		want := Params{Param{"name", "gopher"}}
		ps := FromContextParams(ctx)
		if !reflect.DeepEqual(ps, want) {
			t.Fatalf("wrong wildcard values: want %v, got %v", want, ps)
		}
		return nil
	})

	w := new(mockResponseWriter)

	req, _ := http.NewRequest("GET", "/user/gopher", nil)
	router.ServeHTTP(w, req)

	if !routed {
		t.Fatal("routing failed")
	}
}

func TestRouterContext(t *testing.T) {
	ctx := context.Background()
	ctx = context.WithValue(ctx, 1, "Gopher")
	router := New(ctx)
	router.Use(noopMW())

	routed := false
	router.Handle("GET", "/user", func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
		routed = true

		have, ok := ctx.Value(1).(string)
		if !ok {
			t.Fatal("Cannot extract key 1 from context")
		}
		want := "Gopher"
		if have != want {
			t.Fatal("Have %s\nWant: %s", have, want)
		}

		return nil
	})

	w := new(mockResponseWriter)

	req, _ := http.NewRequest("GET", "/user", nil)
	router.ServeHTTP(w, req)

	if !routed {
		t.Fatal("routing failed")
	}
}

type handlerStruct struct {
	handeled *bool
}

func (h handlerStruct) ServeHTTPContext(_ context.Context, _ http.ResponseWriter, _ *http.Request) error {
	*h.handeled = true
	return nil
}

func TestRouterAPI(t *testing.T) {
	var get, head, options, post, put, patch, delete, handler, handlerFunc bool

	httpHandler := handlerStruct{&handler}

	router := New()
	router.GET("/GET", func(_ context.Context, w http.ResponseWriter, r *http.Request) error {
		get = true
		return nil
	})
	router.HEAD("/GET", func(_ context.Context, w http.ResponseWriter, r *http.Request) error {
		head = true
		return nil
	})
	router.OPTIONS("/GET", func(_ context.Context, w http.ResponseWriter, r *http.Request) error {
		options = true
		return nil
	})
	router.POST("/POST", func(_ context.Context, w http.ResponseWriter, r *http.Request) error {
		post = true
		return nil
	})
	router.PUT("/PUT", func(_ context.Context, w http.ResponseWriter, r *http.Request) error {
		put = true
		return nil
	})
	router.PATCH("/PATCH", func(_ context.Context, w http.ResponseWriter, r *http.Request) error {
		patch = true
		return nil
	})
	router.DELETE("/DELETE", func(_ context.Context, w http.ResponseWriter, r *http.Request) error {
		delete = true
		return nil
	})
	router.Handler("GET", "/Handler", httpHandler)
	router.HandlerFunc("GET", "/HandlerFunc", func(_ context.Context, _ http.ResponseWriter, _ *http.Request) error {
		handlerFunc = true
		return nil
	})

	w := new(mockResponseWriter)

	r, _ := http.NewRequest("GET", "/GET", nil)
	router.ServeHTTP(w, r)
	if !get {
		t.Error("routing GET failed")
	}

	r, _ = http.NewRequest("HEAD", "/GET", nil)
	router.ServeHTTP(w, r)
	if !head {
		t.Error("routing HEAD failed")
	}

	r, _ = http.NewRequest("OPTIONS", "/GET", nil)
	router.ServeHTTP(w, r)
	if !options {
		t.Error("routing OPTIONS failed")
	}

	r, _ = http.NewRequest("POST", "/POST", nil)
	router.ServeHTTP(w, r)
	if !post {
		t.Error("routing POST failed")
	}

	r, _ = http.NewRequest("PUT", "/PUT", nil)
	router.ServeHTTP(w, r)
	if !put {
		t.Error("routing PUT failed")
	}

	r, _ = http.NewRequest("PATCH", "/PATCH", nil)
	router.ServeHTTP(w, r)
	if !patch {
		t.Error("routing PATCH failed")
	}

	r, _ = http.NewRequest("DELETE", "/DELETE", nil)
	router.ServeHTTP(w, r)
	if !delete {
		t.Error("routing DELETE failed")
	}

	r, _ = http.NewRequest("GET", "/Handler", nil)
	router.ServeHTTP(w, r)
	if !handler {
		t.Error("routing Handler failed")
	}

	r, _ = http.NewRequest("GET", "/HandlerFunc", nil)
	router.ServeHTTP(w, r)
	if !handlerFunc {
		t.Error("routing HandlerFunc failed")
	}
}

func TestRouterRoot(t *testing.T) {
	router := New()
	recv := catchPanic(func() {
		router.GET("noSlashRoot", nil)
	})
	if recv == nil {
		t.Fatal("registering path not beginning with '/' did not panic")
	}
}

func TestRouterChaining(t *testing.T) {
	router1 := New()
	router2 := New()
	router1.NotFound = router2

	fooHit := false
	router1.POST("/foo", func(_ context.Context, w http.ResponseWriter, req *http.Request) error {
		fooHit = true
		w.WriteHeader(http.StatusOK)
		return nil
	})

	barHit := false
	router2.POST("/bar", func(_ context.Context, w http.ResponseWriter, req *http.Request) error {
		barHit = true
		w.WriteHeader(http.StatusOK)
		return nil
	})

	r, _ := http.NewRequest("POST", "/foo", nil)
	w := httptest.NewRecorder()
	router1.ServeHTTP(w, r)
	if !(w.Code == http.StatusOK && fooHit) {
		t.Errorf("Regular routing failed with router chaining.")
		t.FailNow()
	}

	r, _ = http.NewRequest("POST", "/bar", nil)
	w = httptest.NewRecorder()
	router1.ServeHTTP(w, r)
	if !(w.Code == http.StatusOK && barHit) {
		t.Errorf("Chained routing failed with router chaining.")
		t.FailNow()
	}

	r, _ = http.NewRequest("POST", "/qax", nil)
	w = httptest.NewRecorder()
	router1.ServeHTTP(w, r)
	if !(w.Code == http.StatusNotFound) {
		t.Errorf("NotFound behavior failed with router chaining.")
		t.FailNow()
	}
}

func TestRouterNotAllowed(t *testing.T) {
	handlerFunc := func(_ context.Context, _ http.ResponseWriter, _ *http.Request) error { return nil }

	router := New()
	router.POST("/path", handlerFunc)

	// Test not allowed
	r, _ := http.NewRequest("GET", "/path", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, r)
	if !(w.Code == http.StatusMethodNotAllowed) {
		t.Errorf("NotAllowed handling failed: Code=%d, Header=%v", w.Code, w.Header())
	}

	w = httptest.NewRecorder()
	responseText := "custom method"
	router.MethodNotAllowed = ctxhttp.HandlerFunc(func(_ context.Context, w http.ResponseWriter, req *http.Request) error {
		w.WriteHeader(http.StatusTeapot)
		w.Write([]byte(responseText))
		return nil
	})
	router.ServeHTTP(w, r)
	if got := w.Body.String(); !(got == responseText) {
		t.Errorf("unexpected response got %q want %q", got, responseText)
	}
	if w.Code != http.StatusTeapot {
		t.Errorf("unexpected response code %d want %d", w.Code, http.StatusTeapot)
	}
}

func TestRouterNotFound(t *testing.T) {
	handlerFunc := func(_ context.Context, _ http.ResponseWriter, _ *http.Request) error { return nil }

	router := New()
	router.GET("/path", handlerFunc)
	router.GET("/dir/", handlerFunc)
	router.GET("/", handlerFunc)

	testRoutes := []struct {
		route  string
		code   int
		header string
	}{
		{"/path/", 301, "map[Location:[/path]]"},   // TSR -/
		{"/dir", 301, "map[Location:[/dir/]]"},     // TSR +/
		{"", 301, "map[Location:[/]]"},             // TSR +/
		{"/PATH", 301, "map[Location:[/path]]"},    // Fixed Case
		{"/DIR/", 301, "map[Location:[/dir/]]"},    // Fixed Case
		{"/PATH/", 301, "map[Location:[/path]]"},   // Fixed Case -/
		{"/DIR", 301, "map[Location:[/dir/]]"},     // Fixed Case +/
		{"/../path", 301, "map[Location:[/path]]"}, // CleanPath
		{"/nope", 404, ""},                         // NotFound
	}
	for _, tr := range testRoutes {
		r, _ := http.NewRequest("GET", tr.route, nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, r)
		if !(w.Code == tr.code && (w.Code == 404 || fmt.Sprint(w.Header()) == tr.header)) {
			t.Errorf("NotFound handling route %s failed: Code=%d, Header=%v", tr.route, w.Code, w.Header())
		}
	}

	// Test custom not found handler
	var notFound bool
	router.NotFound = ctxhttp.HandlerFunc(func(_ context.Context, rw http.ResponseWriter, r *http.Request) error {
		rw.WriteHeader(404)
		notFound = true
		return nil
	})
	r, _ := http.NewRequest("GET", "/nope", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, r)
	if !(w.Code == 404 && notFound == true) {
		t.Errorf("Custom NotFound handler failed: Code=%d, Header=%v", w.Code, w.Header())
	}

	// Test other method than GET (want 307 instead of 301)
	router.PATCH("/path", handlerFunc)
	r, _ = http.NewRequest("PATCH", "/path/", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, r)
	if !(w.Code == 307 && fmt.Sprint(w.Header()) == "map[Location:[/path]]") {
		t.Errorf("Custom NotFound handler failed: Code=%d, Header=%v", w.Code, w.Header())
	}

	// Test special case where no node for the prefix "/" exists
	router = New()
	router.GET("/a", handlerFunc)
	r, _ = http.NewRequest("GET", "/", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, r)
	if !(w.Code == 404) {
		t.Errorf("NotFound handling route / failed: Code=%d", w.Code)
	}
}

func TestRouterPanicHandler(t *testing.T) {
	router := New()
	panicHandled := false

	router.PanicHandler = func(ctx context.Context, rw http.ResponseWriter, r *http.Request) error {
		panicHandled = true
		pa, ok := FromContextParams(ctx).(string)
		if !ok {
			t.Error("PanicFromContext should return a string")
		}
		if pa != "oops!" {
			t.Errorf("Want: oops!\nHave: %s\nPanic %#v\n", pa, FromContextParams(ctx))
		}
		return nil
	}

	router.Handle("PUT", "/user/:name", func(_ context.Context, _ http.ResponseWriter, _ *http.Request) error {
		panic("oops!")
	})

	w := new(mockResponseWriter)
	req, _ := http.NewRequest("PUT", "/user/gopher", nil)

	defer func() {
		if rcv := recover(); rcv != nil {
			t.Fatal("handling panic failed")
		}
	}()

	router.ServeHTTP(w, req)

	if !panicHandled {
		t.Fatal("simulating failed")
	}
}

func TestRouterPanicHandlerError(t *testing.T) {
	router := New()

	router.PanicHandler = func(ctx context.Context, rw http.ResponseWriter, r *http.Request) error {
		return errgo.New("Epic fail")
	}

	router.Handle("PUT", "/user/:name", func(_ context.Context, _ http.ResponseWriter, _ *http.Request) error {
		panic("oops!")
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/user/gopher", nil)

	defer func() {
		if rcv := recover(); rcv != nil {
			t.Fatal("handling panic failed")
		}
	}()

	router.ServeHTTP(w, req)
	if false == strings.Contains(w.Body.String(), "Epic fail\n") {
		t.Errorf("Body should contain Epic fail\nHave: %s", w.Body)
	}
	if w.Code != http.StatusInternalServerError {
		t.Errorf("Have: %d\nWant: %d", w.Code, http.StatusInternalServerError)
	}
}

func TestServeHTTPError(t *testing.T) {
	router := New()

	router.Handle("PUT", "/user/:name", func(_ context.Context, _ http.ResponseWriter, _ *http.Request) error {
		return errgo.New("TestServeHTTPError")
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/user/gopher", nil)

	router.ServeHTTP(w, req)
	if false == strings.Contains(w.Body.String(), "TestServeHTTPError\n") {
		t.Errorf("Body should contain TestServeHTTPError\nHave: %s", w.Body)
	}
	if w.Code != http.StatusInternalServerError {
		t.Errorf("Have: %d\nWant: %d", w.Code, http.StatusInternalServerError)
	}
}

func TestRouterLookup(t *testing.T) {
	routed := false
	wantHandle := func(_ context.Context, _ http.ResponseWriter, _ *http.Request) error {
		routed = true
		return nil
	}
	wantParams := Params{Param{"name", "gopher"}}

	router := New()

	// try empty router first
	handle, _, _, tsr := router.Lookup("GET", "/nope")
	if handle != nil {
		t.Fatalf("Got handle for unregistered pattern: %v", handle)
	}
	if tsr {
		t.Error("Got wrong TSR recommendation!")
	}

	// insert route and try again
	router.GET("/user/:name", wantHandle)

	handle, _, params, tsr := router.Lookup("GET", "/user/gopher")
	if handle == nil {
		t.Fatal("Got no handle!")
	} else {
		handle(nil, nil, nil)
		if !routed {
			t.Fatal("Routing failed!")
		}
	}

	if !reflect.DeepEqual(params, wantParams) {
		t.Fatalf("Wrong parameter values: want %v, got %v", wantParams, params)
	}

	handle, _, _, tsr = router.Lookup("GET", "/user/gopher/")
	if handle != nil {
		t.Fatalf("Got handle for unregistered pattern: %v", handle)
	}
	if !tsr {
		t.Error("Got no TSR recommendation!")
	}

	handle, _, _, tsr = router.Lookup("GET", "/nope")
	if handle != nil {
		t.Fatalf("Got handle for unregistered pattern: %v", handle)
	}
	if tsr {
		t.Error("Got wrong TSR recommendation!")
	}
}

type mockFileSystem struct {
	opened bool
}

func (mfs *mockFileSystem) Open(name string) (http.File, error) {
	mfs.opened = true
	return nil, errors.New("this is just a mock")
}

func TestRouterServeFiles(t *testing.T) {
	router := New()
	mfs := &mockFileSystem{}

	recv := catchPanic(func() {
		router.ServeFiles("/noFilepath", mfs)
	})
	if recv == nil {
		t.Fatal("registering path not ending with '*filepath' did not panic")
	}

	router.ServeFiles("/*filepath", mfs)
	w := new(mockResponseWriter)
	r, _ := http.NewRequest("GET", "/favicon.ico", nil)
	router.ServeHTTP(w, r)
	if !mfs.opened {
		t.Error("serving file failed")
	}
}

// recorderHijacker is a ResponseRecorder that can be hijacked.
type recorderHijacker struct {
	*httptest.ResponseRecorder
	fakeConn *fakeConn
	buf      bytes.Buffer
}

func (rh *recorderHijacker) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	return rh.fakeConn, bufio.NewReadWriter(bufio.NewReader(&rh.buf), bufio.NewWriter(&rh.buf)), nil
}

type fakeConn struct {
	readBuf  bytes.Buffer
	writeBuf bytes.Buffer
}

func (c *fakeConn) LocalAddr() net.Addr                { return nil }
func (c *fakeConn) RemoteAddr() net.Addr               { return nil }
func (c *fakeConn) SetDeadline(t time.Time) error      { return nil }
func (c *fakeConn) SetReadDeadline(t time.Time) error  { return nil }
func (c *fakeConn) SetWriteDeadline(t time.Time) error { return nil }
func (c *fakeConn) Close() error                       { return nil }
func (c *fakeConn) Read(b []byte) (int, error)         { return c.readBuf.Read(b) }
func (c *fakeConn) Write(b []byte) (int, error)        { return c.writeBuf.Write(b) }

func TestWEBSOCKET(t *testing.T) {

	// TODO: setting up a proper websocket server test is quite complicated ...
	// see test files in "golang.org/x/net/websocket"

	router := New()
	router.WEBSOCKET("/ws", func(ctx context.Context, _ http.ResponseWriter, r *http.Request) error {
		ws, ok := FromContextWebsocket(ctx)
		if !ok {
			t.Fatal("Expecting a Websocket but got nothing")
		}
		if ws == nil {
			t.Fatal("Expecting a Websocket but got nil")
		}

		return nil
	})

	w := &recorderHijacker{
		ResponseRecorder: httptest.NewRecorder(),
		fakeConn:         &fakeConn{},
	}
	r, _ := http.NewRequest("GET", "ws://localhost/ws", nil)
	router.ServeHTTP(w, r)

	// Bad Request
	//t.Log(w.Body.String(), w.HeaderMap, w.buf.String())
}
