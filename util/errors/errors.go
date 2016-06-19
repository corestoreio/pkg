// Copyright (c) 2015, Dave Cheney <dave@cheney.net>
// All rights reserved.
//
// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions are met:
//
// * Redistributions of source code must retain the above copyright notice, this
// list of conditions and the following disclaimer.
//
// * Redistributions in binary form must reproduce the above copyright notice,
// this list of conditions and the following disclaimer in the documentation
// and/or other materials provided with the distribution.
//
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
// AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
// IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
// DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE
// FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL
// DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR
// SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER
// CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY,
// OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
// OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.

package errors

import (
	"fmt"
	"io"
)

// _error is an error implementation returned by New and Errorf
// that implements its own fmt.Formatter.
type _error struct {
	msg string
	*stack
}

func (e _error) Error() string { return e.msg }

func (e _error) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		if s.Flag('+') {
			io.WriteString(s, e.msg)
			fmt.Fprintf(s, "%+v", e.StackTrace())
			return
		}
		fallthrough
	case 's':
		io.WriteString(s, e.msg)
	}
}

// New returns an error with the supplied message.
func New(message string) error {
	return _error{
		message,
		callers(),
	}
}

// Errorf formats according to a format specifier and returns the string
// as a value that satisfies error.
func Errorf(format string, args ...interface{}) error {
	return _error{
		fmt.Sprintf(format, args...),
		callers(),
	}
}

type cause struct {
	cause error
	msg   string
}

func (c cause) Error() string { return fmt.Sprintf("%s: %v", c.msg, c.Cause()) }
func (c cause) Cause() error  { return c.cause }

// wrapper is an error implementation returned by Wrap and Wrapf
// that implements its own fmt.Formatter.
type wrapper struct {
	cause
	*stack
}

func (w wrapper) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		if s.Flag('+') {
			fmt.Fprintf(s, "%+v\n", w.Cause())
			fmt.Fprintf(s, "%+v: %s", w.StackTrace()[0], w.msg)
			return
		}
		fallthrough
	case 's':
		io.WriteString(s, w.Error())
	}
}

// Wrap returns an error annotating err with message.
// If err is nil, Wrap returns nil.
func Wrap(err error, message string) error {
	if err == nil {
		return nil
	}
	return wrapper{
		cause: cause{
			cause: err,
			msg:   message,
		},
		stack: callers(),
	}
}

// Wrapf returns an error annotating err with the format specifier.
// If err is nil, Wrapf returns nil.
func Wrapf(err error, format string, args ...interface{}) error {
	if err == nil {
		return nil
	}
	return wrapper{
		cause: cause{
			cause: err,
			msg:   fmt.Sprintf(format, args...),
		},
		stack: callers(),
	}
}

// Cause returns the underlying cause of the error, if possible.
// An error value has a cause if it implements the following
// interface:
//
//     type Causer interface {
//            Cause() error
//     }
//
// If the error does not implement Cause, the original error will
// be returned. If the error is nil, nil will be returned without further
// investigation.
func Cause(err error) error {
	type causer interface {
		Cause() error
	}

	for err != nil {
		cause, ok := err.(causer)
		if !ok {
			break
		}
		err = cause.Cause()
	}
	return err
}
