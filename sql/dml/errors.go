// Copyright 2015-present, Cyrill @ Schumacher.fm and the CoreStore contributors
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

package dml

import (
	"github.com/corestoreio/errors"
	"github.com/go-sql-driver/mysql"
)

func MySQLNumberFromError(err error) uint16 {
	mErr := errors.Cause(err)
	if myErr, ok := mErr.(*mysql.MySQLError); ok {
		return myErr.Number
	}
	return 0
}

func MySQLMessageFromError(err error) string {
	mErr := errors.Cause(err)
	if myErr, ok := mErr.(*mysql.MySQLError); ok {
		return myErr.Message
	}
	return ""
}
