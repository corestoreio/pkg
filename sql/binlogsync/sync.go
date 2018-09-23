package binlogsync

import (
	"context"
	"regexp"
	"time"

	"github.com/corestoreio/errors"
	"github.com/corestoreio/log"
	"github.com/corestoreio/pkg/sql/ddl"
	"github.com/corestoreio/pkg/sql/myreplicator"
)

// Action constants to figure out the type of an event. Those constants will be
// passed to the interface RowsEventHandler.
const (
	UpdateAction = "update"
	InsertAction = "insert"
	DeleteAction = "delete"
)

var (
	expCreateTable = regexp.MustCompile("(?i)^CREATE\\sTABLE(\\sIF\\sNOT\\sEXISTS)?\\s`{0,1}(.*?)`{0,1}\\.{0,1}`{0,1}([^`\\.]+?)`{0,1}\\s.*")
	expAlterTable  = regexp.MustCompile("(?i)^ALTER\\sTABLE\\s.*?`{0,1}(.*?)`{0,1}\\.{0,1}`{0,1}([^`\\.]+?)`{0,1}\\s.*")
	expRenameTable = regexp.MustCompile("(?i)^RENAME\\sTABLE\\s.*?`{0,1}(.*?)`{0,1}\\.{0,1}`{0,1}([^`\\.]+?)`{0,1}\\s{1,}TO\\s.*?")
	expDropTable   = regexp.MustCompile("(?i)^DROP\\sTABLE(\\sIF\\sEXISTS){0,1}\\s`{0,1}(.*?)`{0,1}\\.{0,1}`{0,1}([^`\\.]+?)`{0,1}(?:$|\\s)")
	ddlExpressions = [...]*regexp.Regexp{expCreateTable, expAlterTable, expRenameTable, expDropTable}
)

func extractTableFromQueryEvent(schema, query []byte) (dbName, tableName string) {
	var mb [][]byte

	for _, reg := range ddlExpressions {
		mb = reg.FindSubmatch(query)
		if len(mb) != 0 {
			break
		}
	}
	mbLen := len(mb)
	if mbLen == 0 {
		return
	}

	// the first last is table name, the second last is database name(if exists)
	if len(mb[mbLen-2]) == 0 {
		dbName = string(schema)
	} else {
		dbName = string(mb[mbLen-2])
	}
	tableName = string(mb[mbLen-1])
	return
}

func (c *Canal) clearTableCacheOnDDLStmt(schema, query []byte) {
	if db, tbl := extractTableFromQueryEvent(schema, query); tbl != "" {
		c.ClearTableCache(db, tbl)
		if c.opts.Log.IsInfo() {
			c.opts.Log.Info("[binlogsync] Table structure changed, clear table cache",
				log.String("database", db), log.String("table", tbl), log.String("query", string(query)))
		}
	}
}

func (c *Canal) startSyncBinlog(ctxArg context.Context) error {
	if c.syncer == nil {
		return errors.AlreadyClosed.Newf("[binlogsync] Canal already closed and myreplicator.BinlogSyncer is nil")
	}
	pos := c.masterStatus

	if c.opts.Log.IsInfo() {
		c.opts.Log.Info("[binlogsync] Start syncing of binlog", log.Stringer("position", pos))
	}

	s, err := c.syncer.StartSync(pos)
	if err != nil {
		return errors.Fatal.Newf("[binlogsync] Start sync replication at %s error %v", pos, err)
	}

	timeout := time.Second
	for {
		ctx, cancel := context.WithTimeout(ctxArg, 2*time.Second)
		ev, err := s.GetEvent(ctx)
		cancel()

		if err == context.DeadlineExceeded {
			timeout = 2 * timeout
			continue
		}
		if err != nil {
			return errors.WithStack(err)
		}

		timeout = time.Second

		//next binlog pos
		pos.Position = uint(ev.Header.LogPos)

		switch e := ev.Event.(type) {
		case *myreplicator.RotateEvent:
			if err := c.flushEventHandlers(ctxArg); err != nil {
				// todo maybe better err handling ...
				return errors.WithStack(err)
			}
			pos.File = string(e.NextLogName)
			pos.Position = uint(e.Position)
			// r.ev <- pos

			if c.opts.Log.IsInfo() {
				c.opts.Log.Info("[binlogsync] Rotate binlog to a new position", log.Stringer("position", pos))
			}

			// call event handler OnRotate(e)

		case *myreplicator.RowsEvent:
			// we only focus row based event.
			// NotFound errors get ignores. For example table has been deleted
			// and an old event pops in.
			if err = c.handleRowsEvent(ctxArg, ev); err != nil {
				isNotFound := errors.Is(err, errors.NotFound)
				if c.opts.Log.IsInfo() {
					c.opts.Log.Info("[binlogsync] Rotate binlog to a new position", log.Err(err), log.Stringer("position", pos), log.Bool("ignore_not_found_error", isNotFound))
				}
				if !isNotFound {
					return errors.WithStack(err)
				}
				continue // to not save the master position, not necessary.
			}
		case *myreplicator.XIDEvent:
			// TODO implement
			// if e.GSet != nil {
			// 	c.master.UpdateGTIDSet(e.GSet)
			// }

			// call event OnXID(pos)

		case *myreplicator.MariadbGTIDEvent:
			// call event OnGTID(gtid)

		case *myreplicator.GTIDEvent:
			// call event OnGTID(gtid)

		case *myreplicator.QueryEvent:
			// TODO implement GTID set but review if it makes sense to import siddontang/go-mysql/mysql
			// if e.GSet != nil {
			// 	c.master.UpdateGTIDSet(e.GSet)
			// }

			// handle alert table query
			c.clearTableCacheOnDDLStmt(e.Schema, e.Query)

			// For now really necessary.
			// TODO: call two more event handlers: on OnTableChanged(db,table) and OnDDL(pos, e)

			// save master position, so no continue
		case
			*myreplicator.TableMapEvent,
			*myreplicator.FormatDescriptionEvent:
			// maybe add: *replication.XIDEvent
			// don't update Master with file and position
		default:
			continue
		}

		if err := c.masterSave(pos.File, pos.Position); err != nil {
			c.opts.Log.Info("[binlogsync] startSyncBinlog: Failed to save master position", log.Err(err), log.Stringer("position", pos))
		}
	}
}

