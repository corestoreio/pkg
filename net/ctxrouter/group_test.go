package ctxrouter

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/corestoreio/csfw/net/ctxhttp"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
)

func noopMW() ctxhttp.Middleware {
	return func(hf ctxhttp.HandlerFunc) ctxhttp.HandlerFunc {
		return func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
			return hf(ctx, w, r)
		}
	}
}

func TestGroup(t *testing.T) {
	g := New().Group("/group")

	g.Use(noopMW())

	h := func(context.Context, http.ResponseWriter, *http.Request) error { return nil }
	//g.CONNECT("/", h)
	g.DELETE("/", h)
	g.GET("/", h)
	g.HEAD("/", h)
	g.OPTIONS("/", h)
	g.PATCH("/", h)
	g.POST("/", h)
	g.PUT("/", h)
	g.WEBSOCKET("/ws", h)

	g2 := g.Group("/files")
	mfs := &mockFileSystem{}
	recv := catchPanic(func() {
		g2.ServeFiles("/noFilepath", mfs)
	})
	if recv == nil {
		t.Fatal("registering path not ending with '*filepath' did not panic")
	}
	g2.ServeFiles("/*filepath", mfs)
}

func groupHeader() ctxhttp.Middleware {
	return func(hf ctxhttp.HandlerFunc) ctxhttp.HandlerFunc {
		return func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
			w.Header().Set("X-CoreStore-ID", "Goph3r")
			println(r.RequestURI)
			return hf(ctx, w, r)
		}
	}
}

func TestGroupMiddleware(t *testing.T) {
	r := New()
	g := r.Group("/group", groupHeader())
	h := func(context.Context, http.ResponseWriter, *http.Request) error { return errors.New("Group Error") }
	g.GET("/error", h)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/group/error", nil)
	r.ServeHTTP(w, req)

	assert.Exactly(t, 500, w.Code)
	assert.Exactly(t, "Group Error\n", w.Body.String())
	assert.Exactly(t, "Goph3r", w.Header().Get("X-CoreStore-ID"), "Header key X-CoreStore-ID not found, which has been applied by a middleware")
}

func TestGroupMiddlewareMultipleRoutes(t *testing.T) {
	t.Log("add tests where we have different groups with different middlewarez and they do not each other conflict")
}
