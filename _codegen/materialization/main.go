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

package main

//import (
//	"runtime"
//	"sync"
//
//	"github.com/corestoreio/pkg/codegen"
//	"github.com/corestoreio/pkg/eav"
//	"github.com/corestoreio/pkg/storage/csdb"
//	"github.com/corestoreio/pkg/storage/dbr"
//	"github.com/corestoreio/pkg/util/cstesting"
//	"github.com/corestoreio/errors"
//)
//
//// depends on generated code from tableToStruct
//type context struct {
//	wg  sync.WaitGroup
//	dbc *dbr.ConnPool
//	// will be updated each iteration in materializeAttributes
//	et *eav.TableEntityType
//	// goSrcPath will be used in conjunction with ImportPath to write a file into that directory
//	goSrcPath string
//	// aat = additional attribute table
//	aat *codegen.AddAttrTables
//}
//
//// Connect creates a new database connection from a DSN stored in an
//// environment variable CS_DSN.
//func Connect(opts ...dbr.ConnPoolOption) (*dbr.ConnPool, error) {
//	c, err := dbr.NewConnPool(dbr.WithDSN(csdb.MustGetDSN()))
//	if err != nil {
//		return nil, errors.Wrap(err, "[csdb] dbr.NewConnPool")
//	}
//	if err := c.Options(opts...); err != nil {
//		return nil, errors.Wrap(err, "[csdb] dbr.NewConnPool.Options")
//	}
//	return c, err
//}
//
//func newContext() *context {
//	dbc, err := Connect()
//	codegen.LogFatal(err)
//
//	return &context{
//		wg:        sync.WaitGroup{},
//		dbc:       dbc,
//		goSrcPath: cstesting.RootPath,
//	}
//}
//
func main() {
	//
	//	runtime.GOMAXPROCS(runtime.NumCPU())
	//
	//	ctx := newContext()
	//	defer ctx.dbc.Close()
	//
	//	ctx.wg.Add(1)
	//	go materializeEntityType(ctx)
	//
	//	ctx.wg.Add(1)
	//	go materializeAttributes(ctx)
	//
	//	// EAV -> Create queries for AttributeSets and AttributeGroups
	//	//    ctx.wg.Add(1)
	//	//    go materializeAttributeSets(ctx)
	//	//
	//	//    ctx.wg.Add(1)
	//	//    go materializeAttributeGroups(ctx)
	//
	//	ctx.wg.Wait()
}
