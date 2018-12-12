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

package ddl_test

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"testing"

	"github.com/corestoreio/errors"
	"github.com/corestoreio/pkg/sql/ddl"
	"github.com/corestoreio/pkg/sql/dml"
	"github.com/corestoreio/pkg/sql/dmltest"
	"github.com/corestoreio/pkg/storage/null"
	"github.com/corestoreio/pkg/util/assert"
)

// Check that type adheres to interfaces
var _ fmt.Stringer = (*ddl.Columns)(nil)
var _ fmt.GoStringer = (*ddl.Columns)(nil)
var _ fmt.GoStringer = (*ddl.Column)(nil)
var _ sort.Interface = (*ddl.Columns)(nil)

func TestLoadColumns_Integration(t *testing.T) {
	t.Parallel()

	ctx := context.TODO()
	dbc := dmltest.MustConnectDB(t,
		dml.WithExecSQLOnConnOpen(ctx,
			"DROP TABLE IF EXISTS `core_config_data_test3`;",
			"DROP TABLE IF EXISTS `catalog_category_product_test4`;",
			`CREATE TABLE core_config_data_test3 (
  config_id int(10) unsigned NOT NULL AUTO_INCREMENT COMMENT 'Config Id',
  scope varchar(8) NOT NULL DEFAULT 'default' COMMENT 'Config Scope',
  scope_id int(11) NOT NULL DEFAULT 0 COMMENT 'Config Scope Id',
  path varchar(255) NOT NULL DEFAULT 'general' COMMENT 'Config Path',
  value text DEFAULT NULL COMMENT 'Config Value',
  PRIMARY KEY (config_id),
  UNIQUE KEY CORE_CONFIG_DATA_SCOPE_SCOPE_ID_PATH (scope,scope_id,path)
) ENGINE=InnoDB AUTO_INCREMENT=0 DEFAULT CHARSET=utf8 COMMENT='Config Data';`,

			`CREATE TABLE catalog_category_product_test4 (
  entity_id int(11) NOT NULL AUTO_INCREMENT COMMENT 'Entity ID',
  category_id int(10) unsigned NOT NULL DEFAULT 0 COMMENT 'Category ID',
  product_id int(10) unsigned NOT NULL DEFAULT 0 COMMENT 'Product ID',
  position int(11) NOT NULL DEFAULT 0 COMMENT 'Position',
  PRIMARY KEY (entity_id,category_id,product_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='Catalog Product To Category Linkage Table';
`),
		dml.WithExecSQLOnConnClose(ctx,
			"DROP TABLE IF EXISTS `core_config_data_test3`;",
			"DROP TABLE IF EXISTS `catalog_category_product_test4`;",
		),
	) // create DB
	defer dmltest.Close(t, dbc)

	tests := []struct {
		table          string
		want           string
		wantErrKind    errors.Kind
		wantJoinFields string
	}{
		{"core_config_data_test3",
			"[{\"Field\":\"config_id\",\"Pos\":1,\"Default\":null,\"Null\":\"NO\",\"DataType\":\"int\",\"CharMaxLength\":null,\"Precision\":10,\"Scale\":0,\"ColumnType\":\"int(10) unsigned\",\"Key\":\"PRI\",\"Extra\":\"auto_increment\",\"Comment\":\"Config Id\",\"Generated\":\"NEVER\",\"GenerationExpression\":null,\"Aliases\":null,\"Uniquified\":false,\"StructTag\":\"\"},{\"Field\":\"scope\",\"Pos\":2,\"Default\":\"'default'\",\"Null\":\"NO\",\"DataType\":\"varchar\",\"CharMaxLength\":8,\"Precision\":null,\"Scale\":null,\"ColumnType\":\"varchar(8)\",\"Key\":\"MUL\",\"Extra\":\"\",\"Comment\":\"Config Scope\",\"Generated\":\"NEVER\",\"GenerationExpression\":null,\"Aliases\":null,\"Uniquified\":false,\"StructTag\":\"\"},{\"Field\":\"scope_id\",\"Pos\":3,\"Default\":\"0\",\"Null\":\"NO\",\"DataType\":\"int\",\"CharMaxLength\":null,\"Precision\":10,\"Scale\":0,\"ColumnType\":\"int(11)\",\"Key\":\"\",\"Extra\":\"\",\"Comment\":\"Config Scope Id\",\"Generated\":\"NEVER\",\"GenerationExpression\":null,\"Aliases\":null,\"Uniquified\":false,\"StructTag\":\"\"},{\"Field\":\"path\",\"Pos\":4,\"Default\":\"'general'\",\"Null\":\"NO\",\"DataType\":\"varchar\",\"CharMaxLength\":255,\"Precision\":null,\"Scale\":null,\"ColumnType\":\"varchar(255)\",\"Key\":\"\",\"Extra\":\"\",\"Comment\":\"Config Path\",\"Generated\":\"NEVER\",\"GenerationExpression\":null,\"Aliases\":null,\"Uniquified\":false,\"StructTag\":\"\"},{\"Field\":\"value\",\"Pos\":5,\"Default\":\"NULL\",\"Null\":\"YES\",\"DataType\":\"text\",\"CharMaxLength\":65535,\"Precision\":null,\"Scale\":null,\"ColumnType\":\"text\",\"Key\":\"\",\"Extra\":\"\",\"Comment\":\"Config Value\",\"Generated\":\"NEVER\",\"GenerationExpression\":null,\"Aliases\":null,\"Uniquified\":false,\"StructTag\":\"\"}]",
			errors.NoKind,
			"config_id_scope_scope_id_path_value",
		},
		{"catalog_category_product_test4",
			"[{\"Field\":\"entity_id\",\"Pos\":1,\"Default\":null,\"Null\":\"NO\",\"DataType\":\"int\",\"CharMaxLength\":null,\"Precision\":10,\"Scale\":0,\"ColumnType\":\"int(11)\",\"Key\":\"PRI\",\"Extra\":\"auto_increment\",\"Comment\":\"Entity ID\",\"Generated\":\"NEVER\",\"GenerationExpression\":null,\"Aliases\":null,\"Uniquified\":false,\"StructTag\":\"\"},{\"Field\":\"category_id\",\"Pos\":2,\"Default\":\"0\",\"Null\":\"NO\",\"DataType\":\"int\",\"CharMaxLength\":null,\"Precision\":10,\"Scale\":0,\"ColumnType\":\"int(10) unsigned\",\"Key\":\"PRI\",\"Extra\":\"\",\"Comment\":\"Category ID\",\"Generated\":\"NEVER\",\"GenerationExpression\":null,\"Aliases\":null,\"Uniquified\":false,\"StructTag\":\"\"},{\"Field\":\"product_id\",\"Pos\":3,\"Default\":\"0\",\"Null\":\"NO\",\"DataType\":\"int\",\"CharMaxLength\":null,\"Precision\":10,\"Scale\":0,\"ColumnType\":\"int(10) unsigned\",\"Key\":\"PRI\",\"Extra\":\"\",\"Comment\":\"Product ID\",\"Generated\":\"NEVER\",\"GenerationExpression\":null,\"Aliases\":null,\"Uniquified\":false,\"StructTag\":\"\"},{\"Field\":\"position\",\"Pos\":4,\"Default\":\"0\",\"Null\":\"NO\",\"DataType\":\"int\",\"CharMaxLength\":null,\"Precision\":10,\"Scale\":0,\"ColumnType\":\"int(11)\",\"Key\":\"\",\"Extra\":\"\",\"Comment\":\"Position\",\"Generated\":\"NEVER\",\"GenerationExpression\":null,\"Aliases\":null,\"Uniquified\":false,\"StructTag\":\"\"}]",
			errors.NoKind,
			"entity_id_category_id_product_id_position",
		},
		{"non_existent_table",
			"",
			errors.NotFound,
			"",
		},
	}

	for i, test := range tests {
		tc, err := ddl.LoadColumns(context.TODO(), dbc.DB, test.table)
		cols1 := tc[test.table]
		if !test.wantErrKind.Empty() {
			assert.True(t, test.wantErrKind.Match(err), "Index %d", i)
		} else {
			assert.NoError(t, err, "Index %d\n%+v", i, err)
			data, err := json.Marshal(cols1)
			assert.NoError(t, err)
			assert.Equal(t, test.want, string(data), "Index %d\n%q", i, data)
			assert.Equal(t, test.wantJoinFields, cols1.JoinFields("_"), "Index %d", i)
		}
	}
}

