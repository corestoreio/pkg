// Copyright 2015 CoreStore Authors
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
package csfw contains only go:generate commands to build go code. Some general coding
styles and naming conventions.

Purpose

Why is someone trying to create a Magento database schema compatible online shop in Go?

Because performance :-)

Architecture

...

Names

Generated SQL table structs start with the word "Table". The word "Slice" will
be appended when there is a slice of structs.

Example for generated SQL table structs:

    type (
        // TableStoreSlice contains pointers to TableStore types
        TableStoreSlice []*TableStore
        // TableStore a type for the MySQL table core_store
        TableStore struct {
            ...
        }
    )

Table indexes are iota constants and start with TableIndex[table name].

The word "collection" will be appended to a variable or function when that variable or function contains
a materialized slice or handles it.

Trademarks

Magento is a trademark of MAGENTO, INC. http://www.magentocommerce.com/license/

*/
package csfw
