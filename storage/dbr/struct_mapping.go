package dbr

import (
	"reflect"

	"github.com/corestoreio/csfw/util"
	"github.com/corestoreio/errors"
)

var destDummy interface{}

type fieldMapQueueElement struct {
	Type reflect.Type
	Idxs []int
}

// recordType is the type of a structure
func calculateFieldMap(recordType reflect.Type, columns []string, requireAllColumns bool) ([][]int, error) {

	// each value is either the slice to get to the field via FieldByIndex(index
	// []int) in the record, or nil if we don't want to map it to the structure.
	lenColumns := len(columns)
	fieldMap := make([][]int, lenColumns)

	for i, col := range columns {
		fieldMap[i] = nil

		queue := []fieldMapQueueElement{{Type: recordType, Idxs: nil}}

	QueueLoop:
		for len(queue) > 0 {
			curEntry := queue[0]
			queue = queue[1:]

			curType := curEntry.Type
			curIdxs := curEntry.Idxs
			lenFields := curType.NumField()

			for j := 0; j < lenFields; j++ {
				fieldStruct := curType.Field(j)

				// Skip unexported field
				if len(fieldStruct.PkgPath) != 0 {
					continue
				}

				name := fieldStruct.Tag.Get("db")
				if name != "-" {
					if name == "" {
						name = util.CamelCaseToUnderscore(fieldStruct.Name)
					}
					if name == col {
						fieldMap[i] = append(curIdxs, j)
						break QueueLoop
					}
				}

				if fieldStruct.Type.Kind() == reflect.Struct {
					var idxs2 []int
					copy(idxs2, curIdxs)
					idxs2 = append(idxs2, j)
					queue = append(queue, fieldMapQueueElement{Type: fieldStruct.Type, Idxs: idxs2})
				}
			}
		}

		if requireAllColumns && fieldMap[i] == nil {
			return nil, errors.NewNotFoundf("[dbr] calculateFieldMap: couldn't find match for column %q", col)
		}
	}

	return fieldMap, nil
}

func prepareHolderFor(record reflect.Value, fieldMap [][]int, holder []interface{}) ([]interface{}, error) {
	// Given a query and given a structure (field list), there'ab 2 sets of fields.
	// Take the intersection. We can fill those in. great.
	// For fields in the structure that aren't in the query, we'll let that slide if db:"-"
	// For fields in the structure that aren't in the query but without db:"-", return error
	// For fields in the query that aren't in the structure, we'll ignore them.

	for i, fieldIndex := range fieldMap {
		if fieldIndex == nil {
			holder[i] = &destDummy
		} else {
			field := record.FieldByIndex(fieldIndex)
			holder[i] = field.Addr().Interface()
		}
	}

	return holder, nil
}

//func valuesFor(recordType reflect.Type, record reflect.Value, columns []string) ([]interface{}, error) {
//	fieldMap, err := calculateFieldMap(recordType, columns, true)
//	if err != nil {
//		return nil, errors.Wrap(err, "[dbr] valuesFor.calculateFieldMap")
//	}
//
//	values := make([]interface{}, len(columns))
//	for i, fieldIndex := range fieldMap {
//		if fieldIndex == nil {
//			return nil, errors.NewEmptyf("[dbr] fieldIndex is nil: %#v", fieldMap)
//		}
//		field := record.FieldByIndex(fieldIndex)
//		values[i] = field.Interface()
//	}
//
//	return values, nil
//}
