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

package cfgmodel

import (
	"time"

	"github.com/corestoreio/csfw/config"
	"github.com/corestoreio/csfw/store/scope"
	"github.com/corestoreio/csfw/util/conv"
	"github.com/corestoreio/csfw/util/errors"
)

// Time represents a path in config.Getter which handles time values.
type Time struct{ baseValue }

// NewTime creates a new Time cfgmodel with a given path.
func NewTime(path string, opts ...Option) Time {
	return Time{baseValue: newBaseValue(path, opts...)}
}

// Get returns a time value from ScopedGetter, if empty the
// *Field.Default value will be applied if provided.
// scope.DefaultID will be enforced if *Field.Scopes is empty.
// Get is able to parse available time formats as defined in
// github.com/corestoreio/csfw/util/conv.StringToDate()
func (t Time) Get(sg config.Scoped) (time.Time, scope.Hash, error) {
	// This code must be kept in sync with other Get() functions

	var v time.Time
	var scp = t.initScope().Top()
	if t.Field != nil {
		scp = t.Field.Scopes.Top()
		if d := t.Field.Default; d != nil {
			var err error
			v, err = conv.ToTimeE(d)
			if err != nil {
				return time.Time{}, 0, errors.NewNotValidf("[cfgmodel] ToTimeE: %v", err)
			}
		}
	}

	val, h, err := sg.Time(t.route, scp)
	switch {
	case err == nil: // we found the value in the config service
		v = val
	case !errors.IsNotFound(err):
		err = errors.Wrapf(err, "[cfgmodel] Route %q", t.route)
	default:
		err = nil // a Err(Section|Group|Field)NotFound error and uninteresting, so reset
	}
	return v, h, err
}

// Write writes a time value without validating it against the source.Slice.
func (t Time) Write(w config.Writer, v time.Time, s scope.Scope, scopeID int64) error {
	return t.baseValue.Write(w, v, s, scopeID)
}

// Duration represents a path in config.Getter which handles duration values.
type Duration struct{ Str }

// NewDuration creates a new Duration cfgmodel with a given path.
func NewDuration(path string, opts ...Option) Duration {
	return Duration{Str: NewStr(path, opts...)}
}

// Get returns a duration value from ScopedGetter, if empty the
// *Field.Default value will be applied if provided.
// scope.DefaultID will be enforced if *Field.Scopes is empty.
// A duration string is a possibly signed sequence of
// decimal numbers, each with optional fraction and a unit suffix,
// such as "300ms", "-1.5h" or "2h45m".
// Valid time units are "ns", "us" (or "µs"), "ms", "s", "m", "h".
// Error behaviour: NotValid
func (t Duration) Get(sg config.Scoped) (time.Duration, scope.Hash, error) {
	// This code must be kept in sync with other Get() functions

	var v time.Duration
	var scp = t.initScope().Top()
	if t.Field != nil {
		scp = t.Field.Scopes.Top()
		if d := t.Field.Default; d != nil {
			var err error
			v, err = conv.ToDurationE(d)
			if err != nil {
				return 0, 0, errors.NewNotValidf("[cfgmodel] ToDurationE: %v", err)
			}
		}
	}

	val, h, err := sg.String(t.route, scp)
	switch {
	case err == nil: // we found the value in the config service
		if v, err = conv.ToDurationE(val); err != nil {
			err = errors.NewNotValidf("[cfgmodel] ToDurationE: %v", err)
		}
	case !errors.IsNotFound(err):
		err = errors.Wrapf(err, "[cfgmodel] Route %q", t.route)
	default:
		err = nil // a Err(Section|Group|Field)NotFound error and uninteresting, so reset
	}
	return v, h, err
}

// Write writes a duration value without validating it against the source.Slice.
func (t Duration) Write(w config.Writer, v time.Duration, s scope.Scope, scopeID int64) error {
	return t.Str.Write(w, v.String(), s, scopeID)
}
