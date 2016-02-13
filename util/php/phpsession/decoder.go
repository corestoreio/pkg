package phpsession

import (
	"bytes"
	"io"

	"github.com/corestoreio/csfw/util/bufferpool"
	"github.com/corestoreio/csfw/util/php/phpserialize"
)

type PhpDecoder struct {
	source interface {
		io.Reader
		io.RuneReader
	}
	decoder *phpserialize.UnSerializer
}

func NewPhpDecoder(phpSession []byte) *PhpDecoder {
	decoder := &PhpDecoder{
		source:  bytes.NewReader(phpSession),
		decoder: phpserialize.NewUnSerializer(nil),
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
	buf := bufferpool.Get()
	defer bufferpool.Put(buf)
	for {
		if token, _, err := self.source.ReadRune(); err != nil || token == SEPARATOR_VALUE_NAME {
			return buf.String(), err
		} else {
			buf.WriteRune(token)
		}
	}
	return buf.String(), nil
}
