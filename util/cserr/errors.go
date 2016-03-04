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

// MultiErr represents a container for collecting and printing multiple errors.
// Mostly used for embedding in functional options.
type MultiErr struct {
	errs    []error
	details bool
}

// NewMultiErr creates a new multi error struct.
func NewMultiErr(errs ...error) *MultiErr {
	m := new(MultiErr)
	m.AppendErrors(errs...)
	return m
}

// AppendErrors adds multiple errors to the container. Does not add a location.
func (m *MultiErr) AppendErrors(errs ...error) *MultiErr {
	if m == nil {
		m = new(MultiErr)
	}
	for _, err := range errs {
		if err != nil {
			m.errs = append(m.errs, err)
		}
	}
	return m
}

// HasErrors checks if Multi contains errors.
func (m *MultiErr) HasErrors() bool {
	switch {
	case m == nil:
		return false
	case len(m.errs) > 0:
		return true
	}
	return false
}

// VerboseErrors enables more error details like the location. Use in chaining:
// 		e := NewMultiErr(err1, err2).Details()
func (m *MultiErr) VerboseErrors() *MultiErr {
	m.details = true
	return m
}

// Error returns a string where each error has been separated by a line break.
// The location will be added to the output to show you the file name and line number.
// You should use package github.com/juju/errors.
func (m *MultiErr) Error() string {
	if false == m.HasErrors() {
		return ""
	}
	var buf = bufferpool.Get()
	defer bufferpool.Put(buf)

	var details = errDetail
	if m.details {
		details = errors.Details
	}

	le := len(m.errs) - 1
	for i, e := range m.errs {
		if _, err := buf.WriteString(details(e)); err != nil {
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

var errDetail = func(err error) string {
	if err == nil {
		return ""
	}
	return err.Error()
}
