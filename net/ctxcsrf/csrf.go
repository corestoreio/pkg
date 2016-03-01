// Package csrf (goji/csrf) provides Cross Site Request Forgery
// protection middleware for the Goji microframework (https://goji.io).
package csrf

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"

	"golang.org/x/net/context"

	"goji.io"

	"github.com/gorilla/securecookie"
)

// CSRF token length in bytes.
const tokenLength = 32

// Context/session keys & prefixes
const (
	tokenKey    string = "goji.csrf.Token"
	formKey     string = "goji.csrf.Form"
	errorKey    string = "goji.csrf.Error"
	cookieName  string = "_goji_csrf"
	errorPrefix string = "goji/csrf: "
)

var (
	// The name value used in form fields.
	fieldName = tokenKey
	// The default HTTP request header to inspect
	headerName = "X-CSRF-Token"
	// Idempotent (safe) methods as defined by RFC7231 section 4.2.2.
	safeMethods = []string{"GET", "HEAD", "OPTIONS", "TRACE"}
)

// TemplateTag provides a default template tag - e.g. {{ .csrfField }} - for use
// with the TemplateField function.
var TemplateTag = "csrfField"

var (
	// ErrNoReferer is returned when a HTTPS request provides an empty Referer
	// header.
	ErrNoReferer = errors.New("referer not supplied")
	// ErrBadReferer is returned when the scheme & host in the URL do not match
	// the supplied Referer header.
	ErrBadReferer = errors.New("referer invalid")
	// ErrNoToken is returned if no CSRF token is supplied in the request.
	ErrNoToken = errors.New("CSRF token not found in request")
	// ErrBadToken is returned if the CSRF token in the request does not match
	// the token in the session, or is otherwise malformed.
	ErrBadToken = errors.New("CSRF token invalid")
)

type csrf struct {
	h    goji.Handler
	sc   *securecookie.SecureCookie
	st   store
	opts options
}

// options contains the optional settings for the CSRF middleware.
type options struct {
	MaxAge int
	Domain string
	Path   string
	// Note that the function and field names match the case of the associated
	// http.Cookie field instead of the "correct" HTTPOnly name that golint suggests.
	HttpOnly      bool
	Secure        bool
	RequestHeader string
	FieldName     string
	ErrorHandler  goji.Handler
	CookieName    string
}

// Protect is HTTP middleware that provides Cross-Site Request Forgery
// protection.
//
// It securely generates a masked (unique-per-request) token that
// can be embedded in the HTTP response (e.g. form field or HTTP header).
// The original (unmasked) token is stored in the session, which is inaccessible
// by an attacker (provided you are using HTTPS). Subsequent requests are
// expected to include this token, which is compared against the session token.
// Requests that do not provide a matching token are served with a HTTP 403
// 'Forbidden' error response.
//
// Example:
//	package main
//
//	import (
//	    "html/template"
//	    "net/http"
//
//	    "goji.io"
//	    "github.com/goji/ctx-csrf"
//	    "github.com/zenazn/goji/graceful"
//	)
//
//	func main() {
//	    m := goji.NewMux()
//	    // Add the middleware to your router.
//	    m.UseC(csrf.Protect([]byte("32-byte-long-auth-key")))
//	    m.HandleFuncC(pat.Get("/signup"), ShowSignupForm)
//	    // POST requests without a valid token will return a HTTP 403 Forbidden.
//	    m.HandleFuncC(pat.Post("/signup/post"), SubmitSignupForm)
//
//	    graceful.ListenAndServe(":8000", m)
//	}
//
//	func ShowSignupForm(ctx context.Context, w http.ResponseWriter, r *http.Request) {
//	    // signup_form.tmpl just needs a {{ .csrfField }} template tag for
//	    // csrf.TemplateField to inject the CSRF token into. Easy!
//	    t.ExecuteTemplate(w, "signup_form.tmpl", map[string]interface{
//	        csrf.TemplateTag: csrf.TemplateField(ctx, r),
//	    })
//	    // We could also retrieve the token directly from csrf.Token(c, r) and
//	    // set it in the request header - w.Header.Set("X-CSRF-Token", token)
//	    // This is useful if your sending JSON to clients or a front-end JavaScript
//	    // framework.
//	}
//
//	func SubmitSignupForm(ctx context.Context, w http.ResponseWriter, r *http.Request) {
//	    // We can trust that requests making it this far have satisfied
//	    // our CSRF protection requirements.
//	}
//
func Protect(authKey []byte, opts ...Option) func(goji.Handler) goji.Handler {
	return func(h goji.Handler) goji.Handler {
		cs := parseOptions(h, opts...)

		// Set the defaults if no options have been specified
		if cs.opts.ErrorHandler == nil {
			cs.opts.ErrorHandler = goji.HandlerFunc(unauthorizedHandler)
		}

		if cs.opts.MaxAge < 1 {
			// Default of 12 hours
			cs.opts.MaxAge = 3600 * 12
		}

		if cs.opts.FieldName == "" {
			cs.opts.FieldName = fieldName
		}

		if cs.opts.CookieName == "" {
			cs.opts.CookieName = cookieName
		}

		if cs.opts.RequestHeader == "" {
			cs.opts.RequestHeader = headerName
		}

		// Create an authenticated securecookie instance.
		if cs.sc == nil {
			cs.sc = securecookie.New(authKey, nil)
			// Use JSON serialization (faster than one-off gob encoding)
			cs.sc.SetSerializer(securecookie.JSONEncoder{})
			// Set the MaxAge of the underlying securecookie.
			cs.sc.MaxAge(cs.opts.MaxAge)
		}

		if cs.st == nil {
			// Default to the cookieStore
			cs.st = &cookieStore{
				name:     cs.opts.CookieName,
				maxAge:   cs.opts.MaxAge,
				secure:   cs.opts.Secure,
				httpOnly: cs.opts.HttpOnly,
				path:     cs.opts.Path,
				domain:   cs.opts.Domain,
				sc:       cs.sc,
			}
		}

		return *cs
	}
}

