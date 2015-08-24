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

package csdb_test

import (
	"testing"

	"errors"
	"fmt"

	"github.com/corestoreio/csfw/storage/csdb"
	"github.com/stretchr/testify/assert"
)

func TestGetColumns(t *testing.T) {
	// @todo fix test data retrieving from database ...
	dbc := csdb.MustConnectTest()
	defer dbc.Close()
	sess := dbc.NewSession()

	tests := []struct {
		table   string
		want    string
		wantErr error
	}{
		{"core_config_data",
			`csdb.Columns{csdb.Column{Field:sql.NullString{String:"config_id", Valid:true}, Type:sql.NullString{String:"int(10) unsigned", Valid:true}, Null:sql.NullString{String:"NO", Valid:true}, Key:sql.NullString{String:"PRI", Valid:true}, Default:sql.NullString{String:"", Valid:false}, Extra:sql.NullString{String:"auto_increment", Valid:true}}, csdb.Column{Field:sql.NullString{String:"scope", Valid:true}, Type:sql.NullString{String:"varchar(8)", Valid:true}, Null:sql.NullString{String:"NO", Valid:true}, Key:sql.NullString{String:"MUL", Valid:true}, Default:sql.NullString{String:"default", Valid:true}, Extra:sql.NullString{String:"", Valid:true}}, csdb.Column{Field:sql.NullString{String:"scope_id", Valid:true}, Type:sql.NullString{String:"int(11)", Valid:true}, Null:sql.NullString{String:"NO", Valid:true}, Key:sql.NullString{String:"", Valid:true}, Default:sql.NullString{String:"0", Valid:true}, Extra:sql.NullString{String:"", Valid:true}}, csdb.Column{Field:sql.NullString{String:"path", Valid:true}, Type:sql.NullString{String:"varchar(255)", Valid:true}, Null:sql.NullString{String:"NO", Valid:true}, Key:sql.NullString{String:"", Valid:true}, Default:sql.NullString{String:"general", Valid:true}, Extra:sql.NullString{String:"", Valid:true}}, csdb.Column{Field:sql.NullString{String:"value", Valid:true}, Type:sql.NullString{String:"text", Valid:true}, Null:sql.NullString{String:"YES", Valid:true}, Key:sql.NullString{String:"", Valid:true}, Default:sql.NullString{String:"", Valid:false}, Extra:sql.NullString{String:"", Valid:true}}}` + "\n",
			nil,
		},
		{"catalog_category_product",
			`csdb.Columns{csdb.Column{Field:sql.NullString{String:"category_id", Valid:true}, Type:sql.NullString{String:"int(10) unsigned", Valid:true}, Null:sql.NullString{String:"NO", Valid:true}, Key:sql.NullString{String:"PRI", Valid:true}, Default:sql.NullString{String:"0", Valid:true}, Extra:sql.NullString{String:"", Valid:true}}, csdb.Column{Field:sql.NullString{String:"product_id", Valid:true}, Type:sql.NullString{String:"int(10) unsigned", Valid:true}, Null:sql.NullString{String:"NO", Valid:true}, Key:sql.NullString{String:"PRI", Valid:true}, Default:sql.NullString{String:"0", Valid:true}, Extra:sql.NullString{String:"", Valid:true}}, csdb.Column{Field:sql.NullString{String:"position", Valid:true}, Type:sql.NullString{String:"int(11)", Valid:true}, Null:sql.NullString{String:"NO", Valid:true}, Key:sql.NullString{String:"", Valid:true}, Default:sql.NullString{String:"0", Valid:true}, Extra:sql.NullString{String:"", Valid:true}}}` + "\n",
			nil,
		},
		{"non_existent",
			"",
			errors.New("non_existent"),
		},
	}

	for _, test := range tests {
		cols1, err1 := csdb.GetColumns(sess, test.table)
		if test.wantErr != nil {
			assert.Error(t, err1)
			assert.Contains(t, err1.Error(), test.wantErr.Error())
			//t.Logf("%s\n%#v\n", err1.Error(), err1.(errgo.Locationer).Location())
		} else {
			assert.NoError(t, err1)
			assert.Equal(t, test.want, fmt.Sprintf("%#v\n", cols1))
		}
	}
}

