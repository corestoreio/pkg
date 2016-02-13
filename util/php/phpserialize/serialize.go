package phpserialize

import (
	"bytes"
	"fmt"
	"strconv"

	"github.com/corestoreio/csfw/util/bufferpool"
)

func Serialize(v PhpValue) (string, error) {
	enc := NewSerializer()
	enc.SetSerializedEncodeFunc(SerializedEncodeFunc(Serialize))
	return enc.Encode(v)
}

type Serializer struct {
	lastErr    error
	encodeFunc SerializedEncodeFunc
}

func NewSerializer() *Serializer {
	return &Serializer{}
}

func (s *Serializer) SetSerializedEncodeFunc(f SerializedEncodeFunc) {
	s.encodeFunc = f
}

func (s *Serializer) Encode(v PhpValue) (string, error) {
	buf := bufferpool.Get()
	defer bufferpool.Put(buf)

	switch t := v.(type) {
	default:
		s.saveError(fmt.Errorf("phpserialize: Unknown type %T with value %#v", t, v))
	case nil:
		s.encodeNull(buf)
	case bool:
		s.encodeBool(buf, v)
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64:
		s.encodeNumber(buf, v)
	case string:
		s.encodeString(buf, v, DELIMITER_STRING_LEFT, DELIMITER_STRING_RIGHT, true)
	case PhpArray, map[PhpValue]PhpValue, PhpSlice:
		s.encodeArray(buf, v, true)
	case *PhpObject:
		s.encodeObject(buf, v)
	case *PhpObjectSerialized:
		s.encodeSerialized(buf, v)
	case *PhpSplArray:
		s.encodeSplArray(buf, v)
	}
	return buf.String(), s.lastErr
}

func (s *Serializer) encodeNull(buf *bytes.Buffer) {
	buf.WriteRune(TOKEN_NULL)
	buf.WriteRune(SEPARATOR_VALUES)
}

func (s *Serializer) encodeBool(buf *bytes.Buffer, v PhpValue) {
	buf.WriteRune(TOKEN_BOOL)
	buf.WriteRune(SEPARATOR_VALUE_TYPE)

	var bs = "0"
	if bVal, ok := v.(bool); ok && bVal == true {
		bs = "1"
	}
	buf.WriteString(bs)

	buf.WriteRune(SEPARATOR_VALUES)
}

func (s *Serializer) encodeNumber(buf *bytes.Buffer, v PhpValue) {
	var val string

	isFloat := false

	switch v.(type) {
	default:
		val = "0"
	case int:
		intVal, _ := v.(int)
		val = strconv.FormatInt(int64(intVal), 10)
	case int8:
		intVal, _ := v.(int8)
		val = strconv.FormatInt(int64(intVal), 10)
	case int16:
		intVal, _ := v.(int16)
		val = strconv.FormatInt(int64(intVal), 10)
	case int32:
		intVal, _ := v.(int32)
		val = strconv.FormatInt(int64(intVal), 10)
	case int64:
		intVal, _ := v.(int64)
		val = strconv.FormatInt(int64(intVal), 10)
	case uint:
		intVal, _ := v.(uint)
		val = strconv.FormatUint(uint64(intVal), 10)
	case uint8:
		intVal, _ := v.(uint8)
		val = strconv.FormatUint(uint64(intVal), 10)
	case uint16:
		intVal, _ := v.(uint16)
		val = strconv.FormatUint(uint64(intVal), 10)
	case uint32:
		intVal, _ := v.(uint32)
		val = strconv.FormatUint(uint64(intVal), 10)
	case uint64:
		intVal, _ := v.(uint64)
		val = strconv.FormatUint(uint64(intVal), 10)
	// PHP has precision = 17 by default
	case float32:
		floatVal, _ := v.(float32)
		val = strconv.FormatFloat(float64(floatVal), FORMATTER_FLOAT, FORMATTER_PRECISION, 32)
		isFloat = true
	case float64:
		floatVal, _ := v.(float64)
		val = strconv.FormatFloat(float64(floatVal), FORMATTER_FLOAT, FORMATTER_PRECISION, 64)
		isFloat = true
	}

	var tok = TOKEN_INT
	if isFloat {
		tok = TOKEN_FLOAT
	}
	buf.WriteRune(tok)

	buf.WriteRune(SEPARATOR_VALUE_TYPE)
	buf.WriteString(val)
	buf.WriteRune(SEPARATOR_VALUES)
}

