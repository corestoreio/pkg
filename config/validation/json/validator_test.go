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

package json

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/alecthomas/assert"
	"github.com/corestoreio/errors"
	"github.com/corestoreio/pkg/config"
	"github.com/corestoreio/pkg/config/validation"
	uv "github.com/corestoreio/pkg/util/validation"
)

var (
	_ uv.Validator = (*Validators)(nil)
	_ uv.Validator = (*Validator)(nil)
)

type observerRegistererFake struct {
	t             *testing.T
	wantEvent     uint8
	wantRoute     string
	wantValidator interface{}
	err           error
}

func (orf observerRegistererFake) RegisterObserver(event uint8, route string, o config.Observer) error {
	if orf.err != nil {
		return orf.err
	}
	assert.Exactly(orf.t, orf.wantEvent, event, "Event should be equal")
	assert.Exactly(orf.t, orf.wantRoute, route, "Route should be equal")
	// Pointers are different in the final objects hence they get printed and
	// their structure compared, not the address.
	assert.Exactly(orf.t, fmt.Sprintf("%#v", orf.wantValidator), fmt.Sprintf("%#v", o), "Observer internal types should match")

	return nil
}

func (orf observerRegistererFake) DeregisterObserver(event uint8, route string) error {
	if orf.err != nil {
		return orf.err
	}
	assert.Exactly(orf.t, orf.wantEvent, event, "Event should be equal")
	assert.Exactly(orf.t, orf.wantRoute, route, "Route should be equal")

	return nil
}

func TestDeregisterObservers(t *testing.T) {
	t.Parallel()

	t.Run("JSON malformed", func(t *testing.T) {
		or := observerRegistererFake{
			t: t,
		}

		err := DeregisterObservers(or, bytes.NewBufferString(`[{ 
			"event":before_set, "route":"payment/pp/port", "type":"MinMaxInt64", "condition":{"conditions":[8080,8090]} 
		}]`))
		assert.True(t, errors.BadEncoding.Match(err), "%+v", err)
	})
	t.Run("event not found", func(t *testing.T) {
		or := observerRegistererFake{
			t: t,
		}
		err := DeregisterObservers(or, bytes.NewBufferString(`[{ 
			"event":"before_heck", "route":"payment/pp/port", "type":"MinMaxInt64", "condition":{"conditions":[8080,8090]} 
		}]`))
		assert.True(t, errors.NotFound.Match(err), "%+v", err)
	})

	t.Run("deregistered", func(t *testing.T) {
		or := observerRegistererFake{
			t:         t,
			wantEvent: config.EventOnBeforeSet,
			wantRoute: "payment/pp/port",
		}
		err := DeregisterObservers(or, bytes.NewBufferString(`[{ 
			"event":"before_set", "route":"payment/pp/port" 
		}]`))
		assert.NoError(t, err)
	})

}

