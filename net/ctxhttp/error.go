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

package ctxhttp

import (
	"net/http"

	"github.com/corestoreio/csfw/util/cserr"
)

// Error is a HTTP error with a messgae and a status code
type Error struct {
	// Code displays the status code
	Code int
	// Message contains human readable information
	Message string
	// Verbose enables more error details like file name and line number
	// IF the underlying error has been gift wrapped (masked).
	Verbose bool
	errors  []error
}

// NewError creates a new Error. If msg will be passed one time then
// the StatusText of the code will be overridden.
func NewError(code int, msg ...string) *Error {
	e := &Error{Code: code, Message: http.StatusText(code)}
	if len(msg) > 0 && msg[0] != "" {
		e.Message = msg[0]
	}
	return e
}

// NewErrorFromErrors this lovely name describes that it can create
// a HTTP error from multiple error interfaces. The function util.Errors()
// will be used to extract the errors.Locationer interface.
func NewErrorFromErrors(code int, errs ...error) *Error {
	e := &Error{
		Code:    code,
		Message: http.StatusText(code),
		errors:  errs,
	}
	return e
}

// Error returns message.
func (e *Error) Error() string {
	if len(e.errors) > 0 {
		me := cserr.NewMultiErr(e.errors...)
		if e.Verbose {
			me.VerboseErrors()
		}
		e.Message = me.Error()
	}
	return e.Message
}
