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

// BlackHole logs and does nothing. An empty struct. IsDebug() and IsInfo() will
// always return true.
type BlackHole struct{}

// New returns a new Logger that has this logger's context plus the given context
func (l BlackHole) New(ctx ...interface{}) Logger { return BlackHole{} }

// Debug logs a debug entry.
func (l BlackHole) Debug(msg string, args ...interface{}) {}

// Info logs an info entry.
func (l BlackHole) Info(msg string, args ...interface{}) {}

// Fatal logs a fatal entry then panics.
func (l BlackHole) Fatal(msg string, args ...interface{}) { panic("exit due to fatal error: " + msg) }

// IsDebug determines if this logger logs a debug statement. Returns always true.
func (l BlackHole) IsDebug() bool { return true }

// IsInfo determines if this logger logs an info statement. Returns always true.
func (l BlackHole) IsInfo() bool { return true }

// SetLevel sets the level of this logger.
func (l BlackHole) SetLevel(level int) {}
