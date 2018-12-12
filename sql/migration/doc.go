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

// Package migration provides tools for database schema migrations.
//
// TODO(CyS): https://povilasv.me/2017/02/20/go-schema-migration-tools/
//
// TL;DR If your looking for schema migration tool you can use:
//
// mattes/migrate, SQL defined schema migrations, with a well defined and
// documented API, large database support and a useful CLI tool. This tool is
// actively maintained, has a lot of stars and an A+ from goreport.
//
// rubenv/sql-migrate, go struct based or SQL defined schema migrations, with a
// config file, migration history, prod-dev-test environments. The only drawback
// is that it got B from goreport.
// SQL Schema migration tool for Go. Based on gorp and goose.
// Is better because the API for using it within your code (FOSDEM2016).
//
// https://bitbucket.org/liamstask/goose
// goose is a database migration tool. You can manage your database's evolution by
// creating incremental SQL or Go scripts.
// Cons: https://www.reddit.com/r/golang/comments/2dlbz5/database_migration_handling_in_go/
//
// github.com/mattes/migrate 744 Stars
// A migration helper written in Go. Use it in your existing Golang code or run
// commands via the CLI.
//
// https://github.com/golang-migrate/migrate
package migration
