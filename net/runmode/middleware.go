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

// Options additional customizations for the runMode middleware.
type Options struct {
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
	// Optional. Defaults to type ProcessStoreCodeCookie.
	store.CodeProcessor
	// DisableStoreCodeProcessor set to true and set StoreCodeProcessor to nil
	// to disable store code handling
	DisableStoreCodeProcessor bool
}

// WithRunMode sets for each request the overall runMode aka. scope. The following steps
// will be performed:
//	1. Call to AppRunMode.RunMode.CalculateMode to get the default run mode.
//	2a. Parse Request GET parameter for the store code key (___store).
//	2b. If GET is empty, check cookie for key "store"
//	2c. Lookup CodeToIDMapper.IDbyCode() to get the website/store ID from a website/store code.
//	3. Retrieve all AllowedStoreIDs based on the runMode
//	4. Check if the website/store ID
func WithRunMode(sf store.Finder, o Options) mw.Middleware {

	// todo: code/Magento/Store/App/Request/PathInfoProcessor.php => req.isDirectAccessFrontendName()

	lg := o.Log
	if lg == nil {
		lg = log.BlackHole{} // disabled debug and info logging
	}
	errH := o.ErrorHandler
	if errH == nil {
		errH = mw.ErrorWithStatusCode(http.StatusInternalServerError)
	}
	procCode := o.CodeProcessor
	if procCode == nil {
		procCode = nullCodeProcessor{}
		if !o.DisableStoreCodeProcessor {
			procCode = &ProcessStoreCodeCookie{}
		}
	}
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			// set run mode
			runMode := o.CalculateMode(r)
			r = r.WithContext(scope.WithContextRunMode(r.Context(), runMode))

			// find the default store ID for the runMode
			storeID, websiteID, err := sf.DefaultStoreID(runMode)
			oldStoreID := storeID
			if err != nil {
				if lg.IsDebug() {
					lg.Debug("runmode.WithRunMode.DefaultStoreID.Error", log.Err(err),
						log.Int64("store_id", storeID), log.Int64("website_id", websiteID),
						log.Stringer("run_mode", runMode), log.HTTPRequest("request", r))
				}
				errH(errors.Wrap(err, "[store] WithRunMode.DefaultStoreID")).ServeHTTP(w, r)
				return
			}
			if lg.IsDebug() {
				lg.Debug("runmode.WithRunMode.DefaultStoreID", log.Int64("store_id", storeID), log.Int64("website_id", websiteID),
					log.Stringer("run_mode", runMode), log.HTTPRequest("request", r))
			}

			// extracts the code from GET and/or Cookie or custom implementation and get the
			// new store and website ID.
			if reqCode := procCode.FromRequest(runMode, r); reqCode != "" {
				// convert the code string into its internal ID depending on the scope.
				newStoreID, newWebsiteID, err := sf.StoreIDbyCode(runMode, reqCode)
				if err != nil && !errors.IsNotFound(err) {
					if lg.IsDebug() {
						lg.Debug("runmode.WithRunMode.IDbyCode.Error", log.Err(err), log.String("http_store_code", reqCode),
							log.Int64("store_id", storeID), log.Int64("website_id", websiteID),
							log.Stringer("run_mode", runMode), log.HTTPRequest("request", r))
					}
					errH(errors.Wrap(err, "[store] WithRunMode.IDbyCode")).ServeHTTP(w, r)
					return
				}
				if err == nil {
					storeID = newStoreID
					websiteID = newWebsiteID
				}
				if lg.IsDebug() {
					lg.Debug("runmode.WithRunMode.CodeFromRequest", log.Err(err), log.String("http_store_code", reqCode),
						log.Int64("store_id", storeID), log.Int64("website_id", websiteID),
						log.Stringer("run_mode", runMode), log.HTTPRequest("request", r))
				}
			}

			r = r.WithContext(scope.WithContext(r.Context(), websiteID, storeID))

			// which store IDs are allowed depending on our runMode? Check if the storeID is
			// within the allowed store IDs.
			isStoreAllowed, storeCode, err := sf.IsAllowedStoreID(runMode, storeID)
			if err != nil {
				if lg.IsDebug() {
					lg.Debug("runmode.WithRunMode.AllowedStoreIDs.Error", log.Err(err),
						log.Int64("store_id", storeID), log.Int64("website_id", websiteID),
						log.Stringer("run_mode", runMode), log.HTTPRequest("request", r))
				}
				errH(errors.Wrap(err, "[store] WithRunMode.AllowedStoreIDs")).ServeHTTP(w, r)
				return
			}

			// not found, not active, whatever, we cannot proceed.
			if !isStoreAllowed {
				if lg.IsDebug() {
					lg.Debug("runmode.WithRunMode.StoreNotAllowed",
						log.Bool("is_store_allowed", isStoreAllowed), log.String("store_code", storeCode),
						log.Int64("store_id", storeID), log.Int64("website_id", websiteID),
						log.Stringer("run_mode", runMode), log.HTTPRequest("request", r))
				}
				procCode.ProcessDenied(runMode, oldStoreID, storeID, w, r)
				errH(errors.NewUnauthorizedf("[store] RunMode %s with requested Store ID %d cannot be authorized", runMode, storeID)).ServeHTTP(w, r)
				return
			}

			// if runMode is allowed to change, update the runMode Hash and then put it into the context
			procCode.ProcessAllowed(runMode, oldStoreID, storeID, storeCode, w, r)

			if lg.IsDebug() {
				lg.Debug("runmode.WithRunMode.NextHandler",
					log.Bool("is_store_allowed", isStoreAllowed), log.String("store_code", storeCode),
					log.Int64("store_id", storeID), log.Int64("website_id", websiteID),
					log.Stringer("run_mode", runMode), log.HTTPRequest("request", r))
			}
			next.ServeHTTP(w, r)
		})
	}
}
