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

package cshttp

const (
	// HTTPMethodOverrideHeader represents a commonly used http header to override a request method.
	HTTPMethodOverrideHeader = "X-HTTP-Method-Override"
	// HTTPMethodOverrideFormKey represents a commonly used HTML form key to override a request method.
	HTTPMethodOverrideFormKey = "_method"
)

// HTTPMethodxxx defines the available methods which this framework supports
const (
	HTTPMethodHead    = `HEAD`
	HTTPMethodGet     = "GET"
	HTTPMethodPost    = "POST"
	HTTPMethodPut     = "PUT"
	HTTPMethodPatch   = "PATCH"
	HTTPMethodDelete  = "DELETE"
	HTTPMethodTrace   = "TRACE"
	HTTPMethodOptions = "OPTIONS"
)
