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

package log15w

import (
	"fmt"
	"sync"

	"github.com/corestoreio/csfw/log"
	"github.com/corestoreio/csfw/util/errors"
	"github.com/inconshreveable/log15"
)

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
func (l *Log15) New(ctx ...interface{}) log.Logger {
	return NewLog15(l.Level, l.Wrap.GetHandler(), ctx...)
}

// Fatal exists the app with logging the error
func (l *Log15) Fatal(msg string, fields ...log.Field) {
	l.Wrap.Crit(msg, doLog15FieldWrap(fields...)...)
}

// Info outputs information for users of the app
func (l *Log15) Info(msg string, fields ...log.Field) {
	l.Wrap.Info(msg, doLog15FieldWrap(fields...)...)
}

// Debug outputs information for developers.
func (l *Log15) Debug(msg string, fields ...log.Field) {
	l.Wrap.Debug(msg, doLog15FieldWrap(fields...)...)
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

var log15IFSlicePool = &sync.Pool{
	New: func() interface{} {
		return &log15FieldWrap{
			ifaces: make([]interface{}, 0, 12), // just guessing not more than 12 args / 6 Fields
		}
	},
}

type log15FieldWrap struct {
	ifaces []interface{}
}

func doLog15FieldWrap(fs ...log.Field) []interface{} {
	fw := log15IFSlicePool.Get().(*log15FieldWrap)
	defer log15IFSlicePool.Put(fw)

	if err := log.Fields(fs).AddTo(fw); err != nil {
		fw.AddString(log.ErrorKeyName, fmt.Sprintf("%+v", err))
	}
	return fw.ifaces
}

func (se *log15FieldWrap) append(key string, val interface{}) {
	se.ifaces = append(se.ifaces, key, val)
}

func (se *log15FieldWrap) AddBool(k string, v bool) {
	se.append(k, v)
}
func (se *log15FieldWrap) AddFloat64(k string, v float64) {
	se.append(k, v)
}
func (se *log15FieldWrap) AddInt(k string, v int) {
	se.append(k, v)
}
func (se *log15FieldWrap) AddInt64(k string, v int64) {
	se.append(k, v)
}
func (se *log15FieldWrap) AddMarshaler(k string, v log.Marshaler) error {
	if err := v.MarshalLog(se); err != nil {
		se.AddString(log.ErrorKeyName, fmt.Sprintf("%+v", err))
	}
	return nil
}
func (se *log15FieldWrap) AddObject(k string, v interface{}) {
	se.append(k, v)
}
func (se *log15FieldWrap) AddString(k string, v string) {
	se.append(k, v)
}

func (se *log15FieldWrap) Nest(key string, f func(log.KeyValuer) error) error {
	se.append(key, "nest")
	return errors.Wrap(f(se), "[log15w] log15FieldWrap.Nest.f")
}
