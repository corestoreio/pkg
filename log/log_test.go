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
	"sync"
	"testing"
	"time"

	"github.com/corestoreio/csfw/log"
	"github.com/corestoreio/csfw/log/logw"
	"github.com/stretchr/testify/assert"
)

var (
	_ log.Logger    = (*log.BlackHole)(nil)
	_ log.KeyValuer = (*log.WriteTypes)(nil)
)

func TestWhenDone(t *testing.T) {
	t.Run("Level_Debug", testWhenDone(logw.LevelDebug))
	t.Run("Level_Info", testWhenDone(logw.LevelInfo))
	t.Run("Level_Fatal", testWhenDone(logw.LevelFatal))
}

func testWhenDone(lvl int) func(*testing.T) {
	return func(t *testing.T) {
		buf := &bytes.Buffer{}
		l := logw.NewLog(logw.WithWriter(buf), logw.WithLevel(lvl))
		var wg sync.WaitGroup
		wg.Add(1)
		go func(wg2 *sync.WaitGroup) {
			defer wg2.Done()
			defer log.WhenDone(l).Debug("WhenDoneDebug", log.Int("key1", 123))
			defer log.WhenDone(l).Info("WhenDoneInfo", log.Int("key2", 321))
			time.Sleep(time.Millisecond * 100)
		}(&wg)
		wg.Wait()

		if lvl == logw.LevelDebug {
			assert.Contains(t, buf.String(), `WhenDoneDebug key1: 123 Duration: 10`)
		} else {
			assert.NotContains(t, buf.String(), `WhenDoneDebug key1: 123 Duration: 10`)
		}
		if lvl >= logw.LevelInfo {
			assert.Contains(t, buf.String(), `WhenDoneInfo key2: 321 Duration: 10`)
		} else {
			assert.NotContains(t, buf.String(), `WhenDoneInfo key2: 321 Duration: 10`)
		}
	}
}
