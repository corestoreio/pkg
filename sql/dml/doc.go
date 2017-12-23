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

// Package dml handles the SQL DML for super fast performance,
// type safety and convenience.
//
// Aim: Allow a developer to easily modify a SQL query without type assertion of
// parts of the query. No reflection magic has been used so we must achieve
// type safety with code generation.
//
// This package works only with MySQL and its derivates like MariaDB or Percona.
//
// Abbreviations
//
// DML (https://en.wikipedia.org/wiki/Data_manipulation_language) Select,
// Insert, Update and Delete.
//
// DDL (https://en.wikipedia.org/wiki/Data_definition_language) Create, Drop,
// Alter, and Rename.
//
// DCL (https://en.wikipedia.org/wiki/Data_control_language) Grant and Revoke.
//
// CRUD (https://en.wikipedia.org/wiki/Create,_read,_update_and_delete) Create,
// Read, Update and Delete.
//
// https://mariadb.com/kb/en/mariadb/documentation/
//
// Practical Guide to SQL Transaction Isolation: https://begriffs.com/posts/2017-08-01-practical-guide-sql-isolation.html
//
// NetSPI SQL Injection Wiki: https://sqlwiki.netspi.com/
//
// TODO(CyS) think about named locks:
// https://news.ycombinator.com/item?id=14907679
// https://dev.mysql.com/doc/refman/5.7/en/miscellaneous-functions.html#function_get-lock
// Database locks should not be used by the average developer. Understand
// optimistic concurrency and use serializable isolation.
//
// TODO(CyS) refactor some parts of the code once Go implements generics ;-)
//
// TODO(CyS) implement usage of window functions:
//    - https://mariadb.com/kb/en/library/window-functions/
//    - https://dev.mysql.com/doc/refman/8.0/en/window-functions-usage.html
//    - https://blog.statsbot.co/sql-window-functions-tutorial-b5075b87d129
package dml