func TestColumns(t *testing.T) {
	t.Parallel()
	tests := []struct {
		have  int
		want  int
		haveS string
		wantS string
	}{
		{
			tableMap.MustTable("catalog_category_anc_categs_index_idx").Columns.PrimaryKeys().Len(),
			0,
			tableMap.MustTable("catalog_category_anc_categs_index_idx").Columns.GoString(),
			"ddl.Columns{\n&ddl.Column{Field: \"category_id\", Default: null.MakeString(\"0\"), ColumnType: \"int(10) unsigned\", Key: \"MUL\", Aliases: []string{\"entity_id\"}, Uniquified: true, StructTag: \"json:\\\",omitempty\\\"\", },\n&ddl.Column{Field: \"path\", Null: \"YES\", ColumnType: \"varchar(255)\", Key: \"MUL\", },\n}",
		},
		{
			tableMap.MustTable("catalog_category_anc_categs_index_tmp").Columns.PrimaryKeys().Len(),
			1,
			tableMap.MustTable("catalog_category_anc_categs_index_tmp").Columns.GoString(),
			"ddl.Columns{\n&ddl.Column{Field: \"category_id\", Default: null.MakeString(\"0\"), ColumnType: \"int(10) unsigned\", Key: \"PRI\", },\n&ddl.Column{Field: \"path\", Null: \"YES\", ColumnType: \"varchar(255)\", },\n}",
		},
		{
			tableMap.MustTable("admin_user").Columns.UniqueKeys().Len(), 1,
			tableMap.MustTable("admin_user").Columns.GoString(),
			"ddl.Columns{\n&ddl.Column{Field: \"user_id\", ColumnType: \"int(10) unsigned\", Key: \"PRI\", Extra: \"auto_increment\", },\n&ddl.Column{Field: \"email\", Null: \"YES\", ColumnType: \"varchar(128)\", },\n&ddl.Column{Field: \"first_name\", ColumnType: \"varchar(255)\", },\n&ddl.Column{Field: \"username\", Null: \"YES\", ColumnType: \"varchar(40)\", Key: \"UNI\", },\n}",
		},
		{tableMap.MustTable("admin_user").Columns.PrimaryKeys().Len(), 1, "", ""},
	}

	for i, test := range tests {
		assert.Equal(t, test.want, test.have, "Incorrect length at index %d", i)
		assert.Equal(t, test.wantS, test.haveS, "Index %d", i)
	}

	tsN := tableMap.MustTable("admin_user").Columns.ByField("user_id_not_found")
	assert.NotNil(t, tsN)
	assert.Empty(t, tsN.Field)

	ts4 := tableMap.MustTable("admin_user").Columns.ByField("user_id")
	assert.NotEmpty(t, ts4.Field)
	assert.True(t, ts4.IsAutoIncrement())

	ts4b := tableMap.MustTable("admin_user").Columns.ByField("email")
	assert.NotEmpty(t, ts4b.Field)
	assert.True(t, ts4b.IsNull())

	assert.True(t, tableMap.MustTable("admin_user").Columns.First().IsPK())
	emptyTS := &ddl.Table{}
	assert.False(t, emptyTS.Columns.First().IsPK())
}

