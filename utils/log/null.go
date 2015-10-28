// Copyright 2015, Cyrill @ Schumacher.fm and the CoreStore contributors
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

import "errors"

// Following Code by: https://github.com/mgutz Mario Gutierrez / MIT License

// NullLogger is the default logger for this package.
type NullLogger struct{}

// New returns a new Logger that has this logger's context plus the given context
func (l NullLogger) New(ctx ...interface{}) Logger { return &NullLogger{} }

// Trace logs a trace entry.
func (l NullLogger) Trace(msg string, args ...interface{}) {}

// Debug logs a debug entry.
func (l NullLogger) Debug(msg string, args ...interface{}) {}

// Info logs an info entry.
func (l NullLogger) Info(msg string, args ...interface{}) {}

// Warn logs a warn entry.
func (l NullLogger) Warn(msg string, args ...interface{}) {}

// Error logs an error entry. Returns the first argument as an error OR
// if the 2nd index of args (that is args[1] ;-) ) contains the error
// then that error will be returned.
func (l NullLogger) Error(msg string, args ...interface{}) error {
	lArgs := len(args)
	if lArgs > 0 && lArgs%2 == 0 {
		if err, ok := args[1].(error); ok {
			return err
		}
	}
	return errors.New(msg)
}

// Fatal logs a fatal entry then panics.
func (l NullLogger) Fatal(msg string, args ...interface{}) { panic("exit due to fatal error") }

// Log logs a leveled entry.
func (l NullLogger) Log(level int, msg string, args []interface{}) {}

// IsTrace determines if this logger logs a trace statement.
func (l NullLogger) IsTrace() bool { return false }

// IsDebug determines if this logger logs a debug statement.
func (l NullLogger) IsDebug() bool { return false }

// IsInfo determines if this logger logs an info statement.
func (l NullLogger) IsInfo() bool { return false }

// IsWarn determines if this logger logs a warning statement.
func (l NullLogger) IsWarn() bool { return false }

// SetLevel sets the level of this logger.
func (l NullLogger) SetLevel(level int) {}
