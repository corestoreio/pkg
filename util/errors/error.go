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

package errors

import (
	"fmt"
	"runtime"

	"github.com/corestoreio/csfw/util/bufferpool"
)

const pkgPath = `github.com/corestoreio/csfw/util/errors`

// Errorf creates a new annotated error and records the location that the
// error is created. This should be a drop in replacement for fmt.Errorf.
// No behaviour has been attached to this returned error.
//
// For example:
//    return errors.Errorf("[package name] validation failed: %s", message)
//
func Errorf(format string, args ...interface{}) error {
	pc, _, _, _ := runtime.Caller(1)
	return struct {
		error
		loc
	}{
		fmt.Errorf(format, args...),
		loc(pc),
	}
}

// Wrapf returns an error annotating the cause with a message.
// If cause is nil, Wrap returns nil.
func Wrapf(cause error, format string, args ...interface{}) error {
	if cause == nil {
		return nil
	}
	pc, _, _, _ := runtime.Caller(1)
	return &e{
		cause:   cause,
		message: fmt.Sprintf(format, args...),
		loc:     loc(pc),
	}
}

// PrintLoc prints the error including the location and its stack.
// If the error is nil, returns an empty string.
func PrintLoc(err error) string {
	if err == nil {
		return ""
	}
	var buf = bufferpool.Get()
	defer bufferpool.Put(buf)

	fprint(buf, err)

	return buf.String()
}
