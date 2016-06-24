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

package ratelimit

import (
	"context"
	"net/http"
	"testing"

	"github.com/corestoreio/csfw/util/errors"
	"github.com/stretchr/testify/assert"
)

func TestContextWithError(t *testing.T) {
	var wantErr = errors.New("Contiki Context")
	req, _ := http.NewRequest("GET", "http://localhost", nil)
	req = wrapContextError(req, wantErr)
	assert.NotNil(t, req)

	err := FromContextRateLimit(req.Context())
	assert.EqualError(t, err, wantErr.Error())

	err = FromContextRateLimit(context.TODO())
	assert.True(t, errors.IsNotFound(err), "Error: %s", err)
}
