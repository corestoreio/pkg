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
Package csfw contains the CoreStore FrameWork based on Magento's database structure.
99% compatible to Magento 1 and 2.

The package csfw contains at the moment only go:generate commands to build Go code.

Two skeleton projects (monolith and SOA) are already setup but of course empty.

Purpose

Why is someone trying to create a Magento database schema compatible online shop in Go?

Because performance :-)

Architecture

See the UML diagrams within the Go code. @todo 10km view.

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

Required settings

CS_DSN the environment variable for the MySQL connection.

    $ export CS_DSN='magento1:magento1@tcp(localhost:3306)/magento1'
    $ export CS_DSN='magento2:magento2@tcp(localhost:3306)/magento2'

    $ go get github.com/corestoreio/csfw
    $ export CS_DSN_TEST='see next section'
    $ cd $GOPATH/src/github.com/corestoreio/csfw
    $ go generate ./...

Testing

Setup two databases. One for Magento 1 and one for Magento 2 and fill them with
the provided test data https://github.com/corestoreio/csfw/tree/master/testData

Create a DSN env var CS_DSN_TEST and point it to Magento 1 database. Run the tests.
Change the env var to let it point to Magento 2 database. Rerun the tests.

    $ export CS_DSN_TEST='magento1:magento1@tcp(localhost:3306)/magento1'
    $ export CS_DSN_TEST='magento2:magento2@tcp(localhost:3306)/magento2'

IDE

Currently using the IntelliJ IDEA Community Edition with the https://github.com/go-lang-plugin-org/go-lang-idea-plugin plugin.

At the moment Q2/2015: There are no official jar files for downloading so the go lang plugin will be
compiled on a daily basis. The plugin works very well! Kudos to those developers!

IDEA has been configured with goimports, gofmt, golint, govet and ... with the file watcher plugin.

Why am I not using vim? Because I would only generate passwords ;-|.

UML

Within the Go files there are UML tags from PlantUML http://en.wikipedia.org/wiki/PlantUML because you
can have nice visualisations with the IDEA plugin of PlantUML. To run the PlantUML plugin you must
install:

    $ brew install graphviz

on your OSX. For Unix/Linux ... ?

Contributing

Please have a look at the contribution guidelines https://github.com/corestoreio/corestore/blob/master/CONTRIBUTING.md

Trademarks

Magento is a trademark of MAGENTO, INC. http://www.magentocommerce.com/license/

*/
package csfw
