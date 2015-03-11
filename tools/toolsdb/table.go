// Copyright 2015 CoreStore Authors
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

package toolsdb

import (
	"database/sql"
	"log"

	"github.com/juju/errgo"
)

func GetTables(db *sql.DB, prefix string) ([]string, error) {

	var tableNames = make([]string, 0, 200)
	qry := "SHOW TABLES like '" + prefix + "%'"
	log.Println(qry)

	rows, err := db.Query(qry)
	if err != nil {
		return nil, errgo.Mask(err)
	}
	defer rows.Close()

	for rows.Next() {
		var tableName string
		err := rows.Scan(&tableName)
		if err != nil {
			return nil, errgo.Mask(err)
		}
		tableNames = append(tableNames, tableName)
	}
	err = rows.Err()
	if err != nil {
		return nil, errgo.Mask(err)
	}
	return tableNames, nil
}
