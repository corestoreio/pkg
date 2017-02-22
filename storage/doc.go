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

// Package storage provides everything from MySQL, Redis, BoltDB, file etc
// functions.
//
// DB Schema Go Generation
//
// Will be written on our own as all code generation tools have some serious
// architectural flaw. See package codegen.
//
// bitbucket.org/jatone/genieql
// sql query and code generation program. its purpose is to generate as much of the
// boilerplate code for interacting with database/sql as possible. without putting
// any runtime dependencies into your codebase. it only supports postgresql
// currently. adding additional support is very straight forward, just implement
// the Dialect interface. see the postgresql implementation as the example.
//
// github.com/go-reform/reform
// A better ORM for Go, based on non-empty interfaces and code generation.
// https://gopkg.in/reform.v1
package storage
