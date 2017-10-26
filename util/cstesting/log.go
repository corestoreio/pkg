// Copyright 2015-2017, Cyrill @ Schumacher.fm and the CoreStore contributors
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

package cstesting

import (
	"github.com/corestoreio/log"
)

type logger interface {
	Log(args ...interface{})
}

// NewLogger creates a logger based on the testing.TB.Log function and logs
// level independent for now. Signature of `logger`: Log(args ...interface{})
func NewLogger(l logger, fields ...log.Field) log.Logger {
	return &tLog{
		l:   l,
		ctx: fields,
	}
}

type tLog struct {
	l logger
	// ctx is only set when we act as a child logger
	ctx log.Fields
}

// With returns a new Logger that has this logger's context plus the given
// Fields.
func (l *tLog) With(fields ...log.Field) log.Logger {
	l2 := new(tLog)
	*l2 = *l
	l2.ctx = append(l2.ctx, fields...)
	return l2
}

func (l *tLog) prependCtx(fields log.Fields) log.Fields {
	if ctxl := len(l.ctx); ctxl > 0 {
		all := make(log.Fields, 0, ctxl+len(fields))
		all = append(all, l.ctx...)
		all = append(all, fields...)
		fields = all
	}
	return fields
}

// Debug outputs information for developers including a stack trace.
func (l *tLog) Debug(msg string, fields ...log.Field) {
	l.l.Log("[DEBUG] ", l.prependCtx(fields).ToString(msg))
}

// Info outputs information for users of the app
func (l *tLog) Info(msg string, fields ...log.Field) {
	l.l.Log("[INFO] ", l.prependCtx(fields).ToString(msg))
}

// IsDebug returns true if Debug level is enabled
func (l *tLog) IsDebug() bool {
	return true
}

// IsInfo returns true if Info level is enabled
func (l *tLog) IsInfo() bool {
	return true
}
