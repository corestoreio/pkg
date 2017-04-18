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
	"fmt"

	"github.com/corestoreio/csfw/storage/dbr"
	"github.com/corestoreio/errors"
)

// Make sure that type exampleInsertRecord implements interface
// dbr.ArgumentGenerater.
var _ dbr.ArgumentGenerater = (*exampleInsertRecord)(nil)

// exampleInsertRecord represents just a demo record but can be Products,
// Categories, Sales Orders, etc ...
type exampleInsertRecord struct {
	SomethingID int
	UserID      int64
	Other       bool
}

func (sr exampleInsertRecord) GenerateArguments(statementType byte, columns, condition []string) (dbr.Arguments, error) {
	args := make(dbr.Arguments, 0, 3) // 3 == number of fields in the struct
	// statementType lets you know in which circumstances your function gets
	// called so can return the most suitable arguments. You need to write this
	// boiler plate code only once or let it generate.
	if statementType == dbr.StatementTypeInsert {
		for _, c := range columns {
			switch c {
			case "something_id":
				args = append(args, dbr.ArgInt(sr.SomethingID))
			case "user_id":
				args = append(args, dbr.ArgInt64(sr.UserID))
			case "other":
				args = append(args, dbr.ArgBool(sr.Other))
			default:
				return nil, errors.NewNotFoundf("[dbr_test] Column %q not found", c)
			}
		}
	}
	return args, nil
}

func ExampleInsert_AddRecords() {

	objs := []exampleInsertRecord{{1, 88, false}, {2, 99, true}}

	sqlStr, args, err := dbr.NewInsert("a").AddColumns("something_id", "user_id", "other").
		AddRecords(objs[0]).AddRecords(objs[1]).
		ToSQL()
	if err != nil {
		fmt.Printf("%+v\n", err)
		return
	}

	sqlPre, err := dbr.Preprocess(sqlStr, args...)
	if err != nil {
		fmt.Printf("%+v\n", err)
		return
	}

	fmt.Printf("%s\nArguments: %v\nProcessed: %s\n", sqlStr, args.Interfaces(), sqlPre)
	// Output:
	// INSERT INTO `a` (`something_id`,`user_id`,`other`) VALUES (?,?,?),(?,?,?)
	// Arguments: [1 88 false 2 99 true]
	// Processed: INSERT INTO `a` (`something_id`,`user_id`,`other`) VALUES (1,88,0),(2,99,1)
}
