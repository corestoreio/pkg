# goji/ctx-csrf
[![GoDoc](https://godoc.org/github.com/goji/ctx-csrf?status.svg)](https://godoc.org/github.com/goji/ctx-csrf) [![Build Status](https://travis-ci.org/goji/ctx-csrf.svg?branch=master)](https://travis-ci.org/goji/ctx-csrf)

**ctx-csrf** is a HTTP middleware library that provides [cross-site request
forgery](http://blog.codinghorror.com/preventing-csrf-and-xsrf-attacks/) (CSRF)
protection with support for Go's `net/context` package. It includes:

* The `csrf.Protect` middleware/handler that can be used with `goji.Use` to
  provide CSRF protection on routes attached to a router or a sub-router.
* A `csrf.Token` function that provides the token to pass into your response,
  whether that be a HTML form or a JSON response body.
* ... and a `csrf.TemplateField` helper that you can pass into your `html/template`
  templates to replace a `{{ .csrfField }}` template tag with a hidden input
  field.

This library is designed to work with not just the the
[Goji](https://github.com/goji/goji) micro-framework, but any project that satisfies the
[goji.Handler](https://godoc.org/goji.io#Handler) interface: `ServeHTTPC(context.Context,
http.ResponseWriter, *http.Request)`.

This makes it compatible with other parts of the Go ecosystem. The
`context.Context` request context doesn't rely on a global map, and is therefore
free from contention in a busy web service.

The library also assumes HTTPS by default: sending cookies over vanilla HTTP is
risky and you're likely to get hurt.

*Note*: If you're using Goji v1, the older
[goji/csrf](https://github.com/goji/csrf) still exists.

## Examples

ctx-csrf is easy to use: add the middleware to your stack with the below:

```go
goji.UseC(csrf.Protect([]byte("32-byte-long-auth-key")))
```

... and then collect the token with `csrf.Token(c, r)` before passing it to the
template, JSON body or HTTP header (you pick!). ctx-csrf inspects HTTP headers
(first) and the form body (second) on subsequent POST/PUT/PATCH/DELETE/etc.
requests for the token.

### HTML Forms

Here's the common use-case: HTML forms you want to provide CSRF protection for,
in order to protect malicious POST requests being made:

```go
package main

import (
    "html/template"
    "net/http"

    "goji.io"
    "github.com/goji/ctx-csrf"
    "github.com/zenazn/goji/graceful"
)

func main() {
    m := goji.NewMux()
    // Add the middleware to your router.
    m.UseC(csrf.Protect([]byte("32-byte-long-auth-key")))
    m.HandleFuncC(pat.Get("/signup"), ShowSignupForm)
    // POST requests without a valid token will return a HTTP 403 Forbidden.
    m.HandleFuncC(pat.Post("/signup/post"), SubmitSignupForm)

    graceful.ListenAndServe(":8000", m)
}

func ShowSignupForm(ctx context.Context, w http.ResponseWriter, r *http.Request) {
    // signup_form.tmpl just needs a {{ .csrfField }} template tag for
    // csrf.TemplateField to inject the CSRF token into. Easy!
    t.ExecuteTemplate(w, "signup_form.tmpl", map[string]interface{
        csrf.TemplateTag: csrf.TemplateField(ctx, r),
    })
    // We could also retrieve the token directly from csrf.Token(c, r) and
    // set it in the request header - w.Header.Set("X-CSRF-Token", token)
    // This is useful if your sending JSON to clients or a front-end JavaScript
    // framework.
}

func SubmitSignupForm(ctx context.Context, w http.ResponseWriter, r *http.Request) {
    // We can trust that requests making it this far have satisfied
    // our CSRF protection requirements.
}
```

### JSON Responses

This approach is useful if you're using a front-end JavaScript framework like
Ember or Angular, or are providing a JSON API.

We'll also look at applying selective CSRF protection using Goji's sub-routers,
as we don't handle any POST/PUT/DELETE requests with our top-level router.

```go
package main

import (
    "goji.io"
    "github.com/goji/ctx-csrf"
    "github.com/zenazn/goji/graceful"
)

func main() {
    m := goji.NewMux()
    // Our top-level router doesn't need CSRF protection: it's simple.
    m.HandleFuncC(pat.Get("/"), ShowIndex)

    api := goji.NewMux()
    m.HandleC("/api/*", api)
    // ... but our /api/* routes do, so we add it to the sub-router only.
    api.UseC(csrf.Protect([]byte("32-byte-long-auth-key")))

    api.Get("/api/user/:id", GetUser)
    api.Post("/api/user", PostUser)

    graceful.ListenAndServe(":8000", m)
}

func GetUser(ctx context.Context, w http.ResponseWriter, r *http.Request) {
    // Authenticate the request, get the :id from the route params,
    // and fetch the user from the DB, etc.

    // Get the token and pass it in the CSRF header. Our JSON-speaking client
    // or JavaScript framework can now read the header and return the token in
    // in its own "X-CSRF-Token" request header on the subsequent POST.
    w.Header().Set("X-CSRF-Token", csrf.Token(ctx, r))
    b, err := json.Marshal(user)
    if err != nil {
        http.Error(w, http.StatusText(500), 500)
        return
    }

    w.Write(b)
}
```

### Setting Options

What about providing your own error handler and changing the HTTP header the
package inspects on requests? (i.e. an existing API you're porting to Go). Well,
ctx-csrf provides options for changing these as you see fit:

```go
func main() {
    m := goji.NewMux()
    CSRF := csrf.Protect(
            []byte("a-32-byte-long-key-goes-here"),
            csrf.RequestHeader("Authenticity-Token"),
            csrf.FieldName("authenticity_token"),
            // Note that csrf.ErrorHandler takes a Goji goji.Handler type, else
            // your error handler can't retrieve the error reason from the
            // context.
            // The signature `func UnauthHandler(ctx context.Context, w http.ResponseWriter, r *http.Request)`
            // is a goji.Handler, and the simplest to use if you'd like to serve
            // "pretty" error pages (who doesn't?).
            csrf.ErrorHandler(goji.HandlerFunc(serverError(403))),
        )

    m.UseC(CSRF)
    m.HandleFuncC(pat.Get("/signup"), GetSignupForm)
    m.HandleFuncC(pat.Post("/signup"), PostSignupForm)

    graceful.ListenAndServe(":8000", m)
}
```

Not too bad, right?

If there's something you're confused about or a feature you would like to see
added, open an issue with your code so far.

## Design Notes

Getting CSRF protection right is important, so here's some background:

* This library generates unique-per-request (masked) tokens as a mitigation
  against the [BREACH attack](http://breachattack.com/).
* The 'base' (unmasked) token is stored in the session, which means that
  multiple browser tabs won't cause a user problems as their per-request token
  is compared with the base token.
* Operates on a "whitelist only" approach where safe (non-mutating) HTTP methods
  (GET, HEAD, OPTIONS, TRACE) are the *only* methods where token validation is not
  enforced.
* The design is based on the battle-tested
  [Django](https://docs.djangoproject.com/en/1.8/ref/csrf/) and [Ruby on
  Rails](http://api.rubyonrails.org/classes/ActionController/RequestForgeryProtection.html)
  approaches.
* Cookies are authenticated and based on the [securecookie](https://github.com/gorilla/securecookie)
  library. They're also Secure (issued over HTTPS only) and are HttpOnly
  by default, because sane defaults are important.
* Go's `crypto/rand` library is used to generate the 32 byte (256 bit) tokens
  and the one-time-pad used for masking them.

This library does not seek to be adventurous.

## License

BSD licensed. See the LICENSE file for details.
