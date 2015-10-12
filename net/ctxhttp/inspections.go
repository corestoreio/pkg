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

package ctxhttp

import (
	"net/http"

	"github.com/corestoreio/csfw/config"
	"golang.org/x/net/context"
)

// PathOffloaderHeader defines the header name when a proxy server forwards an already
// terminated TLS request.
const PathOffloaderHeader = "web/secure/offloader_header"

// IsSecure checks if a request has been sent over a TLS connection. Also checks
// if the app runs behind a proxy server and therefore checks the off loader header.
func IsSecure(ctx context.Context, r *http.Request) bool {
	if r.TLS != nil {
		return true
	}

	oh := config.ContextMustReader(ctx).GetString(config.Path(PathOffloaderHeader), config.ScopeDefault())

	h := r.Header.Get(oh)
	hh := r.Header.Get("HTTP_" + oh)

	var isHttps bool
	switch "https" {
	case h, hh:
		isHttps = true
	}
	return isHttps
}

// IsSafeMethod checks if the request method is one of "GET", "HEAD", "TRACE", "OPTIONS"
// which can be considered as "safe".
func IsSafeMethod(r *http.Request) bool {
	// TODD(cs): figure out the usage for that function ...
	switch r.Method {
	case HTTPMethodGet, HTTPMethodHead, HTTPMethodTrace, HTTPMethodOptions:
		return true
	}
	return false
}

// IsAjax checks if the request has been initiated by a XMLHttpRequest
// or a form parameter of ajax or isAjax has been submitted.
func IsAjax(r *http.Request) bool {
	return r.Header.Get("X_REQUESTED_WITH") == "XMLHttpRequest" || r.FormValue("ajax") != "" || r.FormValue("isAjax") != ""
}
