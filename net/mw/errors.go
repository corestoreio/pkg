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

package mw

import (
	"fmt"
	"net/http"

	"github.com/corestoreio/log"
	loghttp "github.com/corestoreio/log/http"
)

// ErrorHandler passes an error to an handler and returns the handler with the
// wrapped error.
type ErrorHandler func(error) http.Handler

// ErrorWithStatusCode wraps an HTTP Status Code into an ErrorHandler. The
// status text message gets printed first followed by the verbose error string.
// This function may leak sensitive information.
func ErrorWithStatusCode(code int) ErrorHandler {
	return func(err error) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			http.Error(w, fmt.Sprintf("%s\n\n%+v", http.StatusText(code), err), code)
		})
	}
}

// LogErrorWithStatusCode same as ErrorWithStatusCode but does not print the
// error message and logs it instead with level debug or info.
func LogErrorWithStatusCode(l log.Logger, code int) ErrorHandler {
	return func(err error) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			http.Error(w, http.StatusText(code), code)

			fields := log.Fields{
				log.Err(err), log.Int("status_code", code),
				loghttp.Request("request", loghttp.ShallowCloneRequest(r)),
			}
			if l.IsDebug() {
				l.Debug("mw.LogErrorWithStatusCode", fields...)
			} else if l.IsInfo() {
				l.Info("mw.LogErrorWithStatusCode", fields...)
			}
		})
	}
}

// ErrorWithPanic implements the ErrorHandler type and panics always. Interesting for
// testing. This function may leak sensitive information.
func ErrorWithPanic(err error) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		panic(fmt.Sprintf(`This handler should not get called!
============================================================
%+v
============================================================
`, err))
	})
}
