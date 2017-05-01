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

package dbr_test

import (
	"bytes"
	"context"
	"encoding/json"
	"testing"

	"github.com/corestoreio/csfw/storage/dbr"
	"github.com/corestoreio/csfw/util/cstesting"
	"github.com/corestoreio/errors"
	"github.com/stretchr/testify/assert"
)

// TableCoreConfigDataSlice represents a collection type for DB table core_config_data
// Generated via tableToStruct.
type TableCoreConfigDataSlice struct {
	columns []string
	Data    []*TableCoreConfigData
}

// TableCoreConfigData represents a type for DB table core_config_data
// Generated via tableToStruct.
type TableCoreConfigData struct {
	ConfigID int64          `json:",omitempty"` // config_id int(10) unsigned NOT NULL PRI  auto_increment
	Scope    string         `json:",omitempty"` // scope varchar(8) NOT NULL MUL DEFAULT 'default'
	ScopeID  int64          `json:",omitempty"` // scope_id int(11) NOT NULL  DEFAULT '0'
	Path     string         `json:",omitempty"` // path varchar(255) NOT NULL  DEFAULT 'general'
	Value    dbr.NullString `json:",omitempty"` // value text NULL
}

func (s *TableCoreConfigDataSlice) ScanArgs(columns []string) []interface{} {
	s.Data = make([]*TableCoreConfigData, 0, 10)
	s.columns = columns
	var c TableCoreConfigData
	return []interface{}{&c.ConfigID, &c.Scope, &c.ScopeID, &c.Path, &c.Value}
}

func (s *TableCoreConfigDataSlice) Row(idx int64, values []interface{}) error {

	// write benchmark to see which is faster. this switch stuff or the if/else
	// blocks.
	//ccd := &TableCoreConfigData{}
	//for i, c := range s.columns {
	//	switch v := values[i].(type) {
	//	case *int64:
	//		switch c {
	//		case "config_id":
	//			ccd.ConfigID = *v
	//		case "scope_id":
	//			ccd.ScopeID = *v
	//		}
	//	case *string:
	//	}
	//}

	c := &TableCoreConfigData{}
	if v, ok := values[0].(*int64); ok {
		c.ConfigID = *v
	} else {
		return errors.NewNotValidf("[pkg] Field ConfigID. Failed to assert to *int64. Got: %#v", values[0])
	}
	if v, ok := values[1].(*string); ok {
		c.Scope = *v
	} else {
		return errors.NewNotValidf("[pkg] Field Scope. Failed to assert to *string. Got: %#v", values[1])
	}
	if v, ok := values[2].(*int64); ok {
		c.ScopeID = *v
	} else {
		return errors.NewNotValidf("[pkg] Field ConfigID. Failed to assert to *int64. Got: %#v", values[2])
	}
	if v, ok := values[3].(*string); ok {
		c.Path = *v
	} else {
		return errors.NewNotValidf("[pkg] Field Path. Failed to assert to *string. Got: %#v", values[3])
	}
	if v, ok := values[4].(*dbr.NullString); ok {
		c.Value = *v
	} else {
		return errors.NewNotValidf("[pkg] Field Path. Failed to assert to *dbr.NullString. Got: %#v", values[4])
	}

	s.Data = append(s.Data, c)

	return nil
}

func TestSelect_LoadX(t *testing.T) {

	dbc, dbMock := cstesting.MockDB(t)
	defer func() {
		dbMock.ExpectClose()
		assert.NoError(t, dbc.Close())
		if err := dbMock.ExpectationsWereMet(); err != nil {
			t.Error("there were unfulfilled expections", err)
		}
	}()

	dbMock.ExpectQuery("SELECT").WillReturnRows(cstesting.MustMockRows(cstesting.WithFile("testdata/core_config_data.csv")))
	s := dbr.NewSelect("*").From("core_config_data")
	s.DB.Querier = dbc.DB

	ccd := &TableCoreConfigDataSlice{}

	_, err := s.LoadX(context.TODO(), ccd)
	assert.NoError(t, err, "%+v", err)

	buf := new(bytes.Buffer)
	je := json.NewEncoder(buf)

	for _, c := range ccd.Data {
		if err := je.Encode(c); err != nil {
			t.Fatalf("%+v", err)
		}
	}
	assert.Equal(t, "{\"ConfigID\":2,\"Scope\":\"default\",\"Path\":\"web/unsecure/base_url\",\"Value\":\"http://mgeto2.local/\"}\n{\"ConfigID\":3,\"Scope\":\"website\",\"ScopeID\":11,\"Path\":\"general/locale/code\",\"Value\":\"en_US\"}\n{\"ConfigID\":4,\"Scope\":\"default\",\"Path\":\"general/locale/timezone\",\"Value\":\"Europe/Berlin\"}\n{\"ConfigID\":5,\"Scope\":\"default\",\"Path\":\"currency/options/base\",\"Value\":\"EUR\"}\n{\"ConfigID\":15,\"Scope\":\"store\",\"ScopeID\":33,\"Path\":\"design/head/includes\",\"Value\":\"\\u003clink  rel=\\\"stylesheet\\\" type=\\\"text/css\\\" href=\\\"{{MEDIA_URL}}styles.css\\\" /\\u003e\"}\n{\"ConfigID\":16,\"Scope\":\"default\",\"Path\":\"admin/security/use_case_sensitive_login\",\"Value\":null}\n{\"ConfigID\":17,\"Scope\":\"default\",\"Path\":\"admin/security/session_lifetime\",\"Value\":\"90000\"}\n",
		buf.String())
}
