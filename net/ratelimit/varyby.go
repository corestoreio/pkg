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

package ratelimit

import (
	"net/http"
	"strings"

	"github.com/corestoreio/pkg/net/request"
	"github.com/corestoreio/pkg/util/bufferpool"
)

// VaryByer is called for each request to generate a key for the limiter. If it
// returns an empty string, all requests use an empty string key ;-). The rate
// limiter checks whether a particular key has exceeded a rate limit.
type VaryByer interface {
	Key(*http.Request) string
}

type emptyVaryBy struct{}

func (emptyVaryBy) Key(_ *http.Request) string {
	return ""
}

// VaryBy defines the criteria to use to group requests.
type VaryBy struct {
	// Vary by the RemoteAddr as specified by the net/http.Request field.
	RemoteAddr bool

	// Vary by the HTTP Method as specified by the net/http.Request field.
	Method bool

	// Vary by the URL's Path as specified by the Path field of the net/http.Request
	// URL field.
	Path bool

	// Vary by this list of header names, read from the net/http.Request Header field.
	Headers []string

	// Vary by this list of parameters, read from the net/http.Request FormValue method.
	Params []string

	// Vary by this list of cookie names, read from the net/http.Request Cookie method.
	Cookies []string

	// Use this separator string to concatenate the various criteria of the VaryBy struct.
	// Defaults to a newline character if empty (\n).
	Separator string

	// SafeUnicode enables the usage of unicode safe strings to lower functions.
	SafeUnicode bool
}

// Key returns the key for this request based on the criteria defined by the VaryBy struct.
func (vb *VaryBy) Key(r *http.Request) string {
	buf := bufferpool.Get()
	defer bufferpool.Put(buf)

	if vb == nil {
		return "" // Special case for no vary-by option
	}

	sep := vb.Separator
	if sep == "" {
		sep = "\n" // Separator defaults to newline
	}
	if vb.RemoteAddr {
		_, _ = buf.WriteString(request.RealIP(r, request.IPForwardedTrust).String())
		_, _ = buf.WriteString(sep)
	}
	if vb.Method {
		_, _ = buf.WriteString(toLower(r.Method, vb.SafeUnicode))
		_, _ = buf.WriteString(sep)
	}
	for _, h := range vb.Headers {
		_, _ = buf.WriteString(toLower(r.Header.Get(h), vb.SafeUnicode))
		_, _ = buf.WriteString(sep)
	}
	if vb.Path {
		_, _ = buf.WriteString(r.URL.Path)
		_, _ = buf.WriteString(sep)
	}
	for _, p := range vb.Params {
		_, _ = buf.WriteString(r.FormValue(p))
		_, _ = buf.WriteString(sep)
	}
	for _, c := range vb.Cookies {
		if ck, err := r.Cookie(c); err == nil {
			_, _ = buf.WriteString(ck.Value)
			_, _ = buf.WriteString(sep)
		}
	}
	return buf.String()
}

func toLower(s string, safeUnicode bool) string {
	if safeUnicode {
		return strings.ToLower(s)
	}
	b := make([]byte, len(s))
	for i := range b {
		c := s[i]
		if c >= 'A' && c <= 'Z' {
			c += 'a' - 'A'
		}
		b[i] = c
	}
	return string(b)
}
