// +build easyjson

package null_test

import (
	"testing"

	"github.com/corestoreio/pkg/storage/null"

	"github.com/corestoreio/pkg/util/assert"
	"github.com/mailru/easyjson"
	"github.com/mailru/easyjson/jlexer"
	"github.com/mailru/easyjson/jwriter"
)

func TestEasyJSON(t *testing.T) {
	runner := func(in easyjson.Marshaler, out easyjson.Unmarshaler, wantJSON string) func(*testing.T) {
		return func(t *testing.T) {
			w := &jwriter.Writer{}
			in.MarshalEasyJSON(w)
			assert.NoError(t, w.Error)
			data := w.Buffer.BuildBytes()
			assert.Exactly(t, wantJSON, string(data))

			l := &jlexer.Lexer{Data: data}
			out.UnmarshalEasyJSON(l)
			assert.NoError(t, l.Error())
			assert.Exactly(t, out, in)
		}
	}

	t.Run("String Raw", runner(&null.String{
		Valid:  true,
		String: "pple",
	}, &null.String{}, "\"\uf8ffpple\""))

	t.Run("StringEmbedded all fields", runner(&StringEmbedded{
		String1: null.MakeString("String1"),
		String2: null.MakeString("pple"),
	}, &StringEmbedded{}, "{\"string_1\":\"String1\",\"string_2\":\"\uf8ffpple\"}"))

	t.Run("StringEmbedded String2", runner(&StringEmbedded{
		String2: null.MakeString("pple"),
	}, &StringEmbedded{}, "{\"string_2\":\"\uf8ffpple\"}"))

	t.Run("StringEmbedded String1", runner(&StringEmbedded{
		String1: null.MakeString(`p"ple`),
	}, &StringEmbedded{}, "{\"string_1\":\"\uf8ffp\\\"ple\",\"string_2\":null}"))

	t.Run("StringEmbedded empty", runner(&StringEmbedded{}, &StringEmbedded{}, "{\"string_2\":null}"))

}

//easyjson:json
type StringEmbedded struct {
	String1 null.String `json:"string_1,omitempty"`
	String2 null.String `json:"string_2"`
}
