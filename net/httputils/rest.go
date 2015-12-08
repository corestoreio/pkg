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

// todo(CS) http://racksburg.com/choosing-an-http-status-code/

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
