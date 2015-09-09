// Copyright 2015, Cyrill @ Schumacher.fm and the CoreStore contributors
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

package log_test

import (
	"bytes"
	"errors"
	std "log"
	"testing"

	"github.com/corestoreio/csfw/utils/log"
	"github.com/stretchr/testify/assert"
)

func TestStdLogger(t *testing.T) {

	var buf bytes.Buffer

	sl := log.NewStdLogger(
		log.SetStdLevel(log.StdLevelTrace),
		log.SetStdTrace(&buf, "TEST-TRACE ", std.LstdFlags),
		log.SetStdDebug(&buf, "TEST-DEBUG ", std.LstdFlags),
		log.SetStdInfo(&buf, "TEST-INFO ", std.LstdFlags),
		log.SetStdWarn(&buf, "TEST-WARN ", std.LstdFlags),
		log.SetStdError(&buf, "TEST-ERROR ", std.LstdFlags),
		log.SetStdFatal(&buf, "TEST-FATAL ", std.LstdFlags),
	)
	sl.SetLevel(log.StdLevelInfo)
	assert.False(t, sl.IsTrace())
	assert.False(t, sl.IsDebug())
	assert.True(t, sl.IsInfo())
	assert.True(t, sl.IsWarn())

	sl.Trace("my trace1")
	sl.Trace("my trace2", "int", 29)
	sl.Debug("my Debug", "float", 3.14152)
	sl.Debug("my Debug2", 2.14152)
	sl.Info("InfoTEST")
	sl.Warn("WarnTEST")
	haveErr := sl.Error("ErrorTEST", "err1a", 1, "err2", 32.4232)
	assert.Contains(t, "ErrorTEST53", sl.Error("ErrorTEST53").Error())
	logs := buf.String()

	assert.EqualError(t, haveErr, "ErrorTEST")
	assert.Contains(t, logs, "InfoTEST")
	assert.Contains(t, logs, "WarnTEST")
	assert.Contains(t, logs, "ErrorTEST")
	assert.NotContains(t, logs, "trace1")
	assert.NotContains(t, logs, "Debug2")

	buf.Reset()
	sl.SetLevel(log.StdLevelTrace)
	assert.True(t, sl.IsTrace())
	assert.True(t, sl.IsDebug())
	assert.True(t, sl.IsInfo())
	assert.True(t, sl.IsWarn())
	sl.Trace("my trace1")
	sl.Trace("my trace2", "int", 29)
	sl.Debug("my Debug", "float", 3.14152)
	sl.Debug("my Debug2", 2.14152)
	sl.Info("InfoTEST")

	logs = buf.String()

	assert.Contains(t, logs, "InfoTEST")
	assert.Contains(t, logs, "trace1")
	assert.Contains(t, logs, "Debug2")

}

func TestStdLoggerGlobals(t *testing.T) {

	var buf bytes.Buffer
	sl := log.NewStdLogger(
		log.SetStdLevel(log.StdLevelDebug),
		log.SetStdWriter(&buf),
		log.SetStdFlag(std.Ldate),
	)
	sl.Trace("my trace1")
	sl.Trace("my trace2", "int", 29)
	sl.Debug("my Debug", "float", 3.14152)
	sl.Debug("my Debug2", 2.14152)
	sl.Info("InfoTEST")

	logs := buf.String()

	assert.NotContains(t, logs, "trace2")
	assert.Contains(t, logs, "InfoTEST")
	assert.NotContains(t, logs, "trace1")
	assert.Contains(t, logs, "Debug2")
}

func TestStdLoggerFormat(t *testing.T) {

	var buf bytes.Buffer
	var bufInfo bytes.Buffer
	sl := log.NewStdLogger(
		log.SetStdLevel(log.StdLevelDebug),
		log.SetStdWriter(&buf),
		log.SetStdInfo(&bufInfo, "TEST-INFO ", std.LstdFlags),
	)

	sl.Debug("my Debug", 3.14152)
	sl.Debug("my Debug2", "", 2.14152)
	sl.Debug("my Debug3", "key3", 3105, 4711, "Hello")
	sl.Info("InfoTEST")
	sl.Info("InfoTEST", "keyI", 117, 2009)

	aTestErr := errors.New("Cannot run PHP code")
	haveErr := sl.Error("ErrorTEST", "myErr", aTestErr)

	logs := buf.String()
	logsInfo := bufInfo.String()

	//	t.Log("", logs)
	//	t.Log("", logsInfo)

	assert.EqualError(t, haveErr, aTestErr.Error())
	assert.Contains(t, logs, "Debug2")
	assert.Contains(t, logs, "BAD_KEY_AT_INDEX_0")
	assert.Contains(t, logs, `key3: 3105 BAD_KEY_AT_INDEX_2: "Hello"`)
	assert.Contains(t, logs, "_: 3.14")

	assert.Contains(t, logsInfo, "InfoTEST")
	assert.Contains(t, logsInfo, `FIX_IMBALANCED_PAIRS: []interface {}{"keyI", 117, 2009}`)
}

func TestStdLoggerNewPanic(t *testing.T) {

	defer func() {
		if r := recover(); r != nil {
			if msg, ok := r.(string); ok {
				assert.EqualValues(t, "Arguments to New() can only be StdOption types!", msg)
			} else {
				t.Error("Expecting a string")
			}
		}
	}()

	var buf bytes.Buffer
	sl := log.NewStdLogger(
		log.SetStdWriter(&buf),
	)
	sl.New(log.SetStdLevel(log.StdLevelDebug), 1)
}