func TestColumnsFilter(t *testing.T) {
	t.Parallel()
	cols := ddl.Columns{
		&ddl.Column{
			Field:      `category_id131`,
			ColumnType: `int10 unsigned`,
			Key:        `PRI`,
			Default:    null.MakeString(`0`),
			Extra:      ``,
		},
		&ddl.Column{
			Field:      `is_root_category180`,
			ColumnType: `smallint2 unsigned`,
			Null:       "YES",
			Key:        ``,
			Default:    null.MakeString(`0`),
			Extra:      ``,
		},
	}
	colsHave := cols.Filter(func(c *ddl.Column) bool {
		return c.Field == "is_root_category180"
	})
	colsWant := ddl.Columns{
		&ddl.Column{Field: `is_root_category180`, ColumnType: `smallint2 unsigned`, Null: "YES", Key: ``, Default: null.MakeString(`0`), Extra: ``},
	}

	assert.Exactly(t, colsWant, colsHave)
}

var adminUserColumns = ddl.Columns{
	&ddl.Column{Field: "user_id", Pos: 1, Default: null.String{}, Null: "NO", DataType: "int", Precision: null.MakeInt64(10), Scale: null.MakeInt64(0), ColumnType: "int(10) unsigned", Key: "PRI", Extra: "auto_increment", Comment: "User ID"},
	&ddl.Column{Field: "firstname", Pos: 2, Default: null.String{}, Null: "YES", DataType: "varchar", CharMaxLength: null.MakeInt64(32), ColumnType: "varchar(32)", Comment: "User First Name"},
	&ddl.Column{Field: "lastname", Pos: 3, Default: null.String{}, Null: "YES", DataType: "varchar", CharMaxLength: null.MakeInt64(32), ColumnType: "varchar(32)", Comment: "User Last Name"},
	&ddl.Column{Field: "email", Pos: 4, Default: null.String{}, Null: "YES", DataType: "varchar", CharMaxLength: null.MakeInt64(128), ColumnType: "varchar(128)", Comment: "User Email"},
	&ddl.Column{Field: "username", Pos: 5, Default: null.String{}, Null: "YES", DataType: "varchar", CharMaxLength: null.MakeInt64(40), ColumnType: "varchar(40)", Key: "UNI", Comment: "User Login"},
	&ddl.Column{Field: "password", Pos: 6, Default: null.String{}, Null: "NO", DataType: "varchar", CharMaxLength: null.MakeInt64(255), ColumnType: "varchar(255)", Comment: "User Password"},
	&ddl.Column{Field: "created", Pos: 7, Default: null.MakeString(`0000-00-00 00:00:00`), Null: "NO", DataType: "timestamp", ColumnType: "timestamp", Comment: "User Created Time"},
	&ddl.Column{Field: "modified", Pos: 8, Default: null.MakeString(`CURRENT_TIMESTAMP`), Null: "NO", DataType: "timestamp", ColumnType: "timestamp", Extra: "on update CURRENT_TIMESTAMP", Comment: "User Modified Time"},
	&ddl.Column{Field: "logdate", Pos: 9, Default: null.String{}, Null: "YES", DataType: "timestamp", ColumnType: "timestamp", Comment: "User Last Login Time"},
	&ddl.Column{Field: "lognum", Pos: 10, Default: null.MakeString(`0`), Null: "NO", DataType: "smallint", Precision: null.MakeInt64(5), Scale: null.MakeInt64(0), ColumnType: "smallint(5) unsigned", Comment: "User Login Number"},
	&ddl.Column{Field: "reload_acl_flag", Pos: 11, Default: null.MakeString(`0`), Null: "NO", DataType: "smallint", Precision: null.MakeInt64(5), Scale: null.MakeInt64(0), ColumnType: "smallint(6)", Comment: "Reload ACL"},
	&ddl.Column{Field: "is_active", Pos: 12, Default: null.MakeString(`1`), Null: "NO", DataType: "smallint", Precision: null.MakeInt64(5), Scale: null.MakeInt64(0), ColumnType: "smallint(6)", Comment: "User Is Active"},
	&ddl.Column{Field: "extra", Pos: 13, Default: null.String{}, Null: "YES", DataType: "text", CharMaxLength: null.MakeInt64(65535), ColumnType: "text", Comment: "User Extra Data"},
	&ddl.Column{Field: "rp_token", Pos: 14, Default: null.String{}, Null: "YES", DataType: "longtext", CharMaxLength: null.MakeInt64(65535), ColumnType: "text", Comment: "Reset Password Link Token"},
	&ddl.Column{Field: "rp_token_created_at", Pos: 15, Default: null.String{}, Null: "YES", DataType: "timestamp", ColumnType: "timestamp", Comment: "Reset Password Link Token Creation Date"},
	&ddl.Column{Field: "interface_locale", Pos: 16, Default: null.MakeString(`en_US`), Null: "NO", DataType: "varchar", CharMaxLength: null.MakeInt64(16), ColumnType: "varchar(16)", Comment: "Backend interface locale"},
	&ddl.Column{Field: "failures_num", Pos: 17, Default: null.MakeString(`0`), Null: "YES", DataType: "smallint", Precision: null.MakeInt64(5), Scale: null.MakeInt64(0), ColumnType: "smallint(6)", Comment: "Failure Number"},
	&ddl.Column{Field: "first_failure", Pos: 18, Default: null.String{}, Null: "YES", DataType: "timestamp", ColumnType: "timestamp", Comment: "First Failure"},
	&ddl.Column{Field: "lock_expires", Pos: 19, Default: null.String{}, Null: "YES", DataType: "timestamp", ColumnType: "timestamp", Comment: "Expiration Lock Dates"},
	&ddl.Column{Field: "virtual_a", Pos: 20, Default: null.String{}, Null: "YES", DataType: "timestamp", ColumnType: "timestamp", Extra: "VIRTUAL GENERATED", Generated: "ALWAYS", GenerationExpression: null.MakeString("`failures_num` MOD 10")},
	&ddl.Column{Field: "stored_b", Pos: 21, Default: null.String{}, Null: "YES", DataType: "timestamp", ColumnType: "timestamp", Extra: "STORED GENERATED", Generated: "ALWAYS", GenerationExpression: null.MakeString("left(`rp_token`,5)")},
	&ddl.Column{Field: "version_ts", Pos: 22, Default: null.String{}, Null: "YES", DataType: "timestamp", ColumnType: "timestamp(6)", Extra: "STORED GENERATED", Comment: "Timestamp Start Versioning", Generated: "ALWAYS", GenerationExpression: null.MakeString("ROW START")},
	&ddl.Column{Field: "version_te", Pos: 23, Default: null.String{}, Null: "YES", DataType: "timestamp", ColumnType: "timestamp(6)", Key: "PRI", Extra: "STORED GENERATED", Comment: "Timestamp End Versioning", Generated: "ALWAYS", GenerationExpression: null.MakeString("ROW END")},
}

