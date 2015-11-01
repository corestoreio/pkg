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
	"errors"
	"net/http"
	"net/url"
	"strings"

	"github.com/corestoreio/csfw/config"
	"github.com/corestoreio/csfw/utils/log"
	"github.com/juju/errgo"
	"golang.org/x/net/context"
)

// ErrBaseUrlDoNotMatch will be returned if the request URL does not match the configured URL.
var ErrBaseUrlDoNotMatch = errors.New("The Base URLs do not match")

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
		if log.IsDebug() {
			log.Debug("net.httputils.IsSecure.FromContextReader", "err", config.ErrContextTypeAssertReaderFailed, "ctx", ctx)
		}
		cr = config.DefaultManager
	}
	oh, err := cr.GetString(config.Path(PathOffloaderHeader), config.ScopeDefault())
	if err != nil {
		if log.IsDebug() {
			log.Debug("net.httputils.IsSecure.FromContextReader.GetString", "err", err, "path", PathOffloaderHeader, "ctx", ctx)
		}
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

// IsBaseUrlCorrect checks if the requested host, scheme and path are same as the servers and
// if the path of the baseURL is included in the request URI.
func IsBaseUrlCorrect(r *http.Request, baseURL string) error {
	uri, err := url.Parse(baseURL)
	if err != nil {
		return errgo.Mask(err)
	}

	if r.Host == uri.Host && r.URL.Host == uri.Host && r.URL.Scheme == uri.Scheme && strings.Contains(r.URL.RequestURI(), uri.Path) {
		return nil
	}
	if log.IsDebug() {
		log.Debug("store.isBaseUrlCorrect.compare", "err", ErrBaseUrlDoNotMatch, "r.Host", r.Host, "baseURL", uri.String(), "requestURL", r.URL.String(), "strings.Contains", []string{r.URL.RequestURI(), uri.Path})
	}
	return errgo.Mask(ErrBaseUrlDoNotMatch)
}
