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
	"github.com/corestoreio/csfw/util/errors"
	"github.com/stretchr/testify/assert"
)

var _ csdb.TableManager = (*csdb.TableService)(nil)
var _ error = (*csdb.TableService)(nil)

func TestNewTableService(t *testing.T) {

	assert.Equal(t, csdb.MustNewTableService().Len(), csdb.Index(0))

	const (
		TableIndexStore   csdb.Index = iota // Table: store
		TableIndexGroup                     // Table: store_group
		TableIndexWebsite                   // Table: store_website
		TableIndexZZZ                       // the maximum index, which is not available.
	)

	tm1 := csdb.MustNewTableService(
		csdb.WithTable(TableIndexStore, "store"),
		csdb.WithTable(TableIndexGroup, "store_group"),
		csdb.WithTable(TableIndexWebsite, "store_website"),
	)
	assert.Equal(t, tm1.Len(), csdb.Index(3))
}
func TestNewTableServicePanic(t *testing.T) {

	defer func() {
		if r := recover(); r != nil {
			err := r.(error)
			assert.True(t, errors.MultiErrContainsAll(err, errors.IsNotValid))
		} else {
			t.Error("Expecting a panic")
		}
	}()

	_ = csdb.MustNewTableService(
		csdb.WithTable(csdb.Index(0), ""),
	)
}

func TestNewTableServiceAppend(t *testing.T) {

	tm0 := csdb.MustNewTableService()
	assert.NotNil(t, tm0)
	assert.True(t, errors.IsFatal(tm0.Append(csdb.Index(0), nil)))
	assert.Equal(t, tm0.Len(), csdb.Index(0))
}

func TestIntegration_NewTableServiceInit(t *testing.T) {

	if _, err := csdb.GetDSN(); errors.IsNotFound(err) {
		t.Skip("Skipping because no DSN found.")
	}

	dbc := csdb.MustConnectTest()
	defer dbc.Close()
	i := csdb.Index(4711)
	tm0 := csdb.MustNewTableService(csdb.WithTable(i, "admin_user"))
	assert.True(t, errors.IsTemporary(tm0.Init(dbc.NewSession(), true)))
	err := tm0.Init(dbc.NewSession())
	assert.NoError(t, err)

	table, err2 := tm0.Structure(i)
	assert.NoError(t, err2)
	assert.Equal(t, 1, table.CountPK)
	assert.Equal(t, 1, table.CountUnique)
	assert.True(t, len(table.Columns.FieldNames()) >= 15)

	assert.Nil(t, tm0.Init(dbc.NewSession()))
}
