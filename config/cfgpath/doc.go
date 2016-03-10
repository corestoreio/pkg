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

// Package cfgpath handles the configuration paths.
//
// A path can either be a short one like a/b/c or a fully qualified path
// like stores/3/a/b/c. The prefix "stores" gets handle by the package store/scope
// and the number 3 represents a Store with ID 3.
// Configuration paths are mainly used in table core_config_data.
// Configuration path attribute can have only three groups of [a-zA-Z0-9_] characters split by '/'.
// Minimal length per part 2 characters. Case sensitive.
//
// Path parts are used as an ID in section, group and field types.
package cfgpath
