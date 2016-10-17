package binlogsync

import (
	"context"
	"time"

	"github.com/corestoreio/csfw/log"
	"github.com/corestoreio/csfw/util/errors"
	"github.com/siddontang/go-mysql/mysql"
	"github.com/siddontang/go-mysql/replication"
)

const (
	UpdateAction = "update"
	InsertAction = "insert"
	DeleteAction = "delete"
)

func (c *Canal) startSyncBinlog() error {
	pos := mysql.Position{Name: c.master.FileName, Pos: c.master.Position}

	if c.Log.IsInfo() {
		c.Log.Info("[binlogsync] Start syncing of binlog", log.Stringer("position", pos))
	}

	s, err := c.syncer.StartSync(pos)
	if err != nil {
		return errors.NewFatalf("[binlogsync] Start sync replication at %s error %v", pos, err)
	}

	timeout := time.Second
	forceSavePos := false
	for {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		ev, err := s.GetEvent(ctx)
		cancel()

		if err == context.DeadlineExceeded {
			timeout = 2 * timeout
			continue
		}
		if err != nil {
			return errors.Wrap(err, "[binlogsync] startSyncBinlog.GetEvent")
		}

		timeout = time.Second

		//next binlog pos
		pos.Pos = ev.Header.LogPos

		forceSavePos = false

		switch e := ev.Event.(type) {
		case *replication.RotateEvent:
			if err := c.flushEventHandlers(); err != nil {
				// todo maybe better err handling ...
				return errors.Wrap(err, "[binlogsync] startSyncBinlog.flushEventHandlers")
			}
			pos.Name = string(e.NextLogName)
			pos.Pos = uint32(e.Position)
			// r.ev <- pos
			forceSavePos = true

			if c.Log.IsInfo() {
				c.Log.Info("[binlogsync] Rotate binlog to a new position", log.Stringer("position", pos))
			}

		case *replication.RowsEvent:
			// we only focus row based event
			if err = c.handleRowsEvent(ev); err != nil {
				if c.Log.IsInfo() {
					c.Log.Info("[binlogsync] Rotate binlog to a new position", log.Err(err), log.Stringer("position", pos))
				}
				return errors.Wrap(err, "[binlogsync] handleRowsEvent")
			}
		case *replication.TableMapEvent:
			continue
			//default:
			//	fmt.Printf("%#v\n\n", e)
		}

		c.master.Update(pos.Name, pos.Pos)
		c.master.Save(forceSavePos)
	}

	return nil
}

func (c *Canal) handleRowsEvent(e *replication.BinlogEvent) error {
	ev, ok := e.Event.(*replication.RowsEvent)
	if !ok {
		return errors.NewFatalf("[binlogsync] handleRowsEvent: Failed to cast to *replication.RowsEvent type")
	}

	// Caveat: table may be altered at runtime.

	if in := string(ev.Table.Schema); c.database != in {
		if c.Log.IsDebug() {
			c.Log.Debug("[binlogsync] Skipping database", log.String("database_have", in), log.String("database_want", c.database), log.Int("table_id", int(ev.TableID)))
		}
		return nil
	}

	table := string(ev.Table.Table)

	t, err := c.GetTable(table)
	if err != nil {
		return errors.Wrapf(err, "[binlogsync] GetTable %q.%q", c.database, table)
	}
	var action string
	switch e.Header.EventType {
	case replication.WRITE_ROWS_EVENTv1, replication.WRITE_ROWS_EVENTv2:
		action = InsertAction
	case replication.DELETE_ROWS_EVENTv1, replication.DELETE_ROWS_EVENTv2:
		action = DeleteAction
	case replication.UPDATE_ROWS_EVENTv1, replication.UPDATE_ROWS_EVENTv2:
		action = UpdateAction
	default:
		return errors.NewNotSupportedf("[binlogsync] EventType %v not yet supported", e.Header.EventType)
	}
	return c.travelRowsEventHandler(action, t, ev.Rows)
}

func (c *Canal) WaitUntilPos(pos mysql.Position, timeout int) error {
	if timeout <= 0 {
		timeout = 60
	}

	timer := time.NewTimer(time.Duration(timeout) * time.Second)
	for {
		select {
		case <-timer.C:
			return errors.NewTimeoutf("[binlogsync] WaitUntilPos wait position %v err", pos)
		default:
			curpos := c.master.Pos()
			if curpos.Compare(pos) >= 0 {
				return nil
			} else {
				time.Sleep(100 * time.Millisecond)
			}
		}
	}

	return nil
}

func (c *Canal) CatchMasterPos(timeout int) error {
	rr, err := c.Execute("SHOW MASTER STATUS")
	if err != nil {
		return errors.Wrap(err, "[binlogsync] CatchMasterPos")
	}

	name, _ := rr.GetString(0, 0)
	pos, _ := rr.GetInt(0, 1)

	return c.WaitUntilPos(mysql.Position{Name: name, Pos: uint32(pos)}, timeout)
}
