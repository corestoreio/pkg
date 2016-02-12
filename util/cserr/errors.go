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

package cserr

import (
	"fmt"

	"github.com/corestoreio/csfw/util/bufferpool"
	"github.com/juju/errors"
)

// Multi represents a container for collecting and printing multiple errors.
// Mostly used for embedding in functional options.
type Multi struct {
	errs []error
}

// NewMulti creates a new multi error struct.
func NewMulti(errs ...error) Multi {
	m := Multi{}
	m.AppendErrors(errs...)
	return m
}

// AppendErrors adds multiple errors to the container. Does not add a location.
func (m *Multi) AppendErrors(errs ...error) {
	for _, err := range errs {
		if err != nil {
			m.errs = append(m.errs, err)
		}
	}
}

// HasErrors checks if Multi contains errors.
func (m Multi) HasErrors() bool {
	return false == (len(m.errs) == 0 || (len(m.errs) == 1 && m.errs[0] == nil))
}

// Error returns a string where each error has been separated by a line break.
// The location will be added to the output to show you the file name and line number.
// You should use package github.com/juju/errors.
func (m Multi) Error() string {
	if len(m.errs) == 0 || (len(m.errs) == 1 && m.errs[0] == nil) {
		return ""
	}
	var buf = bufferpool.Get()
	defer bufferpool.Put(buf)

	le := len(m.errs) - 1
	for i, e := range m.errs {
		if _, err := buf.WriteString(errors.Details(e)); err != nil {
			return fmt.Sprintf("buf.WriteString (1) internal error (%s): %s\n%s", err, e, buf.String())
		}

		if i < le {
			if _, err := buf.WriteString("\n"); err != nil {
				return fmt.Sprintf("buf.WriteString (2) internal error (%s):\n%s", err, buf.String())
			}
		}
	}
	return buf.String()
}
