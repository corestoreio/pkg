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

package httputils

import (
	"net/http"

	"github.com/corestoreio/csfw/config"
	"github.com/corestoreio/csfw/utils/log"
	"golang.org/x/net/context"
)

// PathOffloaderHeader defines the header name when a proxy server forwards an already
// terminated TLS request.
const PathOffloaderHeader = "web/secure/offloader_header"

// IsSecure checks if a request has been sent over a TLS connection. Also checks
// if the app runs behind a proxy server and therefore checks the off loader header.
func IsSecure(ctx context.Context, r *http.Request) bool {
	// due to import cycle this function must be in this package
	if r.TLS != nil {
		return true
	}

	cr, ok := config.FromContextReader(ctx)
	if !ok {
		log.Error("net.httputils.IsSecure.FromContextReader", "err", config.ErrContextTypeAssertReaderFailed, "ctx", ctx)
		cr = config.DefaultManager
	}
	oh, err := cr.GetString(config.Path(PathOffloaderHeader), config.ScopeDefault())
	if err != nil {
		log.Error("net.httputils.IsSecure.FromContextReader.GetString", "err", err, "path", PathOffloaderHeader, "ctx", ctx)
		return false
	}

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
	case MethodGet, MethodHead, MethodTrace, MethodOptions:
		return true
	}
	return false
}

// IsAjax checks if the request has been initiated by a XMLHttpRequest
// or a form parameter of ajax or isAjax has been submitted.
func IsAjax(r *http.Request) bool {
	return r.Header.Get("X_REQUESTED_WITH") == "XMLHttpRequest" || r.FormValue("ajax") != "" || r.FormValue("isAjax") != ""
}
