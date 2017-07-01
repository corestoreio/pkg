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

package csdb

import (
	"strconv"
	"strings"

	"github.com/corestoreio/errors"
)

// MasterStatus This statement provides status information about the binary log
// files of the master. It requires either the SUPER or REPLICATION CLIENT
// privilege.
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

// ToSQL implements dbr.QueryBuilder interface to assemble a SQL string and its
// arguments for query execution.
func (ms *MasterStatus) ToSQL() (string, []interface{}, error) {
	return "SHOW MASTER STATUS", nil, nil
}

// RowScan implements dbr.Scanner interface to scan a row returned from database
// query.
func (ms *MasterStatus) RowScan(_ int64, _ []string, scan func(...interface{}) error) error {
	return errors.WithStack(
		scan(&ms.File, &ms.Position, &ms.BinlogDoDB, &ms.BinlogIgnoreDB, &ms.ExecutedGTIDSet),
	)
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
		return errors.NewNotFoundf("[csdb] MasterStatus FromString: Delimiter semi-colon not found.")
	}

	pos, err := strconv.ParseUint(str[c+1:], 10, 32)
	if err != nil {
		return errors.NewNotValidf("[csdb] MasterStatus FromString: %s", err)
	}
	ms.File = str[:c]
	ms.Position = uint(pos)
	return nil
}
