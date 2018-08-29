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

// +build csall http

package observer_test

import (
	"bytes"
	"fmt"
	gohttp "net/http"
	"net/http/httptest"
	"testing"

	"github.com/corestoreio/pkg/config"
	"github.com/corestoreio/pkg/config/observer"
	"github.com/corestoreio/pkg/util/assert"
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
func TestRegisterObserversFromJSON(t *testing.T) {
	t.Parallel()

	t.Run("StatusCreated", func(t *testing.T) {
		hdnlr := observer.HTTPJSONRegistry(observerRegistererFake{
			t:         t,
			wantEvent: config.EventOnAfterSet,
			wantRoute: "aa/ee/ff",
			wantValidator: observer.MustNewValidator(observer.ValidatorArg{
				Funcs:                   []string{"Locale"},
				AdditionalAllowedValues: []string{"Vulcan"},
				CSVComma:                "|",
			}),
		}, observer.HTTPHandlerOptions{})

		w := httptest.NewRecorder()
		req, err := gohttp.NewRequest("GET", "/", bytes.NewBufferString(`{"Collection":[ { "event":"after_set", "route":"aa/ee/ff", "type":"validator",
		  "condition":{"funcs":["Locale"],"csv_comma":"|","additional_allowed_values":["Vulcan"]}}
		]}`))
		assert.NoError(t, err)
		hdnlr.ServeHTTP(w, req)

		assert.Empty(t, w.Body.String())
		assert.Exactly(t, gohttp.StatusCreated, w.Code, "gohttp.StatusCreated")
	})

	t.Run("custom StatusNonAuthoritativeInfo", func(t *testing.T) {
		hdnlr := observer.HTTPJSONRegistry(observerRegistererFake{
			t:         t,
			wantEvent: config.EventOnAfterSet,
			wantRoute: "aa/ee/ff",
			wantValidator: observer.MustNewValidator(observer.ValidatorArg{
				Funcs: []string{"int"},
			}),
		}, observer.HTTPHandlerOptions{
			StatusCodeOk: gohttp.StatusNonAuthoritativeInfo,
		})

		w := httptest.NewRecorder()
		req, err := gohttp.NewRequest("GET", "/", bytes.NewBufferString(`{"Collection":[ { "event":"after_set", "route":"aa/ee/ff", "type":"validator",
		  "condition":{"funcs":["int"]}}
		]}`))
		assert.NoError(t, err)
		hdnlr.ServeHTTP(w, req)

		assert.Exactly(t, gohttp.StatusNonAuthoritativeInfo, w.Code, "gohttp.StatusNonAuthoritativeInfo")
		assert.Empty(t, w.Body.String())
	})

	t.Run("Request too large StatusNotAcceptable", func(t *testing.T) {
		hdnlr := observer.HTTPJSONRegistry(observerRegistererFake{
			t:         t,
			wantEvent: config.EventOnAfterSet,
			wantRoute: "aa/ee/ff",
			wantValidator: observer.MustNewValidator(observer.ValidatorArg{
				Funcs: []string{"int"},
			}),
		}, observer.HTTPHandlerOptions{
			StatusCodeError: gohttp.StatusNotAcceptable,
			MaxRequestSize:  100,
		})

		w := httptest.NewRecorder()
		req, err := gohttp.NewRequest("GET", "/", bytes.NewBufferString(`{"Collection":[ { "event":"after_set", "route":"aa/ee/ff", "type":"validator",
		  "condition":{"validators":["int"]}}
		]}`))
		assert.NoError(t, err)
		hdnlr.ServeHTTP(w, req)

		assert.Exactly(t, gohttp.StatusNotAcceptable, w.Code, "gohttp.StatusNotAcceptable")
		assert.Contains(t, w.Body.String(), "parse error: EOF reached while skipping array/object or token near offset 100 of ''")
	})
}

func TestDeregisterObserverFromJSON(t *testing.T) {
	t.Parallel()

	t.Run("StatusAccepted", func(t *testing.T) {
		hdnlr := observer.HTTPJSONDeregistry(observerRegistererFake{
			t:         t,
			wantEvent: config.EventOnAfterSet,
			wantRoute: "aa/ee/ff",
		}, observer.HTTPHandlerOptions{})

		w := httptest.NewRecorder()
		req, err := gohttp.NewRequest("GET", "/", bytes.NewBufferString(`{"Collection":[ { "event":"after_set", "route":"aa/ee/ff" }
		]}`))
		assert.NoError(t, err)
		hdnlr.ServeHTTP(w, req)

		assert.Exactly(t, gohttp.StatusAccepted, w.Code, "gohttp.StatusAccepted")
		assert.Empty(t, w.Body.String())
	})

	t.Run("custom StatusNonAuthoritativeInfo", func(t *testing.T) {
		hdnlr := observer.HTTPJSONDeregistry(observerRegistererFake{
			t:         t,
			wantEvent: config.EventOnAfterSet,
			wantRoute: "aa/ee/ff",
		}, observer.HTTPHandlerOptions{
			StatusCodeOk: gohttp.StatusNonAuthoritativeInfo,
		})

		w := httptest.NewRecorder()
		req, err := gohttp.NewRequest("GET", "/", bytes.NewBufferString(`{"Collection":[ { "event":"after_set", "route":"aa/ee/ff"}
		]}`))
		assert.NoError(t, err)
		hdnlr.ServeHTTP(w, req)

		assert.Exactly(t, gohttp.StatusNonAuthoritativeInfo, w.Code, "gohttp.StatusNonAuthoritativeInfo")
		assert.Empty(t, w.Body.String())
	})

	t.Run("Request too large StatusNotAcceptable", func(t *testing.T) {
		hdnlr := observer.HTTPJSONDeregistry(observerRegistererFake{
			t:         t,
			wantEvent: config.EventOnAfterSet,
			wantRoute: "aa/ee/ff",
		}, observer.HTTPHandlerOptions{
			StatusCodeError: gohttp.StatusNotAcceptable,
			MaxRequestSize:  43,
		})

		w := httptest.NewRecorder()
		req, err := gohttp.NewRequest("GET", "/", bytes.NewBufferString(`{"Collection":[ { "event":"after_set", "route":"aa/ee/ff" ]}`))
		assert.NoError(t, err)
		hdnlr.ServeHTTP(w, req)

		assert.Exactly(t, gohttp.StatusNotAcceptable, w.Code, "gohttp.StatusNotAcceptable")
		assert.Contains(t, w.Body.String(), "[config/validation/json] JSON decoding failed")
	})
}
