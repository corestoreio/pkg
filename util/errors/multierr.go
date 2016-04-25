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
	"bytes"
	"strconv"

	"github.com/corestoreio/csfw/util/bufferpool"
)

// MultiErr represents a container for collecting and printing multiple errors.
// Mostly used for embedding in functional options.
type MultiErr struct {
	errs []error
}

// NewMultiErr creates a new multi error struct.
func NewMultiErr(errs ...error) *MultiErr {
	m := new(MultiErr)
	return m.AppendErrors(errs...)
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
			// unwrap MultiErr recursively because in errs can be a MultiErr
			if mErr2, ok := err.(*MultiErr); ok {
				m = m.AppendErrors(mErr2.errs...)
			} else {
				m.errs = append(m.errs, err)
			}
		}
	}
	return m
}

// HasErrors checks if Multi contains errors.
func (m *MultiErr) HasErrors() bool {
	return m != nil && len(m.errs) > 0
}

// Error returns a string where each error has been separated by a line break.
// The location will be added to the output to show you the file name and line number.
func (m *MultiErr) Error() string {
	if !m.HasErrors() {
		return ""
	}
	var buf = bufferpool.Get()
	defer bufferpool.Put(buf)

	for _, e := range m.errs {
		fprint(buf, e)
	}
	return buf.String()
}

// MultiErrContains checks if err contains a behavioral error.
// 1st argument err must be of type (*MultiErr) and validate function vf
// at least one of the many Is*() e.g. IsNotValid().
// More than one validate function will be treated as AND hence
// all validate functions must return true.
func MultiErrContains(err error, vf ...func(error) bool) bool {
	me, ok := err.(*MultiErr)
	if !ok {
		return false
	}

	if len(vf) == 0 || len(me.errs) == 0 {
		return false
	}

	var ec, valids int
	for _, e := range me.errs {
		if e != nil {
			ec++
		}
		for _, f := range vf {
			if f(e) {
				valids++
			}
		}
	}
	return valids == ec || valids == len(vf)
}

// Fprint prints the error to the supplied writer.
// The format of the output is the same as Print.
// If err is nil, nothing is printed.
func fprint(buf *bytes.Buffer, err error) {
	for err != nil {
		location, ok := err.(locationer)
		if ok {
			file, line := location.Location()
			_, _ = buf.WriteString(file)
			_, _ = buf.WriteRune(':')
			_, _ = buf.WriteString(strconv.Itoa(line))
			_, _ = buf.WriteString(": ")
		}
		switch err := err.(type) {
		case *e:
			_, _ = buf.WriteString(err.message)
			_, _ = buf.WriteRune('\n')
		default:
			_, _ = buf.WriteString(err.Error())
			_, _ = buf.WriteRune('\n')
		}

		cause, ok := err.(causer)
		if !ok {
			break
		}
		err = cause.Cause()
	}
}
