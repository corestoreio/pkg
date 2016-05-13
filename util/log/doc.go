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

/*
Package log contains BlackHole, StdLog, Log15 and the Logger interface.

Logging

Interface Logger is used all over the place and there are no other dependencies.
Default Logger is a null logger. You must take care to implement a logger which
is also thread safe.

Convention: Because recording a human-meaningful message is common and good
practice, the first argument to every logging method is the value to the
*implicit* key 'msg'. You may supply any additional context as a set of
key/value pairs to the logging function.

Level guards exists to avoid the cost of building arguments. Get in the
habit of using guards.

	import "github.com/corestoreio/csfw/util/log"

	if log.IsDebug() {
		log.Debug("some message", "key1", expensive())
	}

Standardizes on key-value pair argument sequence:

	import "github.com/corestoreio/csfw/util/log"

	log.Debug("message from inside Fn()", "key1", value1, "key2", value2)

	// instead of this
	log.WithFields(logrus.Fields{"m": "pkg", "key1": value1, "key2": value2}).Debug("inside fn()")

Please consider the key-value pairs when implementing your own logger.

Recommended Loggers are https://github.com/mgutz/logxi and https://github.com/Sirupsen/logrus
and https://github.com/inconshreveable/log15

Standard Logger

CoreStore provides a leveled logger based on Go's standard library without any
dependencies. This StdLog obeys to the interface Logger of this package.

	import "github.com/corestoreio/csfw/util/log"

	func init() {
		log.Set(log.NewStdLog())
	}

log.NewStdLog() accepts a wide range of optional arguments. Please see the functions Std*Option().

Additional Reading

http://dave.cheney.net/2015/11/05/lets-talk-about-logging

https://www.reddit.com/r/golang/comments/3rljir/lets_talk_about_logging/

https://forum.golangbridge.org/t/whats-so-bad-about-the-stdlibs-log-package/1435/2

TODO(cs): http://12factor.net/logs
*/
package log