func TestColumns(t *testing.T) {

	tests := []struct {
		have  int
		want  int
		haveS string
		wantS string
	}{
		{
			mustStructure(table1).Columns.PrimaryKeys().Len(), 0,
			mustStructure(table1).Columns.String(),
			"csdb.Column{\n    Field:   sql.NullString{String:\"category_id\", Valid:true},\n    Type:    sql.NullString{String:\"int(10) unsigned\", Valid:true},\n    Null:    sql.NullString{String:\"NO\", Valid:true},\n    Key:     sql.NullString{String:\"MUL\", Valid:true},\n    Default: sql.NullString{String:\"0\", Valid:true},\n    Extra:   sql.NullString{String:\"\", Valid:true},\n},\ncsdb.Column{\n    Field:   sql.NullString{String:\"path\", Valid:true},\n    Type:    sql.NullString{String:\"varchar(255)\", Valid:true},\n    Null:    sql.NullString{String:\"YES\", Valid:true},\n    Key:     sql.NullString{String:\"MUL\", Valid:true},\n    Default: sql.NullString{},\n    Extra:   sql.NullString{String:\"\", Valid:true},\n}",
		},
		{
			mustStructure(table2).Columns.PrimaryKeys().Len(), 1,
			mustStructure(table2).Columns.String(),
			"csdb.Column{\n    Field:   sql.NullString{String:\"category_id\", Valid:true},\n    Type:    sql.NullString{String:\"int(10) unsigned\", Valid:true},\n    Null:    sql.NullString{String:\"NO\", Valid:true},\n    Key:     sql.NullString{String:\"PRI\", Valid:true},\n    Default: sql.NullString{String:\"0\", Valid:true},\n    Extra:   sql.NullString{String:\"\", Valid:true},\n},\ncsdb.Column{\n    Field:   sql.NullString{String:\"path\", Valid:true},\n    Type:    sql.NullString{String:\"varchar(255)\", Valid:true},\n    Null:    sql.NullString{String:\"YES\", Valid:true},\n    Key:     sql.NullString{},\n    Default: sql.NullString{},\n    Extra:   sql.NullString{String:\"\", Valid:true},\n}",
		},
		{
			mustStructure(table4).Columns.UniqueKeys().Len(), 1,
			mustStructure(table4).Columns.String(),
			"csdb.Column{\n    Field:   sql.NullString{String:\"user_id\", Valid:true},\n    Type:    sql.NullString{String:\"int(10) unsigned\", Valid:true},\n    Null:    sql.NullString{String:\"NO\", Valid:true},\n    Key:     sql.NullString{String:\"PRI\", Valid:true},\n    Default: sql.NullString{},\n    Extra:   sql.NullString{String:\"auto_increment\", Valid:true},\n},\ncsdb.Column{\n    Field:   sql.NullString{String:\"email\", Valid:true},\n    Type:    sql.NullString{String:\"varchar(128)\", Valid:true},\n    Null:    sql.NullString{String:\"YES\", Valid:true},\n    Key:     sql.NullString{String:\"\", Valid:true},\n    Default: sql.NullString{},\n    Extra:   sql.NullString{String:\"\", Valid:true},\n},\ncsdb.Column{\n    Field:   sql.NullString{String:\"username\", Valid:true},\n    Type:    sql.NullString{String:\"varchar(40)\", Valid:true},\n    Null:    sql.NullString{String:\"YES\", Valid:true},\n    Key:     sql.NullString{String:\"UNI\", Valid:true},\n    Default: sql.NullString{},\n    Extra:   sql.NullString{String:\"\", Valid:true},\n}",
		},
		{mustStructure(table4).Columns.PrimaryKeys().Len(), 1, "", ""},
	}

	for i, test := range tests {
		assert.Equal(t, test.have, test.want, "Incorrect length at index %d", i)
		assert.Equal(t, test.haveS, test.wantS)
	}

	tsN := mustStructure(table4).Columns.ByName("user_id_not_found")
	assert.NotNil(t, tsN)
	assert.False(t, tsN.Field.Valid)

	ts4 := mustStructure(table4).Columns.ByName("user_id")
	assert.True(t, ts4.Field.Valid)
	assert.True(t, ts4.IsAutoIncrement())

	ts4b := mustStructure(table4).Columns.ByName("email")
	assert.True(t, ts4b.Field.Valid)
	assert.True(t, ts4b.IsNull())

	assert.True(t, mustStructure(table4).Columns.First().IsPK())
	emptyTS := &csdb.Table{}
	assert.False(t, emptyTS.Columns.First().IsPK())

}
