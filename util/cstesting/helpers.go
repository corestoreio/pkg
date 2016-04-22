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

package cstesting

import (
	"os"

	"github.com/corestoreio/csfw/util/cserr"
)

// Fataler describes the function needed to print the output and stop
// the current running goroutine and hence fail the test.
type Fataler interface {
	Fatal(args ...interface{})
}

// ChangeEnv temporarily changes the environment for non parallel tests.
//		defer cstesting.ChangeEnv(t, "PATH", ".")()
func ChangeEnv(t Fataler, key, value string) func() {
	was := os.Getenv(key)
	FatalIfError(t, os.Setenv(key, value))
	return func() {
		FatalIfError(t, os.Setenv(key, was))
	}
}

// ChangeDir temporarily changes the working directory for non parallel tests.
//			defer cstesting.ChangeDir(t, os.TempDir())()
func ChangeDir(t Fataler, dir string) func() {
	wd, err := os.Getwd()
	FatalIfError(t, err)
	FatalIfError(t, os.Chdir(dir))

	return func() {
		FatalIfError(t, os.Chdir(wd))
	}
}

// FatalIfError fails the tests if an unexpected error occurred.
// If the error is gift wrapped prints the location.
func FatalIfError(t Fataler, err error) {
	if err != nil {
		t.Fatal(cserr.NewMultiErr(err).VerboseErrors())
	}
}
