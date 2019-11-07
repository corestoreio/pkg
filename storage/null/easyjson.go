// +build easyjson

// TODO easyjson does not yet respect build tags to be included when parsing
//  files to generate the code. yet there is a PR which refactores easyjson
//  parser to go/types.

package null

import (
	"github.com/mailru/easyjson/jlexer"
	"github.com/mailru/easyjson/jwriter"
)

// TODO use the athlete struct and run a benchmark comparison between easyjson, stdlib and json-iterator
// TODO add all other types
// TODO fuzzy testing gofuzz

func (a String) MarshalEasyJSON(w *jwriter.Writer) {
	if !a.Valid {
		w.Raw(nil, nil)
		return
	}
	w.String(a.Data)
}

func (a *String) UnmarshalEasyJSON(l *jlexer.Lexer) {
	if l.IsNull() {
		a.Valid = false
		a.Data = ""
		return
	}

	a.Valid = true
	a.Data = l.String()
}

// IsDefined implements easyjson.Optional interface, same as function IsZero of
// this type.
func (a String) IsDefined() bool {
	return a.Valid
}
