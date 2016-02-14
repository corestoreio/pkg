package phpserialize

import (
	"bytes"
	"fmt"
	"io"
	"strconv"

	"github.com/corestoreio/csfw/util/bufferpool"
)

const UnseriazableObjectMaxLen = 10 * 1024 * 1024 * 1024

func UnSerialize(s []byte) (PhpValue, error) {
	dec := NewUnSerializer(s)
	dec.SetSerializedDecodeFunc(SerializedDecodeFunc(UnSerialize))
	return dec.Decode()
}

type UnSerializerReader interface {
	io.Reader
	io.RuneReader
}

type UnSerializer struct {
	source     []byte
	r          UnSerializerReader
	lastErr    error
	decodeFunc SerializedDecodeFunc
}

func NewUnSerializer(data []byte) *UnSerializer {
	return &UnSerializer{
		source: data,
	}
}

func (us *UnSerializer) SetReader(r UnSerializerReader) {
	us.r = r
}

func (us *UnSerializer) SetSerializedDecodeFunc(f SerializedDecodeFunc) {
	us.decodeFunc = f
}

func (us *UnSerializer) Decode() (v PhpValue, err error) {
	if us.r == nil {
		us.r = bytes.NewReader(us.source)
	}

	if token, _, err := us.r.ReadRune(); err == nil {
		switch token {
		default:
			us.saveError(fmt.Errorf("phpserialize: Unknown token %#U", token))
		case TokeNull:
			v = us.decodeNull()
		case TokenBool:
			v = us.decodeBool()
		case TokenInt:
			v = us.decodeNumber(false)
		case TokenFloat:
			v = us.decodeNumber(true)
		case TokenString:
			v = us.decodeString(DelimiterStringLeft, DelimiterStringRight, true)
		case TokenArray:
			v = us.decodeArray()
		case TokenObject:
			v = us.decodeObject()
		case TokenObjectSerialized:
			v = us.decodeSerialized()
		case TokenReference, TOkenReferenceObject:
			v = us.decodeReference()
		case TokenSplArray:
			v = us.decodeSplArray()

		}
	}
	err = us.lastErr
	return
}

func (us *UnSerializer) decodeNull() PhpValue {
	us.expect(SeparatorValues)
	return nil
}

func (us *UnSerializer) decodeBool() PhpValue {
	var (
		raw rune
		err error
	)
	us.expect(SepratorValueTypes)

	if raw, _, err = us.r.ReadRune(); err != nil {
		us.saveError(fmt.Errorf("phpserialize: Error while reading bool value: %v", err))
	}

	us.expect(SeparatorValues)
	return raw == '1'
}

func (us *UnSerializer) decodeNumber(isFloat bool) PhpValue {
	var (
		raw string
		err error
		val PhpValue
	)
	us.expect(SepratorValueTypes)

	if raw, err = us.readUntil(SeparatorValues); err != nil {
		us.saveError(fmt.Errorf("phpserialize: Error while reading number value: %v", err))
	} else {
		if isFloat {
			if val, err = strconv.ParseFloat(raw, 64); err != nil {
				us.saveError(fmt.Errorf("phpserialize: Unable to convert %s to float: %v", raw, err))
			}
		} else {
			if val, err = strconv.Atoi(raw); err != nil {
				us.saveError(fmt.Errorf("phpserialize: Unable to convert %s to int: %v", raw, err))
			}
		}
	}

	return val
}

func (us *UnSerializer) decodeString(left, right rune, isFinal bool) PhpValue {
	var (
		err     error
		val     PhpValue
		strLen  int
		readLen int
	)

	strLen = us.readLen()
	us.expect(left)

	if strLen > 0 {
		buf := make([]byte, strLen, strLen)
		if readLen, err = us.r.Read(buf); err != nil {
			us.saveError(fmt.Errorf("phpserialize: Error while reading string value: %v", err))
		} else {
			if readLen != strLen {
				us.saveError(fmt.Errorf("phpserialize: Unable to read string. Expected %d but have got %d bytes", strLen, readLen))
			} else {
				val = string(buf)
			}
		}
	}

	us.expect(right)
	if isFinal {
		us.expect(SeparatorValues)
	}
	return val
}

