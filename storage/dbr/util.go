package dbr

import (
	"database/sql/driver"

	"github.com/juju/errgo"
)

// NameMapping is the routine to use when mapping column names to struct properties
var NameMapping = camelCaseToSnakeCase

func camelCaseToSnakeCase(name string) string {
	var newstr []rune
	firstTime := true

	for _, chr := range name {
		if isUpper := 'A' <= chr && chr <= 'Z'; isUpper {
			if firstTime == true {
				firstTime = false
			} else {
				newstr = append(newstr, '_')
			}
			chr -= ('A' - 'a')
		}
		newstr = append(newstr, chr)
	}

	return string(newstr)
}

// argsValuer checks if an argument implements driver.Valuer interface. If so
// uses the Value() function to get the correct value.
func argsValuer(args *[]interface{}) error {
	for i, v := range *args {
		if dbVal, ok := v.(driver.Valuer); ok {
			if val, err := dbVal.Value(); err == nil {
				(*args)[i] = val
			} else {
				return errgo.Mask(err)
			}
		}
	}
	return nil
}
