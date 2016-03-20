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

package storenet

import (
	"net/http"

	"github.com/corestoreio/csfw/backend"
	"github.com/corestoreio/csfw/config"
	"github.com/corestoreio/csfw/net/ctxhttp"
	"github.com/corestoreio/csfw/net/ctxjwt"
	"github.com/corestoreio/csfw/net/httputil"
	"github.com/corestoreio/csfw/store"
	"github.com/corestoreio/csfw/store/scope"
	"github.com/juju/errors"
	"golang.org/x/net/context"
)

// HTTPRequestParamStore name of the GET parameter to set a new store in a
// current website/group context
const HTTPRequestParamStore = `___store`

// WithValidateBaseURL is a middleware which checks if the request base URL
// is equal to the one store in the configuration, if not
// i.e. redirect from http://example.com/store/ to http://www.example.com/store/
// @see app/code/Magento/Store/App/FrontController/Plugin/RequestPreprocessor.php
func WithValidateBaseURL(cg config.GetterPubSuber) ctxhttp.Middleware {

	// Having the GetBool command here, means you must restart the app to take
	// changes in effect. @todo refactor and use pub/sub to automatically change
	// the isRedirectToBase value.

	// <todo check logic!>
	cgDefaultScope := cg.NewScoped(0, 0)
	configRedirectCode, err := backend.Backend.WebURLRedirectToBase.Get(cgDefaultScope)
	if err != nil {
		panic(err) // we can panic here because during app start up
	}

	redirectCode := http.StatusMovedPermanently
	if configRedirectCode != redirectCode {
		redirectCode = http.StatusFound
	}
	// </todo check logic>

	return func(hf ctxhttp.HandlerFunc) ctxhttp.HandlerFunc {
		return func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {

			if configRedirectCode > 0 && r.Method != "POST" {

				_, requestedStore, err := FromContextProvider(ctx)
				if err != nil {
					if PkgLog.IsDebug() {
						PkgLog.Debug("ctxhttp.WithValidateBaseUrl.FromContextServiceReader", "err", err, "ctx", ctx)
					}
					return errors.Mask(err)
				}

				baseURL, err := requestedStore.BaseURL(config.URLTypeWeb, requestedStore.IsCurrentlySecure(r))
				if err != nil {
					if PkgLog.IsDebug() {
						PkgLog.Debug("ctxhttp.WithValidateBaseUrl.requestedStore.BaseURL", "err", err, "ctx", ctx)
					}
					return errors.Mask(err)
				}

				if err := httputil.IsBaseURLCorrect(r, &baseURL); err != nil {
					if PkgLog.IsDebug() {
						PkgLog.Debug("store.WithValidateBaseUrl.IsBaseUrlCorrect.error", "err", err, "baseURL", baseURL, "request", r)
					}

					baseURL.Path = r.URL.Path
					baseURL.RawPath = r.URL.RawPath
					baseURL.RawQuery = r.URL.RawQuery
					baseURL.Fragment = r.URL.Fragment
					http.Redirect(w, r, (&baseURL).String(), redirectCode)
					return nil
				}
			}
			return hf(ctx, w, r)
		}
	}
}

