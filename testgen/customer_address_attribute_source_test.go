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

package testgen

import (
	"testing"

	"github.com/corestoreio/csfw/customer/custattr"
	"github.com/corestoreio/csfw/eav"
	"github.com/stretchr/testify/assert"
)

func TestAddressAttributeFrontendLabel(t *testing.T) {
	var err error
	cae, err := eav.GetEntityTypeByCode("customer_address")
	if err != nil {
		t.Error(err)
		return
	}
	aIFs, err := cae.AttributeModel.GetByCode("country_ids")
	assert.Error(t, err)
	assert.Nil(t, aIFs)

	attrIF, err := cae.AttributeModel.GetByCode("country_id")
	if err != nil {
		t.Error(err)
		assert.Error(t, err)
	} else {
		var countryID, ok = attrIF.(custattr.Attributer) // type assertion
		if !ok {
			t.Error("failed to convert countryID into custattr.Attributer type")
		}
		assert.True(t, countryID.SortOrder() > 0)
		assert.Equal(t, "Country", countryID.FrontendLabel())
		// t.Logf("\n%#v\n", countryID.SourceModel())
	}

	apc := cae.AttributeModel.MustGet(CustomerAddressAttributePostcode).(custattr.Attributer)
	assert.Equal(t, "Zip/Postal Code", apc.FrontendLabel())
}

var countryIDFrontendLabel string

// BenchmarkAddressAttributeFrontendLabel	20.000.000	       115 ns/op	       0 B/op	       0 allocs/op
// This is the result for selecting the frontend_label from an attribute assigned to an entity.
// $ac = Mage::getSingleton('eav/config')->getEntityType('customer_address')->getAttributeCollection();
// $ac->addFieldToFilter('attribute_code', ['eq' => 'country_id']);
// $fl = $ac->getFirstItem()->getData('frontend_label');
// PHP 5.5 needs 86600 ns/op (nanosecond) to be fair: with database access but enabled caches with redis and OPcache.
func BenchmarkAddressAttributeFrontendLabel(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		var err error
		cae, err := eav.GetEntityTypeByCode("customer_address")
		if err != nil {
			b.Error(err)
			return
		}

		attrIF, err := cae.AttributeModel.GetByCode("country_id")
		if err != nil {
			b.Error(err)
		} else {
			var countryID, ok = attrIF.(custattr.Attributer) // type assertion
			if !ok {
				b.Error("failed to convert countryID into custattr.Attributer type")
			}
			countryIDFrontendLabel = countryID.FrontendLabel()
		}
	}
}
