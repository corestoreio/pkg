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

// Package ratelimit implements scope based HTTP rate limiting.
//
// Sub-package `backendratelimit` implements the external configuration loading.
// Sub-package `memstore` and `redigostore` provides rate limiting algorithms
// and their storage possibilities. Both packages should be used as either
// functional options to a ratelimit service or as functional option factories
// to the backend type.
package ratelimit
