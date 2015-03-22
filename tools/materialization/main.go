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

// Generates code
package main

import (
	"database/sql"
	"sync"

	"go/build"

	"runtime"

	"github.com/corestoreio/csfw/eav"
	"github.com/corestoreio/csfw/storage/csdb"
	"github.com/corestoreio/csfw/storage/dbr"
	"github.com/corestoreio/csfw/tools"
)

type context struct {
	wg      sync.WaitGroup
	db      *sql.DB
	dbrConn *dbr.Connection
	// will be updated each iteration in materializeAttributes
	et *eav.EntityType
	// goSrcPath will be used in conjunction with ImportPath to write a file into that directory
	goSrcPath string
}

func newContext() *context {
	db, dbrConn, err := csdb.Connect()
	tools.LogFatal(err)

	return &context{
		wg:        sync.WaitGroup{},
		db:        db,
		dbrConn:   dbrConn,
		goSrcPath: build.Default.GOPATH + "/src/",
	}
}

func main() {

	runtime.GOMAXPROCS(runtime.NumCPU())

	ctx := newContext()
	defer ctx.db.Close()

	ctx.wg.Add(1)
	go materializeEntityType(ctx)

	ctx.wg.Add(1)
	go materializeAttributes(ctx)

	// EAV -> Create queries for AttributeSets and AttributeGroups
	//    ctx.wg.Add(1)
	//    go materializeAttributeSets(ctx)
	//
	//    ctx.wg.Add(1)
	//    go materializeAttributeGroups(ctx)

	ctx.wg.Wait()
}
