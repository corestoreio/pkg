// Copyright 2015-2016, Cyrill @ Schumacher.fm and the CoreStore contributors
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

package ctxhttp_test

import (
	"net/http"
	"testing"

	"github.com/corestoreio/csfw/net/ctxhttp"
	"github.com/juju/errors"
	"github.com/stretchr/testify/assert"
)

func TestNewError(t *testing.T) {
	tests := []struct {
		code      int
		msg       []string
		wantError string
	}{
		{http.StatusBadGateway, nil, http.StatusText(http.StatusBadGateway)},
		{http.StatusTeapot, []string{"No coffee pot", "ignored"}, "No coffee pot"},
	}
	for _, test := range tests {
		he := ctxhttp.NewError(test.code, test.msg...)
		assert.Exactly(t, test.code, he.Code)
		assert.Exactly(t, test.wantError, he.Error())
	}
}

func TestNewErrorFromErrors(t *testing.T) {
	tests := []struct {
		code      int
		errs      []error
		wantError string
		verbose   bool
	}{
		{http.StatusBadGateway, nil, http.StatusText(http.StatusBadGateway), false},
		{http.StatusTeapot, []error{errors.New("No coffee pot"), errors.New("Not even a milk pot")}, "No coffee pot\nNot even a milk pot", false},
		{http.StatusConflict, []error{errors.New("Now a coffee pot"), errors.New("Not even close to a milk pot")}, "error_test.go", true},
	}
	for _, test := range tests {
		he := ctxhttp.NewErrorFromErrors(test.code, test.errs...)
		he.Verbose = test.verbose
		assert.Exactly(t, test.code, he.Code)
		assert.Contains(t, he.Error(), test.wantError)
	}
}
