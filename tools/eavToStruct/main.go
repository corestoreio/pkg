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

// Generates code for all EAV types
package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/corestoreio/csfw/eav"
	"github.com/corestoreio/csfw/tools"
	"github.com/corestoreio/csfw/tools/toolsdb"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gocraft/dbr"
	"github.com/juju/errgo"
	"github.com/kr/pretty"
)

var (
	cliDsn     = flag.String("dsn", "test:test@tcp(localhost:3306)/test", "MySQL DSN data source name. Can also be provided via ENV with key CS_DSN")
	pkg        = flag.String("package", "eav", "Package name in template")
	run        = flag.Bool("run", false, "If true program runs")
	outputFile = flag.String("file", "rename_me.go", "Output file name")
	//tableMap = flag.String("tableMap", "", "JSON file for mapping entity_table and attribute_table values to real table names")
)

func main() {
	flag.Parse()

	if false == *run {
		flag.Usage()
		os.Exit(1)
	}

	var dbrConn *dbr.Connection
	db, err := toolsdb.Connect(*cliDsn)
	toolsdb.LogFatal(err)
	defer db.Close()

	dbrConn = dbr.NewConnection(db, nil)
	dbrSess := dbrConn.NewSession(nil)

	type dataContainer struct {
		ETypeData     []string
		Package, Tick string
	}

	etData, err := getEntityTypeData(dbrSess)
	toolsdb.LogFatal(err)

	tplData := &dataContainer{
		ETypeData: etData,
		Package:   *pkg,
		Tick:      "`",
	}

	formatted, err := tools.GenerateCode(tplEav, tplData)
	if err != nil {
		fmt.Printf("\n%s\n", formatted)
		toolsdb.LogFatal(err)
	}

	ioutil.WriteFile(*outputFile, formatted, 0600)
	log.Println("ok")

}

func getEntityTypeData(dbrSess *dbr.Session) ([]string, error) {

	s, err := eav.GetTableStructure(eav.TableEntityType)
	if err != nil {
		return nil, errgo.Mask(err)
	}

	var entityTypeCollection eav.EntityTypeSlice
	_, err = dbrSess.
		Select(s.Columns...).
		From(s.Name).
		LoadStructs(&entityTypeCollection)
	if err != nil {
		return nil, errgo.Mask(err)
	}

	data := make([]string, len(entityTypeCollection))
	// @todo apply mapping of mage1+2 class names to Go types

	for i, eType := range entityTypeCollection {
		data[i] = fmt.Sprintf("%# v", pretty.Formatter(eType))
		data[i] = "&" + data[i][len("&"+*pkg):] // replace "&*pkg." with "&"
	}

	return data, nil
}
