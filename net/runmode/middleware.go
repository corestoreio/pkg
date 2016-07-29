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

package runmode

import (
	"net"
	"net/http"
	"time"

	"github.com/corestoreio/csfw/log"
	"github.com/corestoreio/csfw/net/mw"
	"github.com/corestoreio/csfw/store"
	"github.com/corestoreio/csfw/store/scope"
	"github.com/corestoreio/csfw/util/errors"
)

// WithValidateBaseURL is a middleware which checks if the request base URL is
// equal to the one defined in the configuration, if not i.e. redirect from
// http://example.com/store/ to http://www.example.com/store/ @see
// app/code/Magento/Store/App/FrontController/Plugin/RequestPreprocessor.php
// @todo refactor this whole function as BaseURL must be bound to a store type
//func WithValidateBaseURL(cg config.GetterPubSuber, l log.Logger) mw.Middleware {
//
//	// Having the GetBool command here, means you must restart the app to take
//	// changes in effect. @todo refactor and use pub/sub to automatically change
//	// the isRedirectToBase value.
//
//	// <todo check logic!>
//	cgDefaultScope := cg.NewScoped(0, 0)
//	configRedirectCode, err := backend.Backend.WebURLRedirectToBase.Get(cgDefaultScope) // remove dependency
//	if err != nil {
//		panic(err) // we can panic here because during app start up
//	}
//
//	redirectCode := http.StatusMovedPermanently
//	if configRedirectCode != redirectCode {
//		redirectCode = http.StatusFound
//	}
//	// </todo check logic>
//
//	return func(h http.Handler) http.Handler {
//		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//
//			if configRedirectCode > 0 && r.Method != "POST" {
//
//				requestedStore, err := store.FromContextRequestedStore(r.Context())
//				if err != nil {
//					if l.IsDebug() {
//						l.Debug("http.WithValidateBaseUrl.FromContextServiceReader", log.Err(err), log.Object("request", r))
//					}
//					serveError(h, w, r, errors.Wrap(err, "[runmode] Context"))
//					return
//				}
//
//				baseURL, err := requestedStore.BaseURL(config.URLTypeWeb, requestedStore.IsCurrentlySecure(r))
//				if err != nil {
//					if l.IsDebug() {
//						l.Debug("http.WithValidateBaseUrl.requestedStore.BaseURL", log.Err(err), log.Object("request", r))
//					}
//					serveError(h, w, r, errors.Wrap(err, "[runmode] BaseURL"))
//					return
//				}
//
//				if err := httputil.IsBaseURLCorrect(r, &baseURL); err != nil {
//					if l.IsDebug() {
//						l.Debug("store.WithValidateBaseUrl.IsBaseUrlCorrect.error", log.Err(err), log.Object("request", r), log.Stringer("baseURL", &baseURL))
//					}
//
//					baseURL.Path = r.URL.Path
//					baseURL.RawPath = r.URL.RawPath
//					baseURL.RawQuery = r.URL.RawQuery
//					baseURL.Fragment = r.URL.Fragment
//					http.Redirect(w, r, (&baseURL).String(), redirectCode)
//					return
//				}
//			}
//			h.ServeHTTP(w, r)
//		})
//	}
//}

// AppRunMode initialized for each request the runMode.
type AppRunMode struct {
	// AvailabilityChecker must be set or application panics.
	store.AvailabilityChecker
	// CodeToIDMapper must be set or application panics.
	store.CodeToIDMapper
	// ErrorHandler optional custom error handler. Defaults to sending a HTTP status
	// code 500.
	mw.ErrorHandler

	// Log can be nil, defaults to black hole
	Log log.Logger
	// RunMode optional custom run mode otherwise falls back to
	// scope.DefaultRunMode.
	scope.RunMode
	// CodeExtracter extracts the store code from an HTTP requests. Optional.
	// Defaults to type ExtractCode.
	CodeExtracter
	// CookieTemplate optional pre-configured cookie to set the store code.
	// Expiration time and value will get overwritten.
	CookieTemplate func(*http.Request) *http.Cookie
	// CookieExpiresSet defaults to one year expiration for the store code.
	CookieExpiresSet time.Time // time.Now().AddDate(1, 0, 0) // one year valid
	// CookieExpiresDelete defaults to minus ten years to delete the store code
	// cookie.
	CookieExpiresDelete time.Time // time.Now().AddDate(-10, 0, 0)
}

func (a AppRunMode) getCookie(r *http.Request) *http.Cookie {
	if a.CookieTemplate != nil {
		return a.CookieTemplate(r)
	}
	d, _, err := net.SplitHostPort(r.Host)
	if err != nil {
		d = r.Host // might be a bug ...
	}
	var isSecure bool
	if r.TLS != nil {
		isSecure = true
	}
	return &http.Cookie{
		Name:     FieldName,
		Path:     "/",
		Domain:   d,
		Secure:   isSecure,
		HttpOnly: true, // disable for JavaScript access
	}
}

func (a AppRunMode) getCookieExpiresSet() time.Time {
	if a.CookieExpiresSet.IsZero() {
		return time.Now().AddDate(1, 0, 0) // one year valid
	}
	return a.CookieExpiresSet
}

func (a AppRunMode) getCookieExpiresDelete() time.Time {
	if a.CookieExpiresDelete.IsZero() {
		return time.Now().AddDate(-10, 0, 0) // -10 years
	}
	return a.CookieExpiresDelete
}

