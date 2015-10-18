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

package store

import (
	"net/http"

	"errors"
	"net/url"
	"strings"

	"github.com/corestoreio/csfw/config"
	"github.com/corestoreio/csfw/config/scope"
	"github.com/corestoreio/csfw/net/ctxhttp"
	"github.com/corestoreio/csfw/utils/log"
	"golang.org/x/net/context"
)

// ErrBaseUrlDoNotMatch will be returned if the request URL does not match the configured URL.
var ErrBaseUrlDoNotMatch = errors.New("The Base URLs do not match")

// WithValidateBaseUrl is a middleware which checks if the request base URL
// is equal to the one store in the configuration, if not
// i.e. redirect from http://example.com/store/ to http://www.example.com/store/
// @see app/code/Magento/Store/App/FrontController/Plugin/RequestPreprocessor.php
func WithValidateBaseUrl(cr config.ReaderPubSuber) ctxhttp.Middleware {

	// Having the GetBool command here, means you must restart the app to take
	// changes in effect. @todo refactor and use pub/sub to automatically change
	// the isRedirectToBase value.
	checkBaseURL, err := cr.GetBool(config.Path(PathRedirectToBase)) // scope default
	if config.NotKeyNotFoundError(err) {
		log.Error("ctxhttp.WithValidateBaseUrl.GetBool", "err", err, "path", PathRedirectToBase)
	}

	redirectCode := http.StatusMovedPermanently
	if rc, err := cr.GetInt(config.Path(PathRedirectToBase)); rc != redirectCode && false == config.NotKeyNotFoundError(err) {
		redirectCode = http.StatusFound
	}

	return func(h ctxhttp.Handler) ctxhttp.Handler {
		return ctxhttp.HandlerFunc(func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {

			if checkBaseURL && r.Method != "POST" {

				storeManager, ok := FromContextManagerReader(ctx)
				if !ok {
					return log.Error("ctxhttp.WithValidateBaseUrl.FromContextManagerReader", "err", errors.New("Cannot extract config.Reader from context"), "ctx", ctx)
				}

				activeStore, err := storeManager.Store()
				if err != nil {
					return log.Error("ctxhttp.WithValidateBaseUrl.storeManager.Store", "err", err, "ctx", ctx)
				}

				baseURL := activeStore.BaseURL(config.URLTypeWeb, activeStore.IsCurrentlySecure(r))
				if nil == isBaseUrlCorrect(r, baseURL) {
					redirectURL := baseURL + r.URL.Path
					http.Redirect(w, r, redirectURL, redirectCode)
					return nil
				}
			}
			return h.ServeHTTPContext(ctx, w, r)
		})
	}
}

// isBaseUrlCorrect checks if the requested host, scheme are same as the servers and
// if the path of the baseURL is included in the request URI.
func isBaseUrlCorrect(r *http.Request, baseURL string) error {
	uri, err := url.Parse(baseURL)
	if err != nil {
		return log.Error("store.isBaseUrlCorrect.url.Parse", "err", err)
	}

	if r.Host == uri.Host && r.URL.Host == uri.Host && r.URL.Scheme == uri.Scheme && strings.Contains(r.URL.RequestURI(), uri.Path) {
		return nil
	}
	return log.Error("store.isBaseUrlCorrect.compare", "err", ErrBaseUrlDoNotMatch, "r.Host", r.Host, "baseURL", uri.String(), "requestURL", r.URL.String(), "strings.Contains", []string{r.URL.RequestURI(), uri.Path})
}

// InitByToken returns a Store pointer from a JSON web token. If the store code is invalid,
// this function can return nil,nil. Token argument is equal like jwt.Token.Claim.
func WithInitStoreByToken(scopeType scope.Scope) ctxhttp.Middleware {

	return func(h ctxhttp.Handler) ctxhttp.Handler {
		return ctxhttp.HandlerFunc(func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {

			//			storeManager, ok := FromContextManagerReader(ctx)
			//			if !ok {
			//				return log.Error("ctxhttp.WithValidateBaseUrl.FromContextManagerReader", "err", errors.New("Cannot extract config.Reader from context"), "ctx", ctx)
			//			}

			return nil
		})
	}

	//	scopeOption, err := store.StoreCodeFromClaim(token)
	//	if err == nil {
	//		return sm.GetRequestStore(scopeOption, scopeType)
	//	}

}

// InitByRequest returns a new Store read from a cookie or HTTP request parameter.
// It calls GetRequestStore() to determine the correct store.
// The internal appStore must be set before hand, call Init() before calling this function.
// 1. check cookie store, always a string and the store code
// 2. check for ___store variable, always a string and the store code
// 3. May return nil,nil if nothing is set.
// This function must be used within an HTTP handler.
// The returned new Store must be used in the HTTP context and overrides the appStore.
func WithInitStoreByRequest(scopeType scope.Scope) ctxhttp.Middleware {

	return func(h ctxhttp.Handler) ctxhttp.Handler {
		return ctxhttp.HandlerFunc(func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {

			return nil
		})
	}

	// res http.ResponseWriter, req *http.Request
	//	if sm.appStore == nil {
	//		// that means you must call Init() before executing this function.
	//		return nil, ErrAppStoreNotSet
	//	}
	//
	//	var reqStore *Store
	//	var so scope.Option
	//	var err error
	//	so, err = StoreCodeFromForm(req)
	//	if err != nil { // no cookie set, lets try via form to find the store code
	//
	//		if err == ErrStoreCodeInvalid {
	//			return nil, log.Error("store.Manager.InitByRequest.GetCodeFromForm", "err", err, "req", req, "scopeType", scopeType.String())
	//		}
	//
	//		so, err = StoreCodeFromCookie(req)
	//		switch err {
	//		case ErrStoreCodeEmpty, http.ErrNoCookie:
	//			err = nil
	//		case nil:
	//		// do nothing
	//		default: // err != nil
	//			return nil, log.Error("store.Manager.InitByRequest.GetCodeFromCookie", "err", err, "req", req, "scopeType", scopeType.String())
	//		}
	//	}
	//
	//	// @todo reqStoreCode if number ... cast to int64 because then group id if ScopeGroup is group.
	//	if reqStore, err = sm.GetRequestStore(so, scopeType); err != nil {
	//		return nil, log.Error("store.Manager.InitByRequest.GetRequestStore", "err", err)
	//	}
	//	soStoreCode := so.StoreCode()
	//
	//	// also delete and re-set a new cookie
	//	if reqStore != nil && reqStore.Data.Code.String == soStoreCode {
	//		wds, err := reqStore.Website.DefaultStore()
	//		if err != nil {
	//			return nil, log.Error("store.Manager.InitByRequest.Website.DefaultStore", "err", err, "soStoreCode", soStoreCode)
	//		}
	//		if wds.Data.Code.String == soStoreCode {
	//			reqStore.DeleteCookie(res) // cookie not needed anymore
	//		} else {
	//			reqStore.SetCookie(res) // make sure we force set the new store
	//		}
	//	}
	//
	//	return reqStore, nil // can be nil,nil
}
