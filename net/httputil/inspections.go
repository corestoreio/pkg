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

package httputil

import (
	"errors"
	"net/http"
	"net/url"
	"strings"

	"github.com/corestoreio/csfw/config"
	"github.com/juju/errgo"
	"golang.org/x/net/context"
)

// ErrBaseURLDoNotMatch will be returned if the request URL does not match the
// configured URL.
var ErrBaseURLDoNotMatch = errors.New("The Base URLs do not match")

// PathOffloaderHeader defines the header name when a proxy server forwards an
// already terminated TLS request.
const PathOffloaderHeader = "web/secure/offloader_header"

// CtxIsSecure same as IsSecure() but extract the config.Reader out of the context.
// Wrapper function.
func CtxIsSecure(ctx context.Context, r *http.Request) bool {
	return IsSecure(config.FromContextGetter(ctx), r)
}

// IsSecure checks if a request has been sent over a TLS connection. Also checks
// if the app runs behind a proxy server and therefore checks the off loader header.
func IsSecure(cr config.Getter, r *http.Request) bool {
	// due to import cycle this function must be in this package
	if r.TLS != nil {
		return true
	}

	oh, err := cr.String(config.Path(PathOffloaderHeader), config.ScopeDefault())
	if err != nil {
		if PkgLog.IsDebug() {
			PkgLog.Debug("net.httputil.IsSecure.FromContextReader.String", "err", err, "path", PathOffloaderHeader)
		}
		return false
	}

	h := r.Header.Get(oh)
	hh := r.Header.Get("HTTP_" + oh)

	var isHTTPS bool
	switch "https" {
	case h, hh:
		isHTTPS = true
	}
	return isHTTPS
}

// IsBaseURLCorrect checks if the requested host, scheme and path are same as the servers and
// if the path of the baseURL is included in the request URI.
func IsBaseURLCorrect(r *http.Request, baseURL *url.URL) error {
	if r.Host == baseURL.Host && r.URL.Host == baseURL.Host && r.URL.Scheme == baseURL.Scheme && strings.HasPrefix(r.URL.Path, baseURL.Path) {
		return nil
	}
	if PkgLog.IsDebug() {
		PkgLog.Debug("store.isBaseUrlCorrect.compare", "err", ErrBaseURLDoNotMatch, "r.Host", r.Host, "baseURL", baseURL.String(), "requestURL", r.URL.String(), "strings.Contains", []string{r.URL.RequestURI(), baseURL.Path})
	}
	return errgo.Mask(ErrBaseURLDoNotMatch)
}
