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
	"context"

	"github.com/corestoreio/csfw/util/errors"
)

// MasterStatus: This statement provides status information about the binary log
// files of the master. It requires either the SUPER or REPLICATION CLIENT
// privilege.
type MasterStatus struct {
	File             string
	Position         uint
	Binlog_Do_DB     string
	Binlog_Ignore_DB string
	// Executed_Gtid_Set: When global transaction IDs are in use, Executed_Gtid_Set
	// shows the set of GTIDs for transactions that have been executed on the
	// master. This is the same as the value for the gtid_executed system variable
	// on this server, as well as the value for Executed_Gtid_Set in the output of
	// SHOW SLAVE STATUS on this server.
	Executed_Gtid_Set string
}

// ShowMasterStatus retrieves the current master status.
func ShowMasterStatus(ctx context.Context, db QueryRower) (MasterStatus, error) {
	var ms MasterStatus
	row := db.QueryRowContext(ctx, "SHOW MASTER STATUS")
	if err := row.Scan(&ms.File, &ms.Position, &ms.Binlog_Do_DB, &ms.Binlog_Ignore_DB, &ms.Executed_Gtid_Set); err != nil {
		return ms, errors.Wrap(err, "[csdb] ShowMasterStatus")
	}
	return ms, nil
}
