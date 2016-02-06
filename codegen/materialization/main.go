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

package main

import (
	"go/build"
	"runtime"
	"sync"

	"github.com/corestoreio/csfw/codegen"
	"github.com/corestoreio/csfw/eav"
	"github.com/corestoreio/csfw/storage/csdb"
	"github.com/corestoreio/csfw/storage/dbr"
	"github.com/corestoreio/csfw/util/cstesting"
)

// depends on generated code from tableToStruct
type context struct {
	wg  sync.WaitGroup
	dbc *dbr.Connection
	// will be updated each iteration in materializeAttributes
	et *eav.TableEntityType
	// goSrcPath will be used in conjunction with ImportPath to write a file into that directory
	goSrcPath string
	// aat = additional attribute table
	aat *codegen.AddAttrTables
}

func newContext() *context {
	dbc, err := csdb.Connect()
	codegen.LogFatal(err)

	return &context{
		wg:        sync.WaitGroup{},
		dbc:       dbc,
		goSrcPath: cstesting.RootPath,
	}
}

func main() {

	runtime.GOMAXPROCS(runtime.NumCPU())

	ctx := newContext()
	defer ctx.dbc.Close()

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
