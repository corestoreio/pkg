// Copyright 2015-2017, Cyrill @ Schumacher.fm and the CoreStore contributors
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

// Package mview adds materialized views via events on the MySQL binary log.
//
//
// https://de.slideshare.net/MySQLGeek/flexviews-materialized-views-for-my-sql
// https://github.com/greenlion/swanhart-tools
//
// https://hashrocket.com/blog/posts/materialized-view-strategies-using-postgresql
// Queries returning aggregate, summary, and computed data are frequently used
// in application development. Sometimes these queries are not fast enough.
// Caching query results using Memcached or Redis is a common approach for
// resolving these performance issues. However, these bring their own
// challenges. Before reaching for an external tool it is worth examining what
// techniques PostgreSQL offers for caching query results.
//
// http://www.eschrade.com/page/indexing-in-magento-or-the-wonderful-world-of-materialized-views/
package mview
