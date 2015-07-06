package dbr
import "database/sql/driver"

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

func argsValuer(args *[]interface{}){
	for i,v := range *args {
		if dbVal, ok := v.(driver.Valuer); ok {
			if val, err := dbVal.Value(); err == nil {
				(*args)[i] = val
			} else {
				panic(err) // @todo return error
			}
		}
	}
}
