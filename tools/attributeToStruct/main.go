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

// Generates code for all EAV attribute types
package main

import (
	"database/sql"
	"encoding/json"
	"flag"
	"io/ioutil"
	"os"
	"path"

	"go/build"

	"strings"

	"github.com/corestoreio/csfw/eav"
	"github.com/corestoreio/csfw/materialized"
	"github.com/corestoreio/csfw/storage/csdb"
	"github.com/corestoreio/csfw/storage/dbr"
	"github.com/corestoreio/csfw/tools"
	"github.com/juju/errgo"
)

var (
	run        = flag.Bool("run", false, "If true program runs")
	outputFile = flag.String("o", "", "Output file name")
)

const (
	envModelMap string = "CS_ATTRIBUTE_MODEL_MAP"
)

type context struct {
	db        *sql.DB
	dbrConn   *dbr.Connection
	et        *eav.CSEntityType // will be updated each iteration
	modelMap  tools.AttributeModelMap
	goSrcPath string
}

func newContext() *context {
	db, dbrConn, err := csdb.Connect()
	tools.LogFatal(err)
	modelMap, err := getMapping(os.Getenv(envModelMap), tools.JSONMapAttributeModels)
	tools.LogFatal(err)

	return &context{
		db:        db,
		dbrConn:   dbrConn,
		modelMap:  modelMap,
		goSrcPath: build.Default.GOPATH + "/src/",
	}
}

func main() {
	flag.Parse()

	if false == *run || *outputFile == "" {
		flag.Usage()
		os.Exit(1)
	}

	ctx := newContext()
	defer ctx.db.Close()

	for _, et := range materialized.GetEntityTypeCollection() {
		ctx.et = et
		tools.LogFatal(generateAttributeCode(ctx))

		// EAV -> Create queries for AttributeSets and AttributeGroups
	}
}

// getName
func getName(ctx *context, suffix ...string) string {
	pkg := path.Base(ctx.et.ImportPath)
	structBaseName := ctx.et.EntityTypeCode
	if strings.Contains(ctx.et.EntityTypeCode, "_") {
		structBaseName = strings.Replace(ctx.et.EntityTypeCode, pkg+"_", "", -1)
	}
	return structBaseName + "_" + strings.Join(suffix, "_")
}

func getMapping(fileName string, rawJson []byte) (tools.AttributeModelMap, error) {
	var err error
	if fileName != "" && fileName[len(fileName)-5:] == ".json" { // check if file ext is .json
		rawJson = nil
		rawJson, err = ioutil.ReadFile(fileName)
		if err != nil {
			return nil, errgo.Mask(err)
		}
	}
	mapping := make(tools.AttributeModelMap)
	err = json.Unmarshal(rawJson, &mapping)
	if err != nil {
		return nil, errgo.Mask(err)
	}
	return mapping, nil
}
