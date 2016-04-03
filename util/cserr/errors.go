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
// If *MultiErr is nil it creates a new pointer and returns it.
func (m *MultiErr) AppendErrors(errs ...error) *MultiErr {
	errNilCount := 0
	for _, err := range errs {
		if err == nil {
			errNilCount++
		}
	}
	if errNilCount == len(errs) {
		return m
	}

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
	return m != nil && len(m.errs) > 0
}

// Contains checks if search has been added to the internal error stack.
func (m *MultiErr) Contains(search error) bool {
	if !m.HasErrors() || search == nil {
		return false
	}
	sUnMasked := UnwrapMasked(search)
	for _, err := range m.errs {
		if err != nil && sUnMasked == UnwrapMasked(err) {
			return true
		}
	}
	return false
}

// Contains checks if search1 can be found within search2 and vice versa.
// One or both parameters can be of type *cserr.MultiErr
func Contains(search1, search2 error) bool {
	search1 = UnwrapMasked(search1)
	search2 = UnwrapMasked(search2)
	if search1 != nil && search1 == search2 {
		return true
	}
	me, ok := search1.(*MultiErr)
	if ok && me.Contains(search2) {
		return true
	}

	// flip the search
	me2, ok2 := search2.(*MultiErr)
	if ok2 && me2.Contains(search1) {
		return true
	}

	if !ok || !ok2 {
		return false
	}

	for _, err := range me.errs {
		if err != nil {
			for _, err2 := range me2.errs {
				err = UnwrapMasked(err)
				err2 = UnwrapMasked(err2)
				if err == err2 {
					return true
				}
			}
		}
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
			if _, err := buf.WriteRune('\n'); err != nil {
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
