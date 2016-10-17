package binlogsync

import (
	"github.com/corestoreio/csfw/log"
	"github.com/corestoreio/csfw/util/errors"
	"github.com/siddontang/go-mysql/schema"
)

// RowsEventHandler calls your code when an event gets dispatched.
type RowsEventHandler interface {
	// Do function handles a RowsEvent bound to a specific database. If it
	// returns an error behaviour of "Interrupted", the canal type will stop the
	// syncer. Binlog has three update event version, v0, v1 and v2. For v1 and
	// v2, the rows number must be even. Two rows for one event, format is
	// [before update row, after update row] for update v0, only one row for a
	// event, and we don't support this version yet.
	Do(action string, table *schema.Table, rows [][]interface{}) error
	// Complete runs before a binlog rotation event happens. Same error rules apply
	// here like for function Do().
	Complete() error
	// String returns the name of the handler
	String() string
}

func (c *Canal) RegRowsEventHandler(h RowsEventHandler) {
	c.rsLock.Lock()
	c.rsHandlers = append(c.rsHandlers, h)
	c.rsLock.Unlock()
}

func (c *Canal) travelRowsEventHandler(action string, table *schema.Table, rows [][]interface{}) error {
	c.rsLock.RLock()
	defer c.rsLock.RUnlock()

	var err error
	for _, h := range c.rsHandlers {
		err = h.Do(action, table, rows)
		isInterr := errors.IsInterrupted(err)
		if err != nil && !isInterr {
			c.Log.Info("[binlogsync] Handler.Do error", log.Err(err), log.Stringer("handler_name", h),
				log.String("action", action), log.String("schema", table.Schema), log.String("table", table.Name))
		} else if isInterr {
			c.Log.Info("[binlogsync] Handler.Do Interrupt", log.Err(err), log.Stringer("handler_name", h),
				log.String("action", action), log.String("schema", table.Schema), log.String("table", table.Name))
			return errors.Wrap(err, "[binlogsync] travelRowsEventHandler interrupted")
		}
	}
	return nil
}

func (c *Canal) flushEventHandlers() error {
	c.rsLock.RLock()
	defer c.rsLock.RUnlock()

	var err error
	for _, h := range c.rsHandlers {
		err = h.Complete()
		isInterr := errors.IsInterrupted(err)
		if err != nil && !isInterr {
			c.Log.Info("[binlogsync] flushEventHandlers.Handler.Complete error", log.Err(err), log.Stringer("handler_name", h))
		} else if isInterr {
			c.Log.Info("[binlogsync] flushEventHandlers.Handler.Complete interrupted", log.Err(err), log.Stringer("handler_name", h))
			return errors.Wrap(err, "[binlogsync] flushEventHandlers interrupted")
		}
	}
	return nil
}
