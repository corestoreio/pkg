// Copyright 2015-present, Cyrill @ Schumacher.fm and the CoreStore contributors
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

package request

import (
	"fmt"
	"net/http"
	"strings"
)

// AcceptsJSON returns true if the request requests JSON as a return-format.
func AcceptsJSON(r *http.Request) bool {
	accept := r.Header.Get("Accept")
	if "" == accept {
		return true
	}
	return strings.Contains(accept, "*/*") || strings.Contains(accept, "application/*") || strings.Contains(accept, "application/json")
}

// AcceptsContentType checks if a request requests a specific content type.
func AcceptsContentType(r *http.Request, contentType string) bool {
	accept := r.Header.Get("Accept")
	if "" == accept {
		return true
	}
	if strings.Contains(accept, "*/*") {
		return true
	}
	if strings.Contains(accept, contentType) {
		return true
	}
	typeParts := strings.Split(contentType, "/")
	if len(typeParts) < 2 {
		return false
	}
	return strings.Contains(accept, fmt.Sprintf("%s/*", typeParts[0]))
}
