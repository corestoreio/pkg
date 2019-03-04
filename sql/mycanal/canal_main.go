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

// +build ignore

// This tool helps you debugging and understanding the binary log.

package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/corestoreio/log"
	"github.com/corestoreio/log/logw"
	"github.com/corestoreio/pkg/sql/ddl"
	"github.com/corestoreio/pkg/sql/dml"
	"github.com/corestoreio/pkg/sql/mycanal"
	"github.com/corestoreio/pkg/util/conv"
	"github.com/jroimartin/gocui"
)

var _ mycanal.RowsEventHandler = (*genericEventDump)(nil)

type genericEventDump struct {
	g *gocui.Gui
}

func (ge *genericEventDump) Do(_ context.Context, action string, table *ddl.Table, rows [][]interface{}) error {
	ge.g.Update(func(g *gocui.Gui) error {
		vr, err := g.View("right")
		if err != nil {
			return err
		}
		_, err = fmt.Fprintf(vr, "%q => %q\n", action, table.Name)

		cols := table.Columns
		maxLen := 0
		for _, c := range cols {
			if l := len(c.Field); l > maxLen {
				maxLen = l
			}
		}

		for _, row := range rows {
			for ci, cell := range row {
				_, err = fmt.Fprintf(vr, "%-"+strconv.FormatInt(int64(maxLen), 10)+"s: %T %q\n", cols[ci].Field, cell, conv.ToString(cell))
			}
			fmt.Fprint(vr, "\n")
		}

		fmt.Fprint(vr, "=======================================\n")
		return err
	})

	return nil
}
func (ge *genericEventDump) Complete(_ context.Context) error {
	return nil // errors.NewFatalf("[test] What is incomplete?")
}

func (ge *genericEventDump) String() string {
	return "genericEventDump"
}

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

	c, err := mycanal.NewCanal(os.Getenv(dml.EnvDSN), mycanal.WithMySQL(), &mycanal.Options{
		Log: syncLog,
		OnClose: func(db *dml.ConnPool) (err error) {
			return nil
		},
	})
	mustCheckErr(err)

	ctx, cancel := context.WithCancel(context.Background())

	g, err := gocui.NewGui(gocui.OutputNormal)
	mustCheckErr(err)
	defer g.Close()

	c.RegisterRowsEventHandler(nil, &genericEventDump{g: g})

	g.SetManagerFunc(func(g *gocui.Gui) error {
		maxX, maxY := g.Size()
		if v, err := g.SetView("left", 0, 0, maxX/2-1-10, maxY-1); err != nil && err != gocui.ErrUnknownView {
			return err
		} else {
			v.Wrap = true
			v.Autoscroll = true
		}
		if v, err := g.SetView("right", maxX/2-10, 0, maxX-1, maxY-1); err != nil && err != gocui.ErrUnknownView {
			return err
		} else {
			v.Wrap = true
			v.Autoscroll = true
		}
		g.SetCurrentView("right")
		return nil
	})

	mustCheckErr(g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
		cancel()
		if err := c.Close(); err != nil {
			return err
		}
		return gocui.ErrQuit
	}))
	mustCheckErr(g.SetKeybinding("right", gocui.KeyCtrlW, gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
		v.Clear()
		return nil
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
					_, err = logBuf.WriteTo(vl)
					return err
				})
			}
		}
	}()

	go func() {
		defer wg.Done()
		wg.Add(1)

		if err := c.Start(ctx); err != nil {
			panic("binlog_main_streamer.GetEvent.error " + err.Error())
		}
		// todo terminate this goroutine
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
