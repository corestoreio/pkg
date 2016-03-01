package csrf

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"text/template"

	"goji.io/pat"
	"golang.org/x/net/context"

	"goji.io"
)

var testTemplate = `
<html>
<body>
<form action="/" method="POST">
{{ .csrfField }}
</form>
</body>
</html>
`
var testFieldName = "custom_csrf_field_name"
var testTemplateField = `<input type="hidden" name="%s" value="%s">`

// Test that our form helpers correctly inject a token into the response body.
func TestFormToken(t *testing.T) {
	m := goji.NewMux()
	m.UseC(Protect(testKey))

	// Make the token available outside of the handler for comparison.
	var token string
	m.HandleFuncC(pat.New("/"), func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
		token = Token(ctx, r)
		t := template.Must((template.New("base").Parse(testTemplate)))
		t.Execute(w, map[string]interface{}{
			TemplateTag: TemplateField(ctx, r),
		})
	})

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

	if len(token) != base64.StdEncoding.EncodedLen(tokenLength*2) {
		t.Fatalf("token length invalid: got %v want %v", len(token), base64.StdEncoding.EncodedLen(tokenLength*2))
	}

	if !strings.Contains(rr.Body.String(), token) {
		t.Fatalf("token not in response body: got %v want %v", rr.Body.String(), token)
	}
}

// Test that we can extract a CSRF token from a multipart form.
func TestMultipartFormToken(t *testing.T) {
	m := goji.NewMux()
	m.UseC(Protect(testKey))

	// Make the token available outside of the handler for comparison.
	var token string
	m.HandleFuncC(pat.New("/"), func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
		token = Token(ctx, r)
		t := template.Must((template.New("base").Parse(testTemplate)))
		t.Execute(w, map[string]interface{}{
			TemplateTag: TemplateField(ctx, r),
		})
	})

	r, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	m.ServeHTTP(rr, r)

	// Set up our multipart form
	var b bytes.Buffer
	mp := multipart.NewWriter(&b)
	wr, err := mp.CreateFormField(fieldName)
	if err != nil {
		t.Fatal(err)
	}

	wr.Write([]byte(token))
	mp.Close()

	r, err = http.NewRequest("POST", "/", &b)
	if err != nil {
		t.Fatal(err)
	}

	// Add the multipart header.
	r.Header.Set("Content-Type", mp.FormDataContentType())

	// Send back the issued cookie.
	setCookie(rr, r)

	rr = httptest.NewRecorder()
	m.ServeHTTP(rr, r)

	if rr.Code != http.StatusOK {
		t.Fatalf("middleware failed to pass to the next handler: got %v want %v",
			rr.Code, http.StatusOK)
	}

	if body := rr.Body.String(); !strings.Contains(body, token) {
		t.Fatalf("token not in response body: got %v want %v", body, token)
	}
}

// TestMaskUnmaskTokens tests that a token traversing the mask -> unmask process
// is correctly unmasked to the original 'real' token.
func TestMaskUnmaskTokens(t *testing.T) {
	t.Parallel()

	realToken, err := generateRandomBytes(tokenLength)
	if err != nil {
		t.Fatal(err)
	}

	issued := mask(realToken, nil)
	decoded, err := base64.StdEncoding.DecodeString(issued)
	if err != nil {
		t.Fatal(err)
	}

	unmasked := unmask(decoded)
	if !compareTokens(unmasked, realToken) {
		t.Fatalf("tokens do not match: got %x want %x", unmasked, realToken)
	}
}

// Tests domains that should (or should not) return true for a
// same-origin check.
func TestSameOrigin(t *testing.T) {
	var originTests = []struct {
		o1       string
		o2       string
		expected bool
	}{
		{"https://goji.io/", "https://goji.io", true},
		{"http://golang.org/", "http://golang.org/pkg/net/http", true},
		{"https://goji.io/", "http://goji.io", false},
		{"https://goji.io:3333/", "http://goji.io:4444", false},
	}

	for _, origins := range originTests {
		a, err := url.Parse(origins.o1)
		if err != nil {
			t.Fatal(err)
		}

		b, err := url.Parse(origins.o2)
		if err != nil {
			t.Fatal(err)
		}

		if sameOrigin(a, b) != origins.expected {
			t.Fatalf("origin checking failed: %v and %v, expected %v",
				origins.o1, origins.o2, origins.expected)
		}
	}
}

func TestXOR(t *testing.T) {
	testTokens := []struct {
		a        []byte
		b        []byte
		expected []byte
	}{
		{[]byte("goodbye"), []byte("hello"), []byte{15, 10, 3, 8, 13}},
		{[]byte("gophers"), []byte("clojure"), []byte{4, 3, 31, 2, 16, 0, 22}},
		{nil, []byte("requestToken"), nil},
	}

	for _, token := range testTokens {
		if res := xorToken(token.a, token.b); res != nil {
			if bytes.Compare(res, token.expected) != 0 {
				t.Fatalf("xorBytes failed to return the expected result: got %v want %v",
					res, token.expected)
			}
		}
	}

}

// shortReader provides a broken implementation of io.Reader for testing.
type shortReader struct{}

func (sr shortReader) Read(p []byte) (int, error) {
	return len(p) % 2, io.ErrUnexpectedEOF
}

// TestGenerateRandomBytes tests the (extremely rare) case that crypto/rand does
// not return the expected number of bytes.
func TestGenerateRandomBytes(t *testing.T) {
	// Pioneered from https://github.com/justinas/nosurf
	original := rand.Reader
	rand.Reader = shortReader{}
	defer func() {
		rand.Reader = original
	}()

	b, err := generateRandomBytes(tokenLength)
	if err == nil {
		t.Fatalf("generateRandomBytes did not report a short read: only read %d bytes", len(b))
	}
}

func TestTemplateField(t *testing.T) {
	m := goji.NewMux()
	CSRF := Protect(testKey, FieldName(testFieldName))
	m.UseC(CSRF)

	var token string
	var customTemplateField string
	m.HandleFuncC(pat.New("/"), func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
		token = Token(ctx, r)
		customTemplateField = string(TemplateField(ctx, r))
	})

	r, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	m.ServeHTTP(rr, r)

	expectedTemplateField := fmt.Sprintf(testTemplateField, testFieldName, token)

	if customTemplateField != expectedTemplateField {
		t.Fatalf("templateField not set correctly: got %v want %v",
			customTemplateField, expectedTemplateField)
	}
}
