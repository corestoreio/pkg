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
	"net/http"
	"net/url"
	"strings"

	"github.com/corestoreio/csfw/util/errors"
)

// todo: refactor

//// CheckSecureRequest checks if a request is secure using the SSL offloader header
//type CheckSecureRequest struct {
//	// WebSecureOffloaderHeader => Offloader header.
//	// See package backend.
//	// Path: web/secure/offloader_header
//	WebSecureOffloaderHeader cfgmodel.Str
//	Log                      log.Logger
//}
//
//// NewCeckSecureRequest creates a new SecureRequest type pointer.
//// Requires the correct path to the WebSecureOffloaderHeader configuration.
//func NewCeckSecureRequest(cfgOffloader cfgmodel.Str) *CheckSecureRequest {
//	return &CheckSecureRequest{
//		WebSecureOffloaderHeader: cfgOffloader,
//		Log: log.BlackHole{}, // disabled info and debug logging
//	}
//}
//
//// CtxIs same as IsSecure() but extract the config.ScopedGetter out of the context.
//// Wrapper function.
//func (sr *CheckSecureRequest) CtxIs(ctx context.Context, r *http.Request) bool {
//	sg, ok := config.FromContextScopedGetter(ctx)
//	if !ok {
//		if sr.Log.IsDebug() {
//			sr.Log.Debug("net.httputil.CtxIsSecure.FromContextScopedGetter", log.Bool("ok", ok), log.HTTPRequest("request", r))
//		}
//	}
//	return sr.Is(sg, r)
//}
//
//// Is checks if a request has been sent over a TLS connection. Also checks
//// if the app runs behind a proxy server and therefore checks the off loader header.
//// config.ScopedGetter can be nil.
//func (sr *CheckSecureRequest) Is(sg config.ScopedGetter, r *http.Request) bool {
//
//	if r.TLS != nil {
//		return true
//	}
//
//	if sg == nil {
//		return false
//	}
//
//	oh, err := sr.WebSecureOffloaderHeader.Get(sg)
//	if err != nil {
//		if sr.Log.IsDebug() {
//			sr.Log.Debug("net.httputil.IsSecure.FromContextReader.String", log.Err(err), log.Stringer("path", sr.WebSecureOffloaderHeader.Route()))
//		}
//		return false
//	}
//
//	h := r.Header.Get(oh)
//	hh := r.Header.Get("HTTP_" + oh)
//
//	var isHTTPS bool
//	switch "https" {
//	case h, hh:
//		isHTTPS = true
//	}
//	return isHTTPS
//}

// IsBaseURLCorrect checks if the requested host, scheme and path are same as the servers and
// if the path of the baseURL is included in the request URI.
// Error behaviour: NotValid
func IsBaseURLCorrect(r *http.Request, baseURL *url.URL) error {
	if r.Host == baseURL.Host && r.URL.Host == baseURL.Host && r.URL.Scheme == baseURL.Scheme && strings.HasPrefix(r.URL.Path, baseURL.Path) {
		return nil
	}
	return errors.NewNotValidf("[httputil] Base URLs do not match. BaseURL %q RequestURL %q strings.Contains %v", baseURL.String(), r.URL.String(), []string{r.URL.RequestURI(), baseURL.Path})
}
