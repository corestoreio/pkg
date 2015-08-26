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
	"fmt"
	"runtime"
	"sync"

	"github.com/corestoreio/csfw/codegen"
	"github.com/corestoreio/csfw/storage/csdb"
)

func main() {
	dbc, err := csdb.Connect()
	codegen.LogFatal(err)
	defer dbc.Close()
	var wg sync.WaitGroup
	fmt.Printf("CPUs: %d\tGoroutines: %d\n", runtime.NumCPU(), runtime.NumGoroutine())
	for _, tStruct := range codegen.ConfigTableToStruct {
		go newGenerator(tStruct, dbc, &wg).run()
	}
	fmt.Printf("Goroutines: %d\tGo Version %s\n", runtime.NumGoroutine(), runtime.Version())
	wg.Wait()

	//	for _, ts := range codegen.ConfigTableToStruct {
	//		// due to a race condition the codec generator must run after the newGenerator() calls
	//		runCodec(ts.OutputFile.AppendName("_codec").String(), ts.OutputFile.String())
	//	}

}
