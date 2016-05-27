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

import "time"

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
			if l.IsInfo() {
				l.Info(msg, append(fields, Duration("Duration", time.Since(start)))...)
			}
		},
		Debug: func(msg string, fields ...Field) {
			if l.IsDebug() {
				l.Debug(msg, append(fields, Duration("Duration", time.Since(start)))...)
			}
		},
	}
}
