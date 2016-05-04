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

/*
@todo:
	- routes to implement GH Issue#1
	- authentication
*/

var (
	// RoutePrefix global prefix for this package
	RoutePrefix = "store/"
	// RouteStores defines the REST API endpoints for GET and POST requests.
	RouteStores = RoutePrefix + "stores"
	// RouteStore defines the REST API endpoints for GET, PUT and DELETE a single store.
	RouteStore = RoutePrefix + "stores/:id"
)
