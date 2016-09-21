// Copyright © 2014 Steve Francia <spf@spf13.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package conv

import (
	"fmt"
	"html/template"
	"math"
	"testing"
	"time"

	"github.com/corestoreio/csfw/config/cfgpath"
	"github.com/corestoreio/csfw/storage/text"
	"github.com/corestoreio/csfw/util/errors"
	"github.com/stretchr/testify/assert"
)

func TestToInt(t *testing.T) {

	var eight interface{} = 8
	i_8 := int(8)
	assert.Exactly(t, i_8, ToInt(8))
	assert.Exactly(t, i_8, ToInt(8.31))
	assert.Exactly(t, i_8, ToInt("8"))
	assert.Exactly(t, int(1), ToInt(true))
	assert.Exactly(t, int(0), ToInt(false))
	assert.Exactly(t, i_8, ToInt(eight))
}

func TestToInt64(t *testing.T) {

	var eight interface{} = 8
	var i64_8 = int64(8)
	assert.Exactly(t, i64_8, ToInt64(8))
	assert.Exactly(t, i64_8, ToInt64(8.31))
	assert.Exactly(t, i64_8, ToInt64("8"))
	assert.Exactly(t, int64(1), ToInt64(true))
	assert.Exactly(t, int64(0), ToInt64(false))
	assert.Exactly(t, i64_8, ToInt64(eight))
}

func TestToFloat64(t *testing.T) {

	tests := []struct {
		have    interface{}
		want    float64
		wantErr bool
	}{
		{8, 8, false},
		{math.E, math.E, false},
		{float32(4.56789), 4.567890167236328, false},

		{int64(56789), 56789, false},
		{int32(56789), 56789, false},
		{int16(math.MaxInt16), float64(math.MaxInt16), false},
		{int8(math.MaxInt8), float64(math.MaxInt8), false},
		{int(254), 254, false},

		{uint64(56789), 56789, false},
		{uint32(56789), 56789, false},
		{uint16(math.MaxUint16), float64(math.MaxUint16), false},
		{uint8(math.MaxUint8), float64(math.MaxUint8), false},
		{uint(254), 254, false},

		{fmt.Sprintf("%.10f", math.Phi), 1.6180339887, false},
		{`hello`, 0, true},

		{[]byte(fmt.Sprintf("%.10f", math.Phi)), 1.6180339887, false},
		{[]byte(`hello`), 0, true},
		{nil, 0, true},
	}
	for i, test := range tests {
		gotF64, gotErr := ToFloat64E(test.have)
		if test.wantErr {
			assert.True(t, errors.IsNotValid(gotErr), "Index %d => %s", i, gotErr)
			assert.Exactly(t, float64(0), gotF64, "Index %d", i)
			continue
		}
		assert.NoError(t, gotErr, "Index %d", i)
		assert.Exactly(t, test.want, gotF64, "Index %d", i)
	}
}

func TestToString(t *testing.T) {

	var foo interface{} = "one more time"
	assert.Equal(t, ToString(8), "8")
	assert.Equal(t, ToString(8.12), "8.12")
	assert.Equal(t, ToString([]byte("one time")), "one time")
	assert.Equal(t, ToString(template.HTML("one time")), "one time")
	assert.Equal(t, ToString(template.URL("http://somehost.foo")), "http://somehost.foo")
	assert.Equal(t, ToString(text.Chars("http://somehost.foo")), "http://somehost.foo")
	assert.Equal(t, ToString(cfgpath.NewRoute("http://somehost.foo")), "http://somehost.foo")
	assert.Equal(t, ToString(foo), "one more time")
	assert.Equal(t, ToString(nil), "")
	assert.Equal(t, ToString(true), "true")
	assert.Equal(t, ToString(false), "false")
	assert.Equal(t, ToString(cfgpath.MustNewByParts("aa/bb/cc").BindStore(33)), "stores/33/aa/bb/cc")
}

func TestToByte(t *testing.T) {

	var foo interface{} = []byte("one more time")
	assert.Equal(t, ToByte(8), []byte("8"))
	assert.Equal(t, ToByte(int64(8888)), []byte("8888"))
	assert.Equal(t, ToByte(8.12), []byte("8.12"))
	assert.Equal(t, ToByte([]byte("one time")), []byte("one time"))
	assert.Equal(t, ToByte(template.HTML("one time")), []byte("one time"))
	assert.Equal(t, ToByte(template.URL("http://somehost.foo")), []byte("http://somehost.foo"))
	assert.Equal(t, ToByte(text.Chars("http://somehost.foo")), []byte("http://somehost.foo"))
	assert.Equal(t, ToByte(cfgpath.NewRoute("http://somehost.foo")), []byte("http://somehost.foo"))
	assert.Equal(t, ToByte(foo), []byte("one more time"))
	assert.Equal(t, ToByte(nil), []byte(nil))
	assert.Equal(t, ToByte(true), []byte("true"))
	assert.Equal(t, ToByte(false), []byte("false"))
	assert.Equal(t, ToByte(cfgpath.MustNewByParts("aa/bb/cc").BindStore(33)), []byte("stores/33/aa/bb/cc"))

	b, err := ToByteE(uint8(1))
	assert.Nil(t, b)
	assert.True(t, errors.IsNotValid(err), "Error: %s", err)
}

