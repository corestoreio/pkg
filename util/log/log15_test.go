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
	"testing"

	"github.com/corestoreio/csfw/util/errors"
	"github.com/corestoreio/csfw/util/log"
	"github.com/inconshreveable/log15"
	"github.com/stretchr/testify/assert"
)

var _ log.Logger = (*log.Log15)(nil)

func getLog15(lvl log15.Lvl) string {
	buf := &bytes.Buffer{}
	l := log.NewLog15(lvl, log15.StreamHandler(buf, log15.JsonFormat()), "Hello", "Gophers")

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
	assert.Contains(t, out, `{"Hello":"Gophers","error":"I'm an debug error","lvl":"dbug"`)
	assert.Contains(t, out, `"pi":3.14159`)
	assert.Contains(t, out, `"error":"I'm an info error","lvl":"info"`)
}

func TestNewLog15_Info(t *testing.T) {
	out := getLog15(log15.LvlInfo)
	assert.NotContains(t, out, `{"Hello":"Gophers","error":"I'm an debug error","lvl":"dbug"`)
	assert.Contains(t, out, `"error":"I'm an info error","lvl":"info"`)
	assert.Contains(t, out, `"e":2.7182`)
}
