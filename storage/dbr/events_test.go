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
	"fmt"
	"testing"

	"github.com/corestoreio/csfw/util/errors"
	"github.com/stretchr/testify/assert"
)

var _ fmt.Stringer = (*SelectListeners)(nil)
var _ fmt.Stringer = (*InsertListeners)(nil)
var _ fmt.Stringer = (*UpdateListeners)(nil)
var _ fmt.Stringer = (*DeleteListeners)(nil)

func TestNewListenerBucket(t *testing.T) {

	t.Run("Merge Many", func(t *testing.T) {
		lbOld := MustNewListenerBucket(
			Listen{
				Name:       "Select",
				EventType:  OnBeforeToSQL,
				SelectFunc: func(b *Select) {},
			},
			Listen{
				Name:       "Insert",
				EventType:  OnBeforeToSQL,
				InsertFunc: func(b *Insert) {},
			},
			Listen{
				Name:       "Update",
				EventType:  OnBeforeToSQL,
				UpdateFunc: func(b *Update) {},
			},
			Listen{
				Name:       "Delete",
				EventType:  OnBeforeToSQL,
				DeleteFunc: func(b *Delete) {},
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
				Name:       "Logger",
				EventType:  OnBeforeToSQL,
				SelectFunc: func(b *Select) {},
				InsertFunc: func(b *Insert) {},
				UpdateFunc: func(b *Update) {},
				DeleteFunc: func(b *Delete) {},
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
				assert.True(t, errors.IsEmpty(r.(error)), "%+v", r.(error))
			} else {
				t.Error("Expecting a panic")
			}
		}()
		_ = MustNewListenerBucket(Listen{
			SelectFunc: func(*Select) {},
		})
	})
	t.Run("Error Select", func(t *testing.T) {
		lb, err := NewListenerBucket(Listen{
			SelectFunc: func(*Select) {},
		})
		assert.Nil(t, lb)
		assert.True(t, errors.IsEmpty(err), "%+v", err)
	})

	t.Run("Error Insert", func(t *testing.T) {
		lb, err := NewListenerBucket(Listen{
			InsertFunc: func(*Insert) {},
		})
		assert.Nil(t, lb)
		assert.True(t, errors.IsEmpty(err), "%+v", err)
	})

	t.Run("Error Update", func(t *testing.T) {
		lb, err := NewListenerBucket(Listen{
			UpdateFunc: func(*Update) {},
		})
		assert.Nil(t, lb)
		assert.True(t, errors.IsEmpty(err), "%+v", err)
	})

	t.Run("Error Delete", func(t *testing.T) {
		lb, err := NewListenerBucket(Listen{
			DeleteFunc: func(*Delete) {},
		})
		assert.Nil(t, lb)
		assert.True(t, errors.IsEmpty(err), "%+v", err)
	})

	t.Run("Select Only", func(t *testing.T) {
		called := 0
		lb := MustNewListenerBucket(Listen{
			Name:      "Select",
			EventType: OnBeforeToSQL,
			SelectFunc: func(b *Select) {
				assert.Nil(t, b)
				called++
			},
		})
		err := lb.Select.dispatch(OnBeforeToSQL, nil)
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
			InsertFunc: func(b *Insert) {
				assert.Nil(t, b)
				called++
			},
		})
		err := lb.Insert.dispatch(OnBeforeToSQL, nil)
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
			UpdateFunc: func(b *Update) {
				assert.Nil(t, b)
				called++
			},
		})
		assert.NoError(t, err)
		err = lb.Update.dispatch(OnBeforeToSQL, nil)
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
			DeleteFunc: func(b *Delete) {
				assert.Nil(t, b)
				called++
			},
		})
		assert.NoError(t, err)
		err = lb.Delete.dispatch(OnBeforeToSQL, nil)
		assert.NoError(t, err)
		assert.Exactly(t, 1, called)

		assert.Len(t, lb.Select, 0)
		assert.Len(t, lb.Update, 0)
		assert.Len(t, lb.Insert, 0)
	})

	t.Run("Select Merge", func(t *testing.T) {
		var l1 SelectListeners
		l1.Add(
			Listen{
				Name:       "col1",
				SelectFunc: func(b *Select) {},
			},
		)
		var l2 SelectListeners
		l2.Add(
			Listen{
				Name:       "col2",
				SelectFunc: func(b *Select) {},
			},
		)
		assert.Exactly(t, `col1; col2`, l1.Merge(l2).String())
	})
	t.Run("Insert Merge", func(t *testing.T) {
		var l1 InsertListeners
		l1.Add(
			Listen{
				Name:       "col1",
				InsertFunc: func(b *Insert) {},
			},
		)
		var l2 InsertListeners
		l2.Add(
			Listen{
				Name:       "col2",
				InsertFunc: func(b *Insert) {},
			},
		)
		assert.Exactly(t, `col1; col2`, l1.Merge(l2).String())
	})
	t.Run("Update Merge", func(t *testing.T) {
		var l1 UpdateListeners
		l1.Add(
			Listen{
				Name:       "col1",
				UpdateFunc: func(b *Update) {},
			},
		)
		var l2 UpdateListeners
		l2.Add(
			Listen{
				Name:       "col2",
				UpdateFunc: func(b *Update) {},
			},
		)
		assert.Exactly(t, `col1; col2`, l1.Merge(l2).String())
	})
	t.Run("Delete Merge", func(t *testing.T) {
		var l1 DeleteListeners
		l1.Add(
			Listen{
				Name:       "col1",
				DeleteFunc: func(b *Delete) {},
			},
		)
		var l2 DeleteListeners
		l2.Add(
			Listen{
				Name:       "col2",
				DeleteFunc: func(b *Delete) {},
			},
		)
		assert.Exactly(t, `col1; col2`, l1.Merge(l2).String())
	})

}
