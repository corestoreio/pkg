// Copyright 2015-present, Cyrill @ Schumacher.fm and the CoreStore contributors
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
	"io"
	"os"
)

// fataler describes the function needed to print the output and stop the
// current running goroutine and hence fail the test.
type fataler interface {
	Fatalf(format string, args ...interface{})
}

// ChangeEnv temporarily changes the environment for non parallel tests.
//		defer cstesting.ChangeEnv(t, "PATH", ".")()
func ChangeEnv(t fataler, key, value string) func() {
	was := os.Getenv(key)
	fatalIfError(t, os.Setenv(key, value))
	return func() {
		fatalIfError(t, os.Setenv(key, was))
	}
}

// ChangeDir temporarily changes the working directory for non parallel tests.
//			defer cstesting.ChangeDir(t, os.TempDir())()
func ChangeDir(t fataler, dir string) func() {
	wd, err := os.Getwd()
	fatalIfError(t, err)
	fatalIfError(t, os.Chdir(dir))

	return func() {
		fatalIfError(t, os.Chdir(wd))
	}
}

// fatalIfError fails the tests if an unexpected error occurred. If the error is
// gift wrapped prints the location.
func fatalIfError(t fataler, err error) {
	if err != nil {
		if t != nil {
			t.Fatalf("%+v", err)
		} else {
			panic(err)
		}
	}
}

// Close for usage in conjunction with defer.
// 		defer cstesting.Close(t, con)
func Close(t errorFormatter, c io.Closer) {
	if err := c.Close(); err != nil {
		t.Errorf("%+v", err)
	}
}
