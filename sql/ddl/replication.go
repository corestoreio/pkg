// Copyright 2015-2017, Cyrill @ Schumacher.fm and the CoreStore contributors
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

package ddl

import (
	"strconv"
	"strings"

	"github.com/corestoreio/pkg/sql/dml"
	"github.com/corestoreio/errors"
)

// MasterStatus provides status information about the binary log files of the
// master. It requires either the SUPER or REPLICATION CLIENT privilege. Once a
// MasterStatus pointer variable has been created it can be reused multiple
// times.
type MasterStatus struct {
	File           string
	Position       uint
	BinlogDoDB     string
	BinlogIgnoreDB string
	// ExecutedGTIDSet: When global transaction IDs are in use, ExecutedGTIDSet
	// shows the set of GTIDs for transactions that have been executed on the
	// master. This is the same as the value for the gtid_executed system variable
	// on this server, as well as the value for ExecutedGTIDSet in the output of
	// SHOW SLAVE STATUS on this server.
	ExecutedGTIDSet string
}

// ToSQL implements dml.QueryBuilder interface to assemble a SQL string and its
// arguments for query execution.
func (ms *MasterStatus) ToSQL() (string, []interface{}, error) {
	return "SHOW MASTER STATUS", nil, nil
}

// MapColumns implements dml.ColumnMapper interface to scan a row returned from
// a database query.
func (ms *MasterStatus) MapColumns(rc *dml.ColumnMap) error {
	for rc.Next() {
		switch col := rc.Column(); col {
		case "File":
			rc.String(&ms.File)
		case "Position":
			rc.Uint(&ms.Position)
		case "Binlog_Do_DB":
			rc.String(&ms.BinlogDoDB)
		case "Binlog_Ignore_DB":
			rc.String(&ms.BinlogIgnoreDB)
		case "Executed_Gtid_Set":
			rc.String(&ms.ExecutedGTIDSet)
		default:
			return errors.NewNotFoundf("[ddl] Column %q not found in SHOW MASTER STATUS", col)
		}
	}
	return errors.WithStack(rc.Err())
}

// Compare compares with another MasterStatus. Returns 1 if left hand side is
// bigger, 0 if both are equal and -1 if right hand side is bigger.
func (ms MasterStatus) Compare(other MasterStatus) int {
	switch {
	// First compare binlog name
	case ms.File > other.File:
		return 1
	case ms.File < other.File:
		return -1
	// Same binlog file, compare position
	case ms.Position > other.Position:
		return 1
	case ms.Position < other.Position:
		return -1
	}
	return 0
}

// String converts the file name and the position to a string, separated by a
// semi-colon.
func (ms MasterStatus) String() string {
	if ms.File == "" {
		return ""
	}
	return ms.File + ";" + strconv.FormatUint(uint64(ms.Position), 10)
}

// FromString parses as string in the format: mysql-bin.000002;236423 means
// filename;position.
func (ms *MasterStatus) FromString(str string) error {
	c := strings.IndexByte(str, ';')
	if c < 1 {
		return errors.NewNotFoundf("[ddl] MasterStatus FromString: Delimiter semi-colon not found.")
	}

	pos, err := strconv.ParseUint(str[c+1:], 10, 32)
	if err != nil {
		return errors.NewNotValidf("[ddl] MasterStatus FromString: %s", err)
	}
	ms.File = str[:c]
	ms.Position = uint(pos)
	return nil
}