type foo struct {
	val string
}

func (x foo) String() string {
	return x.val
}

func TestStringerToString(t *testing.T) {

	var x foo
	x.val = "bar"
	assert.Equal(t, "bar", ToString(x))
}

type fu struct {
	val string
}

func (x fu) Error() string {
	return x.val
}

func TestErrorToString(t *testing.T) {

	var x fu
	x.val = "bar"
	assert.Equal(t, "bar", ToString(x))
}

func TestMaps(t *testing.T) {

	var taxonomies = map[interface{}]interface{}{"tag": "tags", "group": "groups"}
	var stringMapBool = map[interface{}]interface{}{"v1": true, "v2": false}

	// ToStringMapString inputs/outputs
	var stringMapString = map[string]string{"key 1": "value 1", "key 2": "value 2", "key 3": "value 3"}
	var stringMapInterface = map[string]interface{}{"key 1": "value 1", "key 2": "value 2", "key 3": "value 3"}
	var interfaceMapString = map[interface{}]string{"key 1": "value 1", "key 2": "value 2", "key 3": "value 3"}
	var interfaceMapInterface = map[interface{}]interface{}{"key 1": "value 1", "key 2": "value 2", "key 3": "value 3"}

	// ToStringMapStringSlice inputs/outputs
	var stringMapStringSlice = map[string][]string{"key 1": {"value 1", "value 2", "value 3"}, "key 2": {"value 1", "value 2", "value 3"}, "key 3": {"value 1", "value 2", "value 3"}}
	var stringMapInterfaceSlice = map[string][]interface{}{"key 1": {"value 1", "value 2", "value 3"}, "key 2": {"value 1", "value 2", "value 3"}, "key 3": {"value 1", "value 2", "value 3"}}
	var stringMapStringSingleSliceFieldsResult = map[string][]string{"key 1": {"value", "1"}, "key 2": {"value", "2"}, "key 3": {"value", "3"}}
	var interfaceMapStringSlice = map[interface{}][]string{"key 1": {"value 1", "value 2", "value 3"}, "key 2": {"value 1", "value 2", "value 3"}, "key 3": {"value 1", "value 2", "value 3"}}
	var interfaceMapInterfaceSlice = map[interface{}][]interface{}{"key 1": {"value 1", "value 2", "value 3"}, "key 2": {"value 1", "value 2", "value 3"}, "key 3": {"value 1", "value 2", "value 3"}}

	var stringMapStringSliceMultiple = map[string][]string{"key 1": {"value 1", "value 2", "value 3"}, "key 2": {"value 1", "value 2", "value 3"}, "key 3": {"value 1", "value 2", "value 3"}}
	var stringMapStringSliceSingle = map[string][]string{"key 1": {"value 1"}, "key 2": {"value 2"}, "key 3": {"value 3"}}

	assert.Equal(t, ToStringMap(taxonomies), map[string]interface{}{"tag": "tags", "group": "groups"})
	assert.Equal(t, ToStringMapBool(stringMapBool), map[string]bool{"v1": true, "v2": false})

	// ToStringMapString tests
	assert.Equal(t, ToStringMapString(stringMapString), stringMapString)
	assert.Equal(t, ToStringMapString(stringMapInterface), stringMapString)
	assert.Equal(t, ToStringMapString(interfaceMapString), stringMapString)
	assert.Equal(t, ToStringMapString(interfaceMapInterface), stringMapString)

	// ToStringMapStringSlice tests
	assert.Equal(t, ToStringMapStringSlice(stringMapStringSlice), stringMapStringSlice)
	assert.Equal(t, ToStringMapStringSlice(stringMapInterfaceSlice), stringMapStringSlice)
	assert.Equal(t, ToStringMapStringSlice(stringMapStringSliceMultiple), stringMapStringSlice)
	assert.Equal(t, ToStringMapStringSlice(stringMapStringSliceMultiple), stringMapStringSlice)
	assert.Equal(t, ToStringMapStringSlice(stringMapString), stringMapStringSliceSingle)
	assert.Equal(t, ToStringMapStringSlice(stringMapInterface), stringMapStringSliceSingle)
	assert.Equal(t, ToStringMapStringSlice(interfaceMapStringSlice), stringMapStringSlice)
	assert.Equal(t, ToStringMapStringSlice(interfaceMapInterfaceSlice), stringMapStringSlice)
	assert.Equal(t, ToStringMapStringSlice(interfaceMapString), stringMapStringSingleSliceFieldsResult)
	assert.Equal(t, ToStringMapStringSlice(interfaceMapInterface), stringMapStringSingleSliceFieldsResult)
}

