package phpsession

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/corestoreio/csfw/util/php/phpserialize"
)

func TestEncodeBooleanValue(t *testing.T) {
	data := PhpSession{
		"login_ok": true,
	}

	encoder := NewPhpEncoder(data)
	if result, err := encoder.Encode(); err != nil {
		t.Errorf("Can not encode boolens value %#v \n", err)
	} else {
		if result != "login_ok|b:1;" {
			t.Errorf("Boolean value was encoded incorrectly %v \n", result)
		}
	}
}

func TestEncodeIntValue(t *testing.T) {
	data := PhpSession{
		"inteiro": 34,
	}

	encoder := NewPhpEncoder(data)
	if result, err := encoder.Encode(); err != nil {
		t.Errorf("Can not encode int value %#v \n", err)
	} else {
		if result != "inteiro|i:34;" {
			t.Errorf("Int value was encoded incorrectly %v \n", result)
		}
	}
}

func TestEncodeFloatValue(t *testing.T) {
	data := PhpSession{
		"float_test": 34.4679999999,
	}

	encoder := NewPhpEncoder(data)
	if result, err := encoder.Encode(); err != nil {
		t.Errorf("Can not encode float value %#v \n", err)
	} else {
		// 34.467999999900002 - PHP has precision = 17 by default
		if result != "float_test|d:34.467999999900002;" {
			t.Errorf("Float value was encoded incorrectly %v\n", result)
		}
	}
}

func TestEncodeStringValue(t *testing.T) {
	data := PhpSession{
		"name": "some text",
	}

	encoder := NewPhpEncoder(data)
	if result, err := encoder.Encode(); err != nil {
		t.Errorf("Can not encode string value %#v \n", err)
	} else {
		if result != "name|s:9:\"some text\";" {
			t.Errorf("String value was encoded incorrectly %v\n", result)
		}
	}
}

func TestEncodeArrayValue(t *testing.T) {
	data := PhpSession{
		"arr": phpserialize.PhpArray{
			// Zero element
			//phpserialize.PhpValue(0): 5,
			0:       5,
			"test":  true,
			"test2": nil,
		},
	}

	encoder := NewPhpEncoder(data)
	if result, err := encoder.Encode(); err != nil {
		t.Errorf("Can not encode array value %#v \n", err)
	} else {
		if !strings.Contains(result, "i:0;i:5;") || !strings.Contains(result, "s:4:\"test\";b:1") || !strings.Contains(result, "s:5:\"test2\";N") {
			t.Errorf("Array value was encoded incorrectly %v\n", result)
		}
	}
}

func TestEncodeObjectValue(t *testing.T) {
	obj := phpserialize.NewPhpObject("TestObject")
	obj.SetPublic("a", 5)
	obj.SetProtected("c", 8)
	obj.SetPrivate("b", "priv")
	data := PhpSession{
		"obj": obj,
	}

	encoder := NewPhpEncoder(data)
	if result, err := encoder.Encode(); err != nil {
		t.Errorf("Can not encode object value %#v \n", err)
	} else {
		if !strings.Contains(result, "s:1:\"a\";i:5") || !strings.Contains(result, "10:\"TestObject\"") || !strings.Contains(result, "s:13:\"\x00TestObject\x00b\";s:4:\"priv\"") || !strings.Contains(result, "s:4:\"\x00*\x00c\";i:8") {
			t.Errorf("Object value was encoded incorrectly %v\n", result)
		}
	}
}

func TestEncodeSerializableObjectValueNoFunc(t *testing.T) {
	obj := phpserialize.NewPhpObjectSerialized("TestObject")
	obj.SetData("a:3:{s:1:\"a\";i:5;s:1:\"b\";s:4:\"priv\";s:1:\"c\";i:8;}")
	data := PhpSession{
		"obj": obj,
	}

	encoder := NewPhpEncoder(data)
	if result, err := encoder.Encode(); err != nil {
		t.Errorf("Can not encode object value %#v \n", err)
	} else {
		if !strings.Contains(result, "a:3:{s:1:\"a\";i:5;s:1:\"b\";s:4:\"priv\";s:1:\"c\";i:8;}") || !strings.Contains(result, "C:10:\"TestObject\"") {
			t.Errorf("Object value was encoded incorrectly %v\n", result)
		}
	}
}

func TestEncodeSerializableObjectValue(t *testing.T) {
	arr := phpserialize.PhpArray{
		"a": 5,
		"b": "priv",
		"c": 8,
	}
	obj := phpserialize.NewPhpObjectSerialized("TestObject")
	obj.SetValue(phpserialize.PhpValue(arr))
	data := PhpSession{
		"obj": obj,
	}

	encoder := NewPhpEncoder(data)
	encoder.SetSerializedEncodeFunc(phpserialize.SerializedEncodeFunc(phpserialize.Serialize))
	if result, err := encoder.Encode(); err != nil {
		t.Errorf("Can not encode object value %#v \n", err)
	} else {
		if !strings.Contains(result, "C:10:\"TestObject\"") {
			t.Errorf("Object value was encoded incorrectly %v\n", result)
		} else if !strings.Contains(result, "s:1:\"a\";i:5;") {
			t.Errorf("Object value was encoded incorrectly %v\n", result)
		} else if !strings.Contains(result, "s:1:\"b\";s:4:\"priv\";") {
			t.Errorf("Object value was encoded incorrectly %v\n", result)
		} else if !strings.Contains(result, "s:1:\"c\";i:8;") {
			t.Errorf("Object value was encoded incorrectly %v\n", result)
		}
	}
}

func TestEncodeSerializableObjectValueJSON(t *testing.T) {
	var f phpserialize.SerializedEncodeFunc
	f = func(v phpserialize.PhpValue) (string, error) {
		res, err := json.Marshal(v)
		return string(res), err
	}

	obj := phpserialize.NewPhpObjectSerialized("Bar")
	obj.SetValue(map[string]string{"public": "public"})
	data := PhpSession{
		"bar": obj,
	}

	encoder := NewPhpEncoder(data)
	encoder.SetSerializedEncodeFunc(f)
	if result, err := encoder.Encode(); err != nil {
		t.Errorf("Can not encode object value %#v \n", err)
	} else {
		if result != "bar|C:3:\"Bar\":19:{{\"public\":\"public\"}}" {
			t.Errorf("Object value was encoded incorrectly %v\n", result)
		}
	}
}