// Implements goji.Handler for the csrf type.
func (cs csrf) ServeHTTPC(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	// Retrieve the token from the session.
	// An error represents either a cookie that failed HMAC validation
	// or that doesn't exist.
	realToken, err := cs.st.Get(r)
	if err != nil || len(realToken) != tokenLength {
		// If there was an error retrieving the token, the token doesn't exist
		// yet, or it's the wrong length, generate a new token.
		// Note that the new token will (correctly) fail validation downstream
		// as it will no longer match the request token.
		realToken, err = generateRandomBytes(tokenLength)
		if err != nil {
			ctx = setEnvError(ctx, err)
			cs.opts.ErrorHandler.ServeHTTPC(ctx, w, r)
			return
		}

		// Save the new (real) token in the session store.
		err = cs.st.Save(realToken, w)
		if err != nil {
			ctx = setEnvError(ctx, err)
			cs.opts.ErrorHandler.ServeHTTPC(ctx, w, r)
			return
		}
	}

	// Save the masked token to the request context
	ctx = context.WithValue(ctx, tokenKey, mask(realToken, r))
	// Save the field name to the request context
	ctx = context.WithValue(ctx, formKey, cs.opts.FieldName)

	// HTTP methods not defined as idempotent ("safe") under RFC7231 require
	// inspection.
	if !contains(safeMethods, r.Method) {
		// Enforce an origin check for HTTPS connections. As per the Django CSRF
		// implementation (https://goo.gl/vKA7GE) the Referer header is almost
		// always present for same-domain HTTP requests.
		if r.URL.Scheme == "https" {
			// Fetch the Referer value. Call the error handler if it's empty or
			// otherwise fails to parse.
			referer, err := url.Parse(r.Referer())
			if err != nil || referer.String() == "" {
				ctx = setEnvError(ctx, ErrNoReferer)
				cs.opts.ErrorHandler.ServeHTTPC(ctx, w, r)
				return
			}

			if sameOrigin(r.URL, referer) == false {
				ctx = setEnvError(ctx, ErrBadReferer)
				cs.opts.ErrorHandler.ServeHTTPC(ctx, w, r)
				return
			}
		}

		// If the token returned from the session store is nil for non-idempotent
		// ("unsafe") methods, call the error handler.
		if realToken == nil {
			ctx = setEnvError(ctx, ErrNoToken)
			cs.opts.ErrorHandler.ServeHTTPC(ctx, w, r)
			return
		}

		// Retrieve the combined token (pad + masked) token and unmask it.
		requestToken := unmask(cs.requestToken(r))

		// Compare the request token against the real token
		if !compareTokens(requestToken, realToken) {
			ctx = setEnvError(ctx, ErrBadToken)
			cs.opts.ErrorHandler.ServeHTTPC(ctx, w, r)
			return
		}

	}

	// Set the Vary: Cookie header to protect clients from caching the response.
	w.Header().Add("Vary", "Cookie")

	// Call the wrapped handler/router on success
	cs.h.ServeHTTPC(ctx, w, r)
}

// unauthorizedhandler sets a HTTP 403 Forbidden status and writes the
// CSRF failure reason to the response.
func unauthorizedHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	http.Error(w, fmt.Sprintf("%s - %s",
		http.StatusText(http.StatusForbidden), FailureReason(ctx, r)),
		http.StatusForbidden)
	return
}
