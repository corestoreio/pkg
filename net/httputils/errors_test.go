// Copyright 2015, Cyrill @ Schumacher.fm and the CoreStore contributors
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

package httputils_test

import (
	"errors"
	"github.com/corestoreio/csfw/net/httputils"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestWriteError(t *testing.T) {
	err := errors.New("My not so cool error")
	rec := httptest.NewRecorder()
	httputils.WriteError(rec, err, http.StatusTeapot)
	assert.EqualValues(t, http.StatusTeapot, rec.Code)
	assert.Equal(t, err.Error()+"\n", rec.Body.String())
}
