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

package logw_test

import (
	"bytes"
	std "log"
	"math"
	"testing"

	"github.com/corestoreio/csfw/log"
	"github.com/corestoreio/csfw/log/logw"
	"github.com/corestoreio/csfw/util/errors"
	"github.com/stretchr/testify/assert"
)

var _ log.Logger = (*logw.Log)(nil)

func TestStdLog(t *testing.T) {

	var buf bytes.Buffer

	sl := logw.NewLog(
		logw.WithLevel(logw.LevelDebug),
		logw.WithDebug(&buf, "TEST-DEBUG ", std.LstdFlags),
		logw.WithInfo(&buf, "TEST-INFO ", std.LstdFlags),
		logw.WithFatal(&buf, "TEST-FATAL ", std.LstdFlags),
	)
	sl.SetLevel(logw.LevelInfo)
	assert.False(t, sl.IsDebug())
	assert.True(t, sl.IsInfo())

	sl.Debug("my Debug", log.Float64("float", 3.14152))
	sl.Debug("my Debug2", log.Float64("float2", 2.14152))
	sl.Info("InfoTEST")

	logs := buf.String()

	assert.Contains(t, logs, "InfoTEST")
	assert.NotContains(t, logs, "Debug2")

	buf.Reset()
	sl.SetLevel(logw.LevelDebug)
	assert.True(t, sl.IsDebug())
	assert.True(t, sl.IsInfo())
	sl.Debug("my Debug", log.Float64("float", 3.14152))
	sl.Debug("my Debug2", log.Float64("float2", 2.14152))
	sl.Info("InfoTEST")

	logs = buf.String()

	assert.Contains(t, logs, "InfoTEST")
	assert.Contains(t, logs, "Debug2")

}

func TestStdLogGlobals(t *testing.T) {

	var buf bytes.Buffer
	sl := logw.NewLog(
		logw.WithLevel(logw.LevelDebug),
		logw.WithWriter(&buf),
		logw.WithFlag(std.Ldate),
	)
	sl.Debug("my Debug", log.Float64("float", 3.14152))
	sl.Debug("my Debug2", log.Float64("float2", 2.14152))
	sl.Info("InfoTEST")

	logs := buf.String()

	assert.NotContains(t, logs, "trace2")
	assert.Contains(t, logs, "InfoTEST")
	assert.NotContains(t, logs, "trace1")
	assert.Contains(t, logs, "Debug2")
}

func TestStdLogFormat(t *testing.T) {

	var buf bytes.Buffer
	var bufInfo bytes.Buffer
	sl := logw.NewLog(
		logw.WithLevel(logw.LevelDebug),
		logw.WithWriter(&buf),
		logw.WithInfo(&bufInfo, "TEST-INFO ", std.LstdFlags),
	)

	sl.Debug("my Debug", log.Float64("float1", 3.14152))
	sl.Debug("my Debug2", log.Float64("", 2.14152))
	sl.Debug("my Debug3", log.Int("key3", 3105), log.Int64("Hello", 4711))
	sl.Info("InfoTEST")
	sl.Info("InfoTEST", log.Int("keyI", 117), log.Int64("year", 2009))
	sl.Info("InfoTEST", log.String("", "Now we have the salad"))

	logs := buf.String()
	logsInfo := bufInfo.String()

	assert.Contains(t, logs, "Debug2")
	assert.NotContains(t, logs, "BAD_KEY_AT_INDEX_0")
	assert.NotContains(t, logs, `key3: 3105 BAD_KEY_AT_INDEX_2: "Hello"`)

	assert.Contains(t, logsInfo, "InfoTEST")
	assert.Contains(t, logsInfo, `_: "Now we have the salad`)
}

type myMarshaler struct {
	string
	float64
	bool
	error
}

func (mm myMarshaler) MarshalLog(kv log.KeyValuer) error {
	kv.AddBool("kvbool", mm.bool)
	kv.AddString("kvstring", mm.string)
	kv.AddFloat64("kvfloat64", mm.float64)
	return mm.error
}

func TestAddMarshaler(t *testing.T) {
	var buf bytes.Buffer
	sl := logw.NewLog(
		logw.WithLevel(logw.LevelDebug),
		logw.WithWriter(&buf),
	)

	sl.Debug("my Debug", log.Float64("float1", math.SqrtE))
	sl.Debug("marshalling", log.Object("anObject", 42), log.Marshaler("myMarshaler", myMarshaler{
		string:  "s1",
		float64: math.Ln2,
		bool:    true,
	}))
	assert.Contains(t, buf.String(), `my Debug float1: 1.6487212707001282`)
	assert.Contains(t, buf.String(), `marshalling anObject: 42 kvbool: true kvstring: "s1" kvfloat64: 0.6931471805599453`)
}

func TestAddMarshaler_Error(t *testing.T) {
	var buf bytes.Buffer
	sl := logw.NewLog(
		logw.WithLevel(logw.LevelDebug),
		logw.WithWriter(&buf),
	)

	sl.Debug("my Debug", log.Float64("float1", math.SqrtE))
	sl.Debug("marshalling", log.Marshaler("myMarshaler", myMarshaler{
		error: errors.New("Whooops"),
	}))
	assert.Contains(t, buf.String(), `marshalling kvbool: false kvstring: "" kvfloat64: 0 Error: github.com/corestoreio/csfw/log/logw/stdLib_test.go:158: Whooops`)
}

func TestStdLogNewPanic(t *testing.T) {

	defer func() {
		if r := recover(); r != nil {
			if msg, ok := r.(string); ok {
				assert.EqualValues(t, "Arguments to New() can only be Option types!", msg)
			} else {
				t.Error("Expecting a string")
			}
		}
	}()

	var buf bytes.Buffer
	sl := logw.NewLog(
		logw.WithWriter(&buf),
	)
	sl.New(logw.WithLevel(logw.LevelDebug), 1)
}

func TestStdLogFatal(t *testing.T) {

	defer func() {
		if r := recover(); r != nil {
			assert.Contains(t, r.(string), "This is sparta")
		}
	}()

	var buf bytes.Buffer
	sl := logw.NewLog(
		logw.WithWriter(&buf),
	)
	sl.Fatal("This is sparta")
}