func (us *UnSerializer) decodeArray() PhpValue {
	var arrLen int
	val := make(PhpArray)

	arrLen = us.readLen()
	us.expect(DelimiterObjectLeft)

	for i := 0; i < arrLen; i++ {
		k, errKey := us.Decode()
		v, errVal := us.Decode()

		if errKey == nil && errVal == nil {
			val[k] = v
			/*switch t := k.(type) {
			default:
				self.saveError(fmt.Errorf("phpserialize: Unexpected key type %T", t))
			case string:
				stringKey, _ := k.(string)
				val[stringKey] = v
			case int:
				intKey, _ := k.(int)
				val[strconv.Itoa(intKey)] = v
			}*/
		} else {
			us.saveError(fmt.Errorf("phpserialize: Error while reading key or(and) value of array"))
		}
	}

	us.expect(DelimiterObjectRight)
	return val
}

func (us *UnSerializer) decodeObject() PhpValue {
	val := &PhpObject{
		className: us.readClassName(),
	}

	rawMembers := us.decodeArray()
	val.members, _ = rawMembers.(PhpArray)

	return val
}

func (us *UnSerializer) decodeSerialized() PhpValue {
	val := &PhpObjectSerialized{
		className: us.readClassName(),
	}

	rawData := us.decodeString(DelimiterObjectLeft, DelimiterObjectRight, false)
	val.data, _ = rawData.(string)

	if us.decodeFunc != nil && val.data != "" {
		var err error
		if val.value, err = us.decodeFunc([]byte(val.data)); err != nil {
			us.saveError(err)
		}
	}

	return val
}

func (us *UnSerializer) decodeReference() PhpValue {
	us.expect(SepratorValueTypes)
	if _, err := us.readUntil(SeparatorValues); err != nil {
		us.saveError(fmt.Errorf("phpserialize: Error while reading reference value: %v", err))
	}
	return nil
}

func (us *UnSerializer) expect(expected rune) {
	token, _, err := us.r.ReadRune()
	switch {
	case err != nil:
		us.saveError(fmt.Errorf("phpserialize: Error while reading expected rune %#U: %v", expected, err))
	case token != expected:
		us.saveError(fmt.Errorf("phpserialize: Expected %#U but have got %#U", expected, token))
	}
}

func (us *UnSerializer) readUntil(stop rune) (string, error) {

	buf := bufferpool.Get()
	defer bufferpool.Put(buf)

	for {
		token, _, err := us.r.ReadRune()
		if err != nil || token == stop {
			return buf.String(), err
		}
		buf.WriteRune(token)
	}
	return buf.String(), nil
}

func (us *UnSerializer) readLen() int {
	var (
		raw string
		err error
		val int
	)
	us.expect(SepratorValueTypes)

	if raw, err = us.readUntil(SepratorValueTypes); err != nil {
		us.saveError(fmt.Errorf("phpserialize: Error while reading lenght of value: %v", err))
	} else {
		if val, err = strconv.Atoi(raw); err != nil {
			us.saveError(fmt.Errorf("phpserialize: Unable to convert %s to int: %v", raw, err))
		} else if val > UnseriazableObjectMaxLen {
			us.saveError(fmt.Errorf("phpserialize: Unserializable object length looks too big(%d). If you are sure you wanna unserialise it, please increase UNSERIAZABLE_OBJECT_MAX_LEN const: %s", val, err))
			val = 0
		}
	}
	return val
}

func (us *UnSerializer) readClassName() (res string) {
	rawClass := us.decodeString(DelimiterStringLeft, DelimiterStringRight, false)
	res, _ = rawClass.(string)
	return
}

func (us *UnSerializer) saveError(err error) {
	if us.lastErr == nil {
		us.lastErr = err
	}
}

func (us *UnSerializer) decodeSplArray() PhpValue {
	var err error
	val := &PhpSplArray{}

	us.expect(SepratorValueTypes)
	us.expect(TokenInt)

	flags := us.decodeNumber(false)
	if flags == nil {
		us.saveError(fmt.Errorf("phpserialize: Unable to read flags of SplArray"))
		return nil
	}
	val.flags = PhpValueInt(flags)

	if val.array, err = us.Decode(); err != nil {
		us.saveError(fmt.Errorf("phpserialize: Can't parse SplArray: %v", err))
		return nil
	}

	us.expect(SeparatorValues)
	us.expect(TokenSplArrayMembers)
	us.expect(SepratorValueTypes)

	if val.properties, err = us.Decode(); err != nil {
		us.saveError(fmt.Errorf("phpserialize: Can't parse properties of SplArray: %v", err))
		return nil
	}

	return val
}
