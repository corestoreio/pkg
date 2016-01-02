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
	"testing"

	"github.com/corestoreio/csfw/util/log"
)

func TestNull(t *testing.T) {
	log.SetLevel(-1000)
	if !log.IsDebug() {
		t.Error("There should be debug logging")
	}
	if !log.IsInfo() {
		t.Error("There should be info logging")
	}
	var args []interface{}
	args = append(args, "key1", 1, "key2", 3.14152)

	log.Debug("Hello World", args...)
	log.Info("Hello World", args...)
}
