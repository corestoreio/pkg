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

/*
Package config handles the configuration and its scopes.

A configuration holds many path.Paths which contains a route, a scope and a scope ID.

A route is defined as a minimum 3 level deep string separated by a slash. For example
catalog/product/enable_flat.

ScopePerm are default, website, group and store. Scope IDs are stored in the core_website,
core_group or core_store tables for M1 and store_website, store_group and store for M2.

Underlying storage can be a simple in memory map (default), MySQL table core_config_data
itself (package config/db) or etcd (package config/etcd) or consul (package todo) or ...

If you use any other configuration storage engine besides config/db package all values
gets bi-directional automatically synchronized (todo).

Elements

The package config/element contains more detailed information.

Scope Values

To get a value from the configuration Service via any type method you have to provide a
path.Path. If you use the ScopedGetter via function MakeScoped() you can only provide a
path.Route to the type methods String(), Int(), Float64(), etc.

The examples show the overall best practices.
*/
// Package cfgpath handles the configuration paths.
//
// It contains two main types: Path and Route:
//
//    +-------+ +-----+ +----------+ +-----+ +----------------+
//    |       | |     | |          | |     | |                |
//    | Scope | |  /  | | Scope ID | |  /  | | Route/to/Value |
//    |       | |     | |          | |     | |                |
//    +-------+ +-----+ +----------+ +-----+ +----------------+
//
//    +                                      +                +
//    |                                      |                |
//    | <--------------+ Path +-----------------------------> |
//    |                                      |                |
//    +                                      + <- Route ----> +
//
// Scope
//
// A scope can only be default, websites or stores. Those three strings are
// defined by constants in package store/scope.
//
// Scope ID
//
// Refers to the database auto increment ID of one of the tables core_website
// and core_store for M1 and store_website plus store for M2.
//
// Type Path
//
// A Path contains always the scope, its scope ID and the route.
// If scope and ID haven't been provided they fall back to scope "default"
// and ID zero (0).
// Configuration paths are mainly used in table core_config_data.
//
// Type Route
//
// A route contains bytes and does not know anything about a scope or an ID.
// In the majority of use cases a route contains three parts to package
// config/element types for building a hierarchical tree structure:
//    element.Section.ID / element.Group.ID / element.Field.ID
// To add little bit more confusion: A route can either be a short one
// like aa/bb/cc or a fully qualified path like
//     scope/scopeID/element.Section.ID/element.Group.ID/element.Field.ID
// But the route always consists of a minimum of three parts.
//
// A route can have only three groups of [a-zA-Z0-9_] characters
// split by '/'. The limitation to [a-zA-Z0-9_] is a M1/M2 thing and can be
// maybe later removed.
// Minimal length per part 2 characters. Case sensitive.
//
// The route parts are used as an ID in element.Section, element.Group and
// element.Field types. See following text.
//
// The following diagram shows the tree structure:
//    +---------+
//    | Section |
//    +---------+
//    |
//    |   +-------+
//    +-->+ Group |
//    |   +-------+
//    |   |
//    |   |    +--------+
//    |   +--->+ Field  |
//    |   |    +--------+
//    |   |
//    |   |    +--------+
//    |   +--->+ Field  |
//    |   |    +--------+
//    |   |
//    |   |    +--------+
//    |   +--->+ Field  |
//    |        +--------+
//    |
//    |   +-------+
//    +-->+ Group |
//    |   +-------+
//    |   |
//    |   |    +--------+
//    |   +--->+ Field  |
//    |   |    +--------+
//    |   |
//    |   |    +--------+
//    |   +--->+ Field  |
//    |   |    +--------+
//    |   |
//    |   |    +--------+
//    |   +--->+ Field  |
//    |        +--------+
//    |
//    http://asciiflow.com/
//
// The three elements Section, Group and Field represents front-end
// configuration fields and more important default values and their permissions.
// A permission is of type scope.Perm and defines which of elements in which
// scope can be shown.
//
// Those three elements represents the tree in function NewConfigStructure()
// which can be found in any package.
//
// Unclear: Your app which includes the cs must merge all
// "PackageConfiguration"s into a single slice. You should submit all default
// values (interface config.Sectioner) to the config.Service.ApplyDefaults()
// function.
//
// The JSON encoding of the three elements Section, Group and Field are intended
// to use on the backend REST API and for debugging and testing. Only used in
// non performance critical parts.package config
package config
