package dmlgen

import (
	"strings"
	"unicode"

	"github.com/corestoreio/pkg/util/strs"
)

// lcFirst transforms the first character of a string to lower case.
func lcFirst(s string) string {
	sr := []rune(s)
	sr[0] = unicode.ToLower(sr[0])
	return string(sr)
}

func collectionName(name string) string {
	tg := strs.ToGoCamelCase(name)
	switch {
	case strings.HasSuffix(name, "y"):
		return tg[:len(tg)-1] + "ies"
	case strings.HasSuffix(name, "ch"):
		return tg + "es"
	case strings.HasSuffix(name, "x"):
		return tg + "es"
	case strings.HasSuffix(name, "us"):
		return tg + "i" // status -> stati
	case strings.HasSuffix(name, "um"):
		return tg + "en" // datum -> daten
	case strings.HasSuffix(name, "s"):
		return tg + "Collection" // stupid case, better ideas?
	default:
		return tg + "s"
	}
}
