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

package backendstore_test

import (
	"bytes"
	"testing"

	"github.com/corestoreio/csfw/config/cfgmock"
	"github.com/corestoreio/csfw/store/backendstore"
)

// backend overall backend models for all tests
var backend *backendstore.Configuration

// this would belong into the test suit setup
func init() {
	cfgStruct, err := backendstore.NewConfigStructure()
	if err != nil {
		panic(err)
	}
	backend = backendstore.New(cfgStruct)
}

func TestConfiguration_FormatAddressText(t *testing.T) {
	var buf = new(bytes.Buffer)
	sg := cfgmock.NewService().NewScoped(3, 4)
	backend.FormatAddressText("", buf, sg)
	t.Log(buf.String())
}
