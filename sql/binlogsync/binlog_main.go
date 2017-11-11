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

// +build ignore

// http://dev.mysql.com/doc/refman/5.7/en/replication-options-binary-log.html

// This tool helps you debugging and understanding the binary log.

package binlogsync

import (
	"bytes"
	"context"
	"fmt"
	"net"

	"github.com/corestoreio/cspkg/sql/ddl"
	"github.com/corestoreio/cspkg/sql/myreplicator"
	"github.com/corestoreio/cspkg/util/conv"
)

func main() {
	// export CS_DSN=mysql://root:PASSWORD@localhost:3306/DATABASE_NAME?BinlogSlaveId=100
	dsn, err := ddl.GetParsedDSN()
	if err != nil {
		panic(fmt.Sprintf("%+v", err))
	}

	host, port, _ := net.SplitHostPort(dsn.Addr)
	cfg := myreplicator.BinlogSyncerConfig{
		ServerID: 100,
		Flavor:   "mysql",
		Host:     host,
		Port:     uint16(conv.ToInt(port)),
		User:     dsn.User,
		Password: dsn.Passwd,
	}
	syncer := myreplicator.NewBinlogSyncer(&cfg)
	defer syncer.Close()

	// mysql.Position change to whatever SHOW MASTER STATUS tells you
	streamer, err := syncer.StartSync(ddl.MasterStatus{File: "mysql-bin.000001", Position: 4})
	if err != nil {
		panic(fmt.Sprintf("%+v", err))
	}

	buf := bytes.Buffer{}
	for {
		ev, err := streamer.GetEvent(context.Background())
		if err != nil {
			panic(fmt.Sprintf("%+v", err))
		}

		ev.Dump(&buf)
		println(buf.String(), "\n")
		buf.Reset()
	}
}
