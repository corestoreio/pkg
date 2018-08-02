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

package http_test

import (
	"bytes"
	"fmt"
	gohttp "net/http"
	"net/http/httptest"
	"testing"

	"github.com/corestoreio/pkg/config"
	"github.com/corestoreio/pkg/config/validation"
	"github.com/corestoreio/pkg/config/validation/http"
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
	if orf.wantEvent != event {
		assert.Exactly(orf.t, orf.wantEvent, event, "Event should be equal: have %d want %d", orf.wantEvent, event)
	}
	if orf.wantRoute != route {
		assert.Exactly(orf.t, orf.wantRoute, route, "Routes")
	}

	// Pointers are different in the final objects hence they get printed and
	// their structure compared, not the address.
	if want, have := fmt.Sprintf("%#v", orf.wantValidator), fmt.Sprintf("%#v", o); want != have {
		assert.Exactly(orf.t, want, have, "Observer internal types should match")
	}
	return nil
}

func (orf observerRegistererFake) DeregisterObserver(event uint8, route string) error {
	if orf.err != nil {
		return orf.err
	}
	if orf.wantEvent != event {
		assert.Exactly(orf.t, orf.wantEvent, event, "Event should be equal: have %d want %d", orf.wantEvent, event)
	}
	if orf.wantRoute != route {
		assert.Exactly(orf.t, orf.wantRoute, route, "Routes")
	}
	return nil
}

func TestRegisterObserversFromJSON(t *testing.T) {
	t.Parallel()

	t.Run("StatusCreated", func(t *testing.T) {
		hdnlr := http.RegisterObserversFromJSON(observerRegistererFake{
			t:         t,
			wantEvent: config.EventOnAfterSet,
			wantRoute: "aa/ee/ff",
			wantValidator: validation.MustNewStrings(validation.Strings{
				Validators:              []string{"Locale"},
				AdditionalAllowedValues: []string{"Vulcan"},
				CSVComma:                "|",
			}),
		}, http.HandlerOptions{})

		w := httptest.NewRecorder()
		req, err := gohttp.NewRequest("GET", "/", bytes.NewBufferString(`[ { "event":"after_set", "route":"aa/ee/ff", "type":"Strings",
		  "condition":{"validators":["Locale"],"csv_comma":"|","additional_allowed_values":["Vulcan"]}}
		]`))
		assert.NoError(t, err)
		hdnlr.ServeHTTP(w, req)

		assert.Exactly(t, w.Code, gohttp.StatusCreated, "gohttp.StatusCreated")
		assert.Empty(t, w.Body.String())
	})

	t.Run("custom StatusNonAuthoritativeInfo", func(t *testing.T) {
		hdnlr := http.RegisterObserversFromJSON(observerRegistererFake{
			t:         t,
			wantEvent: config.EventOnAfterSet,
			wantRoute: "aa/ee/ff",
			wantValidator: validation.MustNewStrings(validation.Strings{
				Validators: []string{"int"},
			}),
		}, http.HandlerOptions{
			StatusCodeOk: gohttp.StatusNonAuthoritativeInfo,
		})

		w := httptest.NewRecorder()
		req, err := gohttp.NewRequest("GET", "/", bytes.NewBufferString(`[ { "event":"after_set", "route":"aa/ee/ff", "type":"Strings",
		  "condition":{"validators":["int"]}}
		]`))
		assert.NoError(t, err)
		hdnlr.ServeHTTP(w, req)

		assert.Exactly(t, w.Code, gohttp.StatusNonAuthoritativeInfo, "gohttp.StatusNonAuthoritativeInfo")
		assert.Empty(t, w.Body.String())
	})

	t.Run("Request too large StatusNotAcceptable", func(t *testing.T) {
		hdnlr := http.RegisterObserversFromJSON(observerRegistererFake{
			t:         t,
			wantEvent: config.EventOnAfterSet,
			wantRoute: "aa/ee/ff",
			wantValidator: validation.MustNewStrings(validation.Strings{
				Validators: []string{"int"},
			}),
		}, http.HandlerOptions{
			StatusCodeError: gohttp.StatusNotAcceptable,
			MaxRequestSize:  100,
		})

		w := httptest.NewRecorder()
		req, err := gohttp.NewRequest("GET", "/", bytes.NewBufferString(`[ { "event":"after_set", "route":"aa/ee/ff", "type":"Strings",
		  "condition":{"validators":["int"]}}
		]`))
		assert.NoError(t, err)
		hdnlr.ServeHTTP(w, req)

		assert.Exactly(t, w.Code, gohttp.StatusNotAcceptable, "gohttp.StatusNotAcceptable")
		assert.Contains(t, w.Body.String(), "parse error: EOF reached while skipping array/object or token near offset 100 of ''")
	})
}

func TestDeregisterObserverFromJSON(t *testing.T) {
	t.Parallel()

	t.Run("StatusAccepted", func(t *testing.T) {
		hdnlr := http.DeregisterObserverFromJSON(observerRegistererFake{
			t:         t,
			wantEvent: config.EventOnAfterSet,
			wantRoute: "aa/ee/ff",
		}, http.HandlerOptions{})

		w := httptest.NewRecorder()
		req, err := gohttp.NewRequest("GET", "/", bytes.NewBufferString(`[ { "event":"after_set", "route":"aa/ee/ff" }
		]`))
		assert.NoError(t, err)
		hdnlr.ServeHTTP(w, req)

		assert.Exactly(t, w.Code, gohttp.StatusAccepted, "gohttp.StatusAccepted")
		assert.Empty(t, w.Body.String())
	})

	t.Run("custom StatusNonAuthoritativeInfo", func(t *testing.T) {
		hdnlr := http.DeregisterObserverFromJSON(observerRegistererFake{
			t:         t,
			wantEvent: config.EventOnAfterSet,
			wantRoute: "aa/ee/ff",
		}, http.HandlerOptions{
			StatusCodeOk: gohttp.StatusNonAuthoritativeInfo,
		})

		w := httptest.NewRecorder()
		req, err := gohttp.NewRequest("GET", "/", bytes.NewBufferString(`[ { "event":"after_set", "route":"aa/ee/ff"}
		]`))
		assert.NoError(t, err)
		hdnlr.ServeHTTP(w, req)

		assert.Exactly(t, w.Code, gohttp.StatusNonAuthoritativeInfo, "gohttp.StatusNonAuthoritativeInfo")
		assert.Empty(t, w.Body.String())
	})

	t.Run("Request too large StatusNotAcceptable", func(t *testing.T) {
		hdnlr := http.DeregisterObserverFromJSON(observerRegistererFake{
			t:         t,
			wantEvent: config.EventOnAfterSet,
			wantRoute: "aa/ee/ff",
		}, http.HandlerOptions{
			StatusCodeError: gohttp.StatusNotAcceptable,
			MaxRequestSize:  43,
		})

		w := httptest.NewRecorder()
		req, err := gohttp.NewRequest("GET", "/", bytes.NewBufferString(`[ { "event":"after_set", "route":"aa/ee/ff" ]`))
		assert.NoError(t, err)
		hdnlr.ServeHTTP(w, req)

		assert.Exactly(t, w.Code, gohttp.StatusNotAcceptable, "gohttp.StatusNotAcceptable")
		assert.Contains(t, w.Body.String(), "[config/validation/json] JSON decoding failed")
	})
}
