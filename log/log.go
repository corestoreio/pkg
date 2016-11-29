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

package log

import (
	"fmt"
	"time"

	"github.com/corestoreio/csfw/util/errors"
)

// ErrorKeyName whenever an error occurs during marshaling that is the
// official key name in the log stream.
const ErrorKeyName = `error`

// Logger defines the minimum requirements for logging. See doc.go for more details.
type Logger interface {
	// New returns a new Logger that has this logger's context plus the given context
	New(ctx ...interface{}) Logger

	// Debug outputs information for developers including a strack trace.
	Debug(msg string, fields ...Field)
	// Info outputs information for users of the app
	Info(msg string, fields ...Field)

	// Fatal exists the app with logging the error
	Fatal(msg string, fields ...Field)

	// SetLevel sets the global log level
	SetLevel(int)
	// IsDebug returns true if Debug level is enabled
	IsDebug() bool
	// IsInfo returns true if Info level is enabled
	IsInfo() bool
}

// AssignmentChar represents the assignment character between key-value pairs
var assignmentChar = ": "

// Separator is the separator to use between key value pairs
var separator = " "

// WriteTypes satisfies the interface KeyValuer. It uses under the hood the
// function Sprintf("%#v", val) to print the values. This costs performance.
type WriteTypes struct {
	// AssignmentChar represents the assignment character between key-value pairs
	AssignmentChar string
	// Separator is the separator to use between key value pairs
	Separator string
	// W used as writer. Must be a pointer. At the moment returned errors
	// are getting ignored.
	W interface {
		WriteString(s string) (n int, err error)
	}
}

// todo: check for errors in WriteString and remove _,_=

func (wt WriteTypes) stdSetKV(key string, value interface{}) {
	if wt.Separator == "" {
		wt.Separator = separator
	}
	_, _ = wt.W.WriteString(wt.Separator)
	if key == "" {
		key = "_"
	}
	_, _ = wt.W.WriteString(key)
	if wt.AssignmentChar == "" {
		wt.AssignmentChar = assignmentChar
	}
	_, _ = wt.W.WriteString(wt.AssignmentChar)
	_, _ = wt.W.WriteString(fmt.Sprintf("%#v", value)) // can be refactored into the different functions
}

func (wt WriteTypes) AddBool(key string, value bool) {
	wt.stdSetKV(key, value)
}
func (wt WriteTypes) AddFloat64(key string, value float64) {
	wt.stdSetKV(key, value)
}
func (wt WriteTypes) AddInt(key string, value int) {
	wt.stdSetKV(key, value)
}
func (wt WriteTypes) AddInt64(key string, value int64) {
	wt.stdSetKV(key, value)
}
func (wt WriteTypes) AddMarshaler(key string, value Marshaler) error {
	if err := value.MarshalLog(wt); err != nil {
		if wt.Separator == "" {
			wt.Separator = separator
		}
		_, _ = wt.W.WriteString(wt.Separator)
		_, _ = wt.W.WriteString(ErrorKeyName)
		if wt.AssignmentChar == "" {
			wt.AssignmentChar = assignmentChar
		}
		_, _ = wt.W.WriteString(wt.AssignmentChar)
		_, _ = wt.W.WriteString(fmt.Sprintf("%+v", err))
	}
	return nil
}
func (wt WriteTypes) AddObject(key string, value interface{}) {
	wt.stdSetKV(key, value)
}
func (wt WriteTypes) AddString(key string, value string) {
	wt.stdSetKV(key, value)
}

// Nest allows the caller to populate a nested object under the provided key.
func (wt WriteTypes) Nest(key string, f func(KeyValuer) error) error {
	if wt.Separator == "" {
		wt.Separator = separator
	}
	_, _ = wt.W.WriteString(wt.Separator)
	if key == "" {
		key = "_"
	}
	_, _ = wt.W.WriteString(key)
	if wt.AssignmentChar == "" {
		wt.AssignmentChar = assignmentChar
	}
	_, _ = wt.W.WriteString(wt.AssignmentChar)
	return errors.Wrap(f(wt), "[log] WriteType.Nest.f")
}

// Deferred defines a logger type which can be used to trace the duration.
// Usage:
//		function main(){
//			var PkgLog = log.NewStdLog()
// 			defer log.WhenDone(PkgLog).Info("Stats", log.String("Package", "main"))
//			...
// 		}
// Outputs the duration for the main action.
type Deferred struct {
	Info  func(msg string, fields ...Field)
	Debug func(msg string, fields ...Field)
}

// WhenDone returns a Logger which tracks the duration
func WhenDone(l Logger) Deferred {
	// @see http://play.golang.org/p/K53LV16F9e from @francesc
	start := time.Now()
	return Deferred{
		Info: func(msg string, fields ...Field) {
			l.Info(msg, append(fields, Duration("Duration", time.Since(start)))...)
		},
		Debug: func(msg string, fields ...Field) {
			l.Debug(msg, append(fields, Duration("Duration", time.Since(start)))...)
		},
	}
}
