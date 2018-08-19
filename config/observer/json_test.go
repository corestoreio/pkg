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

// +build csall json

package observer

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/corestoreio/errors"
	"github.com/corestoreio/pkg/config"
	"github.com/corestoreio/pkg/util/assert"
	uv "github.com/corestoreio/pkg/util/validation"
	"github.com/mailru/easyjson"
)

var (
	_ uv.Validator       = (*Configurations)(nil)
	_ uv.Validator       = (*Configuration)(nil)
	_ easyjson.Marshaler = (*Configuration)(nil)
	_ easyjson.Marshaler = (*Configurations)(nil)
)

func TestDeregisterObservers(t *testing.T) {
	t.Parallel()

	t.Run("JSON malformed", func(t *testing.T) {
		or := observerRegistererFake{
			t: t,
		}

		err := JSONDeregisterObservers(or, bytes.NewBufferString(`{"Collection":[{ 
			"event":before_set, "route":"payment/pp/port", "type":"ValidateMinMaxInt", "condition":{"conditions":[8080,8090]} 
		}]}`))
		assert.True(t, errors.BadEncoding.Match(err), "%+v", err)
	})
	t.Run("event not found", func(t *testing.T) {
		or := observerRegistererFake{
			t: t,
		}
		err := JSONDeregisterObservers(or, bytes.NewBufferString(`{"Collection":[{ 
			"event":"before_heck", "route":"payment/pp/port", "type":"ValidateMinMaxInt", "condition":{"conditions":[8080,8090]} 
		}]}`))
		assert.True(t, errors.NotFound.Match(err), "%+v", err)
	})

	t.Run("deregistered", func(t *testing.T) {
		or := observerRegistererFake{
			t:         t,
			wantEvent: config.EventOnBeforeSet,
			wantRoute: "payment/pp/port",
		}
		err := JSONDeregisterObservers(or, bytes.NewBufferString(`{"Collection":[{ 
			"event":"before_set", "route":"payment/pp/port" 
		}]}`))
		assert.NoError(t, err)
	})

}

