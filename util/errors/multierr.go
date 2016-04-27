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

// MultiErr represents a container for collecting and printing multiple errors.
// Mostly used for embedding in functional options.
type MultiErr struct {
	Errors    []error
	Formatter ErrorFormatFunc
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
				m = m.AppendErrors(mErr2.Errors...)
			} else {
				m.Errors = append(m.Errors, err)
			}
		}
	}
	return m
}

// HasErrors checks if Multi contains errors.
func (m *MultiErr) HasErrors() bool {
	return m != nil && len(m.Errors) > 0
}

// Error returns a string where each error has been separated by a line break.
// The location will be added to the output to show you the file name and line number.
func (m *MultiErr) Error() string {
	if !m.HasErrors() {
		return ""
	}
	if m.Formatter == nil {
		return FormatLineFunc(m.Errors)
	}
	return m.Formatter(m.Errors)
}

// MultiErrContains checks if err contains a behavioral error.
// 1st argument err must be of type (*MultiErr) and validate function vf
// at least one of the many Is*() e.g. IsNotValid(), see type BehaviourFunc.
// More than one validate function will be treated as AND hence
// all validate functions must return true.
// If there are multiple behavioral errors and one BehaviourFunc it will stop
// after all errors matches the BehaviourFunc, not at the first match.
func MultiErrContains(err error, bfs ...BehaviourFunc) bool {
	me, ok := err.(*MultiErr)
	if !ok {
		return false
	}

	if len(bfs) == 0 || len(me.Errors) == 0 {
		return false
	}

	var errCount, validCount int
	for _, e := range me.Errors {
		if e != nil {
			errCount++
		}
		for _, f := range bfs {
			if f(e) {
				validCount++
			}
		}
	}
	return validCount == errCount || validCount == len(bfs)
}
