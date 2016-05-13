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

import "github.com/inconshreveable/log15"

type Log15 struct {
	Level log15.Lvl
	Wrap  log15.Logger
}

// NewLog15 creates a new https://godoc.org/github.com/inconshreveable/log15 logger.
func NewLog15(lvl log15.Lvl, h log15.Handler, ctx ...interface{}) *Log15 {
	l := &Log15{
		Level: lvl,
		Wrap:  log15.New(ctx...),
	}
	l.Wrap.SetHandler(h)
	return l
}

// New creates a new logger with the same level as its parent.
func (l *Log15) New(ctx ...interface{}) Logger {
	return NewLog15(l.Level, l.Wrap.GetHandler(), ctx...)
}

// Fatal exists the app with logging the error
func (l *Log15) Fatal(msg string, args ...interface{}) {
	l.Wrap.Crit(msg, args...)
}

// Info outputs information for users of the app
func (l *Log15) Info(msg string, args ...interface{}) {
	l.Wrap.Info(msg, args...)
}

// Debug outputs information for developers including a strack trace.
func (l *Log15) Debug(msg string, args ...interface{}) {
	l.Wrap.Debug(msg, args...)
}

// SetLevel sets the log level. Panics on incorrect value
func (l *Log15) SetLevel(lvl int) {
	l.Level = log15.Lvl(lvl)
	_, _ = log15.LvlFromString(l.Level.String()) // check for valid setting and panic maybe
}

// IsDebug returns true if Debug level is enabled
func (l *Log15) IsDebug() bool {
	return l.Level >= log15.LvlDebug
}

// IsInfo returns true if Info level is enabled
func (l *Log15) IsInfo() bool {
	return l.Level >= log15.LvlInfo
}
