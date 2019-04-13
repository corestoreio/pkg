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

package ddl

import (
	"bytes"
	"testing"

	"github.com/corestoreio/pkg/util/assert"
)

func TestKeyRelationShips(t *testing.T) {
	krs := &KeyRelationShips{
		relMap: map[string]bool{
			`store.website_id|customer_entity.website_id|MUL`:       true,
			`store.group_id|customer_entity.group_id`:               true,
			`store.website_id|store_group.website_id|MUL`:           true,
			`store.group_id|store_group.group_id|PRI`:               true,
			`store.website_id|store_website.website_id|PRI`:         true,
			`store_group.website_id|customer_entity.website_id|MUL`: true,
			`store_group.website_id|store.website_id|MUL`:           true,
			`store_group.website_id|store_website.website_id|PRI`:   true,
		},
	}
	assert.True(t, krs.IsOneToOne("store_group", "website_id", "store_website", "website_id"))
	assert.True(t, krs.IsOneToOne("store", "group_id", "store_group", "group_id"))
	assert.False(t, krs.IsOneToMany("store", "group_id", "store_group", "group_id"))
	assert.True(t, krs.IsOneToMany("store_group", "website_id", "store", "website_id"))
	assert.False(t, krs.IsOneToOne("store_group", "website_id", "store", "website_id"))

	var buf bytes.Buffer
	krs.Debug(&buf)

	// since Go 1.12 maps are printed sorted
	assert.Exactly(t, `store.group_id|customer_entity.group_id
store.group_id|store_group.group_id|PRI
store.website_id|customer_entity.website_id|MUL
store.website_id|store_group.website_id|MUL
store.website_id|store_website.website_id|PRI
store_group.website_id|customer_entity.website_id|MUL
store_group.website_id|store.website_id|MUL
store_group.website_id|store_website.website_id|PRI
`, buf.String())
}
