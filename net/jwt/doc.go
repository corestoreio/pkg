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

// Package jwt provides a middleware for JSON web token authentication and store
// initialization.
//
// Further reading: https://float-middle.com/json-web-tokens-jwt-vs-sessions/
// and http://cryto.net/~joepie91/blog/2016/06/13/stop-using-jwt-for-sessions/
//
// https://news.ycombinator.com/item?id=11929267 => For people using JWT as a
// substitute for stateful sessions, how do you handle renewal (or revocation)?
package jwt
