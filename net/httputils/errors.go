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

import "net/http"

// WriteError writes an error to the ResponseWriter. Err cannot be nil. This
// variable can be replaced with your customized version.
var WriteError func(w http.ResponseWriter, err error, code int) = httpError

func httpError(w http.ResponseWriter, err error, code int) {
	http.Error(w, err.Error(), code)
}
