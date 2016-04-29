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

// ErrorFormatFunc is a function callback that is called by Error to
// turn the list of errors into a string.
type ErrorFormatFunc func([]error) string

// FormatLineFunc is a basic formatter that outputs the errors
// that occurred along with the filename and the line number of the errors.
// Only if the error has a Location() function.
func FormatLineFunc(errs []error) string {
	var buf = bufferpool.Get()
	defer bufferpool.Put(buf)

	for _, e := range errs {
		fprint(buf, e)
	}
	return buf.String()
}

// Fprint prints the error to the supplied writer.
// The format of the output is the same as Print.
// If err is nil, nothing is printed.
func fprint(buf *bytes.Buffer, err error) {
	type location interface {
		Location() (string, int)
	}
	type message interface {
		Message() string
	}

	for err != nil {
		location, ok := err.(location)
		if ok {
			file, line := location.Location()
			_, _ = buf.WriteString(file)
			_, _ = buf.WriteRune(':')
			_, _ = buf.WriteString(strconv.Itoa(line))
			_, _ = buf.WriteString(": ")
		}
		switch err := err.(type) {
		case message:
			_, _ = buf.WriteString(err.Message())
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
