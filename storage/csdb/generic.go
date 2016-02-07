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

package csdb

import (
	"fmt"

	"github.com/corestoreio/csfw/storage/dbr"
	"github.com/juju/errgo"
)

var _ error = (*Error)(nil)

// Error type with more options
type Error struct {
	msg string
}

// NewError creates a new error from a string
func NewError(text string) *Error {
	return &Error{
		msg: text,
	}
}

// NewErrorf creates a new formatted error. Same arguments as with fmt.Sprintf
func NewErrorf(format string, a ...interface{}) *Error {
	return NewError(fmt.Sprintf(format, a...))
}

// Error returns the error string
func (e *Error) Error() string {
	return e.msg
}

// LoadSlice loads the slice dest with the table structure from tsr TableStructurer and table index ti.
// Returns the number of loaded rows and nil or 0 and an error. Slice must be a pointer to structs.
func LoadSlice(dbrSess dbr.SessionRunner, tsr TableManager, ti Index, dest interface{}, cbs ...dbr.SelectCb) (int, error) {
	ts, err := tsr.Structure(ti)
	if err != nil {
		return 0, errgo.Mask(err)
	}

	sb, err := ts.Select(dbrSess)
	if err != nil {
		return 0, errgo.Mask(err)
	}

	for _, cb := range cbs {
		if cb != nil {
			sb = cb(sb)
		}
	}
	return sb.LoadStructs(dest)
}
