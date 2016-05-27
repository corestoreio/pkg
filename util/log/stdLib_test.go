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

package log_test

import (
	"bytes"
	std "log"
	"testing"

	"github.com/corestoreio/csfw/util/log"
	"github.com/stretchr/testify/assert"
)

func TestStdLog(t *testing.T) {

	var buf bytes.Buffer

	sl := log.NewStdLog(
		log.WithStdLevel(log.StdLevelDebug),
		log.WithStdDebug(&buf, "TEST-DEBUG ", std.LstdFlags),
		log.WithStdInfo(&buf, "TEST-INFO ", std.LstdFlags),
		log.WithStdFatal(&buf, "TEST-FATAL ", std.LstdFlags),
	)
	sl.SetLevel(log.StdLevelInfo)
	assert.False(t, sl.IsDebug())
	assert.True(t, sl.IsInfo())

	sl.Debug("my Debug", log.Float64("float", 3.14152))
	sl.Debug("my Debug2", log.Float64("float2", 2.14152))
	sl.Info("InfoTEST")

	logs := buf.String()

	assert.Contains(t, logs, "InfoTEST")
	assert.NotContains(t, logs, "Debug2")

	buf.Reset()
	sl.SetLevel(log.StdLevelDebug)
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
	sl := log.NewStdLog(
		log.WithStdLevel(log.StdLevelDebug),
		log.WithStdWriter(&buf),
		log.WithStdFlag(std.Ldate),
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
	sl := log.NewStdLog(
		log.WithStdLevel(log.StdLevelDebug),
		log.WithStdWriter(&buf),
		log.WithStdInfo(&bufInfo, "TEST-INFO ", std.LstdFlags),
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

func TestStdLogNewPanic(t *testing.T) {

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
	sl := log.NewStdLog(
		log.WithStdWriter(&buf),
	)
	sl.New(log.WithStdLevel(log.StdLevelDebug), 1)
}

func TestStdLogFatal(t *testing.T) {

	defer func() {
		if r := recover(); r != nil {
			assert.Contains(t, r.(string), "This is sparta")
		}
	}()

	var buf bytes.Buffer
	sl := log.NewStdLog(
		log.WithStdWriter(&buf),
	)
	sl.Fatal("This is sparta")
}
