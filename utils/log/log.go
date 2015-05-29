// Copyright 2015 CoreStore Authors
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

package utils

var (
	logger Logger

	nullLog        = &NullLogger{}
	_       Logger = (*NullLogger)(nil)
)

// Logger defines the minimum requirements for logging. See doc.go for more details.
// Interface may be extended ...
type Logger interface {
	Trace(msg string, args ...interface{})
	Debug(msg string, args ...interface{})
	Info(msg string, args ...interface{})
	Warn(msg string, args ...interface{})
	Error(msg string, args ...interface{}) error
	Fatal(msg string, args ...interface{})
	Log(level int, msg string, args []interface{})

	SetLevel(int)
	IsTrace() bool
	IsDebug() bool
	IsInfo() bool
	IsWarn() bool
	// Error, Fatal not needed, those SHOULD always be logged
}

func init() {
	SetNull()
}

// SetNullLogger resets the logger to the null logger aka. black hole.
func SetNull() {
	if logger != nil {
		Warn("SetNullLogger called to reset the logger")
	}
	logger = nullLog
}

// Set sets your preferred Logger. Default Logger is a null-logger. Panics if called twice.
func Set(l Logger) {
	if logger != nullLog {
		panic("logger already initialized")
	}
	logger = l
}

func Trace(msg string, args ...interface{})         { logger.Trace(msg, args...) }
func Debug(msg string, args ...interface{})         { logger.Debug(msg, args...) }
func Info(msg string, args ...interface{})          { logger.Info(msg, args...) }
func Warn(msg string, args ...interface{})          { logger.Warn(msg, args...) }
func Error(msg string, args ...interface{}) error   { return logger.Error(msg, args...) }
func Fatal(msg string, args ...interface{})         { logger.Fatal(msg, args...) }
func Log(level int, msg string, args []interface{}) { logger.Log(level, msg, args) }

func SetLevel(l int) { logger.SetLevel(l) }
func IsTrace() bool  { return logger.IsTrace() }
func IsDebug() bool  { return logger.IsDebug() }
func IsInfo() bool   { return logger.IsInfo() }
func IsWarn() bool   { return logger.IsWarn() }
