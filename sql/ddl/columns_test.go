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
	"testing"

	"github.com/corestoreio/pkg/storage/null"
	"github.com/corestoreio/pkg/util/assert"
)

var adminUserColumns = Columns{
	&Column{Field: "user_id", Pos: 1, Null: "NO", DataType: "int", Precision: null.MakeInt64(10), Scale: null.MakeInt64(0), ColumnType: "int(10) unsigned", Key: "PRI", Extra: "auto_increment", Comment: "User ID"},
	&Column{Field: "firstname", Pos: 2, Null: "YES", DataType: "varchar", CharMaxLength: null.MakeInt64(32), ColumnType: "varchar(32)", Comment: "User First Name"},
	&Column{Field: "lastname", Pos: 3, Null: "YES", DataType: "varchar", CharMaxLength: null.MakeInt64(32), ColumnType: "varchar(32)", Comment: "User Last Name"},
	&Column{Field: "email", Pos: 4, Null: "YES", DataType: "varchar", CharMaxLength: null.MakeInt64(128), ColumnType: "varchar(128)", Comment: "User Email"},
	&Column{Field: "username", Pos: 5, Null: "YES", DataType: "varchar", CharMaxLength: null.MakeInt64(40), ColumnType: "varchar(40)", Key: "UNI", Comment: "User Login"},
	&Column{Field: "password", Pos: 6, Null: "NO", DataType: "varchar", CharMaxLength: null.MakeInt64(255), ColumnType: "varchar(255)", Comment: "User Password"},
	&Column{Field: "created", Pos: 7, Default: null.MakeString(`CURRENT_TIMESTAMP`), Null: "NO", DataType: "timestamp", ColumnType: "timestamp", Comment: "User Created Time"},
	&Column{Field: "modified", Pos: 8, Null: "NO", DataType: "timestamp", ColumnType: "timestamp", Extra: "on update CURRENT_TIMESTAMP", Comment: "User Modified Time"},
	&Column{Field: "logdate", Pos: 9, Null: "YES", DataType: "timestamp", ColumnType: "timestamp", Comment: "User Last Login Time"},
	&Column{Field: "lognum", Pos: 10, Default: null.MakeString(`0`), Null: "NO", DataType: "smallint", Precision: null.MakeInt64(5), Scale: null.MakeInt64(0), ColumnType: "smallint(5) unsigned", Comment: "User Login Number"},
	&Column{Field: "reload_acl_flag", Pos: 11, Default: null.MakeString(`0`), Null: "NO", DataType: "smallint", Precision: null.MakeInt64(5), Scale: null.MakeInt64(0), ColumnType: "smallint(6)", Comment: "Reload ACL"},
	&Column{Field: "is_active", Pos: 12, Default: null.MakeString(`1`), Null: "NO", DataType: "smallint", Precision: null.MakeInt64(5), Scale: null.MakeInt64(0), ColumnType: "smallint(6)", Comment: "User Is Active"},
	&Column{Field: "extra", Pos: 13, Null: "YES", DataType: "text", CharMaxLength: null.MakeInt64(65535), ColumnType: "text", Comment: "User Extra Data"},
	&Column{Field: "rp_token", Pos: 14, Null: "YES", DataType: "longtext", CharMaxLength: null.MakeInt64(65535), ColumnType: "text", Comment: "Reset Password Link Token"},
	&Column{Field: "rp_token_created_at", Pos: 15, Null: "YES", DataType: "timestamp", ColumnType: "timestamp", Comment: "Reset Password Link Token Creation Date"},
	&Column{Field: "interface_locale", Pos: 16, Default: null.MakeString(`en_US`), Null: "NO", DataType: "varchar", CharMaxLength: null.MakeInt64(16), ColumnType: "varchar(16)", Comment: "Backend interface locale"},
	&Column{Field: "failures_num", Pos: 17, Default: null.MakeString(`0`), Null: "YES", DataType: "smallint", Precision: null.MakeInt64(5), Scale: null.MakeInt64(0), ColumnType: "smallint(6)", Comment: "Failure Number"},
	&Column{Field: "first_failure", Pos: 18, Null: "YES", DataType: "timestamp", ColumnType: "timestamp", Comment: "First Failure"},
	&Column{Field: "lock_expires", Pos: 19, Null: "YES", DataType: "timestamp", ColumnType: "timestamp", Comment: "Expiration Lock Dates"},
	&Column{Field: "virtual_a", Pos: 20, Null: "YES", DataType: "timestamp", ColumnType: "timestamp", Extra: "VIRTUAL GENERATED", Generated: "ALWAYS", GenerationExpression: null.MakeString("`failures_num` MOD 10")},
	&Column{Field: "stored_b", Pos: 21, Null: "YES", DataType: "timestamp", ColumnType: "timestamp", Extra: "STORED GENERATED", Generated: "ALWAYS", GenerationExpression: null.MakeString("left(`rp_token`,5)")},
	&Column{Field: "version_ts", Pos: 22, Null: "YES", DataType: "timestamp", ColumnType: "timestamp(6)", Extra: "STORED GENERATED", Comment: "Timestamp Start Versioning", Generated: "ALWAYS", GenerationExpression: null.MakeString("ROW START")},
	&Column{Field: "version_te", Pos: 23, Null: "YES", DataType: "timestamp", ColumnType: "timestamp(6)", Key: "PRI", Extra: "STORED GENERATED", Comment: "Timestamp End Versioning", Generated: "ALWAYS", GenerationExpression: null.MakeString("ROW END")},
}

func TestColumnsIsEligibleForUpsert(t *testing.T) {
	assert.False(t, columnsIsEligibleForUpsert(adminUserColumns.ByField("user_id")), "PK")
	assert.False(t, columnsIsEligibleForUpsert(adminUserColumns.ByField("version_ts")), "versioned")
	assert.False(t, columnsIsEligibleForUpsert(adminUserColumns.ByField("virtual_a")), "virtual")
	assert.False(t, columnsIsEligibleForUpsert(adminUserColumns.ByField("stored_b")), "stored_b")
	assert.False(t, columnsIsEligibleForUpsert(adminUserColumns.ByField("created")), "timestamp created")
	assert.False(t, columnsIsEligibleForUpsert(adminUserColumns.ByField("modified")), "timestamp modified")
	assert.True(t, columnsIsEligibleForUpsert(adminUserColumns.ByField("password")), "timestamp modified")
}
