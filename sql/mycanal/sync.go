package mycanal

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
	expCreateTable   = regexp.MustCompile("(?i)^CREATE\\sTABLE(\\sIF\\sNOT\\sEXISTS)?\\s`{0,1}(.*?)`{0,1}\\.{0,1}`{0,1}([^`\\.]+?)`{0,1}\\s.*")
	expAlterTable    = regexp.MustCompile("(?i)^ALTER\\sTABLE\\s.*?`{0,1}(.*?)`{0,1}\\.{0,1}`{0,1}([^`\\.]+?)`{0,1}\\s.*")
	expRenameTable   = regexp.MustCompile("(?i)^RENAME\\sTABLE\\s.*?`{0,1}(.*?)`{0,1}\\.{0,1}`{0,1}([^`\\.]+?)`{0,1}\\s{1,}TO\\s.*?")
	expDropTable     = regexp.MustCompile("(?i)^DROP\\sTABLE(\\sIF\\sEXISTS){0,1}\\s`{0,1}(.*?)`{0,1}\\.{0,1}`{0,1}([^`\\.]+?)`{0,1}(?:$|\\s)")
	expTruncateTable = regexp.MustCompile("(?i)^TRUNCATE\\s+(?:TABLE\\s+)?(?:`?([^`\\s]+)`?\\.`?)?([^`\\s]+)`?")
	ddlExpressions   = [...]*regexp.Regexp{expCreateTable, expAlterTable, expRenameTable, expDropTable, expTruncateTable}
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
	defer log.WhenDone(c.opts.Log).Info("myCanal.clearTableCacheOnDDLStmt")
	if db, tbl := extractTableFromQueryEvent(schema, query); tbl != "" {
		c.ClearTableCache(db, tbl)
		if c.opts.Log.IsDebug() {
			c.opts.Log.Debug("[mycanal] Table structure changed, clear table cache",
				log.String("database", db), log.String("table", tbl), log.String("query", string(query)))
		}
	}
}

func (c *Canal) startStream() (s *myreplicator.BinlogStream, err error) {
	if c.syncer == nil {
		return nil, errors.AlreadyClosed.Newf("[mycanal] Canal already closed and myreplicator.BinlogSyncer is nil")
	}

	c.masterMu.RLock()
	defer c.masterMu.RUnlock()

	gset := c.masterGTID
	var gsetStr string
	if gset == nil {
		s, err = c.syncer.StartSync(c.masterStatus)
	} else {
		gsetStr = gset.String()
		s, err = c.syncer.StartSyncGTID(gset)
	}
	if c.opts.Log.IsDebug() {
		c.opts.Log.Debug("myCanal.startStream.start", log.Err(err), log.Stringer("position", c.masterStatus), log.String("gtid", gsetStr))
	}
	if err != nil {
		err = errors.WithStack(err)
	}
	return
}

