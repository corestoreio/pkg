// Copyright 2015-present, Cyrill @ Schumacher.fm and the CoreStore contributors
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

//go:generate easyjson $GOFILE

package json

import (
	"encoding/json"

	"github.com/corestoreio/errors"
	"github.com/corestoreio/pkg/config"
)

// Validator defines the data retrieved from the outside as JSON to add a new
// validator for a specific route and event.
//easyjson:json
type Validator struct {
	// Route defines at least three parts: e.g. general/information/store
	Route string `json:"route,omitempty"`
	// Event can be before_set, after_set, before_get or after_get. See
	// config.MakeEvent.
	Event string `json:"event,omitempty"`
	// Type name of struct to decode and specifies the type of the validator.
	// Case sensitive.
	Type string `json:"type,omitempty"`
	// Condition contains the JSON object for a type like: MinMaxInt64, UUID, etc
	Condition json.RawMessage `json:"condition,omitempty"`
}

// Validators a list of Validator types.
//easyjson:json
type Validators []*Validator

type unmarshalableObserver interface { // can be used to implement custom types, but later.
	config.Observer
	UnmarshalJSON(data []byte) error
}

func (v *Validator) NewFromJSON() (event uint8, route string, _ config.Observer, err error) {
	var ov unmarshalableObserver
	switch v.Type {
	case "MinMaxInt64":
		ov = new(MinMaxInt64)
	case "UUID":
		ov = new(UUID)
	case "ISO3166Alpha2":
		ov, _ = NewISO3166Alpha2()
	default:
		return 0, "", nil, errors.NotFound.Newf("[config/validation] Validator type %q not found in list.", v.Type)
	}

	if err := config.Route(v.Route).IsValid(); err != nil {
		return 0, "", nil, errors.Wrapf(err, "[config/validation] Invalid route: %#v", v)
	}

	if len(v.Condition) == 0 {
		return 0, "", nil, errors.Empty.Newf("[config/validation] Data for %q is empty. %#v", v.Type, v)
	}

	if err := ov.UnmarshalJSON(v.Condition); err != nil {
		return 0, "", nil, errors.Wrapf(err, "[config/validation] Failed to decode: %#v", v)
	}

	if event, err = config.MakeEvent(v.Event); err != nil {
		return 0, "", nil, errors.WithStack(err)
	}

	return event, v.Route, ov, nil
}

func (v *Validator) MakeEventRouteFromJSON() (event uint8, route string, err error) {

	if err := config.Route(v.Route).IsValid(); err != nil {
		return event, "", errors.Wrapf(err, "[config/validation] Invalid route: %#v", v)
	}

	event, err = config.MakeEvent(v.Event)
	err = errors.WithStack(err)
	route = v.Route
	return
}

func (m *Validator) Size() (n int) {
	var l int
	_ = l
	l = len(m.Route)
	if l > 0 {
		n += 1 + l + sovValidator(uint64(l))
	}
	l = len(m.Event)
	if l > 0 {
		n += 1 + l + sovValidator(uint64(l))
	}
	l = len(m.Type)
	if l > 0 {
		n += 1 + l + sovValidator(uint64(l))
	}
	l = len(m.Condition)
	if l > 0 {
		n += 1 + l + sovValidator(uint64(l))
	}
	return n
}

func sovValidator(x uint64) (n int) {
	for {
		n++
		x >>= 7
		if x == 0 {
			break
		}
	}
	return n
}

func RegisterObserversFromJSON(or config.ObserverRegisterer, jsonData []byte) error {

	vs := make(Validators, 0, 5)
	if err := vs.UnmarshalJSON(jsonData); err != nil {
		return errors.WithStack(err)
	}

	for _, v := range vs {
		event, route, o, err := v.NewFromJSON()
		if err != nil {
			return errors.Wrapf(err, "[config/validation] Data: %#v", v)
		}
		if err := or.RegisterObserver(event, route, o); err != nil {
			return errors.WithStack(err)
		}
	}
	return nil
}
