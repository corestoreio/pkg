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

// Package mycanal adds event listener to a MySQL/MariaDB compatible binlog,
// based on pkg myreplicator to provide triggerless events.
//
// Overview
//
// Triggers are stored routines which are invoked on a per-row operation upon
// INSERT, DELETE, UPDATE on a table. They were introduced in MySQL 5.0. A
// trigger may contain a set of queries, and these queries run in the same
// transaction space as the query that manipulates the table. This makes for an
// atomicity of both the original operation on the table and the trigger-invoked
// operations.
//
// Triggers, overhead
//
// A trigger in MySQL is a stored routine. MySQL stored routines are
// interpreted, never compiled. With triggers, for every INSERT, DELETE, UPDATE
// on the often busy tables, it pays the necessary price of the additional
// write, but also the price of interpreting the trigger body.
//
// We know this to be a visible overhead on very busy or very large tables.
//
// Triggers, locks
//
// When a table with triggers is concurrently being written to, the triggers,
// being in same transaction space as the incoming queries, are also executed
// concurrently. While concurrent queries compete for resources via locks (e.g.
// the auto_increment value), the triggers need to simultaneously compete for
// their own locks (e.g., likewise on the auto_increment value on the ghost
// table, in a synchronous solution). These competitions are non-coordinated.
//
// We have evidenced near or complete lock downs in production, to the effect of
// rendering the table or the entire database inaccessible due to lock
// contention.
//
// Thus, triggers must keep operating. On busy servers, we have seen that even
// as the online operation throttles, the master is brought down by the load of
// the triggers.
package mycanal
