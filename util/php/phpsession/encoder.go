package phpsession

import (
	"fmt"

	"github.com/corestoreio/csfw/util/bufferpool"
	"github.com/corestoreio/csfw/util/php/phpserialize"
)

type PhpEncoder struct {
	data    PhpSession
	encoder *phpserialize.Serializer
}

func NewPhpEncoder(data PhpSession) *PhpEncoder {
	return &PhpEncoder{
		data:    data,
		encoder: phpserialize.NewSerializer(),
	}
}

func (self *PhpEncoder) SetSerializedEncodeFunc(f phpserialize.SerializedEncodeFunc) {
	self.encoder.SetSerializedEncodeFunc(f)
}

func (self *PhpEncoder) Encode() (string, error) {
	if self.data == nil {
		return "", nil
	}
	buf := bufferpool.Get()
	defer bufferpool.Put(buf)

	for k, v := range self.data {
		buf.WriteString(k)
		buf.WriteRune(SEPARATOR_VALUE_NAME)
		val, err := self.encoder.Encode(v)
		if err != nil {
			return buf.String(), fmt.Errorf("php_session: error during encode value for %q: %v", k, err)
		}
		buf.WriteString(val)
	}

	return buf.String(), nil
}
