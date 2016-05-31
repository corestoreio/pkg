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

package zapw

import (
	"github.com/corestoreio/csfw/log"
	"github.com/corestoreio/csfw/util/errors"
	"github.com/uber-go/zap"
)

// Wrap creates a new zap Logger. Their options cannot be applied as an argument
// because this interface https://godoc.org/github.com/uber-go/zap#Option has a
// private function :-( So use:
//		l := zapw.Wrap{
//			Zap: zap.NewJSON(zap.Option ... )
// 		}
type Wrap struct {
	Zap zap.Logger
}

// New should create a new zap Logger but returns nil
// because https://godoc.org/github.com/uber-go/zap#Option has a private
// function. Maybe this will change in the futur or someone has another
// idea.
func (l Wrap) New(_ ...interface{}) log.Logger {
	return nil
}

// Fatal exists the app with logging the error
func (l Wrap) Fatal(msg string, fields ...log.Field) {
	l.Zap.Fatal(msg, doFieldWrap(fields...)...)
}

// Info outputs information for users of the app
func (l Wrap) Info(msg string, fields ...log.Field) {
	l.Zap.Info(msg, doFieldWrap(fields...)...)
}

// Debug outputs information for developers.
func (l Wrap) Debug(msg string, fields ...log.Field) {
	l.Zap.Debug(msg, doFieldWrap(fields...)...)
}

// SetLevel sets the log level. Panics on incorrect value
func (l Wrap) SetLevel(lvl int) {
	l.Zap.SetLevel(zap.Level(lvl))
}

// IsDebug returns true if Debug level is enabled
func (l Wrap) IsDebug() bool {
	return l.Zap.Level() <= zap.Debug
}

// IsInfo returns true if Info level is enabled
func (l Wrap) IsInfo() bool {
	return l.Zap.Level() <= zap.Info
}

type zapFieldWrap struct {
	zf []zap.Field
}

func doFieldWrap(fs ...log.Field) []zap.Field {
	fw := &zapFieldWrap{
		zf: make([]zap.Field, 0, len(fs)),
	}

	if err := log.Fields(fs).AddTo(fw); err != nil {
		fw.AddString(log.ErrorKeyName, errors.PrintLoc(err))
	}
	return fw.zf
}

func (se *zapFieldWrap) AddBool(k string, v bool) {
	se.zf = append(se.zf, zap.Bool(k, v))
}
func (se *zapFieldWrap) AddFloat64(k string, v float64) {
	se.zf = append(se.zf, zap.Float64(k, v))
}
func (se *zapFieldWrap) AddInt(k string, v int) {
	se.zf = append(se.zf, zap.Int(k, v))
}
func (se *zapFieldWrap) AddInt64(k string, v int64) {
	se.zf = append(se.zf, zap.Int64(k, v))
}
func (se *zapFieldWrap) AddMarshaler(k string, v log.Marshaler) error {
	if err := v.MarshalLog(se); err != nil {
		se.AddString(log.ErrorKeyName, errors.PrintLoc(err))
	}
	return nil
}
func (se *zapFieldWrap) AddObject(k string, v interface{}) {
	se.zf = append(se.zf, zap.Object(k, v))
}
func (se *zapFieldWrap) AddString(k string, v string) {
	se.zf = append(se.zf, zap.String(k, v))
}

func (se *zapFieldWrap) Nest(key string, f func(log.KeyValuer) error) error {
	// not that nice ...
	se.zf = append(se.zf, zap.String(key, "StartNest"))
	err := errors.Wrap(f(se), "[zapw] Nest")
	se.zf = append(se.zf, zap.String(key, "EndNest"))
	return err
}
