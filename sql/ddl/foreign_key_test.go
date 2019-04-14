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

	"github.com/corestoreio/pkg/storage/null"
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

func TestReverseKeyColumnUsage(t *testing.T) {
	kcuc := map[string]KeyColumnUsageCollection{
		"x910cms_block_store": {
			Data: []*KeyColumnUsage{
				{
					ConstraintCatalog:          "def",
					ConstraintSchema:           "test",
					ConstraintName:             "CMS_BLOCK_STORE_BLOCK_ID_CMS_BLOCK_BLOCK_ID",
					TableCatalog:               "def",
					TableSchema:                "test",
					TableName:                  "x910cms_block_store",
					ColumnName:                 "block_id",
					OrdinalPosition:            1,
					PositionInUniqueConstraint: null.MakeInt64(1),
					ReferencedTableSchema:      null.MakeString(`test`),
					ReferencedTableName:        null.MakeString(`x910cms_block`),
					ReferencedColumnName:       null.MakeString(`block_id`),
				},
				{
					ConstraintCatalog:          "def",
					ConstraintSchema:           "test",
					ConstraintName:             "CMS_BLOCK_STORE_STORE_ID_STORE_STORE_ID",
					TableCatalog:               "def",
					TableSchema:                "test",
					TableName:                  "x910cms_block_store",
					ColumnName:                 "store_id",
					OrdinalPosition:            1,
					PositionInUniqueConstraint: null.MakeInt64(1),
					ReferencedTableSchema:      null.MakeString(`test`),
					ReferencedTableName:        null.MakeString(`store`),
					ReferencedColumnName:       null.MakeString(`store_id`),
				},
			},
			BeforeMapColumns: nil,
			AfterMapColumns:  nil,
		},
		"x910cms_page_store": {
			Data: []*KeyColumnUsage{
				{
					ConstraintCatalog:          "def",
					ConstraintSchema:           "test",
					ConstraintName:             "CMS_PAGE_STORE_PAGE_ID_CMS_PAGE_PAGE_ID",
					TableCatalog:               "def",
					TableSchema:                "test",
					TableName:                  "x910cms_page_store",
					ColumnName:                 "page_id",
					OrdinalPosition:            1,
					PositionInUniqueConstraint: null.MakeInt64(1),
					ReferencedTableSchema:      null.MakeString(`test`),
					ReferencedTableName:        null.MakeString(`x910cms_page`),
					ReferencedColumnName:       null.MakeString(`page_id`),
				},
				{
					ConstraintCatalog:          "def",
					ConstraintSchema:           "test",
					ConstraintName:             "CMS_PAGE_STORE_STORE_ID_STORE_STORE_ID",
					TableCatalog:               "def",
					TableSchema:                "test",
					TableName:                  "x910cms_page_store",
					ColumnName:                 "store_id",
					OrdinalPosition:            1,
					PositionInUniqueConstraint: null.MakeInt64(1),
					ReferencedTableSchema:      null.MakeString(`test`),
					ReferencedTableName:        null.MakeString(`store`),
					ReferencedColumnName:       null.MakeString(`store_id`),
				},
			},
			BeforeMapColumns: nil,
			AfterMapColumns:  nil,
		},
		"catalog_category_entity": {
			Data: []*KeyColumnUsage{
				{
					ConstraintCatalog:          "def",
					ConstraintSchema:           "test",
					ConstraintName:             "CAT_CTGR_ENTT_ENTT_ID_SEQUENCE_CAT_CTGR_SEQUENCE_VAL",
					TableCatalog:               "def",
					TableSchema:                "test",
					TableName:                  "catalog_category_entity",
					ColumnName:                 "entity_id",
					OrdinalPosition:            1,
					PositionInUniqueConstraint: null.MakeInt64(1),
					ReferencedTableSchema:      null.MakeString(`test`),
					ReferencedTableName:        null.MakeString(`sequence_catalog_category`),
					ReferencedColumnName:       null.MakeString(`sequence_value`),
				},
			},
			BeforeMapColumns: nil,
			AfterMapColumns:  nil,
		},
		"customer_address_entity": {
			Data: []*KeyColumnUsage{
				{
					ConstraintCatalog:          "def",
					ConstraintSchema:           "test",
					ConstraintName:             "CUSTOMER_ADDRESS_ENTITY_PARENT_ID_CUSTOMER_ENTITY_ENTITY_ID",
					TableCatalog:               "def",
					TableSchema:                "test",
					TableName:                  "customer_address_entity",
					ColumnName:                 "parent_id",
					OrdinalPosition:            1,
					PositionInUniqueConstraint: null.MakeInt64(1),
					ReferencedTableSchema:      null.MakeString(`test`),
					ReferencedTableName:        null.MakeString(`customer_entity`),
					ReferencedColumnName:       null.MakeString(`entity_id`),
				},
			},
			BeforeMapColumns: nil,
			AfterMapColumns:  nil,
		},
		"customer_entity": {
			Data: []*KeyColumnUsage{
				{
					ConstraintCatalog:          "def",
					ConstraintSchema:           "test",
					ConstraintName:             "CUSTOMER_ENTITY_STORE_ID_STORE_STORE_ID",
					TableCatalog:               "def",
					TableSchema:                "test",
					TableName:                  "customer_entity",
					ColumnName:                 "store_id",
					OrdinalPosition:            1,
					PositionInUniqueConstraint: null.MakeInt64(1),
					ReferencedTableSchema:      null.MakeString(`test`),
					ReferencedTableName:        null.MakeString(`store`),
					ReferencedColumnName:       null.MakeString(`store_id`),
				},
				{
					ConstraintCatalog:          "def",
					ConstraintSchema:           "test",
					ConstraintName:             "CUSTOMER_ENTITY_WEBSITE_ID_STORE_WEBSITE_WEBSITE_ID",
					TableCatalog:               "def",
					TableSchema:                "test",
					TableName:                  "customer_entity",
					ColumnName:                 "website_id",
					OrdinalPosition:            1,
					PositionInUniqueConstraint: null.MakeInt64(1),
					ReferencedTableSchema:      null.MakeString(`test`),
					ReferencedTableName:        null.MakeString(`store_website`),
					ReferencedColumnName:       null.MakeString(`website_id`),
				},
			},
			BeforeMapColumns: nil,
			AfterMapColumns:  nil,
		},
		"store": {
			Data: []*KeyColumnUsage{
				{
					ConstraintCatalog:          "def",
					ConstraintSchema:           "test",
					ConstraintName:             "STORE_GROUP_ID_STORE_GROUP_GROUP_ID",
					TableCatalog:               "def",
					TableSchema:                "test",
					TableName:                  "store",
					ColumnName:                 "group_id",
					OrdinalPosition:            1,
					PositionInUniqueConstraint: null.MakeInt64(1),
					ReferencedTableSchema:      null.MakeString(`test`),
					ReferencedTableName:        null.MakeString(`store_group`),
					ReferencedColumnName:       null.MakeString(`group_id`),
				},
				{
					ConstraintCatalog:          "def",
					ConstraintSchema:           "test",
					ConstraintName:             "STORE_WEBSITE_ID_STORE_WEBSITE_WEBSITE_ID",
					TableCatalog:               "def",
					TableSchema:                "test",
					TableName:                  "store",
					ColumnName:                 "website_id",
					OrdinalPosition:            1,
					PositionInUniqueConstraint: null.MakeInt64(1),
					ReferencedTableSchema:      null.MakeString(`test`),
					ReferencedTableName:        null.MakeString(`store_website`),
					ReferencedColumnName:       null.MakeString(`website_id`),
				},
			},
			BeforeMapColumns: nil,
			AfterMapColumns:  nil,
		},
		"store_group": {
			Data: []*KeyColumnUsage{
				{
					ConstraintCatalog:          "def",
					ConstraintSchema:           "test",
					ConstraintName:             "STORE_GROUP_WEBSITE_ID_STORE_WEBSITE_WEBSITE_ID",
					TableCatalog:               "def",
					TableSchema:                "test",
					TableName:                  "store_group",
					ColumnName:                 "website_id",
					OrdinalPosition:            1,
					PositionInUniqueConstraint: null.MakeInt64(1),
					ReferencedTableSchema:      null.MakeString(`test`),
					ReferencedTableName:        null.MakeString(`store_website`),
					ReferencedColumnName:       null.MakeString(`website_id`),
				},
			},
			BeforeMapColumns: nil,
			AfterMapColumns:  nil,
		},
		"x859admin_passwords": {
			Data: []*KeyColumnUsage{
				{
					ConstraintCatalog:          "def",
					ConstraintSchema:           "test",
					ConstraintName:             "ADMIN_PASSWORDS_USER_ID_ADMIN_USER_USER_ID",
					TableCatalog:               "def",
					TableSchema:                "test",
					TableName:                  "x859admin_passwords",
					ColumnName:                 "user_id",
					OrdinalPosition:            1,
					PositionInUniqueConstraint: null.MakeInt64(1),
					ReferencedTableSchema:      null.MakeString(`test`),
					ReferencedTableName:        null.MakeString(`x859admin_user`),
					ReferencedColumnName:       null.MakeString(`user_id`),
				},
			},
			BeforeMapColumns: nil,
			AfterMapColumns:  nil,
		},
	}

	wantKcucRev := map[string]KeyColumnUsageCollection{
		"x859admin_user": {
			Data: []*KeyColumnUsage{
				{
					ConstraintCatalog:          "def",
					ConstraintSchema:           "test",
					ConstraintName:             "ADMIN_PASSWORDS_USER_ID_ADMIN_USER_USER_ID",
					TableCatalog:               "def",
					TableSchema:                "test",
					TableName:                  "x859admin_user",
					ColumnName:                 "user_id",
					OrdinalPosition:            1,
					PositionInUniqueConstraint: null.MakeInt64(1),
					ReferencedTableSchema:      null.MakeString(`test`),
					ReferencedTableName:        null.MakeString(`x859admin_passwords`),
					ReferencedColumnName:       null.MakeString(`user_id`),
				},
			},
			BeforeMapColumns: nil,
			AfterMapColumns:  nil,
		},
		"x910cms_block": {
			Data: []*KeyColumnUsage{
				{
					ConstraintCatalog:          "def",
					ConstraintSchema:           "test",
					ConstraintName:             "CMS_BLOCK_STORE_BLOCK_ID_CMS_BLOCK_BLOCK_ID",
					TableCatalog:               "def",
					TableSchema:                "test",
					TableName:                  "x910cms_block",
					ColumnName:                 "block_id",
					OrdinalPosition:            1,
					PositionInUniqueConstraint: null.MakeInt64(1),
					ReferencedTableSchema:      null.MakeString(`test`),
					ReferencedTableName:        null.MakeString(`x910cms_block_store`),
					ReferencedColumnName:       null.MakeString(`block_id`),
				},
			},
			BeforeMapColumns: nil,
			AfterMapColumns:  nil,
		},
		"x910cms_page": {
			Data: []*KeyColumnUsage{
				{
					ConstraintCatalog:          "def",
					ConstraintSchema:           "test",
					ConstraintName:             "CMS_PAGE_STORE_PAGE_ID_CMS_PAGE_PAGE_ID",
					TableCatalog:               "def",
					TableSchema:                "test",
					TableName:                  "x910cms_page",
					ColumnName:                 "page_id",
					OrdinalPosition:            1,
					PositionInUniqueConstraint: null.MakeInt64(1),
					ReferencedTableSchema:      null.MakeString(`test`),
					ReferencedTableName:        null.MakeString(`x910cms_page_store`),
					ReferencedColumnName:       null.MakeString(`page_id`),
				},
			},
			BeforeMapColumns: nil,
			AfterMapColumns:  nil,
		},
		"sequence_catalog_category": {
			Data: []*KeyColumnUsage{
				{
					ConstraintCatalog:          "def",
					ConstraintSchema:           "test",
					ConstraintName:             "CAT_CTGR_ENTT_ENTT_ID_SEQUENCE_CAT_CTGR_SEQUENCE_VAL",
					TableCatalog:               "def",
					TableSchema:                "test",
					TableName:                  "sequence_catalog_category",
					ColumnName:                 "sequence_value",
					OrdinalPosition:            1,
					PositionInUniqueConstraint: null.MakeInt64(1),
					ReferencedTableSchema:      null.MakeString(`test`),
					ReferencedTableName:        null.MakeString(`catalog_category_entity`),
					ReferencedColumnName:       null.MakeString(`entity_id`),
				},
			},
			BeforeMapColumns: nil,
			AfterMapColumns:  nil,
		},
		"customer_entity": {
			Data: []*KeyColumnUsage{
				{
					ConstraintCatalog:          "def",
					ConstraintSchema:           "test",
					ConstraintName:             "CUSTOMER_ADDRESS_ENTITY_PARENT_ID_CUSTOMER_ENTITY_ENTITY_ID",
					TableCatalog:               "def",
					TableSchema:                "test",
					TableName:                  "customer_entity",
					ColumnName:                 "entity_id",
					OrdinalPosition:            1,
					PositionInUniqueConstraint: null.MakeInt64(1),
					ReferencedTableSchema:      null.MakeString(`test`),
					ReferencedTableName:        null.MakeString(`customer_address_entity`),
					ReferencedColumnName:       null.MakeString(`parent_id`),
				},
			},
			BeforeMapColumns: nil,
			AfterMapColumns:  nil,
		},
		"store": {
			Data: []*KeyColumnUsage{
				{
					ConstraintCatalog:          "def",
					ConstraintSchema:           "test",
					ConstraintName:             "CMS_BLOCK_STORE_STORE_ID_STORE_STORE_ID",
					TableCatalog:               "def",
					TableSchema:                "test",
					TableName:                  "store",
					ColumnName:                 "store_id",
					OrdinalPosition:            1,
					PositionInUniqueConstraint: null.MakeInt64(1),
					ReferencedTableSchema:      null.MakeString(`test`),
					ReferencedTableName:        null.MakeString(`x910cms_block_store`),
					ReferencedColumnName:       null.MakeString(`store_id`),
				},
				{
					ConstraintCatalog:          "def",
					ConstraintSchema:           "test",
					ConstraintName:             "CMS_PAGE_STORE_STORE_ID_STORE_STORE_ID",
					TableCatalog:               "def",
					TableSchema:                "test",
					TableName:                  "store",
					ColumnName:                 "store_id",
					OrdinalPosition:            1,
					PositionInUniqueConstraint: null.MakeInt64(1),
					ReferencedTableSchema:      null.MakeString(`test`),
					ReferencedTableName:        null.MakeString(`x910cms_page_store`),
					ReferencedColumnName:       null.MakeString(`store_id`),
				},
				{
					ConstraintCatalog:          "def",
					ConstraintSchema:           "test",
					ConstraintName:             "CUSTOMER_ENTITY_STORE_ID_STORE_STORE_ID",
					TableCatalog:               "def",
					TableSchema:                "test",
					TableName:                  "store",
					ColumnName:                 "store_id",
					OrdinalPosition:            1,
					PositionInUniqueConstraint: null.MakeInt64(1),
					ReferencedTableSchema:      null.MakeString(`test`),
					ReferencedTableName:        null.MakeString(`customer_entity`),
					ReferencedColumnName:       null.MakeString(`store_id`),
				},
			},
			BeforeMapColumns: nil,
			AfterMapColumns:  nil,
		},
		"store_website": {
			Data: []*KeyColumnUsage{
				{
					ConstraintCatalog:          "def",
					ConstraintSchema:           "test",
					ConstraintName:             "CUSTOMER_ENTITY_WEBSITE_ID_STORE_WEBSITE_WEBSITE_ID",
					TableCatalog:               "def",
					TableSchema:                "test",
					TableName:                  "store_website",
					ColumnName:                 "website_id",
					OrdinalPosition:            1,
					PositionInUniqueConstraint: null.MakeInt64(1),
					ReferencedTableSchema:      null.MakeString(`test`),
					ReferencedTableName:        null.MakeString(`customer_entity`),
					ReferencedColumnName:       null.MakeString(`website_id`),
				},
				{
					ConstraintCatalog:          "def",
					ConstraintSchema:           "test",
					ConstraintName:             "STORE_GROUP_WEBSITE_ID_STORE_WEBSITE_WEBSITE_ID",
					TableCatalog:               "def",
					TableSchema:                "test",
					TableName:                  "store_website",
					ColumnName:                 "website_id",
					OrdinalPosition:            1,
					PositionInUniqueConstraint: null.MakeInt64(1),
					ReferencedTableSchema:      null.MakeString(`test`),
					ReferencedTableName:        null.MakeString(`store_group`),
					ReferencedColumnName:       null.MakeString(`website_id`),
				},
				{
					ConstraintCatalog:          "def",
					ConstraintSchema:           "test",
					ConstraintName:             "STORE_WEBSITE_ID_STORE_WEBSITE_WEBSITE_ID",
					TableCatalog:               "def",
					TableSchema:                "test",
					TableName:                  "store_website",
					ColumnName:                 "website_id",
					OrdinalPosition:            1,
					PositionInUniqueConstraint: null.MakeInt64(1),
					ReferencedTableSchema:      null.MakeString(`test`),
					ReferencedTableName:        null.MakeString(`store`),
					ReferencedColumnName:       null.MakeString(`website_id`),
				},
			},
			BeforeMapColumns: nil,
			AfterMapColumns:  nil,
		},
		"store_group": {
			Data: []*KeyColumnUsage{
				{
					ConstraintCatalog:          "def",
					ConstraintSchema:           "test",
					ConstraintName:             "STORE_GROUP_ID_STORE_GROUP_GROUP_ID",
					TableCatalog:               "def",
					TableSchema:                "test",
					TableName:                  "store_group",
					ColumnName:                 "group_id",
					OrdinalPosition:            1,
					PositionInUniqueConstraint: null.MakeInt64(1),
					ReferencedTableSchema:      null.MakeString(`test`),
					ReferencedTableName:        null.MakeString(`store`),
					ReferencedColumnName:       null.MakeString(`group_id`),
				},
			},
			BeforeMapColumns: nil,
			AfterMapColumns:  nil,
		},
	}

	kcucRev := ReverseKeyColumnUsage(kcuc)
	// maps are printed in Go 1.12 in an ordered way, otherwise the test would randomly fail.
	assert.Exactly(t, len(wantKcucRev), len(kcucRev))

	for key := range wantKcucRev {
		wantKcucRev[key].Sort()
		kcucRev[key].Sort()
		assert.Exactly(t, wantKcucRev[key], kcucRev[key], "Mismatch at key %q", key)
	}
}
