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
	"testing"

	"github.com/corestoreio/csfw/storage/dbr"
	"github.com/corestoreio/errors"
)

var _ dbr.Binder = (*someRecord)(nil)

type someRecord struct {
	SomethingID int
	UserID      int64
	Other       bool
}

func (sr someRecord) appendBind(args dbr.Arguments, column string) (_ dbr.Arguments, err error) {
	switch column {
	case "something_id":
		args = args.Int(sr.SomethingID)
	case "user_id":
		args = args.Int64(sr.UserID)
	case "other":
		args = args.Bool(sr.Other)
	default:
		err = errors.NewNotFoundf("[dbr_test] Column %q not found", column)
	}
	return args, err
}

func (sr someRecord) AppendBind(args dbr.Arguments, columns []string) (_ dbr.Arguments, err error) {
	l := len(columns)
	if l == 1 {
		return sr.appendBind(args, columns[0])
	}
	if l == 0 {
		return args.Int(sr.SomethingID).Int64(sr.UserID).Bool(sr.Other), nil // except auto inc column ;-)
	}
	for _, col := range columns {
		if args, err = sr.appendBind(args, col); err != nil {
			return nil, errors.WithStack(err)
		}
	}
	return args, err
}

func TestInsert_Bind(t *testing.T) {
	t.Parallel()
	objs := []someRecord{{1, 88, false}, {2, 99, true}, {3, 101, true}}
	wantArgs := []interface{}{int64(1), int64(88), false, int64(2), int64(99), true, int64(3), int64(101), true, int64(99)}

	t.Run("valid with multiple records", func(t *testing.T) {
		compareToSQL(t,
			dbr.NewInsert("a").
				AddColumns("something_id", "user_id", "other").
				Bind(objs[0]).Bind(objs[1], objs[2]).
				AddOnDuplicateKey(
					dbr.Column("something_id").Int64(99),
					dbr.Column("user_id").Values(),
				),
			nil,
			"INSERT INTO `a` (`something_id`,`user_id`,`other`) VALUES (?,?,?),(?,?,?),(?,?,?) ON DUPLICATE KEY UPDATE `something_id`=?, `user_id`=VALUES(`user_id`)",
			"INSERT INTO `a` (`something_id`,`user_id`,`other`) VALUES (1,88,0),(2,99,1),(3,101,1) ON DUPLICATE KEY UPDATE `something_id`=99, `user_id`=VALUES(`user_id`)",
			wantArgs...,
		)
	})
	t.Run("without columns, all columns requested", func(t *testing.T) {
		compareToSQL(t,
			dbr.NewInsert("a").
				SetRecordValueCount(3).
				Bind(objs[0]).Bind(objs[1], objs[2]).
				AddOnDuplicateKey(
					dbr.Column("something_id").Int64(99),
					dbr.Column("user_id").Values(),
				),
			nil,
			"INSERT INTO `a` VALUES (?,?,?),(?,?,?),(?,?,?) ON DUPLICATE KEY UPDATE `something_id`=?, `user_id`=VALUES(`user_id`)",
			"INSERT INTO `a` VALUES (1,88,0),(2,99,1),(3,101,1) ON DUPLICATE KEY UPDATE `something_id`=99, `user_id`=VALUES(`user_id`)",
			wantArgs...,
		)
	})
	t.Run("column not found", func(t *testing.T) {
		objs := []someRecord{{1, 88, false}, {2, 99, true}}
		compareToSQL(t,
			dbr.NewInsert("a").AddColumns("something_it", "user_id", "other").Bind(objs[0]).Bind(objs[1]),
			errors.IsNotFound,
			"",
			"",
		)
	})
}
