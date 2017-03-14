// Copyright 2015-2017, Cyrill @ Schumacher.fm and the CoreStore contributors
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

package eav_test

import (
	"testing"

	"github.com/corestoreio/csfw/eav"
	"github.com/stretchr/testify/assert"
)

func TestAttributeSource(t *testing.T) {
	a := eav.NewAttributeSource(
		// temporary because later these values comes from another slice/container/database
		func(as *eav.AttributeSource) {
			as.Source = []string{
				"BAY", "Bavaria",
				"BAW", "Baden-Würstchenberg",
				"HAM", "Hamburg",
				"BER", "Bärlin",
			}
		},
	)
	assert.Equal(
		t,
		eav.AttributeSourceOptions{eav.AttributeSourceOption{Value: "BAY", Label: "Bavaria"}, eav.AttributeSourceOption{Value: "BAW", Label: "Baden-Würstchenberg"}, eav.AttributeSourceOption{Value: "HAM", Label: "Hamburg"}, eav.AttributeSourceOption{Value: "BER", Label: "Bärlin"}},
		a.GetAllOptions(),
	)

}