func TestColumns_UniqueColumns(t *testing.T) {
	t.Parallel()
	assert.Exactly(t, []string{"user_id", "username"}, adminUserColumns.UniqueColumns().FieldNames())
}

func TestColumnsSort(t *testing.T) {
	t.Parallel()
	sort.Sort(adminUserColumns)
	assert.Exactly(t, `user_id`, adminUserColumns.First().Field)
}

func TestColumn_GoComment(t *testing.T) {
	t.Parallel()

	assert.Exactly(t, "// user_id int(10) unsigned NOT NULL PRI  auto_increment \"User ID\"",
		adminUserColumns.First().GoComment())
	assert.Exactly(t, "// firstname varchar(32) NULL    \"User First Name\"",
		adminUserColumns.ByField("firstname").GoComment())

}

func TestColumn_IsUnsigned(t *testing.T) {
	t.Parallel()
	assert.True(t, adminUserColumns.ByField("lognum").IsUnsigned())
	assert.False(t, adminUserColumns.ByField("reload_acl_flag").IsUnsigned())
}

func TestColumn_IsCurrentTimestamp(t *testing.T) {
	t.Parallel()
	assert.True(t, adminUserColumns.ByField("modified").IsCurrentTimestamp())
	assert.False(t, adminUserColumns.ByField("reload_acl_flag").IsCurrentTimestamp())
}

