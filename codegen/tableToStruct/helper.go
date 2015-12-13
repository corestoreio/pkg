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
	"errors"
	"fmt"
	"regexp"
	"sync"

	"github.com/corestoreio/csfw/codegen"
	"github.com/corestoreio/csfw/codegen/tableToStruct/codecgen"
	"github.com/corestoreio/csfw/storage/csdb"
	"github.com/corestoreio/csfw/storage/dbr"
	"github.com/corestoreio/csfw/util"
	"github.com/corestoreio/csfw/util/log"
)

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
	for n := range dc.dups {
		ret = ret + n + ", "
	}
	return ret
}

// runCodec generates the codecs to be used later in JSON or msgpack or etc
func runCodec(pkg, outfile, readfile string) {
	defer log.WhenDone(PkgLog).Info("Stats", "Package", pkg, "Step", "runCodec")
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

func detectMagentoVersion(dbrSess dbr.SessionRunner) (MageOne, MageTwo bool) {
	defer log.WhenDone(PkgLog).Info("Stats", "Package", "DetectMagentoVersion")
	allTables, err := codegen.GetTables(dbrSess)
	codegen.LogFatal(err)
	MageOne, MageTwo = util.MagentoVersion(codegen.TablePrefix, allTables)

	if MageOne == MageTwo {
		codegen.LogFatal(errors.New("Cannot detect your Magento version"))
	}
	return
}

// findBy is a template function used in runTable()
func findBy(s string) string {
	return "FindBy" + util.UnderscoreCamelize(s)
}

// dbrType is a template function used in runTable()
func dbrType(c csdb.Column) string {
	switch {
	// order of the c.Is* functions matters ... :-|
	case false == c.IsNull():
		return ""
	case c.IsBool():
		return ".Bool" // dbr.NullBool
	case c.IsString():
		return ".String" // dbr.NullString
	case c.IsMoney():
		return "" // money.Money
	case c.IsFloat():
		return ".Float64" // dbr.NullFloat64
	case c.IsInt():
		return ".Int64" // dbr.NullInt64
	case c.IsDate():
		return ".Time" // dbr.NullTime
	}
	return ""
}
