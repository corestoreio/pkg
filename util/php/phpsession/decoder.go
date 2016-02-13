package phpsession

import (
	"bytes"
	"io"
	"strings"

	"github.com/corestoreio/csfw/util/php/phpserialize"
)

type PhpDecoder struct {
	source  *strings.Reader
	decoder *phpserialize.UnSerializer
}

func NewPhpDecoder(phpSession string) *PhpDecoder {
	decoder := &PhpDecoder{
		source:  strings.NewReader(phpSession),
		decoder: phpserialize.NewUnSerializer(""),
	}
	decoder.decoder.SetReader(decoder.source)
	return decoder
}

func (self *PhpDecoder) SetSerializedDecodeFunc(f phpserialize.SerializedDecodeFunc) {
	self.decoder.SetSerializedDecodeFunc(f)
}

func (self *PhpDecoder) Decode() (PhpSession, error) {
	var (
		name  string
		err   error
		value phpserialize.PhpValue
	)
	res := make(PhpSession)

	for {
		if name, err = self.readName(); err != nil {
			break
		}
		if value, err = self.decoder.Decode(); err != nil {
			break
		}
		res[name] = value
	}

	if err == io.EOF {
		err = nil
	}
	return res, err
}

func (self *PhpDecoder) readName() (string, error) {
	var (
		token rune
		err   error
	)
	buf := bytes.NewBuffer([]byte{})
	for {
		if token, _, err = self.source.ReadRune(); err != nil || token == SEPARATOR_VALUE_NAME {
			break
		} else {
			buf.WriteRune(token)
		}
	}
	return buf.String(), err
}
