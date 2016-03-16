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

// Package store implements the handling of websites, groups and stores.
//
// The following shows a hierarchical diagram of the structure:
//        +---------------------+
//        |  Website            |
//        |   ID     <-----------------+---+
//        |   Code              |      |   |
//     +----+ Default Group ID  |      |   |
//     |  |   Is Default        |      |   |
//     |  +---------------------+      |   |
//     |                               |   |
//     |    +----------------------+   |   |
//     |    |  Group               |   |   |
//     +------> ID                 |   |   |
//          |   Website ID +-----------+   |
//          |   Root Category ID   |       |
//     +------+ Default Store ID   |       |
//     |    +----------------------+       |
//     |                                   |
//     |      +---------------+            |
//     |      |  Store        |            |
//     |      |   ID          |            |
//     |      |   Code        |            |
//     +--------> Group ID    |            |
//            |   Website ID +-------------+
//            |   Is Active   |
//            +---------------+
//     http://asciiflow.com
//
// Those three objects also represents the tables in the database.
//
// Sub package Scope
//
// The subpackage scope depends on these structure except that the group has
// been removed and a default scope has been introduced.
//
// More explanation @todo
package store