func TestRegisterObservers(t *testing.T) {
	t.Parallel()

	t.Run("RegisterObservers JSON malformed", func(t *testing.T) {
		or := observerRegistererFake{
			t: t,
		}

		err := RegisterObservers(or, bytes.NewBufferString(`[{ 
			"event":before_set, "route":"payment/pp/port", "type":"MinMaxInt64", "condition":{"conditions":[8080,8090]} 
		}]`))
		assert.True(t, errors.BadEncoding.Match(err), "%+v", err)
	})

	t.Run("MinMaxInt64 OK", func(t *testing.T) {
		or := observerRegistererFake{
			t:             t,
			wantEvent:     config.EventOnBeforeSet,
			wantRoute:     "payment/pp/port",
			wantValidator: validation.MinMaxInt64{Conditions: []int64{8080, 8090}},
		}

		err := RegisterObservers(or, bytes.NewBufferString(`[{ 
			"event":"before_set", "route":"payment/pp/port", "type":"MinMaxInt64", "condition":{"conditions":[8080,8090]} 
		}]`))
		assert.NoError(t, err)
	})
	t.Run("MinMaxInt64 Empty conditions", func(t *testing.T) {
		or := observerRegistererFake{
			t:             t,
			wantEvent:     config.EventOnBeforeSet,
			wantRoute:     "payment/pp/port",
			wantValidator: validation.MinMaxInt64{Conditions: []int64{}},
		}

		err := RegisterObservers(or, bytes.NewBufferString(`[{ 
			"event":"before_set", "route":"payment/pp/port", "type":"MinMaxInt64", "condition":{"conditions":[]} 
		}]`))
		assert.NoError(t, err)
	})
	t.Run("MinMaxInt64 empty condition", func(t *testing.T) {
		or := observerRegistererFake{
			t: t,
		}
		err := RegisterObservers(or, bytes.NewBufferString(`[{ 
			"event":"before_set", "route":"payment/pp/port", "type":"MinMaxInt64" 
		}]`))
		assert.True(t, errors.Empty.Match(err), "%+v", err)
	})
	t.Run("MinMaxInt64 invalid route", func(t *testing.T) {
		or := observerRegistererFake{
			t: t,
		}
		err := RegisterObservers(or, bytes.NewBufferString(`[{ 
			"event":"before_set", "route":"pay", "type":"MinMaxInt64" 
		}]`))
		assert.True(t, errors.NotValid.Match(err), "%+v", err)
	})
	t.Run("MinMaxInt64 invalid event", func(t *testing.T) {
		or := observerRegistererFake{
			t: t,
		}
		err := RegisterObservers(or, bytes.NewBufferString(`[{ 
			"event":"before_sunrise", "route":"payment/pp/port", "type":"MinMaxInt64", "condition":{"conditions":[3]}
		}]`))
		assert.True(t, errors.NotFound.Match(err), "%+v", err)
	})
	t.Run("MinMaxInt64 malformed condition JSON", func(t *testing.T) {
		or := observerRegistererFake{
			t: t,
		}
		err := RegisterObservers(or, bytes.NewBufferString(`[{ 
			"event":"before_set", "route":"payment/pp/port", "type":"MinMaxInt64", "condition":{"conditions":[x]}
		}]`))
		assert.True(t, errors.BadEncoding.Match(err), "%+v", err)
	})

	t.Run("Strings success", func(t *testing.T) {
		or := observerRegistererFake{
			t:         t,
			wantEvent: config.EventOnAfterSet,
			wantRoute: "aa/ee/ff",
			wantValidator: validation.MustNewStrings(validation.Strings{
				Validators:              []string{"Locale"},
				CSVComma:                "|",
				AdditionalAllowedValues: []string{"Vulcan"},
			}),
		}

		err := RegisterObservers(or, bytes.NewBufferString(`[ { "event":"after_set", "route":"aa/ee/ff", "type":"Strings",
		  "condition":{"validators":["Locale"],"csv_comma":"|","additional_allowed_values":["Vulcan"]}}
		]`))
		assert.NoError(t, err)
	})

	t.Run("Strings condition JSON malformed", func(t *testing.T) {
		or := observerRegistererFake{
			t: t,
		}

		err := RegisterObservers(or, bytes.NewBufferString(`[ { "event":"after_set", "route":"aa/ee/ff", "type":"Strings",
		  "condition":{"validators":["Locale"],"csv_comma":|,"additional_allowed_values":["Vulcan"]}}
		]`))
		assert.True(t, errors.BadEncoding.Match(err), "%+v", err)
	})
	t.Run("Strings condition unsupported validator", func(t *testing.T) {
		or := observerRegistererFake{
			t: t,
		}

		err := RegisterObservers(or, bytes.NewBufferString(`[ { "event":"after_set", "route":"aa/ee/ff", "type":"Strings",
		  "condition":{"validators":["IsPHP"],"additional_allowed_values":["Vulcan"]}}
		]`))
		assert.True(t, errors.NotSupported.Match(err), "%+v", err)
	})

	t.Run("customObserverRegistry success", func(t *testing.T) {
		wantConditionJSON := []byte(`{"validators":["IsPHP"],"additional_allowed_values":["Vulcan"]}`)

		or := observerRegistererFake{
			t:         t,
			wantEvent: config.EventOnAfterGet,
			wantRoute: "bb/ee/ff",
			wantValidator: xmlValidator{
				t:        t,
				wantJSON: wantConditionJSON,
			},
		}
		RegisterCustomObserver("XMLValidation", xmlValidator{
			t:        t,
			wantJSON: wantConditionJSON,
		})

		err := RegisterObservers(or, bytes.NewBufferString(`[ { "event":"after_get", "route":"bb/ee/ff", "type":"XMLValidation",
		  "condition":{"validators":["IsPHP"],"additional_allowed_values":["Vulcan"]}}
		]`))
		assert.NoError(t, err)
	})

	t.Run("customObserverRegistry malformed condition JSON", func(t *testing.T) {
		or := observerRegistererFake{
			t: t,
		}
		RegisterCustomObserver("XMLValidation", xmlValidator{
			t:     t,
			ujErr: errors.New("Ups"),
		})

		err := RegisterObservers(or, bytes.NewBufferString(`[ { "event":"after_get", "route":"bb/ee/ff", "type":"XMLValidation",
		  "condition":{"validators":IsPHP,"additional_allowed_values":["Vulcan"]}}
		]`))
		assert.True(t, errors.BadEncoding.Match(err), "%+v", err)
	})

	t.Run("observer not found", func(t *testing.T) {
		or := observerRegistererFake{
			t: t,
		}
		err := RegisterObservers(or, bytes.NewBufferString(`[ { "event":"after_get", "route":"bb/ee/ff", "type":"YAMLValidation",
		  "condition":{ }}
		]`))
		assert.True(t, errors.NotFound.Match(err), "%+v", err)
	})

}

