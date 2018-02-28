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
	"fmt"
	"sort"
	"testing"

	"github.com/corestoreio/errors"
	"github.com/corestoreio/pkg/sql/ddl"
	"github.com/corestoreio/pkg/sql/dml"
	"github.com/corestoreio/pkg/sql/dmltest"
	"github.com/stretchr/testify/assert"
)

// Check that type adheres to interfaces
var _ fmt.Stringer = (*ddl.Columns)(nil)
var _ fmt.GoStringer = (*ddl.Columns)(nil)
var _ fmt.GoStringer = (*ddl.Column)(nil)
var _ sort.Interface = (*ddl.Columns)(nil)

func TestLoadColumns_Integration_Mage(t *testing.T) {
	t.Parallel()

	dbc := dmltest.MustConnectDB(t)
	defer dmltest.Close(t, dbc)

	tests := []struct {
		table          string
		want           string
		wantErrKind    errors.Kind
		wantJoinFields string
	}{
		{"core_config_data",
			"ddl.Columns{\n&ddl.Column{Field: \"config_id\", Pos: 1, Null: \"NO\", DataType: \"int\", Precision: dml.MakeNullInt64(10), Scale: dml.MakeNullInt64(0), ColumnType: \"int(10) unsigned\", Key: \"PRI\", Extra: \"auto_increment\", Comment: \"Config Id\", },\n&ddl.Column{Field: \"scope\", Pos: 2, Default: dml.MakeNullString(\"'default'\"), Null: \"NO\", DataType: \"varchar\", CharMaxLength: dml.MakeNullInt64(8), ColumnType: \"varchar(8)\", Key: \"MUL\", Comment: \"Config Scope\", },\n&ddl.Column{Field: \"scope_id\", Pos: 3, Default: dml.MakeNullString(\"0\"), Null: \"NO\", DataType: \"int\", Precision: dml.MakeNullInt64(10), Scale: dml.MakeNullInt64(0), ColumnType: \"int(11)\", Comment: \"Config Scope Id\", },\n&ddl.Column{Field: \"path\", Pos: 4, Default: dml.MakeNullString(\"'general'\"), Null: \"NO\", DataType: \"varchar\", CharMaxLength: dml.MakeNullInt64(255), ColumnType: \"varchar(255)\", Comment: \"Config Path\", },\n&ddl.Column{Field: \"value\", Pos: 5, Default: dml.MakeNullString(\"NULL\"), Null: \"YES\", DataType: \"text\", CharMaxLength: dml.MakeNullInt64(65535), ColumnType: \"text\", Comment: \"Config Value\", },\n}\n",
			errors.NoKind,
			"config_id_scope_scope_id_path_value",
		},
		{"catalog_category_product",
			"ddl.Columns{\n&ddl.Column{Field: \"entity_id\", Pos: 1, Null: \"NO\", DataType: \"int\", Precision: dml.MakeNullInt64(10), Scale: dml.MakeNullInt64(0), ColumnType: \"int(11)\", Key: \"PRI\", Extra: \"auto_increment\", Comment: \"Entity ID\", },\n&ddl.Column{Field: \"category_id\", Pos: 2, Default: dml.MakeNullString(\"0\"), Null: \"NO\", DataType: \"int\", Precision: dml.MakeNullInt64(10), Scale: dml.MakeNullInt64(0), ColumnType: \"int(10) unsigned\", Key: \"PRI\", Comment: \"Category ID\", },\n&ddl.Column{Field: \"product_id\", Pos: 3, Default: dml.MakeNullString(\"0\"), Null: \"NO\", DataType: \"int\", Precision: dml.MakeNullInt64(10), Scale: dml.MakeNullInt64(0), ColumnType: \"int(10) unsigned\", Key: \"PRI\", Comment: \"Product ID\", },\n&ddl.Column{Field: \"position\", Pos: 4, Default: dml.MakeNullString(\"0\"), Null: \"NO\", DataType: \"int\", Precision: dml.MakeNullInt64(10), Scale: dml.MakeNullInt64(0), ColumnType: \"int(11)\", Comment: \"Position\", },\n}\n",
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
			assert.NoError(t, err, "Index %d", i)
			assert.Equal(t, test.want, fmt.Sprintf("%#v\n", cols1), "Index %d", i)
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
			"ddl.Columns{\n&ddl.Column{Field: \"category_id\", Default: dml.MakeNullString(\"0\"), ColumnType: \"int(10) unsigned\", Key: \"MUL\", Aliases: []string{\"entity_id\"}, Uniquified: true, StructTag: \"json:\\\",omitempty\\\"\", },\n&ddl.Column{Field: \"path\", Null: \"YES\", ColumnType: \"varchar(255)\", Key: \"MUL\", },\n}",
		},
		{
			tableMap.MustTable("catalog_category_anc_categs_index_tmp").Columns.PrimaryKeys().Len(),
			1,
			tableMap.MustTable("catalog_category_anc_categs_index_tmp").Columns.GoString(),
			"ddl.Columns{\n&ddl.Column{Field: \"category_id\", Default: dml.MakeNullString(\"0\"), ColumnType: \"int(10) unsigned\", Key: \"PRI\", },\n&ddl.Column{Field: \"path\", Null: \"YES\", ColumnType: \"varchar(255)\", },\n}",
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
			Default:    dml.MakeNullString(`0`),
			Extra:      ``,
		},
		&ddl.Column{
			Field:      `is_root_category180`,
			ColumnType: `smallint2 unsigned`,
			Null:       "YES",
			Key:        ``,
			Default:    dml.MakeNullString(`0`),
			Extra:      ``,
		},
	}
	colsHave := cols.Filter(func(c *ddl.Column) bool {
		return c.Field == "is_root_category180"
	})
	colsWant := ddl.Columns{
		&ddl.Column{Field: `is_root_category180`, ColumnType: `smallint2 unsigned`, Null: "YES", Key: ``, Default: dml.MakeNullString(`0`), Extra: ``},
	}

	assert.Exactly(t, colsWant, colsHave)
}

