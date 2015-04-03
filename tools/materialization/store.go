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

package main

import (
	"fmt"
	"io/ioutil"

	"github.com/corestoreio/csfw/store"
	"github.com/corestoreio/csfw/tools"
	"github.com/juju/errgo"
)

type (
	storeTplData struct {
		PackageName string
		Stores      store.TableStoreSlice
		Groups      store.TableGroupSlice
		Websites    store.TableWebsiteSlice
	}
)

func newStoreTplData(ctx *context) *storeTplData {
	tplData := &storeTplData{
		PackageName: tools.ConfigMaterializationStore.Package,
	}

	rowCount, err := tplData.Stores.Load(ctx.dbrConn.NewSession(nil))
	tools.LogFatal(errgo.Mask(err))
	if rowCount < 1 {
		tools.LogFatal(errgo.New("There are no stores in the database!"))
	}

	rowCount, err = tplData.Groups.Load(ctx.dbrConn.NewSession(nil))
	tools.LogFatal(errgo.Mask(err))
	if rowCount < 1 {
		tools.LogFatal(errgo.New("There are no groups in the database!"))
	}

	rowCount, err = tplData.Websites.Load(ctx.dbrConn.NewSession(nil))
	tools.LogFatal(errgo.Mask(err))
	if rowCount < 1 {
		tools.LogFatal(errgo.New("There are no groups in the database!"))
	}
	return tplData
}

// materializeStore writes the data from store, store_group and store_website.
// Depends on generated code from tableToStruct.
func materializeStore(ctx *context) {
	defer ctx.wg.Done()

	tplData := newStoreTplData(ctx)
	formatted, err := tools.GenerateCode(tools.ConfigMaterializationStore.Package, tplMaterializationStore, tplData, nil)
	if err != nil {
		fmt.Printf("\n%s\n", formatted)
		tools.LogFatal(err)
	}

	tools.LogFatal(ioutil.WriteFile(tools.ConfigMaterializationStore.OutputFile, formatted, 0600))
}
