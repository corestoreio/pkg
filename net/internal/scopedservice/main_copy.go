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

// +build ignore

package main

import (
	"bytes"
	"io/ioutil"
	"os"
	"path/filepath"
)

// This file copies the *_scopedservice.go files to the packages via go:generate

func main() {

	pkgName := os.Args[1]

	files, err := filepath.Glob("../../scopedservice/*_generic.go")
	if err != nil {
		panic(err)
	}

	for _, f := range files {
		content, err := ioutil.ReadFile(f)
		if err != nil {
			panic(err)
		}

		content = bytes.Replace(content, []byte(`scopedservice`), []byte(pkgName), -1)
		base := filepath.Base(f)
		if err := ioutil.WriteFile(base, content, 0644); err != nil {
			panic(err)
		}
		println(f, "[scopedservice copier] copied to", base, "for package", pkgName)
	}
}
