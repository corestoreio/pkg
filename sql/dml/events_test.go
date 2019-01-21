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

package dml

import (
	"fmt"
	"testing"

	"github.com/corestoreio/errors"
	"github.com/corestoreio/pkg/util/assert"
)

var _ fmt.Stringer = (*ListenersSelect)(nil)
var _ fmt.Stringer = (*ListenersInsert)(nil)
var _ fmt.Stringer = (*ListenersUpdate)(nil)
var _ fmt.Stringer = (*ListenersDelete)(nil)

func TestNewListenerBucket(t *testing.T) {

	t.Run("Merge Many", func(t *testing.T) {
		lbOld := MustNewListenerBucket(
			Listen{
				Name:           "Select",
				EventType:      OnBeforeToSQL,
				ListenSelectFn: func(b *Select) {},
			},
			Listen{
				Name:           "Insert",
				EventType:      OnBeforeToSQL,
				ListenInsertFn: func(b *Insert) {},
			},
			Listen{
				Name:           "Update",
				EventType:      OnBeforeToSQL,
				ListenUpdateFn: func(b *Update) {},
			},
			Listen{
				Name:           "Delete",
				EventType:      OnBeforeToSQL,
				ListenDeleteFn: func(b *Delete) {},
			},
		)

		lbNew := MustNewListenerBucket().Merge(lbOld)
		assert.Len(t, lbNew.Select, 1)
		assert.Len(t, lbNew.Insert, 1)
		assert.Len(t, lbNew.Update, 1)
		assert.Len(t, lbNew.Delete, 1)

		assert.Exactly(t, `Select`, lbNew.Select.String())
		assert.Exactly(t, `Insert`, lbNew.Insert.String())
		assert.Exactly(t, `Update`, lbNew.Update.String())
		assert.Exactly(t, `Delete`, lbNew.Delete.String())
	})
	t.Run("Merge One", func(t *testing.T) {
		lbOld := MustNewListenerBucket(
			Listen{
				Name:           "Logger",
				EventType:      OnBeforeToSQL,
				ListenSelectFn: func(b *Select) {},
				ListenInsertFn: func(b *Insert) {},
				ListenUpdateFn: func(b *Update) {},
				ListenDeleteFn: func(b *Delete) {},
			},
		)

		lbNew := MustNewListenerBucket().Merge(lbOld)
		assert.Len(t, lbNew.Select, 1)
		assert.Len(t, lbNew.Insert, 1)
		assert.Len(t, lbNew.Update, 1)
		assert.Len(t, lbNew.Delete, 1)

		assert.Exactly(t, `Logger`, lbNew.Select.String())
		assert.Exactly(t, `Logger`, lbNew.Insert.String())
		assert.Exactly(t, `Logger`, lbNew.Update.String())
		assert.Exactly(t, `Logger`, lbNew.Delete.String())
	})
	t.Run("panic", func(t *testing.T) {
		defer func() {
			if r := recover(); r != nil {
				assert.ErrorIsKind(t, errors.Empty, r.(error))
			} else {
				t.Error("Expecting a panic")
			}
		}()
		_ = MustNewListenerBucket(Listen{
			ListenSelectFn: func(*Select) {},
		})
	})
	t.Run("Error Select", func(t *testing.T) {
		lb, err := NewListenerBucket(Listen{
			ListenSelectFn: func(*Select) {},
		})
		assert.Nil(t, lb)
		assert.ErrorIsKind(t, errors.Empty, err)
	})
	t.Run("Error Insert", func(t *testing.T) {
		lb, err := NewListenerBucket(Listen{
			ListenInsertFn: func(*Insert) {},
		})
		assert.Nil(t, lb)
		assert.ErrorIsKind(t, errors.Empty, err)
	})
	t.Run("Error Update", func(t *testing.T) {
		lb, err := NewListenerBucket(Listen{
			ListenUpdateFn: func(*Update) {},
		})
		assert.Nil(t, lb)
		assert.ErrorIsKind(t, errors.Empty, err)
	})
	t.Run("Error Delete", func(t *testing.T) {
		lb, err := NewListenerBucket(Listen{
			ListenDeleteFn: func(*Delete) {},
		})
		assert.Nil(t, lb)
		assert.ErrorIsKind(t, errors.Empty, err)
	})
	t.Run("Select Only", func(t *testing.T) {
		called := 0
		lb := MustNewListenerBucket(Listen{
			Name:      "Select",
			EventType: OnBeforeToSQL,
			ListenSelectFn: func(b *Select) {
				called++
			},
		})
		err := lb.Select.dispatch(OnBeforeToSQL, &Select{})
		assert.NoError(t, err)
		assert.Exactly(t, 1, called)

		assert.Len(t, lb.Insert, 0)
		assert.Len(t, lb.Update, 0)
		assert.Len(t, lb.Delete, 0)
	})
	t.Run("Insert Only", func(t *testing.T) {
		called := 0
		lb := MustNewListenerBucket(Listen{
			Name:      "Insert",
			EventType: OnBeforeToSQL,
			ListenInsertFn: func(b *Insert) {
				called++
			},
		})
		err := lb.Insert.dispatch(OnBeforeToSQL, &Insert{})
		assert.NoError(t, err)
		assert.Exactly(t, 1, called)

		assert.Len(t, lb.Select, 0)
		assert.Len(t, lb.Update, 0)
		assert.Len(t, lb.Delete, 0)
	})
	t.Run("Update Only", func(t *testing.T) {
		called := 0
		lb, err := NewListenerBucket(Listen{
			Name:      "Update",
			EventType: OnBeforeToSQL,
			ListenUpdateFn: func(b *Update) {
				called++
			},
		})
		assert.NoError(t, err)
		err = lb.Update.dispatch(OnBeforeToSQL, &Update{})
		assert.NoError(t, err)
		assert.Exactly(t, 1, called)

		assert.Len(t, lb.Select, 0)
		assert.Len(t, lb.Insert, 0)
		assert.Len(t, lb.Delete, 0)
	})
	t.Run("Delete Only", func(t *testing.T) {
		called := 0
		lb, err := NewListenerBucket(Listen{
			Name:      "Delete",
			EventType: OnBeforeToSQL,
			ListenDeleteFn: func(b *Delete) {
				called++
			},
		})
		assert.NoError(t, err)
		err = lb.Delete.dispatch(OnBeforeToSQL, &Delete{})
		assert.NoError(t, err)
		assert.Exactly(t, 1, called)

		assert.Len(t, lb.Select, 0)
		assert.Len(t, lb.Update, 0)
		assert.Len(t, lb.Insert, 0)
	})

	t.Run("Select Merge", func(t *testing.T) {
		var l1 ListenersSelect
		l1.Add(
			Listen{
				Name:           "col1",
				ListenSelectFn: func(b *Select) {},
			},
		)
		var l2 ListenersSelect
		l2.Add(
			Listen{
				Name:           "col2",
				ListenSelectFn: func(b *Select) {},
			},
		)
		assert.Exactly(t, `col1; col2`, l1.Merge(l2).String())
	})
	t.Run("Insert Merge", func(t *testing.T) {
		var l1 ListenersInsert
		l1.Add(
			Listen{
				Name:           "col1",
				ListenInsertFn: func(b *Insert) {},
			},
		)
		var l2 ListenersInsert
		l2.Add(
			Listen{
				Name:           "col2",
				ListenInsertFn: func(b *Insert) {},
			},
		)
		assert.Exactly(t, `col1; col2`, l1.Merge(l2).String())
	})
	t.Run("Update Merge", func(t *testing.T) {
		var l1 ListenersUpdate
		l1.Add(
			Listen{
				Name:           "col1",
				ListenUpdateFn: func(b *Update) {},
			},
		)
		var l2 ListenersUpdate
		l2.Add(
			Listen{
				Name:           "col2",
				ListenUpdateFn: func(b *Update) {},
			},
		)
		assert.Exactly(t, `col1; col2`, l1.Merge(l2).String())
	})
	t.Run("Delete Merge", func(t *testing.T) {
		var l1 ListenersDelete
		l1.Add(
			Listen{
				Name:           "col1",
				ListenDeleteFn: func(b *Delete) {},
			},
		)
		var l2 ListenersDelete
		l2.Add(
			Listen{
				Name:           "col2",
				ListenDeleteFn: func(b *Delete) {},
			},
		)
		assert.Exactly(t, `col1; col2`, l1.Merge(l2).String())
	})

}
