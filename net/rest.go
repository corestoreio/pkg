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

// net provides functions and configuration values for the network.
package net

import (
	"encoding/json"
	"net/http"

	"github.com/corestoreio/csfw/utils/log"
	"github.com/juju/errgo"
	"github.com/labstack/echo"
)

const (
	CharsetUTF8 = "charset=utf-8"

	ContentLength = "Content-Length"
	ContentType   = "Content-Type"

	ApplicationJSON            = "application/json"
	ApplicationJSONCharsetUTF8 = ApplicationJSON + "; " + CharsetUTF8
)

// APIRoute defines the current API version
const APIRoute apiVersion = "/V1/"

type apiVersion string

// Versionize prepends the API version as defined in constant APIRoute to a route.
func (a apiVersion) Versionize(r string) string {
	if len(r) > 0 && r[:1] == "/" {
		r = r[1:]
	}
	return string(a) + r
}

// String returns the current version and not the full route
func (a apiVersion) String() string {
	return string(a)
}

// WriteJSON encodes v into JSON and sets the appropriate header except Length header
func WriteJSON(w http.ResponseWriter, v interface{}) error {
	w.Header().Set(ContentType, ApplicationJSONCharsetUTF8)
	// w.Header().Set(ContentLength, strconv.Itoa(len( ??? ))) @todo ...
	w.WriteHeader(http.StatusOK)
	return errgo.Mask(json.NewEncoder(w).Encode(v))
}

// RESTErrorHandler default REST error handler ... @todo remove echo dependency
func RESTErrorHandler(err error, c *echo.Context) {
	code := http.StatusInternalServerError
	msg := http.StatusText(code)
	if he, ok := err.(*echo.HTTPError); ok {
		code = he.Code()
		msg = he.Error()
	}
	if log.IsDebug() {
		log.Error("net.RESTErrorHandler", "err", err)
	}
	msg = err.Error()

	http.Error(c.Response(), msg, code)

}
