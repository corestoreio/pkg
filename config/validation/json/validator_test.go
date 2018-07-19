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
	"testing"

	"github.com/corestoreio/pkg/config"
	"github.com/corestoreio/pkg/config/validation"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
	assert.Exactly(orf.t, orf.wantValidator, o, "Observer internal types should match")

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

func TestRegisterObserversFromJSON(t *testing.T) {

	t.Run("MinMaxInt64", func(t *testing.T) {
		var payLoad = []byte(`[ { "event":"before_set", "route":"payment/pp/port", "type":"MinMaxInt64", "condition":{ "min":8080, "max":8090 } } ]`)
		or := observerRegistererFake{
			t:             t,
			wantEvent:     config.EventOnBeforeSet,
			wantRoute:     "payment/pp/port",
			wantValidator: &validation.MinMaxInt64{Min: 8080, Max: 8090},
		}

		err := validation.RegisterObserversFromJSON(or, payLoad)
		require.NoError(t, err)
	})

	t.Run("UUID", func(t *testing.T) {
		var payLoad = []byte(`[ { "event":"after_get", "route":"aa/bb/cc", "type":"UUID", "condition":{ "version": 4 } } ]`)
		or := observerRegistererFake{
			t:             t,
			wantEvent:     config.EventOnAfterGet,
			wantRoute:     "aa/bb/cc",
			wantValidator: &validation.UUID{Version: 4},
		}

		err := validation.RegisterObserversFromJSON(or, payLoad)
		require.NoError(t, err)
	})

	t.Run("Strings", func(t *testing.T) {
		isoC2Val, err := validation.NewStrings("DE")
		require.NoError(t, err)
		isoC2Val.CSVComma = ";"

		var payLoad = []byte(`[ { "event":"after_set", "route":"aa/ee/ff", "type":"Strings", "condition":{ "csv_comma": ";", "restrict_to_codes":{"DE":true } } } ]`)
		or := observerRegistererFake{
			t:             t,
			wantEvent:     config.EventOnAfterSet,
			wantRoute:     "aa/ee/ff",
			wantValidator: isoC2Val,
		}

		err = validation.RegisterObserversFromJSON(or, payLoad)
		require.NoError(t, err)
	})

}