func TestRegisterObservers(t *testing.T) {
	t.Parallel()

	t.Run("JSONRegisterObservers JSON malformed", func(t *testing.T) {
		or := observerRegistererFake{
			t: t,
		}

		err := JSONRegisterObservers(or, bytes.NewBufferString(`{"Collection":[{ 
			"event":before_set, "route":"payment/pp/port", "type":"ValidateMinMaxInt", "condition":{"conditions":[8080,8090]} 
		}]}`))
		assert.True(t, errors.BadEncoding.Match(err), "%+v", err)
	})

	t.Run("ValidateMinMaxInt OK", func(t *testing.T) {
		or := observerRegistererFake{
			t:             t,
			wantEvent:     config.EventOnBeforeSet,
			wantRoute:     "payment/pp/port",
			wantValidator: &ValidateMinMaxInt{Conditions: []int64{8080, 8090}},
		}

		err := JSONRegisterObservers(or, bytes.NewBufferString(`{"Collection":[{ 
			"event":"before_set", "route":"payment/pp/port", "type":"ValidateMinMaxInt", "condition":{"conditions":[8080,8090]} 
		}]}`))
		assert.NoError(t, err)
	})
	t.Run("ValidateMinMaxInt Empty conditions", func(t *testing.T) {
		or := observerRegistererFake{
			t:             t,
			wantEvent:     config.EventOnBeforeSet,
			wantRoute:     "payment/pp/port",
			wantValidator: &ValidateMinMaxInt{Conditions: []int64{}},
		}

		err := JSONRegisterObservers(or, bytes.NewBufferString(`{"Collection":[{ 
			"event":"before_set", "route":"payment/pp/port", "type":"ValidateMinMaxInt", "condition":{"conditions":[]} 
		}]}`))
		assert.NoError(t, err)
	})
	t.Run("ValidateMinMaxInt empty condition", func(t *testing.T) {
		or := observerRegistererFake{
			t: t,
		}
		err := JSONRegisterObservers(or, bytes.NewBufferString(`{"Collection":[{ 
			"event":"before_set", "route":"payment/pp/port", "type":"ValidateMinMaxInt" 
		}]}`))
		assert.True(t, errors.Empty.Match(err), "%+v", err)
	})
	t.Run("ValidateMinMaxInt invalid route", func(t *testing.T) {
		or := observerRegistererFake{
			t: t,
		}
		err := JSONRegisterObservers(or, bytes.NewBufferString(`{"Collection":[{ 
			"event":"before_set", "route":"pay", "type":"ValidateMinMaxInt" 
		}]}`))
		assert.True(t, errors.NotValid.Match(err), "%+v", err)
	})
	t.Run("ValidateMinMaxInt invalid event", func(t *testing.T) {
		or := observerRegistererFake{
			t: t,
		}
		err := JSONRegisterObservers(or, bytes.NewBufferString(`{"Collection":[{ 
			"event":"before_sunrise", "route":"payment/pp/port", "type":"ValidateMinMaxInt", "condition":{"conditions":[3]}
		}]}`))
		assert.True(t, errors.NotFound.Match(err), "%+v", err)
	})
	t.Run("ValidateMinMaxInt malformed condition JSON", func(t *testing.T) {
		or := observerRegistererFake{
			t: t,
		}
		err := JSONRegisterObservers(or, bytes.NewBufferString(`{"Collection":[{ 
			"event":"before_set", "route":"payment/pp/port", "type":"ValidateMinMaxInt", "condition":{"conditions":[x]}
		}]}`))
		assert.True(t, errors.BadEncoding.Match(err), "%+v", err)
	})

	t.Run("ValidatorArg success", func(t *testing.T) {
		or := observerRegistererFake{
			t:         t,
			wantEvent: config.EventOnAfterSet,
			wantRoute: "aa/ee/ff",
			wantValidator: MustNewValidator(ValidatorArg{
				Funcs:                   []string{"Locale"},
				CSVComma:                "|",
				AdditionalAllowedValues: []string{"Vulcan"},
			}),
		}

		err := JSONRegisterObservers(or, bytes.NewBufferString(`{"Collection":[ { "event":"after_set", "route":"aa/ee/ff", "type":"ValidatorArg",
		  "condition":{"validators":["Locale"],"csv_comma":"|","additional_allowed_values":["Vulcan"]}}
		]}`))
		assert.NoError(t, err)
	})

	t.Run("ValidatorArg condition JSON malformed", func(t *testing.T) {
		or := observerRegistererFake{
			t: t,
		}

		err := JSONRegisterObservers(or, bytes.NewBufferString(`{"Collection":[ { "event":"after_set", "route":"aa/ee/ff", "type":"ValidatorArg",
		  "condition":{"validators":["Locale"],"csv_comma":|,"additional_allowed_values":["Vulcan"]}}
		]}`))
		assert.True(t, errors.BadEncoding.Match(err), "%+v", err)
	})
	t.Run("ValidatorArg condition unsupported validator", func(t *testing.T) {
		or := observerRegistererFake{
			t: t,
		}

		err := JSONRegisterObservers(or, bytes.NewBufferString(`{"Collection":[ { "event":"after_set", "route":"aa/ee/ff", "type":"ValidatorArg",
		  "condition":{"validators":["IsPHP"],"additional_allowed_values":["Vulcan"]}}
		]}`))
		assert.True(t, errors.NotSupported.Match(err), "%+v", err)
	})

	t.Run("customObserverRegistry success", func(t *testing.T) {
		wantConditionJSON := []byte(`{"validators":["IsPHP"],"additional_allowed_values":["Vulcan"]}`)

		or := observerRegistererFake{
			t:         t,
			wantEvent: config.EventOnAfterGet,
			wantRoute: "bb/ee/ff",
			wantValidator: xmlValidator{
				wantJSON: wantConditionJSON,
			},
		}
		RegisterCustom("XMLValidationOK", func(data json.RawMessage) (config.Observer, error) {
			assert.Exactly(t, wantConditionJSON, data)
			return xmlValidator{wantJSON: wantConditionJSON}, nil
		})

		err := JSONRegisterObservers(or, bytes.NewBufferString(`{"Collection":[ { "event":"after_get", "route":"bb/ee/ff", "type":"XMLValidationOK",
		  "condition":{"validators":["IsPHP"],"additional_allowed_values":["Vulcan"]}}
		]}`))
		assert.NoError(t, err)
	})

	t.Run("customObserverRegistry new validator throws error", func(t *testing.T) {
		or := observerRegistererFake{
			t: t,
		}
		RegisterCustom("XMLValidationErr01", func(data json.RawMessage) (config.Observer, error) {
			assert.Exactly(t, []byte("{\"validators\":IsPHP,\"additional_allowed_values\":[\"Vulcan\"]}"), data)
			return nil, errors.Blocked.Newf("Ups")
		})

		err := JSONRegisterObservers(or, bytes.NewBufferString(`{"Collection":[ { "event":"after_get", "route":"bb/ee/ff", "type":"XMLValidationErr01",
		  "condition":{"validators":IsPHP,"additional_allowed_values":["Vulcan"]}}
		]}`))
		assert.True(t, errors.Blocked.Match(err), "%+v", err)
	})

	t.Run("observer not found", func(t *testing.T) {
		or := observerRegistererFake{
			t: t,
		}
		err := JSONRegisterObservers(or, bytes.NewBufferString(`{"Collection":[ { "event":"after_get", "route":"bb/ee/ff", "type":"YAMLValidation",
		  "condition":{ }}
		]}`))
		assert.True(t, errors.NotFound.Match(err), "%+v", err)
	})

}

