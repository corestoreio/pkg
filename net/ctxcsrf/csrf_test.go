package csrf

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"goji.io/pat"

	"golang.org/x/net/context"

	"goji.io"
)

var testKey = []byte("keep-it-secret-keep-it-safe-----")
var testHandler = goji.HandlerFunc(func(ctx context.Context, w http.ResponseWriter, r *http.Request) {})

// TestProtect is a high-level test to make sure the middleware returns the
// wrapped handler with a 200 OK status.
func TestProtect(t *testing.T) {
	m := goji.NewMux()
	m.UseC(Protect(testKey))
	m.HandleFuncC(pat.Get("/"), testHandler)

	r, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	m.ServeHTTP(rr, r)

	if rr.Code != http.StatusOK {
		t.Fatalf("middleware failed to pass to the next handler: got %v want %v",
			rr.Code, http.StatusOK)
	}

	if rr.Header().Get("Set-Cookie") == "" {
		t.Fatalf("cookie not set: got %q", rr.Header().Get("Set-Cookie"))
	}
}

// Test that idempotent methods return a 200 OK status and that non-idempotent
// methods return a 403 Forbidden status when a CSRF cookie is not present.
func TestMethods(t *testing.T) {
	m := goji.NewMux()
	m.UseC(Protect(testKey))
	m.HandleFuncC(pat.New("/"), testHandler)

	// Test idempontent ("safe") methods
	for _, method := range safeMethods {
		r, err := http.NewRequest(method, "/", nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		m.ServeHTTP(rr, r)

		if rr.Code != http.StatusOK {
			t.Fatalf("middleware failed to pass to the next handler: got %v want %v",
				rr.Code, http.StatusOK)
		}

		if rr.Header().Get("Set-Cookie") == "" {
			t.Fatalf("cookie not set: got %q", rr.Header().Get("Set-Cookie"))
		}
	}

	// Test non-idempotent methods (should return a 403 without a cookie set)
	nonIdempotent := []string{"POST", "PUT", "DELETE", "PATCH"}
	for _, method := range nonIdempotent {
		r, err := http.NewRequest(method, "/", nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		m.ServeHTTP(rr, r)

		if rr.Code != http.StatusForbidden {
			t.Fatalf("middleware failed to pass to the next handler: got %v want %v",
				rr.Code, http.StatusForbidden)
		}

		if rr.Header().Get("Set-Cookie") == "" {
			t.Fatalf("cookie not set: got %q", rr.Header().Get("Set-Cookie"))
		}
	}

}

// TestBadCookie tests for failure when a cookie header is modified (malformed).
func TestBadCookie(t *testing.T) {
	m := goji.NewMux()
	m.UseC(Protect(testKey))

	var token string
	m.HandleFuncC(pat.New("/"), func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
		token = Token(ctx, r)
	})

	// Obtain a CSRF cookie via a GET request.
	r, err := http.NewRequest("GET", "http://www.gorillatoolkit.org/", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	m.ServeHTTP(rr, r)

	// POST the token back in the header.
	r, err = http.NewRequest("POST", "http://www.gorillatoolkit.org/", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Replace the cookie prefix
	badHeader := strings.Replace("_csrfToken=", rr.Header().Get("Set-Cookie"), "_badCookie", -1)
	r.Header.Set("Cookie", badHeader)
	r.Header.Set("X-CSRF-Token", token)
	r.Header.Set("Referer", "http://www.gorillatoolkit.org/")

	rr = httptest.NewRecorder()
	m.ServeHTTP(rr, r)

	if rr.Code != http.StatusForbidden {
		t.Fatalf("middleware failed to reject a bad cookie: got %v want %v",
			rr.Code, http.StatusForbidden)
	}

}

func TestErrorHandler(t *testing.T) {
	m := goji.NewMux()
	errorHandler := goji.HandlerFunc(func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
		if ctx.Value(errorKey) == nil {
			t.Errorf("should have access to a goji error")
		}
		http.Error(w, "", http.StatusTeapot)
	})
	m.UseC(Protect(testKey, ErrorHandler(errorHandler)))

	r, _ := http.NewRequest("POST", "/", nil)

	rr := httptest.NewRecorder()
	m.ServeHTTPC(context.Background(), rr, r)

	if rr.Code != http.StatusTeapot {
		t.Fatalf("custom error handler was not called: got %v want %v",
			rr.Code, http.StatusTeapot)
	}
}

// Responses should set a "Vary: Cookie" header to protect client/proxy caching.
func TestVaryHeader(t *testing.T) {
	m := goji.NewMux()
	m.UseC(Protect(testKey))
	m.HandleFuncC(pat.Get("/"), testHandler)

	r, err := http.NewRequest("HEAD", "https://www.golang.org/", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	m.ServeHTTP(rr, r)

	if rr.Code != http.StatusOK {
		t.Fatalf("middleware failed to pass to the next handler: got %v want %v",
			rr.Code, http.StatusOK)
	}

	if rr.Header().Get("Vary") != "Cookie" {
		t.Fatalf("vary header not set: got %q want %q", rr.Header().Get("Vary"), "Cookie")
	}
}

// Requests with no Referer header should fail.
func TestNoReferer(t *testing.T) {
	m := goji.NewMux()
	m.UseC(Protect(testKey))
	m.HandleFuncC(pat.Get("/"), testHandler)

	r, err := http.NewRequest("POST", "https://golang.org/", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	m.ServeHTTP(rr, r)

	if rr.Code != http.StatusForbidden {
		t.Fatalf("middleware failed to pass to the next handler: got %v want %v",
			rr.Code, http.StatusForbidden)
	}
}

// TestBadReferer checks that HTTPS requests with a Referer that does not
// match the request URL correctly fail CSRF validation.
func TestBadReferer(t *testing.T) {
	m := goji.NewMux()
	m.UseC(Protect(testKey))

	var token string
	m.HandleFuncC(pat.New("/"), func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
		token = Token(ctx, r)
	})

	// Obtain a CSRF cookie via a GET request.
	r, err := http.NewRequest("GET", "https://www.gorillatoolkit.org/", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	m.ServeHTTP(rr, r)

	// POST the token back in the header.
	r, err = http.NewRequest("POST", "https://www.gorillatoolkit.org/", nil)
	if err != nil {
		t.Fatal(err)
	}

	setCookie(rr, r)
	r.Header.Set("X-CSRF-Token", token)

	// Set a non-matching Referer header.
	r.Header.Set("Referer", "http://goji.io")

	rr = httptest.NewRecorder()
	m.ServeHTTP(rr, r)

	if rr.Code != http.StatusForbidden {
		t.Fatalf("middleware failed to pass to the next handler: got %v want %v",
			rr.Code, http.StatusForbidden)
	}
}

// Requests with a valid Referer should pass.
func TestWithReferer(t *testing.T) {
	m := goji.NewMux()
	m.UseC(Protect(testKey))

	var token string
	m.HandleFuncC(pat.New("/"), func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
		token = Token(ctx, r)
	})

	// Obtain a CSRF cookie via a GET request.
	r, err := http.NewRequest("GET", "http://www.gorillatoolkit.org/", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	m.ServeHTTP(rr, r)

	// POST the token back in the header.
	r, err = http.NewRequest("POST", "http://www.gorillatoolkit.org/", nil)
	if err != nil {
		t.Fatal(err)
	}

	setCookie(rr, r)
	r.Header.Set("X-CSRF-Token", token)
	r.Header.Set("Referer", "http://www.gorillatoolkit.org/")

	rr = httptest.NewRecorder()
	m.ServeHTTP(rr, r)

	if rr.Code != http.StatusOK {
		t.Fatalf("middleware failed to pass to the next handler: got %v want %v",
			rr.Code, http.StatusOK)
	}
}

// TestFormField tests that a token in the form field takes precedence over a
// token in the HTTP header.
// TODO(matt): Finish this test.
func TestFormField(t *testing.T) {

}

func setCookie(rr *httptest.ResponseRecorder, r *http.Request) {
	r.Header.Set("Cookie", rr.Header().Get("Set-Cookie"))
}