func (s *Serializer) encodeString(buf *bytes.Buffer, v PhpValue, left, right rune, isFinal bool) {
	val, _ := v.(string)

	if isFinal {
		buf.WriteRune(TOKEN_STRING)
	}

	buf.WriteString(s.prepareLen(len(val)))
	buf.WriteRune(left)
	buf.WriteString(val)
	buf.WriteRune(right)

	if isFinal {
		buf.WriteRune(SEPARATOR_VALUES)
	}
}

func (s *Serializer) encodeArray(buf *bytes.Buffer, v PhpValue, isFinal bool) {
	var (
		arrLen int
		str    string
	)

	if isFinal {
		buf.WriteRune(TOKEN_ARRAY)
	}

	switch v.(type) {
	case PhpArray:
		arrVal, _ := v.(PhpArray)
		arrLen = len(arrVal)

		buf.WriteString(s.prepareLen(arrLen))
		buf.WriteRune(DELIMITER_OBJECT_LEFT)

		for k, v := range arrVal {
			str, _ = s.Encode(k)
			buf.WriteString(str)
			str, _ = s.Encode(v)
			buf.WriteString(str)
		}

	case map[PhpValue]PhpValue:
		arrVal, _ := v.(map[PhpValue]PhpValue)
		arrLen = len(arrVal)

		buf.WriteString(s.prepareLen(arrLen))
		buf.WriteRune(DELIMITER_OBJECT_LEFT)

		for k, v := range arrVal {
			str, _ = s.Encode(k)
			buf.WriteString(str)
			str, _ = s.Encode(v)
			buf.WriteString(str)
		}
	case PhpSlice:
		arrVal, _ := v.(PhpSlice)
		arrLen = len(arrVal)

		buf.WriteString(s.prepareLen(arrLen))
		buf.WriteRune(DELIMITER_OBJECT_LEFT)

		for k, v := range arrVal {
			str, _ = s.Encode(k)
			buf.WriteString(str)
			str, _ = s.Encode(v)
			buf.WriteString(str)
		}
	}

	buf.WriteRune(DELIMITER_OBJECT_RIGHT)
}

func (s *Serializer) encodeObject(buf *bytes.Buffer, v PhpValue) {
	obj, _ := v.(*PhpObject)
	buf.WriteRune(TOKEN_OBJECT)
	s.prepareClassName(buf, obj.className)
	s.encodeArray(buf, obj.members, false)
	return
}

func (s *Serializer) encodeSerialized(buf *bytes.Buffer, v PhpValue) {
	var serialized string

	obj, _ := v.(*PhpObjectSerialized)
	buf.WriteRune(TOKEN_OBJECT_SERIALIZED)
	s.prepareClassName(buf, obj.className)

	if s.encodeFunc == nil {
		serialized = obj.GetData()
	} else {
		var err error
		if serialized, err = s.encodeFunc(obj.GetValue()); err != nil {
			s.saveError(err)
		}
	}

	s.encodeString(buf, serialized, DELIMITER_OBJECT_LEFT, DELIMITER_OBJECT_RIGHT, false)
}

func (s *Serializer) encodeSplArray(buf *bytes.Buffer, v PhpValue) {

	obj, _ := v.(*PhpSplArray)

	buf.WriteRune(TOKEN_SPL_ARRAY)
	buf.WriteRune(SEPARATOR_VALUE_TYPE)

	s.encodeNumber(buf, obj.flags)

	data, _ := s.Encode(obj.array)
	buf.WriteString(data)

	buf.WriteRune(SEPARATOR_VALUES)
	buf.WriteRune(TOKEN_SPL_ARRAY_MEMBERS)
	buf.WriteRune(SEPARATOR_VALUE_TYPE)

	data, _ = s.Encode(obj.properties)
	buf.WriteString(data)

}

func (s *Serializer) prepareLen(l int) string {
	return string(SEPARATOR_VALUE_TYPE) + strconv.Itoa(l) + string(SEPARATOR_VALUE_TYPE)
}

func (s *Serializer) prepareClassName(buf *bytes.Buffer, name string) {
	s.encodeString(buf, name, DELIMITER_STRING_LEFT, DELIMITER_STRING_RIGHT, false)
}

func (s *Serializer) saveError(err error) {
	if s.lastErr == nil {
		s.lastErr = err
	}
}
