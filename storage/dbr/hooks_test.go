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

package dbr_test

import (
	"testing"

	"github.com/corestoreio/csfw/storage/dbr"
	"github.com/stretchr/testify/assert"
)

func TestNewHook(t *testing.T) {

	hSel := func(*dbr.Select) {}
	hIns := func(*dbr.Insert) {}
	hUpd := func(*dbr.Update) {}
	hDel := func(*dbr.Delete) {}

	h := dbr.NewHook()
	h.AddSelectAfter(hSel, hSel)
	h.AddInsertAfter(hIns, hIns)
	h.AddUpdateAfter(hUpd, hUpd)
	h.AddDeleteAfter(hDel, hDel)

	assert.Len(t, h.BeforeToSQL.SelectHooks, 2)
	assert.Len(t, h.BeforeToSQL.InsertHooks, 2)
	assert.Len(t, h.BeforeToSQL.UpdateHooks, 2)
	assert.Len(t, h.BeforeToSQL.DeleteHooks, 2)

	h2 := dbr.NewHook()
	h2.AddSelectAfter(hSel, hSel)
	h2.AddInsertAfter(hIns, hIns)
	h2.AddUpdateAfter(hUpd, hUpd)
	h2.AddDeleteAfter(hDel, hDel)

	assert.Len(t, h2.BeforeToSQL.SelectHooks, 2)
	assert.Len(t, h2.BeforeToSQL.InsertHooks, 2)
	assert.Len(t, h2.BeforeToSQL.UpdateHooks, 2)
	assert.Len(t, h2.BeforeToSQL.DeleteHooks, 2)

	h.Merge(h2, h2)
	assert.Len(t, h.BeforeToSQL.SelectHooks, 6)
	assert.Len(t, h.BeforeToSQL.InsertHooks, 6)
	assert.Len(t, h.BeforeToSQL.UpdateHooks, 6)
	assert.Len(t, h.BeforeToSQL.DeleteHooks, 6)

	h.BeforeToSQL.SelectHooks.Apply(nil)
	h.BeforeToSQL.InsertHooks.Apply(nil)
	h.BeforeToSQL.UpdateHooks.Apply(nil)
	h.BeforeToSQL.DeleteHooks.Apply(nil)
}
