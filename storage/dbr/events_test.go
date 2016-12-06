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

package dbr

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewHook(t *testing.T) {

	hSel := func(*Select) {}
	hIns := func(*Insert) {}
	hUpd := func(*Update) {}
	hDel := func(*Delete) {}

	h := NewEvents()
	h.Select.AddBeforeToSQL(hSel, hSel)
	h.Insert.AddBeforeToSQL(hIns, hIns)
	h.Update.AddBeforeToSQL(hUpd, hUpd)
	h.Delete.AddBeforeToSQL(hDel, hDel)

	assert.Len(t, h.Select.receivers[eventToSQLBefore], 2)
	assert.Len(t, h.Insert.receivers[eventToSQLBefore], 2)
	assert.Len(t, h.Update.receivers[eventToSQLBefore], 2)
	assert.Len(t, h.Delete.receivers[eventToSQLBefore], 2)

	h2 := NewEvents()
	h2.Select.AddBeforeToSQL(hSel, hSel)
	h2.Insert.AddBeforeToSQL(hIns, hIns)
	h2.Update.AddBeforeToSQL(hUpd, hUpd)
	h2.Delete.AddBeforeToSQL(hDel, hDel)

	assert.Len(t, h.Select.receivers[eventToSQLBefore], 2)
	assert.Len(t, h.Insert.receivers[eventToSQLBefore], 2)
	assert.Len(t, h.Update.receivers[eventToSQLBefore], 2)
	assert.Len(t, h.Delete.receivers[eventToSQLBefore], 2)

	h.Merge(h2, h2)
	assert.Len(t, h.Select.receivers[eventToSQLBefore], 6)
	assert.Len(t, h.Insert.receivers[eventToSQLBefore], 6)
	assert.Len(t, h.Update.receivers[eventToSQLBefore], 6)
	assert.Len(t, h.Delete.receivers[eventToSQLBefore], 6)

	h.Select.dispatch(eventToSQLBefore, nil)
	h.Insert.dispatch(eventToSQLBefore, nil)
	h.Update.dispatch(eventToSQLBefore, nil)
	h.Delete.dispatch(eventToSQLBefore, nil)
}