func TestColumn_IsGenerated(t *testing.T) {
	t.Parallel()
	assert.True(t, adminUserColumns.ByField("virtual_a").IsGenerated())
	assert.True(t, adminUserColumns.ByField("stored_b").IsGenerated())
	assert.False(t, adminUserColumns.ByField("reload_acl_flag").IsGenerated())
}

func TestColumn_IsSystemVersioned(t *testing.T) {
	t.Parallel()
	assert.True(t, adminUserColumns.ByField("version_ts").IsSystemVersioned())
	assert.True(t, adminUserColumns.ByField("version_te").IsSystemVersioned())
	assert.False(t, adminUserColumns.ByField("reload_acl_flag").IsSystemVersioned())
	assert.False(t, adminUserColumns.ByField("stored_b").IsSystemVersioned())
	assert.False(t, adminUserColumns.ByField("virtual_a").IsSystemVersioned())
}

func TestColumn_IsString(t *testing.T) {
	t.Parallel()
	assert.False(t, adminUserColumns.ByField("version_ts").IsString())
	assert.True(t, adminUserColumns.ByField("firstname").IsString())
	assert.True(t, adminUserColumns.ByField("extra").IsString())
}

func TestColumn_IsBlobDataType(t *testing.T) {
	t.Parallel()
	assert.False(t, adminUserColumns.ByField("version_ts").IsBlobDataType(), "version_ts")
	assert.False(t, adminUserColumns.ByField("firstname").IsBlobDataType(), "firstname")
	assert.True(t, adminUserColumns.ByField("extra").IsBlobDataType(), "extra")
	assert.True(t, adminUserColumns.ByField("rp_token").IsBlobDataType(), "rp_token")

}
