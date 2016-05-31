// Copyright 2015-2016, Cyrill @ Schumacher.fm and the CoreStore contributors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package log15w_test

import (
	"bytes"
	"testing"

	"math"

	"github.com/corestoreio/csfw/log"
	"github.com/corestoreio/csfw/log/log15w"
	"github.com/corestoreio/csfw/util/errors"
	"github.com/inconshreveable/log15"
	"github.com/stretchr/testify/assert"
)

var _ log.Logger = (*log15w.Log15)(nil)

func getLog15(lvl log15.Lvl) string {
	buf := &bytes.Buffer{}
	l := log15w.NewLog15(lvl, log15.StreamHandler(buf, log15.JsonFormat()), "Hello", "Gophers")

	if l.IsDebug() {
		l.Debug("log_15_debug", log.Err(errors.New("I'm a debug error")), log.Float64("pi", 3.14159))
	}
	if l.IsInfo() {
		l.Info("log_15_info", log.Err(errors.New("I'm an info error")), log.Float64("e", 2.7182))
	}
	return buf.String()
}

func TestNewLog15_Debug(t *testing.T) {
	out := getLog15(log15.LvlDebug)
	assert.Contains(t, out, `"error":"I'm a debug error"`)
	assert.Contains(t, out, `"Hello":"Gophers"`)
	assert.Contains(t, out, `"lvl":"dbug"`)
	assert.Contains(t, out, `"msg":"log_15_debug"`)
	assert.Contains(t, out, `"pi":3.14159`)
	assert.Contains(t, out, `"error":"I'm an info error"`)
	assert.Contains(t, out, `"Hello":"Gophers"`)
	assert.Contains(t, out, `"lvl":"info"`)
	assert.Contains(t, out, `"msg":"log_15_info"`)
}

func TestNewLog15_Info(t *testing.T) {
	out := getLog15(log15.LvlInfo)
	assert.NotContains(t, out, `{"Hello":"Gophers","error":"I'm an debug error","lvl":"dbug"`)
	assert.Contains(t, out, `"error":"I'm an info error"`)
	assert.Contains(t, out, `"Hello":"Gophers"`)
	assert.Contains(t, out, `"lvl":"info"`)
	assert.Contains(t, out, `"e":2.7182`)
}

type marshalMock struct {
	string
	float64
	bool
	error
}

func (mm marshalMock) MarshalLog(kv log.KeyValuer) error {
	kv.AddBool("kvbool", mm.bool)
	kv.AddString("kvstring", mm.string)
	kv.AddFloat64("kvfloat64", mm.float64)
	kv.Nest("startNest", func(kv2 log.KeyValuer) error {
		kv2.AddInt64("nestedInt64", 4711)
		return nil
	})
	return mm.error
}

func TestAddMarshaler(t *testing.T) {
	buf := &bytes.Buffer{}
	l := log15w.NewLog15(log15.LvlDebug, log15.StreamHandler(buf, log15.JsonFormat()), "Hello", "Gophers")

	l.Debug("log_15_debug", log.Err(errors.New("I'm an debug error")), log.Float64("pi", 3.14159))

	l.Debug("log_15_marshalling", log.Object("anObject", 42), log.Marshal("marshalLogMock", marshalMock{
		string:  "s1",
		float64: math.Ln2,
		bool:    true,
	}))
	assert.Contains(t, buf.String(), `"anObject":42`)
	assert.Contains(t, buf.String(), `"e":2.7182`)
	assert.Contains(t, buf.String(), `"kvbool":"true"`)
	assert.Contains(t, buf.String(), `"kvfloat64":0.6931471805599453`)
	assert.Contains(t, buf.String(), `"kvstring":"s1"`)
	assert.Contains(t, buf.String(), `"nestedInt64":4711`)
}

func TestAddMarshaler_Error(t *testing.T) {
	buf := &bytes.Buffer{}
	l := log15w.NewLog15(log15.LvlDebug, log15.StreamHandler(buf, log15.JsonFormat()), "Hello", "Gophers")

	l.Debug("marshalling", log.Marshal("marshalLogMock", marshalMock{
		error: errors.New("Whooops"),
	}))
	assert.Contains(t, buf.String(), `"error":"github.com/corestoreio/csfw/log/log15w/log15_test.go:110: Whooops\n"`)
	assert.Contains(t, buf.String(), `"kvbool":"false","kvfloat64":0,"kvstring":""`)
}
