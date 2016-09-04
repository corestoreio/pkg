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
Package storage provides everything from MySQL, Redis, BoltDB, file, etc functions.

DB Schema Go Generation

bitbucket.org/jatone/genieql
sql query and code generation program. its purpose is to generate as much of the
boilerplate code for interacting with database/sql as possible. without putting
any runtime dependencies into your codebase. it only supports postgresql
currently. adding additional support is very straight forward, just implement
the Dialect interface. see the postgresql implementation as the example.

github.com/go-reform/reform
A better ORM for Go, based on non-empty interfaces and code generation.
https://gopkg.in/reform.v1


DB Migrations

github.com/rubenv/sql-migrate 500 Stars
SQL Schema migration tool for Go. Based on gorp and goose.
Is better because the API for using it within your code (FOSDEM2016).

https://bitbucket.org/liamstask/goose
goose is a database migration tool. You can manage your database's evolution by
creating incremental SQL or Go scripts.
Cons: https://www.reddit.com/r/golang/comments/2dlbz5/database_migration_handling_in_go/

github.com/mattes/migrate 744 Stars
A migration helper written in Go. Use it in your existing Golang code or run
commands via the CLI.


*/
package storage
