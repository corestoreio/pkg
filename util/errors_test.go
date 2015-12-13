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

package util_test

import (
	"errors"
	"testing"

	"github.com/corestoreio/csfw/util"
	"github.com/juju/errgo"
	"github.com/stretchr/testify/assert"
)

func TestErrors(t *testing.T) {
	assert.Equal(t, "Err1\nErr2\nErr3", util.Errors(
		errors.New("Err1"),
		errors.New("Err2"),
		errors.New("Err3"),
	))

	err := util.Errors(
		errgo.New("Err1"),
		errgo.New("Err2"),
		errors.New("Err3"),
	)
	assert.Contains(t, err, "corestoreio/csfw/util/errors_test.go:34\nErr2")
}
