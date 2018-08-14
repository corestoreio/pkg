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

package modification

import (
	"encoding/json"
	"io"

	"github.com/corestoreio/errors"
	"github.com/corestoreio/pkg/config"
)

// Observer defines the data retrieved from the outside as JSON to add a new
// modificator for a specific route and event.
//easyjson:json
type Observer struct {
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

// Observers a list of Validator types.
//easyjson:json
type Observers struct {
	Collection []*Observer
}

func (o *Observer) Validate() error {

	return nil
}

// Validate validates the current slice and used in GRPC middleware.
func (m Observers) Validate() error {
	for _, v := range m.Collection {
		if err := v.Validate(); err != nil {
			return errors.WithStack(err)
		}
	}
	return nil
}

// RegisterObservers reads all JSON byte data from r into memory, parses it,
// creates the appropriate observers and registers them with the config.Service.
func RegisterObservers(or config.ObserverRegisterer, r io.Reader) error {

	// jsonData, err := ioutil.ReadAll(r)
	// if err != nil {
	// 	return errors.ReadFailed.New(err, "[config/validation/json] Reading failed")
	// }
	//
	// vs := make(Observers, 0, 5)
	// if err := vs.UnmarshalJSON(jsonData); err != nil {
	// 	return errors.BadEncoding.New(err, "[config/validation/json] JSON decoding failed")
	// }
	//
	// for _, v := range vs {
	// 	event, route, o, err := v.MakeObserver()
	// 	if err != nil {
	// 		return errors.Wrapf(err, "[config/validation] Data: %#v", v)
	// 	}
	// 	if err := or.RegisterObserver(event, route, o); err != nil {
	// 		return errors.WithStack(err)
	// 	}
	// }
	return nil
}

// DeregisterObservers reads all JSON byte data from r into memory, parses it,
// removes the appropriate observers which matches the route and the event.
func DeregisterObservers(or config.ObserverRegisterer, r io.Reader) error {

	// jsonData, err := ioutil.ReadAll(r)
	// if err != nil {
	// 	return errors.ReadFailed.New(err, "[config/validation/json] Reading failed")
	// }
	//
	// vs := make(Observers, 0, 5)
	// if err := vs.UnmarshalJSON(jsonData); err != nil {
	// 	return errors.BadEncoding.New(err, "[config/validation/json] JSON decoding failed")
	// }
	//
	// for _, v := range vs {
	// 	event, route, err := v.MakeEventRoute()
	// 	if err != nil {
	// 		return errors.Wrapf(err, "[config/validation] Data: %#v", v)
	// 	}
	// 	if err := or.DeregisterObserver(event, route); err != nil {
	// 		return errors.WithStack(err)
	// 	}
	// }
	return nil
}
