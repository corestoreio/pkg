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

// +build !json
// +build !csall

package store_test

import (
	"io/ioutil"
	"testing"

	storemock "github.com/corestoreio/pkg/store/mock"
	"github.com/corestoreio/pkg/util/assert"
)

func TestService_Sorting_JSON(t *testing.T) {

	t.Run("EuroW11G11S19", func(t *testing.T) {
		srv := storemock.NewServiceEuroW11G11S19()
		haveData := toJSON(srv)

		goldenData, err := ioutil.ReadFile("testdata/sort_euroW11G11S19.golden.json")
		assert.NoError(t, err)
		assert.Exactly(t, goldenData, haveData)
	})

	t.Run("EuroOZ", func(t *testing.T) {
		srv := storemock.NewServiceEuroOZ()
		haveData := toJSON(srv)

		goldenData, err := ioutil.ReadFile("testdata/sort_euroOZ.golden.json")
		assert.NoError(t, err)
		assert.Exactly(t, goldenData, haveData)
	})
}
