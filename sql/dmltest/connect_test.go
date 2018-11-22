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

package dmltest_test

import (
	"database/sql"
	"testing"

	"github.com/corestoreio/errors"
	"github.com/corestoreio/pkg/sql/dmltest"
	"github.com/corestoreio/pkg/util/assert"
)

func TestMockDB(t *testing.T) {
	dbc, mockDB := dmltest.MockDB(t)
	defer dmltest.MockClose(t, dbc, mockDB)
	assert.NotNil(t, dbc)
	assert.NotNil(t, mockDB)
}

func TestMustConnectDB(t *testing.T) {
	t.Parallel()

	defer func() {
		if r := recover(); r != nil {
			assert.NotNil(t, r, "There should be no panic")
		}
	}()

	db := dmltest.MustConnectDB(t)
	assert.NotNil(t, db)
}

type tErrorLID struct {
	msg string
	err errors.Kind
	t   *testing.T
}

func (t tErrorLID) Errorf(format string, args ...interface{}) {
	for _, a := range args {
		switch at := a.(type) {
		case string:
			assert.Exactly(t.t, t.msg, at)
		case error:
			assert.True(t.t, t.err.Match(at), "%+v", at)
		case nil:
		default:
			t.t.Fatalf("Type %#v not supported", a)
		}
	}
}

type sqlResult struct {
	err error
	id  int64
}

func (sr sqlResult) LastInsertId() (int64, error) { return sr.id, sr.err }
func (sr sqlResult) RowsAffected() (int64, error) { return sr.id, sr.err }

func newLID(id int64, err error) (sql.Result, error) {
	return sqlResult{id: id}, err
}

func newLIDErr(id int64, err error) (sql.Result, error) {
	return sqlResult{id: id, err: err}, nil
}

func TestCheckLastInsertID(t *testing.T) {

	t.Run("happy path", func(t *testing.T) {
		lid := dmltest.CheckLastInsertID(tErrorLID{t: t})(newLID(123, nil))
		assert.Exactly(t, int64(123), lid)
	})
	t.Run("error in newLID", func(t *testing.T) {
		lid := dmltest.CheckLastInsertID(tErrorLID{t: t, err: errors.NotAllowed})(newLID(123, errors.NotAllowed.Newf("upss")))
		assert.Exactly(t, int64(0), lid)
	})
	t.Run("error in newLID msg", func(t *testing.T) {
		lid := dmltest.CheckLastInsertID(tErrorLID{t: t, err: errors.NotAllowed, msg: "Hello"}, "Hello")(newLID(123, errors.NotAllowed.Newf("upss")))
		assert.Exactly(t, int64(0), lid)
	})

	t.Run("error in newLIDErr", func(t *testing.T) {
		lid := dmltest.CheckLastInsertID(tErrorLID{t: t, err: errors.NotAllowed})(newLIDErr(123, errors.NotAllowed.Newf("upss")))
		assert.Exactly(t, int64(0), lid)
	})
	t.Run("error in newLIDErr msg", func(t *testing.T) {
		lid := dmltest.CheckLastInsertID(tErrorLID{t: t, err: errors.NotAllowed, msg: "Hello"}, "Hello")(newLIDErr(123, errors.NotAllowed.Newf("upss")))
		assert.Exactly(t, int64(0), lid)
	})
}