func TestValidator_MakeEventRoute(t *testing.T) {
	t.Parallel()

	t.Run("All Valid", func(t *testing.T) {
		v := &Validator{
			Route: "general/stores/information",
			Event: "after_set",
		}
		e, r, err := v.MakeEventRoute()
		assert.NoError(t, err)
		assert.Exactly(t, "general/stores/information", r)
		assert.Exactly(t, config.EventOnAfterSet, e)
	})
	t.Run("event not found", func(t *testing.T) {
		v := &Validator{
			Route: "general/stores/information",
			Event: "after_bet",
		}
		e, r, err := v.MakeEventRoute()
		assert.True(t, errors.NotFound.Match(err))
		assert.Empty(t, r)
		assert.Empty(t, e)
	})
	t.Run("invalid route", func(t *testing.T) {
		v := &Validator{
			Route: "d/f",
			Event: "after_set",
		}
		e, r, err := v.MakeEventRoute()
		assert.True(t, errors.NotValid.Match(err))
		assert.Empty(t, r)
		assert.Empty(t, e)
	})
}

var _ UnmarshallableObserver = (*xmlValidator)(nil)

type xmlValidator struct {
	t        *testing.T
	wantJSON []byte
	ujErr    error
}

func (xv xmlValidator) UnmarshalJSON(data []byte) error {
	if xv.ujErr != nil {
		return xv.ujErr
	}
	assert.Exactly(xv.t, xv.wantJSON, data)
	return nil
}

func (xv xmlValidator) Observe(p config.Path, rawData []byte, found bool) (newRawData []byte, err error) {
	return rawData, nil
}

func TestValidators_JSON(t *testing.T) {
	t.Parallel()

	data := []byte(`[ 
	{ "event":"after_get", "route":"gg/ee/ff", "type":"Strings",
		  "condition":{"validators":["Locale"],"additional_allowed_values":["Rhomulan"]}},
	{ "event":"after_set", "route":"aa/ee/ff", "type":"Strings",
		  "condition":{"validators":["Locale"],"csv_comma":"|","additional_allowed_values":["Vulcan"]}}
		]`)
	valis := new(Validators)
	assert.NoError(t, valis.UnmarshalJSON(data))

	newData, err := valis.MarshalJSON()
	assert.NoError(t, err)

	assert.Exactly(t, `[{"route":"gg/ee/ff","event":"after_get","type":"Strings","condition":{"validators":["Locale"],"additional_allowed_values":["Rhomulan"]}},{"route":"aa/ee/ff","event":"after_set","type":"Strings","condition":{"validators":["Locale"],"csv_comma":"|","additional_allowed_values":["Vulcan"]}}]`,
		string(newData))
}

func TestValidators_Validate(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		valis := Validators{
			&Validator{
				Route:     "aa/bb/cc",
				Event:     "before_set",
				Type:      "strings",
				Condition: []byte(`{"validators":["Locale"],"csv_comma":"|","additional_allowed_values":["Vulcan"]}`),
			},
		}
		assert.NoError(t, valis.Validate())
	})

	t.Run("an errors", func(t *testing.T) {
		valis := Validators{
			&Validator{
				Route: "aa/bb/cc",
				Event: "before_set",
				Type:  "strings",
			},
		}
		assert.True(t, errors.Empty.Match(valis.Validate()))
	})

}
