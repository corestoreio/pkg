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
	std "log"
	"testing"

	"github.com/corestoreio/csfw/utils/log"
	"github.com/stretchr/testify/assert"
)

func TestStdLogger(t *testing.T) {

	var buf bytes.Buffer

	sl := log.NewStdLogger(
		log.SetStdLevel(log.StdLevelDebug),
		log.SetStdDebug(&buf, "TEST-DEBUG ", std.LstdFlags),
		log.SetStdInfo(&buf, "TEST-INFO ", std.LstdFlags),
		log.SetStdFatal(&buf, "TEST-FATAL ", std.LstdFlags),
	)
	sl.SetLevel(log.StdLevelInfo)
	assert.False(t, sl.IsDebug())
	assert.True(t, sl.IsInfo())

	sl.Debug("my Debug", "float", 3.14152)
	sl.Debug("my Debug2", 2.14152)
	sl.Info("InfoTEST")

	logs := buf.String()

	assert.Contains(t, logs, "InfoTEST")
	assert.NotContains(t, logs, "Debug2")

	buf.Reset()
	sl.SetLevel(log.StdLevelDebug)
	assert.True(t, sl.IsDebug())
	assert.True(t, sl.IsInfo())
	sl.Debug("my Debug", "float", 3.14152)
	sl.Debug("my Debug2", 2.14152)
	sl.Info("InfoTEST")

	logs = buf.String()

	assert.Contains(t, logs, "InfoTEST")
	assert.Contains(t, logs, "Debug2")

}

func TestStdLoggerGlobals(t *testing.T) {

	var buf bytes.Buffer
	sl := log.NewStdLogger(
		log.SetStdLevel(log.StdLevelDebug),
		log.SetStdWriter(&buf),
		log.SetStdFlag(std.Ldate),
	)
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
	sl.Info("InfoTEST", "Now we have the salad")

	logs := buf.String()
	logsInfo := bufInfo.String()

	assert.Contains(t, logs, "Debug2")
	assert.Contains(t, logs, "BAD_KEY_AT_INDEX_0")
	assert.Contains(t, logs, `key3: 3105 BAD_KEY_AT_INDEX_2: "Hello"`)

	assert.Contains(t, logsInfo, "InfoTEST")
	assert.Contains(t, logsInfo, `_: "Now we have the salad`)
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

func TestStdLoggerFatal(t *testing.T) {

	defer func() {
		if r := recover(); r != nil {
			assert.Contains(t, r.(string), "This is sparta")
		}
	}()

	var buf bytes.Buffer
	sl := log.NewStdLogger(
		log.SetStdWriter(&buf),
	)
	sl.Fatal("This is sparta")
}
