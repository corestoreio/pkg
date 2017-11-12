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

package directory_test

import (
	std "log"

	"github.com/corestoreio/pkg/directory"
	"github.com/corestoreio/pkg/util/log"
)

var debugLogBuf *log.MutexBuffer
var infoLogBuf *log.MutexBuffer

func init() {
	debugLogBuf = new(log.MutexBuffer)
	infoLogBuf = new(log.MutexBuffer)

	directory.PkgLog = log.NewStdLog(
		log.WithStdDebug(debugLogBuf, "testDebug: ", std.Lshortfile),
		log.WithStdInfo(infoLogBuf, "testInfo: ", std.Lshortfile),
	)
	directory.PkgLog.SetLevel(log.StdLevelDebug)
}