// WithInitStoreByToken is a middleware which initializes a request based store
// via a JSON Web Token.
// Extracts the store.Provider and jwt.Token from context.Context. If the requested
// store is different than the initialized requested store than the new requested
// store will be saved in the context.
func WithInitStoreByToken() ctxhttp.Middleware {

	return func(hf ctxhttp.HandlerFunc) ctxhttp.HandlerFunc {
		return func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {

			storeService, requestedStore, err := FromContextProvider(ctx)
			if err != nil {
				if PkgLog.IsDebug() {
					PkgLog.Debug("store.WithInitStoreByToken.FromContextServiceReader", "err", err, "ctx", ctx)
				}
				return errors.Mask(err)
			}

			token, err := ctxjwt.FromContext(ctx)
			if err != nil {
				if PkgLog.IsDebug() {
					PkgLog.Debug("store.WithInitStoreByToken.ctxjwt.FromContext.err", "err", err, "ctx", ctx)
				}
				return errors.Mask(err)
			}

			scopeOption, err := CodeFromClaim(token.Claims)
			if err != nil {
				if PkgLog.IsDebug() {
					PkgLog.Debug("store.WithInitStoreByToken.StoreCodeFromClaim", "err", err, "token", token, "ctx", ctx)
				}
				return errors.Mask(err)
			}

			newRequestedStore, err := storeService.RequestedStore(scopeOption)
			if err != nil {
				if PkgLog.IsDebug() {
					PkgLog.Debug("store.WithInitStoreByToken.RequestedStore", "err", err, "token", token, "scopeOption", scopeOption, "ctx", ctx)
				}
				return errors.Mask(err)
			}

			if newRequestedStore.StoreID() != requestedStore.StoreID() {
				// this may lead to a bug because the previously set storeService and requestedStore
				// will still exists and have not been removed.
				ctx = WithContextProvider(ctx, storeService, newRequestedStore)
			}

			return hf.ServeHTTPContext(ctx, w, r)
		}
	}
}

// WithInitStoreByFormCookie reads from a GET parameter or cookie the store code.
// Checks if the store code is valid and allowed. If so it adjusts the context.Context
// to provide the new requestedStore.
//
// It calls Getter.RequestedStore() to determine the correct store.
// 		1. check cookie store, always a string and the store code
// 		2. check for GET ___store variable, always a string and the store code
func WithInitStoreByFormCookie() ctxhttp.Middleware {
	return func(hf ctxhttp.HandlerFunc) ctxhttp.HandlerFunc {
		return func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {

			storeService, requestedStore, err := FromContextProvider(ctx)
			if err != nil {
				if PkgLog.IsDebug() {
					PkgLog.Debug("store.WithInitStoreByToken.FromContextServiceReader", "err", err, "ctx", ctx)
				}
				return errors.Mask(err)
			}

			var reqSO scope.Option

			reqSO, err = CodeFromRequestGET(r)
			if err != nil {
				if PkgLog.IsDebug() {
					PkgLog.Debug("store.WithInitStoreByFormCookie.StoreCodeFromRequestGET", "err", err, "req", r, "scope", reqSO)
				}

				reqSO, err = CodeFromCookie(r)
				if err != nil {
					// ignore further processing because all codes are invalid or not found
					if PkgLog.IsDebug() {
						PkgLog.Debug("store.WithInitStoreByFormCookie.StoreCodeFromCookie", "err", err, "req", r, "scope", reqSO)
					}
					return hf.ServeHTTPContext(ctx, w, r)
				}
			}

			var newRequestedStore *store.Store
			if newRequestedStore, err = storeService.RequestedStore(reqSO); err != nil {
				if PkgLog.IsDebug() {
					PkgLog.Debug("store.WithInitStoreByFormCookie.storeService.RequestedStore", "err", err, "req", r, "scope", reqSO)
				}
				return errors.Mask(err)
			}

			soStoreCode := reqSO.StoreCode()

			// delete and re-set a new cookie, adjust context.Context
			if newRequestedStore != nil && newRequestedStore.Data.Code.String == soStoreCode {
				wds, err := newRequestedStore.Website.DefaultStore()
				if err != nil {
					if PkgLog.IsDebug() {
						PkgLog.Debug("store.WithInitStoreByFormCookie.Website.DefaultStore", "err", err, "soStoreCode", soStoreCode)
					}
					return errors.Mask(err)
				}
				keks := Cookie{Store: newRequestedStore}
				if wds.Data.Code.String == soStoreCode {
					keks.Delete(w) // cookie not needed anymore
				} else {
					keks.Set(w) // make sure we force set the new store

					if newRequestedStore.StoreID() != requestedStore.StoreID() {
						// this may lead to a bug because the previously set storeService and requestedStore
						// will still exists and have not been removed.
						ctx = WithContextProvider(ctx, storeService, newRequestedStore)
					}
				}
			}

			return hf(ctx, w, r)

		}
	}
}
