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
	"io"
	"io/ioutil"
	"sync"

	"github.com/corestoreio/errors"
	"github.com/corestoreio/pkg/config"
	"github.com/corestoreio/pkg/config/validation"
)

// RegisterObservers reads JSON byte data from r, parses it, creates the
// appropriate observers and registers them with the config.Service.
func RegisterObservers(or config.ObserverRegisterer, r io.Reader) error {

	jsonData, err := ioutil.ReadAll(r)
	if err != nil {
		return errors.ReadFailed.New(err, "[config/validation/json] Reading failed")
	}

	vs := make(Validators, 0, 5)
	if err := vs.UnmarshalJSON(jsonData); err != nil {
		return errors.BadEncoding.New(err, "[config/validation/json] JSON decoding failed")
	}

	for _, v := range vs {
		event, route, o, err := v.makeObserver()
		if err != nil {
			return errors.Wrapf(err, "[config/validation] Data: %#v", v)
		}
		if err := or.RegisterObserver(event, route, o); err != nil {
			return errors.WithStack(err)
		}
	}
	return nil
}

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
	// Case sensitive. Supported names are `strings` or `min_max_int64` or TBC.
	Type string `json:"type,omitempty"`
	// Condition contains the JSON object for a type in this package like:
	// `Strings` or `MinMaxInt64` or TBC.
	Condition json.RawMessage `json:"condition,omitempty"`
}

// Validators a list of Validator types.
//easyjson:json
type Validators []*Validator

func (v *Validator) makeObserver() (event uint8, route string, _ config.Observer, err error) {
	if err := config.Route(v.Route).IsValid(); err != nil {
		return 0, "", nil, errors.Wrapf(err, "[config/validation] Invalid route: %#v", v)
	}

	if len(v.Condition) == 0 {
		return 0, "", nil, errors.Empty.Newf("[config/validation] Data for %q is empty. %#v", v.Type, v)
	}

	if event, err = config.MakeEvent(v.Event); err != nil {
		return 0, "", nil, errors.WithStack(err)
	}

	switch v.Type {
	case "min_max_int64", "minmaxint64", "MinMaxInt64":
		mm := new(MinMaxInt64)
		if err := mm.UnmarshalJSON(v.Condition); err != nil {
			return 0, "", nil, errors.BadEncoding.New(err, "[config/validation] Failed to decode: %#v", v)
		}
		return event, v.Route, mm.MinMaxInt64, nil

	case "strings", "Strings":
		str := new(Strings)
		if err := str.UnmarshalJSON(v.Condition); err != nil {
			return 0, "", nil, errors.BadEncoding.New(err, "[config/validation] Failed to decode: %#v", v)
		}
		vStr, err := validation.NewStrings(str.Strings)
		if err != nil {
			return 0, "", nil, errors.WithStack(err)
		}
		return event, v.Route, vStr, nil

	default:
		customObserverRegistry.RLock()
		defer customObserverRegistry.RUnlock()
		if uo, ok := customObserverRegistry.pool[v.Type]; ok {
			if err := uo.UnmarshalJSON(v.Condition); err != nil {
				return 0, "", nil, errors.BadEncoding.New(err, "[config/validation] Failed to decode: %#v", v)
			}
			return event, v.Route, uo, nil
		}
		return 0, "", nil, errors.NotFound.Newf("[config/validation] Validator type %q not found in list.", v.Type)
	}
}

// UnmarshallableObserver allows to implement custom validation for function
// RegisterObservers. Don't abuse this >:-|
type UnmarshallableObserver interface {
	json.Unmarshaler
	config.Observer
}

type customObservers struct {
	sync.RWMutex
	pool map[string]UnmarshallableObserver
}

var customObserverRegistry = &customObservers{
	pool: make(map[string]UnmarshallableObserver),
}

// RegisterCustomObserver adds a custom observer to the global registry.
func RegisterCustomObserver(typeName string, uo UnmarshallableObserver) {
	customObserverRegistry.Lock()
	defer customObserverRegistry.Unlock()
	customObserverRegistry.pool[typeName] = uo
}
