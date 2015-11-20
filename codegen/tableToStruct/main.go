// Copyright 2015, Cyrill @ Schumacher.fm and the CoreStore contributors
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
	"sync"

	"github.com/corestoreio/csfw/codegen"
	"github.com/corestoreio/csfw/storage/csdb"
	"github.com/corestoreio/csfw/utils/log"
)

// PkgLog global package based logger
var PkgLog log.Logger = log.NewStdLogger()

func init() {
	codegen.PkgLog = PkgLog
	log.PkgLog = PkgLog
}

func main() {
	defer log.WhenDone(PkgLog).Info("Stats")
	dbc, err := csdb.Connect()
	codegen.LogFatal(err)
	defer dbc.Close()
	var wg sync.WaitGroup

	mageV1, mageV2 := detectMagentoVersion(dbc.NewSession())

	for _, tStruct := range codegen.ConfigTableToStruct {
		go newGenerator(tStruct, dbc, &wg).setMagentoVersion(mageV1, mageV2).run()
	}

	wg.Wait()

	//	for _, ts := range codegen.ConfigTableToStruct {
	//		// due to a race condition the codec generator must run after the newGenerator() calls
	//		// TODO(cs) fix https://github.com/ugorji/go/issues/92#issuecomment-140410732
	//		runCodec(ts.Package, ts.OutputFile.AppendName("_codec").String(), ts.OutputFile.String())
	//	}
}
