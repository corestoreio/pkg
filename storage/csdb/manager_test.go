// Copyright 2015-2016, Cyrill @ Schumacher.fm and the CoreStore contributors
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

	"github.com/corestoreio/csfw/storage/csdb"
	"github.com/stretchr/testify/assert"
)

func TestNewTableManager(t *testing.T) {
	assert.Equal(t, csdb.NewTableManager().Len(), csdb.Index(0))

	const (
		TableIndexStore   csdb.Index = iota // Table: store
		TableIndexGroup                     // Table: store_group
		TableIndexWebsite                   // Table: store_website
		TableIndexZZZ                       // the maximum index, which is not available.
	)

	tm1 := csdb.NewTableManager(
		csdb.AddTableByName(TableIndexStore, "store"),
		csdb.AddTableByName(TableIndexGroup, "store_group"),
		csdb.AddTableByName(TableIndexWebsite, "store_website"),
	)
	assert.Equal(t, tm1.Len(), csdb.Index(3))
}
func TestNewTableManagerPanic(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			assert.Contains(t, r.(string), csdb.ErrManagerIncorrectValue.Error())
		}
	}()

	tm0 := csdb.NewTableManager(
		csdb.AddTableByName(csdb.Index(0), ""),
	)
	assert.NotNil(t, tm0)
	assert.Equal(t, tm0.Len(), csdb.Index(0))
}

func TestNewTableManagerAppend(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			assert.Contains(t, r.(string), "Table pointer cannot be nil for Index")
		}
	}()

	tm0 := csdb.NewTableManager()
	tm0.Append(csdb.Index(0), nil)
	assert.NotNil(t, tm0)
	assert.Equal(t, tm0.Len(), csdb.Index(0))
}

func TestNewTableManagerInit(t *testing.T) {
	dbc := csdb.MustConnectTest()
	defer dbc.Close()
	i := csdb.Index(4711)
	tm0 := csdb.NewTableManager(csdb.AddTableByName(i, "admin_user"))
	assert.EqualError(t, tm0.Init(dbc.NewSession(), true), csdb.ErrManagerInitReload.Error())
	err := tm0.Init(dbc.NewSession())
	assert.NoError(t, err)

	table, err2 := tm0.Structure(i)
	assert.NoError(t, err2)
	assert.Equal(t, 1, table.CountPK)
	assert.Equal(t, 1, table.CountUnique)
	assert.True(t, len(table.Columns.FieldNames()) >= 15)

	assert.Nil(t, tm0.Init(dbc.NewSession()))
}
