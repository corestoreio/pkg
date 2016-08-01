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
	"net/http"

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

// AppRunMode initialized for each request the main important runMode.
type AppRunMode struct {
	// AvailabilityChecker must be set or middleware panics.
	store.AvailabilityChecker
	// CodeToIDMapper must be set or middleware panics.
	store.CodeToIDMapper

	// Following fields are optional

	// ErrorHandler optional custom error handler. Defaults to sending an HTTP
	// status code 500 and exposing the real error including full paths.
	mw.ErrorHandler
	// Log can be nil, defaults to black hole.
	Log log.Logger
	// RunMode optional custom runMode otherwise falls back to
	// scope.DefaultRunMode which selects the default website with its default
	// store. To use the admin area enable scope.Store and ID 0.
	scope.RunMode
	// StoreCodeProcessor extracts the store code from an HTTP requests.
	// Optional. Defaults to type ProcessStoreCode.
	StoreCodeProcessor
	// DisableStoreCodeProcessor set to true and set StoreCodeProcessor to nil
	// to disable store code handling
	DisableStoreCodeProcessor bool
}

func (a AppRunMode) checkStoreIDAllowed(runMode scope.Hash, newStoreID int64) (allowedStoreIDs []int64, isStoreAllowed bool, err error) {

	// which store IDs are allowed depending on our runMode? Check if the newStoreID is
	// within the allowed store IDs.
	allowedStoreIDs, err = a.AllowedStoreIDs(runMode)
	if err != nil {
		return nil, false, errors.Wrap(err, "[store] WithRunMode.checkStoreIDAllowed")
	}

	for _, s := range allowedStoreIDs {
		if s == newStoreID {
			isStoreAllowed = true
			return
		}
	}
	return
}

// WithRunMode sets for each request the overall runMode aka. scope. The following steps
// will be performed:
//	1. Call to AppRunMode.RunMode.CalculateMode to get the default run mode.
//	2a. Parse Request GET parameter for the store code key (___store).
//	2b. If GET is empty, check cookie for key "store"
//	2c. Lookup CodeToIDMapper.IDbyCode() to get the website/store ID from a website/store code.
//	3. Retrieve all AllowedStoreIDs based on the runMode
//	4. Check if the website/store ID
func (a AppRunMode) WithRunMode() mw.Middleware {

	// todo: code/Magento/Store/App/Request/PathInfoProcessor.php => req.isDirectAccessFrontendName()

	aLog := a.Log
	if aLog == nil {
		aLog = log.BlackHole{} // disabled debug and info logging
	}
	aErrH := a.ErrorHandler
	if aErrH == nil {
		aErrH = mw.ErrorWithStatusCode(http.StatusInternalServerError)
	}
	aGetCode := a.StoreCodeProcessor
	if aGetCode == nil {
		aGetCode = nullCodeProcessor{}
		if !a.DisableStoreCodeProcessor {
			aGetCode = &ProcessStoreCode{}
		}
	}
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			// set run mode
			runMode := a.CalculateMode(w, r)

			// find the default store ID for the runMode
			newStoreID, err := a.DefaultStoreID(runMode)
			if err != nil {
				if aLog.IsDebug() {
					aLog.Debug("runmode.WithRunMode.DefaultStoreID.Error", log.Err(err),
						log.Int64("store_id", newStoreID), log.Stringer("run_mode", runMode), log.HTTPRequest("request", r))
				}
				aErrH(errors.Wrap(err, "[store] WithRunMode.DefaultStoreID")).ServeHTTP(w, r)
				return
			}
			if aLog.IsDebug() {
				aLog.Debug("runmode.WithRunMode.DefaultStoreID",
					log.Int64("store_id", newStoreID), log.Stringer("run_mode", runMode), log.HTTPRequest("request", r))
			}

			// extracts the code from GET and/or Cookie or custom implementation and get the
			// new runID.
			var reqStoreCode string
			reqStoreCode = a.StoreCodeProcessor.FromRequest(r)
			if reqStoreCode != "" {
				var err error
				// convert the code string into its internal ID depending on the scope.
				newStoreID, err = a.StoreIDbyCode(runMode, reqStoreCode)
				if err != nil && !errors.IsNotFound(err) {
					if aLog.IsDebug() {
						aLog.Debug("runmode.WithRunMode.IDbyCode.Error", log.Err(err),
							log.String("http_store_code", reqStoreCode), log.Int64("store_id", newStoreID),
							log.Stringer("run_mode", runMode), log.HTTPRequest("request", r))
					}
					aErrH(errors.Wrap(err, "[store] WithRunMode.IDbyCode")).ServeHTTP(w, r)
					return
				}
				if aLog.IsDebug() {
					aLog.Debug("runmode.WithRunMode.CodeFromRequest", log.String("http_store_code", reqStoreCode),
						log.Int64("store_id", newStoreID), log.Stringer("run_mode", runMode), log.HTTPRequest("request", r))
				}
			} // ignore everything else

			// which store IDs are allowed depending on our runMode? Check if the newStoreID is
			// within the allowed store IDs.
			allowedStoreIDs, isStoreAllowed, err := a.checkStoreIDAllowed(runMode, newStoreID)
			if err != nil {
				if aLog.IsDebug() {
					aLog.Debug("runmode.WithRunMode.AllowedStoreIDs.Error", log.Err(err),
						log.Int64("store_id", newStoreID), log.Stringer("run_mode", runMode), log.HTTPRequest("request", r))
				}
				aErrH(errors.Wrap(err, "[store] WithRunMode.AllowedStoreIDs")).ServeHTTP(w, r)
				return
			}

			// not found, not active, whatever, we cannot proceed.
			if !isStoreAllowed {
				if aLog.IsDebug() {
					aLog.Debug("runmode.WithRunMode.StoreNotAllowed",
						log.Bool("is_store_allowed", isStoreAllowed), log.Int64s("allowed_store_IDs", allowedStoreIDs...),
						log.Int64("store_id", newStoreID), log.Stringer("run_mode", runMode), log.HTTPRequest("request", r))
				}
				a.StoreCodeProcessor.ProcessDenied(runMode, newStoreID, w, r)
				aErrH(errors.NewUnauthorizedf("[store] RunMode %s with requested Store ID %d cannot be authorized", runMode, newStoreID)).ServeHTTP(w, r)
				return
			}

			// if runMode is allowed to change, update the runMode Hash and then put it into the context
			a.StoreCodeProcessor.ProcessAllowed(runMode, newStoreID, w, r)
			previousRunMode := runMode
			if isStoreAllowed && newStoreID != runMode.ID() {
				runMode = scope.NewHash(runMode.Scope(), newStoreID)
			}
			if aLog.IsDebug() {
				aLog.Debug("runmode.WithRunMode.NextHandler",
					log.Bool("is_store_allowed", isStoreAllowed), log.Int64s("allowed_store_IDs", allowedStoreIDs...),
					log.Int64("store_id", newStoreID), log.Stringer("run_mode", runMode),
					log.Stringer("previous_run_mode", previousRunMode), log.HTTPRequest("request", r))
			}
			next.ServeHTTP(w, r.WithContext(scope.WithContextRunMode(r.Context(), runMode)))
		})
	}
}
