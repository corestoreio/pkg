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
	"regexp"
	"strings"
	"sync"

	"github.com/corestoreio/csfw/codegen"
	"github.com/corestoreio/csfw/codegen/tableToStruct/internal/codecgen"
)

var mapm1m2Mu sync.Mutex // protects map TableMapMagento1To2
const MethodRecvPrefix = "parent"

// TypePrefix of the generated types e.g. TableStoreSlice, TableStore ...
// If you change this you must change all "Table" in the template.
const TypePrefix = "Table"

// generatedFunctions: If a package has already such a function
// then prefix MethodRecvPrefix will be appended to the generated function
// so that in our code we can refer to the "parent" function. No composition possible.
// var generatedFunctions = map[string]bool{"Load": false, "Len": false, "Filter": false}

type duplicateChecker struct {
	dups map[string]bool
	mu   sync.RWMutex
}

func newDuplicateChecker(names ...string) *duplicateChecker {
	dc := &duplicateChecker{
		dups: make(map[string]bool),
		mu:   sync.RWMutex{},
	}
	dc.add(names...)
	return dc
}

func (dc *duplicateChecker) has(name string) bool {
	dc.mu.RLock()
	defer dc.mu.RUnlock()
	return dc.dups[name]
}

func (dc *duplicateChecker) add(names ...string) {
	dc.mu.Lock()
	defer dc.mu.Unlock()
	for _, n := range names {
		dc.dups[n] = true
	}
}

func (dc *duplicateChecker) debug() string {
	dc.mu.RLock()
	defer dc.mu.RUnlock()
	var ret string
	for n, _ := range dc.dups {
		ret = ret + n + ", "
	}
	return ret
}

// runCodec generates the codecs to be used later in JSON or msgpack or etc
func runCodec(outfile, readfile string) {

	if err := codecgen.Generate(
		outfile, // outfile
		"",      // buildTag
		codecgen.GenCodecPath,
		false, // use unsafe
		"",
		regexp.MustCompile(TypePrefix+".*"), // Prefix of generated structs and slices
		true,     // delete temp files
		readfile, // read from file
	); err != nil {
		fmt.Println("codecgen.Generate Error:")
		codegen.LogFatal(err)
	}
}

// isDuplicate slow duplicate checker ...
func isDuplicate(sl []string, st string) bool {
	for _, s := range sl {
		if s == st {
			return true
		}
	}
	return false
}

func getTableName(table string) (name string) {
	mapm1m2Mu.Lock()
	name = table
	if mappedName, ok := codegen.TableMapMagento1To2[strings.Replace(table, codegen.TablePrefix, "", 1)]; ok {
		name = mappedName
	}
	mapm1m2Mu.Unlock()
	return
}
