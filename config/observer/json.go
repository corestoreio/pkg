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

// +build csall json http proto

//go:generate easyjson -build_tags "csall json http proto" $GOFILE

package observer

import (
	"encoding/json"
	"io"
	"io/ioutil"

	"github.com/corestoreio/errors"
	"github.com/corestoreio/pkg/config"
)

func init() {
	RegisterFactory("validateMinMaxInt", func(rawJSON []byte) (config.Observer, error) {
		mm := new(ValidateMinMaxInt)
		if err := mm.UnmarshalJSON(rawJSON); err != nil {
			return nil, errors.BadEncoding.New(err, "[config/observer] Failed to decode: %q", string(rawJSON))
		}
		return mm, nil
	})

	RegisterFactory("validator", func(rawJSON []byte) (config.Observer, error) {
		var va ValidatorArg
		if err := va.UnmarshalJSON(rawJSON); err != nil {
			return nil, errors.BadEncoding.New(err, "[config/observer] Failed to decode: %q", rawJSON)
		}
		o, err := NewValidator(va)
		return o, errors.WithStack(err)
	})

	RegisterFactory("modifier", func(rawJSON []byte) (config.Observer, error) {
		var ma ModifierArg
		if err := ma.UnmarshalJSON(rawJSON); err != nil {
			return nil, errors.BadEncoding.New(err, "[config/observer] Failed to decode: %q", string(rawJSON))
		}
		o, err := NewModifier(ma)
		return o, errors.WithStack(err)
	})
}

// Configuration defines the data retrieved from the outside as JSON to add a
// new observer for a specific route and event. For example an HTTP requests
// contains this Configurtion data.
//easyjson:json
type Configuration struct {
	// Route defines the route for which an event gets dispatched and hence the
	// observer triggered. The route can be a route prefix (e.g. `general`) or a
	// full route (e.g. `general/information/store`) In the above case the
	// observer for prefix route `general` gets dispatched every time a route
	// with that prefix gets called.
	Route string `json:"route,omitempty"`
	// Event can be before_set, after_set, before_get or after_get. See
	// config.MakeEvent.
	Event string `json:"event,omitempty"`
	// Type specifies the kind of the observer which should be created. Case
	// sensitive. Supported names are: "ValidateMinMaxInt", "validator",
	// "modifier" and the keys registered via function RegisterFactory.
	Type string `json:"type,omitempty"`
	// Condition contains the JSON object for a type in this package like:
	// `ValidatorArg` or `ValidateMinMaxInt` or TBC.
	// Depends on function RegisterFactory.
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
func NewConfigurations(c ...*Configuration) *Configurations {
	vs := &Configurations{
		Collection: c,
	}
	if c == nil {
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
func (m *Configuration) Validate() error {
	if m == nil {
		return nil
	}
	if m.Route == "" {
		return errors.Empty.Newf("[config/observer] Route is empty for type %q", m.Type)
	}

	if len(m.Condition) == 0 {
		return errors.Empty.Newf("[config/observer] Data for %q is empty. %#v", m.Type, m)
	}
	if _, err := config.MakeEvent(m.Event); err != nil {
		return errors.WithStack(err)
	}

	if _, ok := lookupFactory(m.Type); !ok {
		return errors.NotFound.Newf("[config/observer] Configuration type %q not found in list %v.", m.Type, availableFactories())
	}
	return nil
}

// MakeEventRoute extracts a validated event and a route from the data.
func (m *Configuration) MakeEventRoute() (event uint8, route string, err error) {
	if event, err = config.MakeEvent(m.Event); err != nil {
		return 0, "", errors.WithStack(err)
	}

	if m.Route == "" {
		return 0, "", errors.Empty.Newf("[config/observer] Route is empty for type %q", m.Type)
	}

	return event, m.Route, nil
}

// MakeObserver transforms and validates the Configuration data into a functional
// observer for an event and a specific route.
func (m *Configuration) MakeObserver() (event uint8, route string, _ config.Observer, err error) {
	if err := m.Validate(); err != nil {
		return 0, "", nil, errors.WithStack(err)
	}

	if newObsFn, ok := lookupFactory(m.Type); ok {
		co, err := newObsFn(m.Condition)
		if err != nil {
			return 0, "", nil, errors.Wrapf(err, "[config/observer] Failed to decode: %q Route: %q Condition: %q", m.Type, m.Route, string(m.Condition))
		}
		event, _ = config.MakeEvent(m.Event)
		return event, m.Route, co, nil
	}
	return 0, "", nil, errors.Fatal.Newf("[config/observer] A programmer made an error. This can never happen.")
}

// RegisterWithJSON reads all JSON byte data from r into memory, parses it,
// creates the appropriate observers and registers them with the config.Service.
// The data in io.Reader must have the structure of type `Configurations`.
func RegisterWithJSON(or config.ObserverRegisterer, r io.Reader) error {

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
			return errors.Wrapf(err, "[config/observer] Data: %#v", v)
		}
		if err := or.RegisterObserver(event, route, o); err != nil {
			return errors.WithStack(err)
		}
	}
	return nil
}

// DeregisterWithJSON reads all JSON byte data from r into memory, parses it,
// removes the appropriate observers which matches the route and the event.
// The data in io.Reader must have the structure of type `Configurations`.
func DeregisterWithJSON(or config.ObserverRegisterer, r io.Reader) error {

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
			return errors.Wrapf(err, "[config/observer] Data: %#v", v)
		}
		if err := or.DeregisterObserver(event, route); err != nil {
			return errors.WithStack(err)
		}
	}
	return nil
}
