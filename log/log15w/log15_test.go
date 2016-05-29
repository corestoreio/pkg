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
		l.Debug("log_15_debug", log.Err(errors.New("I'm an debug error")), log.Float64("pi", 3.14159))
	}
	if l.IsInfo() {
		l.Info("log_15_info", log.Err(errors.New("I'm an info error")), log.Float64("e", 2.7182))
	}
	return buf.String()
}

func TestNewLog15_Debug(t *testing.T) {
	out := getLog15(log15.LvlDebug)
	assert.Contains(t, out, `{"Error":"I'm an debug error","Hello":"Gophers","lvl":"dbug","msg":"log_15_debug","pi":3.14159`)
	assert.Contains(t, out, `"pi":3.14159`)
	assert.Contains(t, out, `{"Error":"I'm an info error","Hello":"Gophers","e":2.7182,"lvl":"info","msg":"log_15_info"`)
}

func TestNewLog15_Info(t *testing.T) {
	out := getLog15(log15.LvlInfo)
	assert.NotContains(t, out, `{"Hello":"Gophers","Error":"I'm an debug error","lvl":"dbug"`)
	assert.Contains(t, out, `{"Error":"I'm an info error","Hello":"Gophers","e":2.7182,"lvl":"info",`)
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
	return mm.error
}

func TestAddMarshaler(t *testing.T) {
	buf := &bytes.Buffer{}
	l := log15w.NewLog15(log15.LvlDebug, log15.StreamHandler(buf, log15.JsonFormat()), "Hello", "Gophers")

	l.Debug("log_15_debug", log.Err(errors.New("I'm an debug error")), log.Float64("pi", 3.14159))

	l.Debug("log_15_marshalling", log.Object("anObject", 42), log.Marshaler("marshalLogMock", marshalMock{
		string:  "s1",
		float64: math.Ln2,
		bool:    true,
	}))
	assert.Contains(t, buf.String(), `"anObject":42,"e":2.7182,"kvbool":"true","kvfloat64":0.6931471805599453,"kvstring":"s1",`)
}

func TestAddMarshaler_Error(t *testing.T) {
	buf := &bytes.Buffer{}
	l := log15w.NewLog15(log15.LvlDebug, log15.StreamHandler(buf, log15.JsonFormat()), "Hello", "Gophers")

	l.Debug("marshalling", log.Marshaler("marshalLogMock", marshalMock{
		error: errors.New("Whooops"),
	}))
	assert.Contains(t, buf.String(), `{"Error":"github.com/corestoreio/csfw/log/log15w/log15_test.go:93: Whooops\n","Hello":"Gophers","anObject":42,"e":2.7182,"kvbool":"false","kvfloat64":0,"kvstring":""`)
}
