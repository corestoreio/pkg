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

// Package csMake replaces the Makefile and is only used via go:generate.
package main

import (
	"errors"
	"go/build"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/corestoreio/csfw/storage/csdb"
	"github.com/mgutz/ansi"
	// "github.com/pquerna/ffjson" @todo
)

type (
	aCommand struct {
		name string
		args []string
		rm   bool
	}
)

var (
	goCmd = build.Default.GOROOT + "/bin/go"
	pwd   = build.Default.GOPATH + "/src/github.com/corestoreio/csfw/"
)

// getCommands returns list of shell commands to be execute in its specific order
func getCommands() []aCommand {
	// @todo make it configurable if your catalog or customer tables different names.
	// @todo also include table_prefix for the whole database
	return []aCommand{
		aCommand{
			name: goCmd,
			args: []string{"build", "-a", "github.com/corestoreio/csfw/codegen/tableToStruct"},
			rm:   false,
		},
		aCommand{
			name: "find",
			args: []string{pwd, "-name", "generated_*.go", "-delete"},
			rm:   false,
		},
		aCommand{
			name: pwd + "tableToStruct",
			rm:   true,
		},
		aCommand{
			// this commands depends on the generated source from tableToStruct 8-)
			name: goCmd,
			args: []string{"build", "-a", "github.com/corestoreio/csfw/codegen/materialization"},
			rm:   false,
		},
		aCommand{
			name: pwd + "materialization",
			rm:   true,
		},
	}
}

// checkEnv verifies if all env vars are set
func checkEnv() error {
	dsn, err := csdb.GetDSN()
	if dsn == "" || err != nil {
		return errors.New(
			ansi.Color("Missing environment variable CS_DSN.", "red") +
				"Please see https://github.com/corestoreio/csfw#usage",
		)
	}
	return nil
}

func main() {

	if err := checkEnv(); err != nil {
		log.Fatal(err)
	}

	for _, cmd := range getCommands() {
		out, err := exec.Command(cmd.name, cmd.args...).CombinedOutput()
		if cmd.rm {
			defer os.Remove(cmd.name)
		}
		args := strings.Join(cmd.args, " ")
		if err != nil {
			log.Fatalf(ansi.Color("Failed:\n%s %s => %s\n%s", "red"), cmd.name, args, err, out)
		}

		log.Printf("%s %s", cmd.name, args)
		if nil != out && len(out) > 0 {
			log.Printf("%s", out)
		}
	}
	log.Println(ansi.Color("Done", "green"))
}
