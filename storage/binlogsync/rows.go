package binlogsync

import (
	"github.com/corestoreio/csfw/util/errors"
	"github.com/siddontang/go-mysql/schema"
)

// Get primary keys in one row for a table, a table may use multi fields as the PK
func GetPKValues(table *schema.Table, row []interface{}) ([]interface{}, error) {
	indexes := table.PKColumns
	if len(indexes) == 0 {
		return nil, errors.NewNotFoundf("[binlogsync] Table %q has no primary key", table)
	} else if len(table.Columns) != len(row) {
		return nil, errors.NewNotValidf("[binlogsync] Table %q has %d columns, but row data %v len is %d",
			table, len(table.Columns), row, len(row))
	}

	values := make([]interface{}, 0, len(indexes))

	for _, index := range indexes {
		keyPart := row[index]
		if keyPart == nil { // todo fix bug
			return nil, errors.NewNotValidf("[binlogsync] Row in table %q has no primary key: %v", table, row)
		}
		values = append(values, row[index])
	}

	return values, nil
}
