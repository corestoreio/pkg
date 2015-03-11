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
	"fmt"
	"strings"

	"github.com/juju/errgo"
)

type (
	column struct {
		Field, Type, Null, Key, Default, Extra sql.NullString
		GoType, GoName                         string
	}
)

func (c column) Comment() string {
	sqlNull := "NOT NULL"
	if c.Null.String == "YES" {
		sqlNull = "NULL"
	}
	sqlDefault := ""
	if c.Default.String != "" {
		sqlDefault = "DEFAULT '" + c.Default.String + "'"
	}
	return fmt.Sprintf(
		"// %s %s %s %s %s %s",
		c.Field.String,
		c.Type.String,
		sqlNull,
		c.Key.String,
		sqlDefault,
		c.Extra.String,
	)
}

func GetColumns(db *sql.DB, table string) ([]column, error) {
	var columns = make([]column, 0, 200)
	rows, err := db.Query("SHOW COLUMNS FROM `" + table + "`")
	if err != nil {
		return nil, errgo.Mask(err)
	}
	defer rows.Close()

	// Drop unused column entity_type_id in customer__* and catalog_* tables
	isEntityTypeIdFree := strings.Index(table, "catalog_") >= 0 || strings.Index(table, "customer_") >= 0

	for rows.Next() {
		col := column{}
		err := rows.Scan(&col.Field, &col.Type, &col.Null, &col.Key, &col.Default, &col.Extra)
		if err != nil {
			return nil, errgo.Mask(err)
		}

		if isEntityTypeIdFree && col.Field.String == "entity_type_id" {
			continue
		}

		updateGoType(&col)
		columns = append(columns, col)
	}
	err = rows.Err()
	if err != nil {
		return nil, errgo.Mask(err)
	}
	return columns, nil
}

func updateGoType(col *column) {
	// dbr relates to github.com/gocraft/dbr
	col.GoType = "undefined"
	col.GoName = Camelize(col.Field.String)
	if strings.Index(col.Field.String, "is_") == 0 {
		col.GoType = "bool"
		if col.Null.String == "YES" {
			col.GoType = "dbr.NullBool"
		}
	} else if strings.Contains(col.Type.String, "int") {
		col.GoType = "int64"
		if col.Null.String == "YES" {
			col.GoType = "dbr.NullInt64"
		}
	} else if strings.Contains(col.Type.String, "varchar") || strings.Contains(col.Type.String, "text") {
		col.GoType = "string"
		if col.Null.String == "YES" {
			col.GoType = "dbr.NullString"
		}
	} else if strings.Contains(col.Type.String, "decimal") || strings.Contains(col.Type.String, "float") {
		col.GoType = "float64"
		if col.Null.String == "YES" {
			col.GoType = "dbr.NullFloat64"
		}
	} else if strings.Contains(col.Type.String, "timestamp") || strings.Contains(col.Type.String, "date") {
		col.GoType = "time.Time"
		if col.Null.String == "YES" {
			col.GoType = "dbr.NullTime"
		}
	}
}