func TestSlices(t *testing.T) {

	assert.Equal(t, []string{"a", "b"}, ToStringSlice([]string{"a", "b"}))
	assert.Equal(t, []string{"1", "3"}, ToStringSlice([]interface{}{1, 3}))
	assert.Equal(t, []int{1, 3}, ToIntSlice([]int{1, 3}))
	assert.Equal(t, []int{1, 3}, ToIntSlice([]interface{}{1.2, 3.2}))
	assert.Equal(t, []int{2, 3}, ToIntSlice([]string{"2", "3"}))
	assert.Equal(t, []int{2, 3}, ToIntSlice([2]string{"2", "3"}))
}

func TestToBool(t *testing.T) {

	assert.Equal(t, ToBool(0), false)
	assert.Equal(t, ToBool(int64(0)), false)
	assert.Equal(t, ToBool(nil), false)
	assert.Equal(t, ToBool("false"), false)
	assert.Equal(t, ToBool("FALSE"), false)
	assert.Equal(t, ToBool("False"), false)
	assert.Equal(t, ToBool("f"), false)
	assert.Equal(t, ToBool("F"), false)
	assert.Equal(t, ToBool(false), false)
	assert.Equal(t, ToBool("foo"), false)

	assert.Equal(t, ToBool("true"), true)
	assert.Equal(t, ToBool("TRUE"), true)
	assert.Equal(t, ToBool("True"), true)
	assert.Equal(t, ToBool("t"), true)
	assert.Equal(t, ToBool("T"), true)
	assert.Equal(t, ToBool(1), true)
	assert.Equal(t, ToBool(int64(2)), true)
	assert.Equal(t, ToBool(true), true)
	assert.Equal(t, ToBool(-1), true)
	assert.Equal(t, ToBool(int64(-1)), true)
}

func TestIndirectPointers(t *testing.T) {

	x := 13
	y := &x
	z := &y

	assert.Equal(t, ToInt(y), 13)
	assert.Equal(t, ToInt(z), 13)
}

func TestToDuration(t *testing.T) {

	a := time.Second * 5
	ai := int64(a)
	b := time.Second * 5
	bf := float64(b)

	dai, err := ToDurationE(ai)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, dai, a)

	dbf, err := ToDurationE(bf)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, dbf, b)
}

func getMockTime(format string) time.Time {
	nowS := time.Now().Format(format)
	t, err := time.ParseInLocation(format, nowS, time.Local)
	if err != nil {
		panic(err)
	}
	return t
}

func TestStringToDate(t *testing.T) {

	for i, f := range TimeFormats {
		now := getMockTime(f)
		nowS := now.Format(f)

		haveT, haveErr := StringToDate(nowS, time.Local)
		if haveErr != nil {
			t.Fatal("Index", i, "Error", haveErr)
		}
		assert.Exactly(t, now.Unix(), haveT.Unix(), "Index %d => Format %s", i, f)
	}
}

func TestToTimeE(t *testing.T) {

	now := time.Now()

	fUnix := float64(now.Unix()) + (float64(now.Nanosecond()) / 1e9)

	tests := []struct {
		arg      interface{}
		wantUnix int64
		wantErr  bool
	}{
		{now, now.Unix(), false},
		{'r', 0, true},
		{now.Unix(), now.Unix(), false},
		{fUnix, now.Unix(), false},
	}
	for i, test := range tests {
		haveT, haveErr := ToTimeE(test.arg)
		if test.wantErr {
			assert.True(t, errors.IsNotValid(haveErr), "Index %d => %s", i, haveErr)
			continue
		}
		if haveErr != nil {
			t.Fatal("Index", i, " => ", haveErr)
		}
		assert.Exactly(t, test.wantUnix, haveT.Unix(), "Index %d", i)
	}
}

func TestToTimeSpecific(t *testing.T) {

	const have = "2012-08-23 09:20:13"
	tm, err := ToTimeE(have)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(tm.String())
}
