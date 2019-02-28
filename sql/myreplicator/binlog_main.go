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

package main

import (
	"context"
	"flag"
	"fmt"
	"net"
	"os"
	"sync"
	"time"

	"github.com/corestoreio/errors"
	"github.com/corestoreio/log"
	"github.com/corestoreio/log/logw"
	"github.com/corestoreio/pkg/sql/ddl"
	"github.com/corestoreio/pkg/sql/dml"
	"github.com/corestoreio/pkg/sql/myreplicator"
	"github.com/corestoreio/pkg/util/conv"
	"github.com/go-sql-driver/mysql"
	"github.com/jroimartin/gocui"
)

var _ = context.TODO() // because sometimes code commented out 8-)

// this program creates a terminal UI with two panels. on the left you see the
// debug logging outout and on the right panel you see the incoming binary log
// data. Use Ctrl+C to exit the program. During exit it will panic to print the last log entries from the log buffer.
func main() {
	// export CS_DSN=mysql://root:PASSWORD@localhost:3306/DATABASE_NAME?BinlogSlaveId=100
	var masterStatusFile string
	var masterStatusPosition uint
	flag.StringVar(&masterStatusFile, "master_file", "", "File name from the SHOW MASTER STATUS command")
	flag.UintVar(&masterStatusPosition, "master_pos", 4, "Position from the SHOW MASTER STATUS command")
	flag.Parse()

	logBuf := new(log.MutexBuffer)
	syncLog := logw.NewLog(
		logw.WithWriter(logBuf),
		logw.WithLevel(logw.LevelDebug),
	)
	dsn, err := mysql.ParseDSN(os.Getenv(dml.EnvDSN))
	mustCheckErr(err)

	host, port, err := net.SplitHostPort(dsn.Addr)
	mustCheckErr(err)
	cfg := myreplicator.BinlogSyncerConfig{
		ServerID: 100,
		Flavor:   "mysql",
		Host:     host,
		Port:     uint16(conv.ToInt(port)),
		User:     dsn.User,
		Password: dsn.Passwd,
		Log:      syncLog,
	}
	binSync := myreplicator.NewBinlogSyncer(&cfg)
	// mysql.Position change to whatever SHOW MASTER STATUS tells you
	streamer, err := binSync.StartSync(ddl.MasterStatus{File: masterStatusFile, Position: masterStatusPosition})
	mustCheckErr(err)

	ctx, cancel := context.WithCancel(context.Background())

	g, err := gocui.NewGui(gocui.OutputNormal)
	mustCheckErr(err)
	defer g.Close()
	g.SetManagerFunc(func(g *gocui.Gui) error {
		maxX, maxY := g.Size()
		if _, err := g.SetView("left", 0, 0, maxX/2-1, maxY-1); err != nil && err != gocui.ErrUnknownView {
			return err
		}
		if _, err := g.SetView("right", maxX/2, 0, maxX-1, maxY-1); err != nil && err != gocui.ErrUnknownView {
			return err
		}
		return nil
	})

	mustCheckErr(g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
		if err := binSync.Close(); err != nil {
			return err
		}
		cancel()
		return gocui.ErrQuit
	}))

	var wg sync.WaitGroup
	go func() {
		defer wg.Done()
		wg.Add(1)
		for {
			select {
			case <-ctx.Done():
				return
			case <-time.After(750 * time.Millisecond):
				g.Update(func(g *gocui.Gui) error {
					vl, err := g.View("left")
					if err != nil {
						return err
					}
					vl.Wrap = true
					vl.Autoscroll = true
					// vl.Clear()
					_, err = logBuf.WriteTo(vl)
					return err
				})
			}
		}
	}()

	go func() {
		defer wg.Done()
		wg.Add(1)

		for {
			ev, err := streamer.GetEvent(ctx)
			if err != nil {
				if errors.Cause(err) != context.Canceled {
					panic("binlog_main_streamer.GetEvent.error " + err.Error())
				}
				return
			}
			g.Update(func(g *gocui.Gui) error {
				vr, err := g.View("right")
				if err != nil {
					return err
				}
				vr.Wrap = true
				vr.Autoscroll = true
				// vl.Clear()
				ev.Dump(vr)
				return err
			})
		}
	}()

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		panic(fmt.Sprintf("%+v", err))
	}

	wg.Wait()
	// we have to panic here to show the following output, when using println
	// gocui would just clear the screen and the output is gone.
	panic("All goroutines terminated.\nDEBUG BUFFER\n" + logBuf.String())
}

func mustCheckErr(err error) {
	if err != nil {
		panic(fmt.Sprintf("%+v", err))
	}
}