func (c *Canal) startSyncBinlog(ctxArg context.Context) error {

	s, err := c.startStream()
	if err != nil {
		return err
	}

	c.masterMu.RLock()
	pos := c.masterStatus
	c.masterMu.RUnlock()
	for {
		ev, err := s.GetEvent(ctxArg)
		if err != nil {
			if errors.Cause(err) == context.Canceled {
				if c.opts.Log.IsDebug() {
					c.opts.Log.Debug("myCanal.startSyncBinlog.context.canceled")
				}
				return nil
			}
			return errors.WithStack(err)
		}

		currentPos := pos
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

			if c.opts.Log.IsDebug() {
				c.opts.Log.Debug("myCanal.startSyncBinlog.rotateEvent.newPosition", log.Stringer("position", pos))
			}

			// call event handler OnRotate(e)

		case *myreplicator.RowsEvent:
			// we only focus row based event.
			// NotFound errors get ignores. For example table has been deleted
			// and an old event pops in.
			if err = c.handleRowsEvent(ctxArg, ev); err != nil {
				isNotFound := errors.Is(err, errors.NotFound)
				if c.opts.Log.IsDebug() {
					c.opts.Log.Debug("myCanal.startSyncBinlog.rowsEvent.newPosition", log.Err(err),
						log.Stringer("header_position", pos), log.Stringer("current_position", currentPos),
						log.Bool("ignore_not_found_error", isNotFound))
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

		case *myreplicator.MariadbGTIDEvent:
			// TODO implement
			// gtid, err := mysql.ParseMariadbGTIDSet(e.GTID.String())
			// if err != nil {
			// 	return errors.WithStack(err)
			// }
			// if err := c.eventHandler.OnGTID(gtid); err != nil {
			// 	return errors.Trace(err)
			// }

		case *myreplicator.GTIDEvent:
			// TODO implement
			// u, _ := uuid.FromBytes(e.SID)
			// gtid, err := mysql.ParseMysqlGTIDSet(fmt.Sprintf("%s:%d", u.String(), e.GNO))
			// if err != nil {
			// 	return errors.Trace(err)
			// }
			// if err := c.eventHandler.OnGTID(gtid); err != nil {
			// 	return errors.Trace(err)
			// }

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
			// don't update Master with file and position
		default:
			if c.opts.Log.IsDebug() {
				c.opts.Log.Debug("myCanal.startSyncBinlog.unknown.event", log.ObjectTypeOf("event_type", ev.Event), log.Stringer("position", pos))
			}
			continue
		}

		if err := c.masterSave(pos.File, pos.Position); err != nil {
			if c.opts.Log.IsInfo() {
				c.opts.Log.Info("myCanal.startSyncBinlog.master.position.Failed", log.Err(err), log.Stringer("position", pos))
			}
		}
	}
}

// handleRowsEvent handles an event on the rows and calls all registered rows
// event handler. can return different error behaviours.
func (c *Canal) handleRowsEvent(ctx context.Context, e *myreplicator.BinlogEvent) error {
	if c.opts.Log.IsDebug() {
		defer log.WhenDone(c.opts.Log).Debug("myCanal.handleRowsEvent")
	}
	ev, ok := e.Event.(*myreplicator.RowsEvent)
	if !ok {
		return errors.Fatal.Newf("[mycanal] handleRowsEvent: Failed to cast to *myreplicator.RowsEvent type")
	}

	// Caveat: table may be altered at runtime.
	schemaName := string(ev.Table.Schema)
	if c.dsn.DBName != schemaName {
		if c.opts.Log.IsDebug() {
			c.opts.Log.Debug("myCanal.handleRowsEvent.Skipping.database", log.String("database_have", schemaName),
				log.String("database_want", c.dsn.DBName), log.Int("table_id", int(ev.TableID)))
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
			c.opts.Log.Debug("myCanal.handleRowsEvent.Skipping.not_allowed_table", log.String("database", schemaName),
				log.String("table", table), log.Int("table_id", int(ev.TableID)))
		}
		return nil
	default:
		return errors.Wrapf(err, "[mycanal] handleRowsEvent %q.%q", c.dsn.DBName, table)
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
		return errors.NotSupported.Newf("[mycanal] EventType %v not yet supported. Table %q.%q", e.Header.EventType, c.dsn.DBName, table)
	}
	return c.processRowsEventHandler(ctx, a, t, ev.Rows)
}

func (c *Canal) GetMasterPos() (ms ddl.MasterStatus, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), c.opts.MasterStatusQueryTimeout)
	defer cancel()
	if _, err = c.dbcp.WithQueryBuilder(&ms).Load(ctx, &ms); err != nil {
		return ms, errors.WithStack(err)
	}
	return
}

// FlushBinlog executes FLUSH BINARY LOGS.
func (c *Canal) FlushBinlog() error {
	_, err := c.dbcp.DB.Exec("FLUSH BINARY LOGS")
	return errors.WithStack(err)
}

// WaitUntilPos flushes the binary logs until we've reached the desired position.
func (c *Canal) WaitUntilPos(pos ddl.MasterStatus, timeout time.Duration) error {
	timer := time.NewTimer(timeout)
	for {
		select {
		case <-timer.C:
			return errors.Timeout.Newf("[mycanal] Waited for position %v too long > %s", pos, timeout.String())
		default:
			if err := c.FlushBinlog(); err != nil {
				return errors.WithStack(err)
			}
			curPos := c.SyncedPosition()
			if curPos.Compare(pos) >= 0 {
				return nil
			} else {
				if c.opts.Log.IsDebug() {
					c.opts.Log.Debug("myCanal.WaitUntilPos",
						log.String("current_pos", curPos.String()), log.String("waiting_for_post", pos.String()))
				}
				time.Sleep(100 * time.Millisecond)
			}
		}
	}
}

// CatchMasterPos reads the current master position and waits until we reached
// it.
func (c *Canal) CatchMasterPos(timeout time.Duration) error {
	pos, err := c.GetMasterPos()
	if err != nil {
		return errors.WithStack(err)
	}

	return c.WaitUntilPos(pos, timeout)
}
