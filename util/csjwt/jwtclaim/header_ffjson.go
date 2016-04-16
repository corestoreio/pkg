// DO NOT EDIT!
// Code generated by ffjson <https://github.com/pquerna/ffjson>
// source: header.go
// DO NOT EDIT!

package jwtclaim

import (
	"bytes"
	"fmt"

	fflib "github.com/pquerna/ffjson/fflib/v1"
)

func (mj *HeadSegments) MarshalJSON() ([]byte, error) {
	var buf fflib.Buffer
	if mj == nil {
		buf.WriteString("null")
		return buf.Bytes(), nil
	}
	err := mj.MarshalJSONBuf(&buf)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
func (mj *HeadSegments) MarshalJSONBuf(buf fflib.EncodingBuffer) error {
	if mj == nil {
		buf.WriteString("null")
		return nil
	}
	var err error
	var obj []byte
	_ = obj
	_ = err
	buf.WriteString(`{ `)
	if len(mj.Algorithm) != 0 {
		buf.WriteString(`"alg":`)
		fflib.WriteJsonString(buf, string(mj.Algorithm))
		buf.WriteByte(',')
	}
	if len(mj.Type) != 0 {
		buf.WriteString(`"typ":`)
		fflib.WriteJsonString(buf, string(mj.Type))
		buf.WriteByte(',')
	}
	if len(mj.JKU) != 0 {
		buf.WriteString(`"jku":`)
		fflib.WriteJsonString(buf, string(mj.JKU))
		buf.WriteByte(',')
	}
	if len(mj.KID) != 0 {
		buf.WriteString(`"kid":`)
		fflib.WriteJsonString(buf, string(mj.KID))
		buf.WriteByte(',')
	}
	if len(mj.X5U) != 0 {
		buf.WriteString(`"x5u":`)
		fflib.WriteJsonString(buf, string(mj.X5U))
		buf.WriteByte(',')
	}
	if len(mj.X5T) != 0 {
		buf.WriteString(`"x5t":`)
		fflib.WriteJsonString(buf, string(mj.X5T))
		buf.WriteByte(',')
	}
	buf.Rewind(1)
	buf.WriteByte('}')
	return nil
}

const (
	ffj_t_HeadSegmentsbase = iota
	ffj_t_HeadSegmentsno_such_key

	ffj_t_HeadSegments_Algorithm

	ffj_t_HeadSegments_Type

	ffj_t_HeadSegments_JKU

	ffj_t_HeadSegments_KID

	ffj_t_HeadSegments_X5U

	ffj_t_HeadSegments_X5T
)

var ffj_key_HeadSegments_Algorithm = []byte("alg")

var ffj_key_HeadSegments_Type = []byte("typ")

var ffj_key_HeadSegments_JKU = []byte("jku")

var ffj_key_HeadSegments_KID = []byte("kid")

var ffj_key_HeadSegments_X5U = []byte("x5u")

var ffj_key_HeadSegments_X5T = []byte("x5t")

func (uj *HeadSegments) UnmarshalJSON(input []byte) error {
	fs := fflib.NewFFLexer(input)
	return uj.UnmarshalJSONFFLexer(fs, fflib.FFParse_map_start)
}

func (uj *HeadSegments) UnmarshalJSONFFLexer(fs *fflib.FFLexer, state fflib.FFParseState) error {
	var err error = nil
	currentKey := ffj_t_HeadSegmentsbase
	_ = currentKey
	tok := fflib.FFTok_init
	wantedTok := fflib.FFTok_init

mainparse:
	for {
		tok = fs.Scan()
		//	println(fmt.Sprintf("debug: tok: %v  state: %v", tok, state))
		if tok == fflib.FFTok_error {
			goto tokerror
		}

		switch state {

		case fflib.FFParse_map_start:
			if tok != fflib.FFTok_left_bracket {
				wantedTok = fflib.FFTok_left_bracket
				goto wrongtokenerror
			}
			state = fflib.FFParse_want_key
			continue

		case fflib.FFParse_after_value:
			if tok == fflib.FFTok_comma {
				state = fflib.FFParse_want_key
			} else if tok == fflib.FFTok_right_bracket {
				goto done
			} else {
				wantedTok = fflib.FFTok_comma
				goto wrongtokenerror
			}

		case fflib.FFParse_want_key:
			// json {} ended. goto exit. woo.
			if tok == fflib.FFTok_right_bracket {
				goto done
			}
			if tok != fflib.FFTok_string {
				wantedTok = fflib.FFTok_string
				goto wrongtokenerror
			}

			kn := fs.Output.Bytes()
			if len(kn) <= 0 {
				// "" case. hrm.
				currentKey = ffj_t_HeadSegmentsno_such_key
				state = fflib.FFParse_want_colon
				goto mainparse
			} else {
				switch kn[0] {

				case 'a':

					if bytes.Equal(ffj_key_HeadSegments_Algorithm, kn) {
						currentKey = ffj_t_HeadSegments_Algorithm
						state = fflib.FFParse_want_colon
						goto mainparse
					}

				case 'j':

					if bytes.Equal(ffj_key_HeadSegments_JKU, kn) {
						currentKey = ffj_t_HeadSegments_JKU
						state = fflib.FFParse_want_colon
						goto mainparse
					}

				case 'k':

					if bytes.Equal(ffj_key_HeadSegments_KID, kn) {
						currentKey = ffj_t_HeadSegments_KID
						state = fflib.FFParse_want_colon
						goto mainparse
					}

				case 't':

					if bytes.Equal(ffj_key_HeadSegments_Type, kn) {
						currentKey = ffj_t_HeadSegments_Type
						state = fflib.FFParse_want_colon
						goto mainparse
					}

				case 'x':

					if bytes.Equal(ffj_key_HeadSegments_X5U, kn) {
						currentKey = ffj_t_HeadSegments_X5U
						state = fflib.FFParse_want_colon
						goto mainparse

					} else if bytes.Equal(ffj_key_HeadSegments_X5T, kn) {
						currentKey = ffj_t_HeadSegments_X5T
						state = fflib.FFParse_want_colon
						goto mainparse
					}

				}

				if fflib.AsciiEqualFold(ffj_key_HeadSegments_X5T, kn) {
					currentKey = ffj_t_HeadSegments_X5T
					state = fflib.FFParse_want_colon
					goto mainparse
				}

				if fflib.AsciiEqualFold(ffj_key_HeadSegments_X5U, kn) {
					currentKey = ffj_t_HeadSegments_X5U
					state = fflib.FFParse_want_colon
					goto mainparse
				}

				if fflib.EqualFoldRight(ffj_key_HeadSegments_KID, kn) {
					currentKey = ffj_t_HeadSegments_KID
					state = fflib.FFParse_want_colon
					goto mainparse
				}

				if fflib.EqualFoldRight(ffj_key_HeadSegments_JKU, kn) {
					currentKey = ffj_t_HeadSegments_JKU
					state = fflib.FFParse_want_colon
					goto mainparse
				}

				if fflib.SimpleLetterEqualFold(ffj_key_HeadSegments_Type, kn) {
					currentKey = ffj_t_HeadSegments_Type
					state = fflib.FFParse_want_colon
					goto mainparse
				}

				if fflib.SimpleLetterEqualFold(ffj_key_HeadSegments_Algorithm, kn) {
					currentKey = ffj_t_HeadSegments_Algorithm
					state = fflib.FFParse_want_colon
					goto mainparse
				}

				currentKey = ffj_t_HeadSegmentsno_such_key
				state = fflib.FFParse_want_colon
				goto mainparse
			}

		case fflib.FFParse_want_colon:
			if tok != fflib.FFTok_colon {
				wantedTok = fflib.FFTok_colon
				goto wrongtokenerror
			}
			state = fflib.FFParse_want_value
			continue
		case fflib.FFParse_want_value:

			if tok == fflib.FFTok_left_brace || tok == fflib.FFTok_left_bracket || tok == fflib.FFTok_integer || tok == fflib.FFTok_double || tok == fflib.FFTok_string || tok == fflib.FFTok_bool || tok == fflib.FFTok_null {
				switch currentKey {

				case ffj_t_HeadSegments_Algorithm:
					goto handle_Algorithm

				case ffj_t_HeadSegments_Type:
					goto handle_Type

				case ffj_t_HeadSegments_JKU:
					goto handle_JKU

				case ffj_t_HeadSegments_KID:
					goto handle_KID

				case ffj_t_HeadSegments_X5U:
					goto handle_X5U

				case ffj_t_HeadSegments_X5T:
					goto handle_X5T

				case ffj_t_HeadSegmentsno_such_key:
					err = fs.SkipField(tok)
					if err != nil {
						return fs.WrapErr(err)
					}
					state = fflib.FFParse_after_value
					goto mainparse
				}
			} else {
				goto wantedvalue
			}
		}
	}

handle_Algorithm:

	/* handler: uj.Algorithm type=string kind=string quoted=false*/

	{

		{
			if tok != fflib.FFTok_string && tok != fflib.FFTok_null {
				return fs.WrapErr(fmt.Errorf("cannot unmarshal %s into Go value for string", tok))
			}
		}

		if tok == fflib.FFTok_null {

		} else {

			outBuf := fs.Output.Bytes()

			uj.Algorithm = string(string(outBuf))

		}
	}

	state = fflib.FFParse_after_value
	goto mainparse

handle_Type:

	/* handler: uj.Type type=string kind=string quoted=false*/

	{

		{
			if tok != fflib.FFTok_string && tok != fflib.FFTok_null {
				return fs.WrapErr(fmt.Errorf("cannot unmarshal %s into Go value for string", tok))
			}
		}

		if tok == fflib.FFTok_null {

		} else {

			outBuf := fs.Output.Bytes()

			uj.Type = string(string(outBuf))

		}
	}

	state = fflib.FFParse_after_value
	goto mainparse

handle_JKU:

	/* handler: uj.JKU type=string kind=string quoted=false*/

	{

		{
			if tok != fflib.FFTok_string && tok != fflib.FFTok_null {
				return fs.WrapErr(fmt.Errorf("cannot unmarshal %s into Go value for string", tok))
			}
		}

		if tok == fflib.FFTok_null {

		} else {

			outBuf := fs.Output.Bytes()

			uj.JKU = string(string(outBuf))

		}
	}

	state = fflib.FFParse_after_value
	goto mainparse

handle_KID:

	/* handler: uj.KID type=string kind=string quoted=false*/

	{

		{
			if tok != fflib.FFTok_string && tok != fflib.FFTok_null {
				return fs.WrapErr(fmt.Errorf("cannot unmarshal %s into Go value for string", tok))
			}
		}

		if tok == fflib.FFTok_null {

		} else {

			outBuf := fs.Output.Bytes()

			uj.KID = string(string(outBuf))

		}
	}

	state = fflib.FFParse_after_value
	goto mainparse

handle_X5U:

	/* handler: uj.X5U type=string kind=string quoted=false*/

	{

		{
			if tok != fflib.FFTok_string && tok != fflib.FFTok_null {
				return fs.WrapErr(fmt.Errorf("cannot unmarshal %s into Go value for string", tok))
			}
		}

		if tok == fflib.FFTok_null {

		} else {

			outBuf := fs.Output.Bytes()

			uj.X5U = string(string(outBuf))

		}
	}

	state = fflib.FFParse_after_value
	goto mainparse

handle_X5T:

	/* handler: uj.X5T type=string kind=string quoted=false*/

	{

		{
			if tok != fflib.FFTok_string && tok != fflib.FFTok_null {
				return fs.WrapErr(fmt.Errorf("cannot unmarshal %s into Go value for string", tok))
			}
		}

		if tok == fflib.FFTok_null {

		} else {

			outBuf := fs.Output.Bytes()

			uj.X5T = string(string(outBuf))

		}
	}

	state = fflib.FFParse_after_value
	goto mainparse

wantedvalue:
	return fs.WrapErr(fmt.Errorf("wanted value token, but got token: %v", tok))
wrongtokenerror:
	return fs.WrapErr(fmt.Errorf("ffjson: wanted token: %v, but got token: %v output=%s", wantedTok, tok, fs.Output.String()))
tokerror:
	if fs.BigError != nil {
		return fs.WrapErr(fs.BigError)
	}
	err = fs.Error.ToError()
	if err != nil {
		return fs.WrapErr(err)
	}
	panic("ffjson-generated: unreachable, please report bug.")
done:
	return nil
}
