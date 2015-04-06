// Copyright 2015 CoreStore Authors
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

package testgen

import (
	"testing"

	"github.com/corestoreio/csfw/customer/custattr"
	"github.com/corestoreio/csfw/eav"
	"github.com/stretchr/testify/assert"
)

func init() {

}

func TestAddressAttributeSource(t *testing.T) {
	var err error
	cae, err := eav.GetEntityTypeByCode("customer_address")
	if err != nil {
		t.Error(err)
		return
	}

	countryID, err := cae.AttributeModel.GetByCode("country_ids")
	if err != nil {
		t.Error(err)
		assert.Error(t, err)
	} else {
		var ok bool
		if countryID, ok = countryID.(custattr.Attributer); !ok {
			t.Error("failed to convert countryID into custattr.Attributer")
		}

		t.Logf("\n%#v\n", countryID)
		//		caac.
		//			assert.Equal(
		//			t,
		//			eav.AttributeSourceOptions{eav.AttributeSourceOption{Value: "AU", Label: "Straya"}, eav.AttributeSourceOption{Value: "NZ", Label: "Kiwi land"}, eav.AttributeSourceOption{Value: "DE", Label: "Autobahn"}, eav.AttributeSourceOption{Value: "SE", Label: "Smørrebrød"}},
		//			attr.SourceModel().GetAllOptions(),
		//		)
	}
}
