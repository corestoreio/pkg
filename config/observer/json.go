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

// +build csall json proto

package observer

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"sync"

	"github.com/corestoreio/errors"
	"github.com/corestoreio/pkg/config"
)

// NewFunc allows to implement a custom observer which gets created based on the
// raw JSON. The function gets called in Configuration.MakeObserver or in
// JSONRegisterObservers.
type NewFunc func(data json.RawMessage) (config.Observer, error)

type customObservers struct {
	sync.RWMutex
	pool map[string]NewFunc
}

var customObserverRegistry = &customObservers{
	pool: make(map[string]NewFunc),
}

// RegisterCustom adds a custom observer to the global registry. A
// custom observer can be accessed via Configuration.MakeObserver or via
// JSONRegisterObservers.
func RegisterCustom(typeName string, fn NewFunc) {
	customObserverRegistry.Lock()
	defer customObserverRegistry.Unlock()
	customObserverRegistry.pool[typeName] = fn
}

// Configuration defines the data retrieved from the outside as JSON to add a
// new observer for a specific route and event.
//easyjson:json
type Configuration struct {
	// Route defines at least three parts: e.g. general/information/store
	Route string `json:"route,omitempty"`
	// Event can be before_set, after_set, before_get or after_get. See
	// config.MakeEvent.
	Event string `json:"event,omitempty"`
	// Type specifies the kind of the observer which should be created. Case
	// sensitive. Supported names are: "ValidateMinMaxInt", "validator",
	// "modificator" and the keys registered via function RegisterCustom.
	Type string `json:"type,omitempty"`
	// Condition contains the JSON object for a type in this package like:
	// `ValidatorArg` or `ValidateMinMaxInt` or TBC.
	Condition json.RawMessage `json:"condition,omitempty"`
}

// Configurations a list of Configuration types.
//easyjson:json
type Configurations struct {
	// Collection list of available validators. The Collection name has been
	// chosen because of protobuf.
	Collection []*Configuration
}

// NewConfigurations creates a new Configurations object.
func NewConfigurations(v ...*Configuration) *Configurations {
	vs := &Configurations{
		Collection: v,
	}
	if v == nil {
		vs.Collection = make([]*Configuration, 0, 5)
	}
	return vs
}

// Validate validates the current slice and used in GRPC middleware.
func (m Configurations) Validate() error {
	for _, v := range m.Collection {
		if err := v.Validate(); err != nil {
			return errors.WithStack(err)
		}
	}
	return nil
}

// Validate checks if the data is confirm to the business logic. Returns nil on success.
// Also used by github.com/grpc-ecosystem/go-grpc-middleware/validator
func (v *Configuration) Validate() error {
	if v == nil {
		return nil
	}
	if err := config.Route(v.Route).IsValid(); err != nil {
		return errors.Wrapf(err, "[config/validation] Invalid route: %#v", v)
	}

	if len(v.Condition) == 0 {
		return errors.Empty.Newf("[config/validation] Data for %q is empty. %#v", v.Type, v)
	}
	if _, err := config.MakeEvent(v.Event); err != nil {
		return errors.WithStack(err)
	}

	// This logic is duplicated but for now ok. No need to add another abstraction.
	switch v.Type {
	case "ValidateMinMaxInt", "validator", "modificator":
		// ok
	default:
		customObserverRegistry.RLock()
		defer customObserverRegistry.RUnlock()
		if _, ok := customObserverRegistry.pool[v.Type]; !ok {
			return errors.NotFound.Newf("[config/validation] Configuration type %q not found in list.", v.Type)
		}
	}
	return nil
}

// MakeEventRoute extracts a validated event and a route from the data.
func (v *Configuration) MakeEventRoute() (event uint8, route string, err error) {
	if event, err = config.MakeEvent(v.Event); err != nil {
		return 0, "", errors.WithStack(err)
	}

	if err := config.Route(v.Route).IsValid(); err != nil {
		return 0, "", errors.Wrapf(err, "[config/validation] Invalid route: %#v", v)
	}

	return event, v.Route, nil
}

