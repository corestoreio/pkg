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

import (
	"errors"
	"time"
)

var (
	ErrLoggerSet = errors.New("Logger already initialized")
	logger       Logger

	nullLog        = &NullLogger{}
	_       Logger = (*NullLogger)(nil)
	_       Logger = (*StdLogger)(nil)
)

// Logger defines the minimum requirements for logging. See doc.go for more details.
type Logger interface {
	// New returns a new Logger that has this logger's context plus the given context
	New(ctx ...interface{}) Logger

	// Debug outputs information for developers including a strack trace.
	Debug(msg string, args ...interface{})
	// Info outputs information for users of the app
	Info(msg string, args ...interface{})

	// Fatal exists the app with logging the error
	Fatal(msg string, args ...interface{})

	// SetLevel sets the global log level
	SetLevel(int)
	// IsDebug returns true if Debug level is enabled
	IsDebug() bool
	// IsInfo returns true if Info level is enabled
	IsInfo() bool
}

func init() {
	SetNull()
}

// SetNullLogger resets the logger to the null logger aka. black hole.
func SetNull() {
	logger = nullLog
}

// Set sets your preferred Logger to be used in CoreStore. Default Logger is
// a null-logger. Panics if called twice.
func Set(l Logger) {
	if logger != nullLog {
		panic(ErrLoggerSet)
	}
	logger = l
}

func Debug(msg string, args ...interface{}) { logger.Debug(msg, args...) }
func Info(msg string, args ...interface{})  { logger.Info(msg, args...) }
func Fatal(msg string, args ...interface{}) { logger.Fatal(msg, args...) }

func SetLevel(l int) { logger.SetLevel(l) }
func IsDebug() bool  { return logger.IsDebug() }
func IsInfo() bool   { return logger.IsInfo() }

// Deferred defines a logger type which can be used to trace the duration.
// Usage:
//		function main(){
//			var PkgLog = log.NewStdLogger()
// 			defer log.WhenDone(PkgLog).Info("Stats", "Package", "main")
//			...
// 		}
// Outputs the duration for the main action.
type Deferred struct {
	Info  func(msg string, args ...interface{})
	Debug func(msg string, args ...interface{})
}

// WhenDone returns a Logger which tracks the duration
func WhenDone(l Logger) Deferred {
	// @see http://play.golang.org/p/K53LV16F9e from @francesc
	start := time.Now()
	return Deferred{
		Info: func(msg string, args ...interface{}) {
			if l.IsInfo() {
				l.Info(msg, append(args, "Duration", time.Since(start).String())...)
			}
		},
		Debug: func(msg string, args ...interface{}) {
			if l.IsDebug() {
				l.Debug(msg, append(args, "Duration", time.Since(start).String())...)
			}
		},
	}
}
