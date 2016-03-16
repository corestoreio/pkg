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
	"github.com/juju/errors"
)

// Time represents a path in config.Getter which handles time values.
type Time struct{ baseValue }

// NewTime creates a new Time cfgmodel with a given path.
func NewTime(path string, opts ...Option) Time {
	return Time{baseValue: NewValue(path, opts...)}
}

// Get returns a time value from ScopedGetter, if empty the
// *Field.Default value will be applied if provided.
// scope.DefaultID will be enforced if *Field.Scopes is empty.
// Get is able to parse available time formats as defined in
// github.com/corestoreio/csfw/util/conv.StringToDate()
func (t Time) Get(sg config.ScopedGetter) (time.Time, error) {
	// This code must be kept in sync with other Get() functions

	var v time.Time
	var scp = scope.DefaultID
	if t.Field != nil {
		scp = t.Field.Scopes.Top()
		var err error
		v, err = conv.ToTimeE(t.Field.Default)
		if err != nil {
			return time.Time{}, errors.Mask(err)
		}
	}

	val, err := sg.Time(t.route, scp)
	switch {
	case err == nil: // we found the value in the config service
		v = val
	case config.NotKeyNotFoundError(err):
		err = errors.Maskf(err, "Route %s", t.route)
	default:
		err = nil // a Err(Section|Group|Field)NotFound error and uninteresting, so reset
	}
	return v, err
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
func (t Duration) Get(sg config.ScopedGetter) (time.Duration, error) {
	// This code must be kept in sync with other Get() functions

	var v time.Duration
	var scp = scope.DefaultID
	if t.Field != nil {
		scp = t.Field.Scopes.Top()
		var err error
		v, err = conv.ToDurationE(t.Field.Default)
		if err != nil {
			return 0, errors.Mask(err)
		}
	}

	val, err := sg.String(t.route, scp)
	switch {
	case err == nil: // we found the value in the config service
		v, err = conv.ToDurationE(val)
	case config.NotKeyNotFoundError(err):
		err = errors.Maskf(err, "Route %s", t.route)
	default:
		err = nil // a Err(Section|Group|Field)NotFound error and uninteresting, so reset
	}
	return v, err
}

// Write writes a duration value without validating it against the source.Slice.
func (t Duration) Write(w config.Writer, v time.Duration, s scope.Scope, scopeID int64) error {
	return t.Str.Write(w, v.String(), s, scopeID)
}
