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

// Package objcache transcodes arbitrary Go types to bytes and stores them in
// a cache reducing GC.
//
// A Cache can be either in memory or a persistent one. Cache adapters are
// available for bigcache or Redis. To enable the cache adapter use build tags
// "bigcache" or "redis" or "csall".
//
// Use case: Caching millions of Go types as a byte slice reduces the pressure
// to the GC.
//
// For more details regarding bigcache: https://godoc.org/github.com/allegro/bigcache
package objcache