func TestValidator_MakeEventRoute(t *testing.T) {
	t.Parallel()

	t.Run("All Valid", func(t *testing.T) {
		v := &Configuration{
			Route: "general/stores/information",
			Event: "after_set",
		}
		e, r, err := v.MakeEventRoute()
		assert.NoError(t, err)
		assert.Exactly(t, "general/stores/information", r)
		assert.Exactly(t, config.EventOnAfterSet, e)
	})
	t.Run("event not found", func(t *testing.T) {
		v := &Configuration{
			Route: "general/stores/information",
			Event: "after_bet",
		}
		e, r, err := v.MakeEventRoute()
		assert.True(t, errors.NotFound.Match(err))
		assert.Empty(t, r)
		assert.Empty(t, e)
	})
	t.Run("invalid route", func(t *testing.T) {
		v := &Configuration{
			Route: "d/f",
			Event: "after_set",
		}
		e, r, err := v.MakeEventRoute()
		assert.True(t, errors.NotValid.Match(err))
		assert.Empty(t, r)
		assert.Empty(t, e)
	})
}

type xmlValidator struct {
	wantJSON []byte
	err      error
}

func (xv xmlValidator) Observe(p config.Path, rawData []byte, found bool) (newRawData []byte, err error) {
	if xv.err != nil {
		return nil, xv.err
	}
	return rawData, nil
}

func TestValidators_JSON(t *testing.T) {
	t.Parallel()

	data := []byte(`{"Collection":[ 
	{ "event":"after_get", "route":"gg/ee/ff", "type":"ValidatorArg",
		  "condition":{"validators":["Locale"],"additional_allowed_values":["Rhomulan"]}},
	{ "event":"after_set", "route":"aa/ee/ff", "type":"ValidatorArg",
		  "condition":{"validators":["Locale"],"csv_comma":"|","additional_allowed_values":["Vulcan"]}}
		]}`)
	valis := new(Configurations)
	assert.NoError(t, valis.UnmarshalJSON(data))

	newData, err := valis.MarshalJSON()
	assert.NoError(t, err)

	assert.Exactly(t, `{"Collection":[{"route":"gg/ee/ff","event":"after_get","type":"ValidatorArg","condition":{"validators":["Locale"],"additional_allowed_values":["Rhomulan"]}},{"route":"aa/ee/ff","event":"after_set","type":"ValidatorArg","condition":{"validators":["Locale"],"csv_comma":"|","additional_allowed_values":["Vulcan"]}}]}`,
		string(newData))
}

func TestValidators_Validate(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		valis := Configurations{
			Collection: []*Configuration{
				{
					Route:     "aa/bb/cc",
					Event:     "before_set",
					Type:      "strings",
					Condition: []byte(`{"validators":["Locale"],"csv_comma":"|","additional_allowed_values":["Vulcan"]}`),
				},
			},
		}
		assert.NoError(t, valis.Validate())
	})

	t.Run("an errors", func(t *testing.T) {
		valis := Configurations{
			Collection: []*Configuration{
				{
					Route: "aa/bb/cc",
					Event: "before_set",
					Type:  "strings",
				},
			},
		}
		assert.True(t, errors.Empty.Match(valis.Validate()))
	})

}