// MakeObserver transforms and validates the Configuration data into a functional
// observer for an event and a specific route.
func (v Configuration) MakeObserver() (event uint8, route string, _ config.Observer, err error) {
	if err := v.Validate(); err != nil {
		return 0, "", nil, errors.WithStack(err)
	}

	event, _ = config.MakeEvent(v.Event)

	switch v.Type {
	case "ValidateMinMaxInt":
		mm := new(ValidateMinMaxInt)
		if err := mm.UnmarshalJSON(v.Condition); err != nil {
			return 0, "", nil, errors.BadEncoding.New(err, "[config/validation] Failed to decode: %#v", v)
		}
		return event, v.Route, mm, nil

	case "validator":
		var va ValidatorArg
		if err := va.UnmarshalJSON(v.Condition); err != nil {
			return 0, "", nil, errors.BadEncoding.New(err, "[config/validation] Failed to decode: %#v", v)
		}
		vStr, err := NewValidator(va)
		if err != nil {
			return 0, "", nil, errors.WithStack(err)
		}
		return event, v.Route, vStr, nil

	case "modificator":
		var ma ModificatorArg
		if err := ma.UnmarshalJSON(v.Condition); err != nil {
			return 0, "", nil, errors.BadEncoding.New(err, "[config/validation] Failed to decode: %#v", v)
		}
		vStr, err := NewModificator(ma)
		if err != nil {
			return 0, "", nil, errors.WithStack(err)
		}
		return event, v.Route, vStr, nil
	}

	customObserverRegistry.RLock()
	defer customObserverRegistry.RUnlock()
	newObFn := customObserverRegistry.pool[v.Type]

	co, err := newObFn(v.Condition)
	if err != nil {
		return 0, "", nil, errors.Wrapf(err, "[config/validation] Failed to decode: %#v", v)
	}
	return event, v.Route, co, nil
}

// JSONRegisterObservers reads all JSON byte data from r into memory, parses it,
// creates the appropriate observers and registers them with the config.Service.
func JSONRegisterObservers(or config.ObserverRegisterer, r io.Reader) error {

	jsonData, err := ioutil.ReadAll(r)
	if err != nil {
		return errors.ReadFailed.New(err, "[config/validation/json] Reading failed")
	}
	_ = jsonData // easyjson code generation hack

	vs := NewConfigurations()
	if err := vs.UnmarshalJSON(jsonData); err != nil {
		return errors.BadEncoding.New(err, "[config/validation/json] JSON decoding failed")
	}

	for _, v := range vs.Collection {
		event, route, o, err := v.MakeObserver()
		if err != nil {
			return errors.Wrapf(err, "[config/validation] Data: %#v", v)
		}
		if err := or.RegisterObserver(event, route, o); err != nil {
			return errors.WithStack(err)
		}
	}
	return nil
}

// JSONDeregisterObservers reads all JSON byte data from r into memory, parses it,
// removes the appropriate observers which matches the route and the event.
func JSONDeregisterObservers(or config.ObserverRegisterer, r io.Reader) error {

	jsonData, err := ioutil.ReadAll(r)
	if err != nil {
		return errors.ReadFailed.New(err, "[config/validation/json] Reading failed")
	}
	_ = jsonData // easyjson code generation hack

	vs := NewConfigurations()
	if err := vs.UnmarshalJSON(jsonData); err != nil {
		return errors.BadEncoding.New(err, "[config/validation/json] JSON decoding failed")
	}

	for _, v := range vs.Collection {
		event, route, err := v.MakeEventRoute()
		if err != nil {
			return errors.Wrapf(err, "[config/validation] Data: %#v", v)
		}
		if err := or.DeregisterObserver(event, route); err != nil {
			return errors.WithStack(err)
		}
	}
	return nil
}
