// Copyright 2015 CoreStore Authors
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
	"testing"

	"github.com/corestoreio/csfw/utils/log"
	"github.com/stretchr/testify/assert"
)

func TestSet(t *testing.T) {
	nn1 := &log.NullLogger{}
	log.Set(nn1)
	defer func() {
		if r := recover(); r != nil {
			if es, ok := r.(error); ok {
				assert.EqualError(t, log.ErrLoggerSet, es.Error())
			} else {
				t.Error("Expected a cast to error interface")
			}
		}
	}()
	nn2 := &log.NullLogger{}
	log.Set(nn2)
}

func TestNull(t *testing.T) {
	log.SetNull()
	log.SetLevel(-1000)
	if log.IsTrace() {
		t.Error("There should be no trace")
	}
	if log.IsDebug() {
		t.Error("There should be no debug")
	}
	if log.IsInfo() {
		t.Error("There should be no info")
	}
	if log.IsWarn() {
		t.Error("There should be no warn")
	}
	var args []interface{}
	args = append(args, "key1", 1, "key2", 3.14152)

	log.Trace("Hello World", args...)
	log.Debug("Hello World", args...)
	log.Info("Hello World", args...)
	log.Warn("Hello World", args...)
	log.Error("Hello World", args...)
	log.Log(1, "Hello World", args)
}
