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

/*
Package config handles the configuration and its scopes.

A configuration holds many path.Paths which contains a route, a scope and a scope ID.

A route is defined as a minimum 3 level deep string separated by a slash. For example
catalog/product/enable_flat.

Scopes are default, website, group and store. Scope IDs are stored in the core_website,
core_group or core_store tables for M1 and store_website, store_group and store for M2.

Underlying storage can be a simple in memory map (default), MySQL table core_config_data
itself (package config/db) or etcd (package config/etcd) or consul (package todo) or ...

If you use any other configuration storage engine besides config/db package all values
gets bi-directional automatically synchronized (todo).

Elements

The package config/element contains more detailed information.

Scope Values

To get a value from the configuration Service via any type method you have to provide a
path.Path. If you use the ScopedGetter via function NewScoped() you can only provide a
path.Route to the type methods String(), Int(), Float64(), etc.

The examples show the overall best practices.
*/
package config