// handleRowsEvent handles an event on the rows and calls all registered rows
// event handler. can return different error behaviours.
func (c *Canal) handleRowsEvent(ctx context.Context, e *myreplicator.BinlogEvent) error {
	ev, ok := e.Event.(*myreplicator.RowsEvent)
	if !ok {
		return errors.Fatal.Newf("[binlogsync] handleRowsEvent: Failed to cast to *myreplicator.RowsEvent type")
	}

	// Caveat: table may be altered at runtime.
	schemaName := string(ev.Table.Schema)
	if c.dsn.DBName != schemaName {
		if c.opts.Log.IsDebug() {
			c.opts.Log.Debug("[binlogsync.Canal.handleRowsEvent.Skipping.database", log.String("database_have", schemaName), log.String("database_want", c.dsn.DBName), log.Int("table_id", int(ev.TableID)))
		}
		return nil
	}

	table := string(ev.Table.Table)

	t, err := c.FindTable(ctx, table)
	switch {
	case err == nil:
		// noop
	case errors.NotAllowed.Match(err):
		// do not execute the processRowsEventHandler function
		if c.opts.Log.IsDebug() {
			c.opts.Log.Debug("binlogsync.Canal.handleRowsEvent.Skipping.not_allowed_table", log.String("database", schemaName), log.String("table", table), log.Int("table_id", int(ev.TableID)))
		}
		return nil
	default:
		return errors.Wrapf(err, "[binlogsync] handleRowsEvent %q.%q", c.dsn.DBName, table)
	}

	var a string
	switch e.Header.EventType {
	case myreplicator.WRITE_ROWS_EVENTv1, myreplicator.WRITE_ROWS_EVENTv2:
		a = InsertAction
	case myreplicator.DELETE_ROWS_EVENTv1, myreplicator.DELETE_ROWS_EVENTv2:
		a = DeleteAction
	case myreplicator.UPDATE_ROWS_EVENTv1, myreplicator.UPDATE_ROWS_EVENTv2:
		a = UpdateAction
	default:
		return errors.NotSupported.Newf("[binlogsync] EventType %v not yet supported. Table %q.%q", e.Header.EventType, c.dsn.DBName, table)
	}
	return c.processRowsEventHandler(ctx, a, t, ev.Rows)
}

// FlushBinlog executes FLUSH BINARY LOGS.
func (c *Canal) FlushBinlog() error {
	_, err := c.dbcp.DB.Exec("FLUSH BINARY LOGS")
	return errors.WithStack(err)
}

func (c *Canal) WaitUntilPos(pos ddl.MasterStatus, timeout time.Duration) error {
	timer := time.NewTimer(timeout)
	for {
		select {
		case <-timer.C:
			return errors.Timeout.Newf("[binlogsync] Waited for position %v too long > %s", pos, timeout.String())
		default:
			if err := c.FlushBinlog(); err != nil {
				return errors.WithStack(err)
			}
			curPos := c.SyncedPosition()
			if curPos.Compare(pos) >= 0 {
				return nil
			} else {
				if c.opts.Log.IsDebug() {
					c.opts.Log.Debug("binlogsync.Canal.WaitUntilPos",
						log.String("current_pos", curPos.String()), log.String("waiting_for_post", pos.String()))
				}
				time.Sleep(100 * time.Millisecond)
			}
		}
	}
	return nil
}

//func (c *Canal) CatchMasterPos(timeout int) error {
//	rr, err := c.Execute("SHOW MASTER STATUS")
//	if err != nil {
//		return errors.Wrap(err, "[binlogsync] CatchMasterPos")
//	}
//
//	name, _ := rr.GetString(0, 0)
//	pos, _ := rr.GetInt(0, 1)
//
//	return c.WaitUntilPos(mysql.Position{Name: name, Pos: uint32(pos)}, timeout)
//}