var adminUserColumns = ddl.Columns{
	&ddl.Column{Field: "user_id", Pos: 1, Default: dml.NullString{}, Null: "NO", DataType: "int", CharMaxLength: dml.NullInt64{}, Precision: dml.MakeNullInt64(10), Scale: dml.MakeNullInt64(0), ColumnType: "int(10) unsigned", Key: "PRI", Extra: "auto_increment", Comment: "User ID"},
	&ddl.Column{Field: "firstname", Pos: 2, Default: dml.NullString{}, Null: "YES", DataType: "varchar", CharMaxLength: dml.MakeNullInt64(32), Precision: dml.NullInt64{}, Scale: dml.NullInt64{}, ColumnType: "varchar(32)", Key: "", Extra: "", Comment: "User First Name"},
	&ddl.Column{Field: "lastname", Pos: 3, Default: dml.NullString{}, Null: "YES", DataType: "varchar", CharMaxLength: dml.MakeNullInt64(32), Precision: dml.NullInt64{}, Scale: dml.NullInt64{}, ColumnType: "varchar(32)", Key: "", Extra: "", Comment: "User Last Name"},
	&ddl.Column{Field: "email", Pos: 4, Default: dml.NullString{}, Null: "YES", DataType: "varchar", CharMaxLength: dml.MakeNullInt64(128), Precision: dml.NullInt64{}, Scale: dml.NullInt64{}, ColumnType: "varchar(128)", Key: "", Extra: "", Comment: "User Email"},
	&ddl.Column{Field: "username", Pos: 5, Default: dml.NullString{}, Null: "YES", DataType: "varchar", CharMaxLength: dml.MakeNullInt64(40), Precision: dml.NullInt64{}, Scale: dml.NullInt64{}, ColumnType: "varchar(40)", Key: "UNI", Extra: "", Comment: "User Login"},
	&ddl.Column{Field: "password", Pos: 6, Default: dml.NullString{}, Null: "NO", DataType: "varchar", CharMaxLength: dml.MakeNullInt64(255), Precision: dml.NullInt64{}, Scale: dml.NullInt64{}, ColumnType: "varchar(255)", Key: "", Extra: "", Comment: "User Password"},
	&ddl.Column{Field: "created", Pos: 7, Default: dml.MakeNullString(`0000-00-00 00:00:00`), Null: "NO", DataType: "timestamp", CharMaxLength: dml.NullInt64{}, Precision: dml.NullInt64{}, Scale: dml.NullInt64{}, ColumnType: "timestamp", Key: "", Extra: "", Comment: "User Created Time"},
	&ddl.Column{Field: "modified", Pos: 8, Default: dml.MakeNullString(`CURRENT_TIMESTAMP`), Null: "NO", DataType: "timestamp", CharMaxLength: dml.NullInt64{}, Precision: dml.NullInt64{}, Scale: dml.NullInt64{}, ColumnType: "timestamp", Key: "", Extra: "on update CURRENT_TIMESTAMP", Comment: "User Modified Time"},
	&ddl.Column{Field: "logdate", Pos: 9, Default: dml.NullString{}, Null: "YES", DataType: "timestamp", CharMaxLength: dml.NullInt64{}, Precision: dml.NullInt64{}, Scale: dml.NullInt64{}, ColumnType: "timestamp", Key: "", Extra: "", Comment: "User Last Login Time"},
	&ddl.Column{Field: "lognum", Pos: 10, Default: dml.MakeNullString(`0`), Null: "NO", DataType: "smallint", CharMaxLength: dml.NullInt64{}, Precision: dml.MakeNullInt64(5), Scale: dml.MakeNullInt64(0), ColumnType: "smallint(5) unsigned", Key: "", Extra: "", Comment: "User Login Number"},
	&ddl.Column{Field: "reload_acl_flag", Pos: 11, Default: dml.MakeNullString(`0`), Null: "NO", DataType: "smallint", CharMaxLength: dml.NullInt64{}, Precision: dml.MakeNullInt64(5), Scale: dml.MakeNullInt64(0), ColumnType: "smallint(6)", Key: "", Extra: "", Comment: "Reload ACL"},
	&ddl.Column{Field: "is_active", Pos: 12, Default: dml.MakeNullString(`1`), Null: "NO", DataType: "smallint", CharMaxLength: dml.NullInt64{}, Precision: dml.MakeNullInt64(5), Scale: dml.MakeNullInt64(0), ColumnType: "smallint(6)", Key: "", Extra: "", Comment: "User Is Active"},
	&ddl.Column{Field: "extra", Pos: 13, Default: dml.NullString{}, Null: "YES", DataType: "text", CharMaxLength: dml.MakeNullInt64(65535), Precision: dml.NullInt64{}, Scale: dml.NullInt64{}, ColumnType: "text", Key: "", Extra: "", Comment: "User Extra Data"},
	&ddl.Column{Field: "rp_token", Pos: 14, Default: dml.NullString{}, Null: "YES", DataType: "text", CharMaxLength: dml.MakeNullInt64(65535), Precision: dml.NullInt64{}, Scale: dml.NullInt64{}, ColumnType: "text", Key: "", Extra: "", Comment: "Reset Password Link Token"},
	&ddl.Column{Field: "rp_token_created_at", Pos: 15, Default: dml.NullString{}, Null: "YES", DataType: "timestamp", CharMaxLength: dml.NullInt64{}, Precision: dml.NullInt64{}, Scale: dml.NullInt64{}, ColumnType: "timestamp", Key: "", Extra: "", Comment: "Reset Password Link Token Creation Date"},
	&ddl.Column{Field: "interface_locale", Pos: 16, Default: dml.MakeNullString(`en_US`), Null: "NO", DataType: "varchar", CharMaxLength: dml.MakeNullInt64(16), Precision: dml.NullInt64{}, Scale: dml.NullInt64{}, ColumnType: "varchar(16)", Key: "", Extra: "", Comment: "Backend interface locale"},
	&ddl.Column{Field: "failures_num", Pos: 17, Default: dml.MakeNullString(`0`), Null: "YES", DataType: "smallint", CharMaxLength: dml.NullInt64{}, Precision: dml.MakeNullInt64(5), Scale: dml.MakeNullInt64(0), ColumnType: "smallint(6)", Key: "", Extra: "", Comment: "Failure Number"},
	&ddl.Column{Field: "first_failure", Pos: 18, Default: dml.NullString{}, Null: "YES", DataType: "timestamp", CharMaxLength: dml.NullInt64{}, Precision: dml.NullInt64{}, Scale: dml.NullInt64{}, ColumnType: "timestamp", Key: "", Extra: "", Comment: "First Failure"},
	&ddl.Column{Field: "lock_expires", Pos: 19, Default: dml.NullString{}, Null: "YES", DataType: "timestamp", CharMaxLength: dml.NullInt64{}, Precision: dml.NullInt64{}, Scale: dml.NullInt64{}, ColumnType: "timestamp", Key: "", Extra: "", Comment: "Expiration Lock Dates"},
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