// WithRunMode sets for each request the overall runMode aka. scope. The following steps
// will be performed:
//	1. Call to AppRunMode.RunMode.CalculateMode to get the default run mode.
//	2a. Parse Request GET parameter for the store code key (___store).
//	2b. If GET is empty, check cookie for key "store"
//	2c. Lookup CodeToIDMapper.IDbyCode() to get the store ID from a store code.
func (a AppRunMode) WithRunMode() mw.Middleware {

	// todo check if store is not active anymore, and if inactive call error handler
	// todo: code/Magento/Store/App/Request/PathInfoProcessor.php => req.isDirectAccessFrontendName()

	aLog := a.Log
	if aLog == nil {
		aLog = log.BlackHole{} // disabled debug and info logging
	}
	aErrH := a.ErrorHandler
	if aErrH == nil {
		aErrH = mw.ErrorWithStatusCode(http.StatusInternalServerError)
	}
	aGetCode := a.CodeExtracter
	if aGetCode == nil {
		aGetCode = ExtractCode{}
	}
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			// set run mode
			var runMode scope.Hash
			runMode = a.CalculateMode(w, r)
			runID := runMode.ID()

			// extracts the code from GET and/or Cookie or custom implementation and get the
			// new runID.
			reqCode, reqCodeValid := a.CodeExtracter.FromRequest(r)
			if reqCodeValid > 0 {
				var err error
				// convert the code string into its internal ID depending on the scope.
				runID, err = a.IDbyCode(runMode.Scope(), reqCode)
				if err != nil && !errors.IsNotFound(err) {
					if aLog.IsDebug() {
						aLog.Debug("runmode.WithRunMode.IDbyCode.Error", log.Err(err), log.String("http_store_code", reqCode),
							log.Int64("run_id", runID), log.Stringer("run_mode", runMode), log.HTTPRequest("request", r))
					}
					aErrH(errors.Wrap(err, "[store] WithRunMode.IDbyCode")).ServeHTTP(w, r)
					return
				}
				if aLog.IsDebug() {
					aLog.Debug("runmode.WithRunMode.CodeFromRequest", log.String("http_store_code", reqCode),
						log.Int64("run_id", runID), log.Stringer("run_mode", runMode), log.HTTPRequest("request", r))
				}
			} // ignore everything else

			// which store IDs are allowed depending on our runMode? Check if the runID is
			// within the allowed store IDs.
			var isStoreAllowed bool
			allowedStoreIDs, err := a.AllowedStoreIDs(runMode)
			if err != nil {
				if aLog.IsDebug() {
					aLog.Debug("runmode.WithRunMode.AllowedStoreIDs.Error", log.Err(err),
						log.Int64("run_id", runID), log.Stringer("run_mode", runMode), log.HTTPRequest("request", r))
				}
				aErrH(errors.Wrap(err, "[store] WithRunMode.AllowedStoreIDs")).ServeHTTP(w, r)
				return
			}
			for _, s := range allowedStoreIDs {
				if s == runID {
					isStoreAllowed = true
				}
			}

			// if the runID is zero, it can be the admin scope OR if our store ID
			// cannot be found in the list of allowed store IDs we must look up the
			// DefaultStoreID for the runMode.
			if !isStoreAllowed {
				if aLog.IsDebug() {
					aLog.Debug("runmode.WithRunMode.StoreNotAllowed",
						log.Bool("is_store_allowed", isStoreAllowed),
						log.Int64s("allowed_store_IDs", allowedStoreIDs...),
						log.Int64("run_id", runID), log.Stringer("run_mode", runMode), log.HTTPRequest("request", r))
				}

				// not found, not active, whatever, we cannot proceed.
				if runID > 0 {
					aErrH(errors.NewUnauthorizedf("[store] Scope %q with requested ID %d cannot be authorized", runMode.Scope(), runID)).ServeHTTP(w, r)
					return
				}

				// todo: admin scope which is website,group,store are 0
				var err error
				runID, err = a.DefaultStoreID(runMode)
				if err != nil {
					if aLog.IsDebug() {
						aLog.Debug("runmode.WithRunMode.DefaultStoreID.Error", log.Err(err),
							log.Bool("is_store_allowed", isStoreAllowed),
							log.Int64("run_id", runID), log.Stringer("run_mode", runMode), log.HTTPRequest("request", r))
					}
					aErrH(errors.Wrap(err, "[store] WithRunMode.DefaultStoreID")).ServeHTTP(w, r)
					return
				}
				if aLog.IsDebug() {
					aLog.Debug("runmode.WithRunMode.DefaultStoreID",
						log.Bool("is_store_allowed", isStoreAllowed),
						log.Int64("run_id", runID), log.Stringer("run_mode", runMode), log.HTTPRequest("request", r))
				}
				// runID can still be zero
			}

			// if store code found in cookie and not valid anymore, delete the cookie.
			if reqCodeValid == 20 && !isStoreAllowed {
				keks := a.getCookie(r)
				keks.Expires = a.getCookieExpiresDelete()
				http.SetCookie(w, keks)
			}

			// if runMode is allowed to change, update the runMode Hash and then put it into the context
			if isStoreAllowed && runID != runMode.ID() {
				runMode = scope.NewHash(runMode.Scope(), runID)

				if reqCodeValid < 20 { // no cookie found but the code changed
					// set cookie once with the new code
					keks := a.getCookie(r)
					keks.Expires = a.getCookieExpiresSet()
					http.SetCookie(w, keks)
				}
			}

			next.ServeHTTP(w, r.WithContext(scope.WithContextRunMode(r.Context(), runMode)))
		})
	}
}
