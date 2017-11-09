package binlogsync

import (
	"context"

	"github.com/corestoreio/csfw/sql/ddl"
	"github.com/corestoreio/errors"
	"github.com/corestoreio/log"
	"golang.org/x/sync/errgroup"
)

// TODO(CyS) investigate what would happen in case of transaction? should all
// the events be gathered together once a transaction starts? because on
// RollBack all events must be invalidated or better RowsEventHandler should not
// be called at all.

// RowsEventHandler calls your code when an event gets dispatched.
type RowsEventHandler interface {
	// Do function handles a RowsEvent bound to a specific database. If it
	// returns an error behaviour of "Interrupted", the canal type will stop the
	// syncer. Binlog has three update event version, v0, v1 and v2. For v1 and
	// v2, the rows number must be even. Two rows for one event, format is
	// [before update row, after update row] for update v0, only one row for a
	// event, and we don't support this version yet. The Do function will run in
	// its own Goroutine.
	Do(ctx context.Context, action string, t ddl.Table, rows [][]interface{}) error
	// Complete runs before a binlog rotation event happens. Same error rules
	// apply here like for function Do(). The Complete function will run in its
	// own Goroutine.
	Complete(context.Context) error
	// String returns the name of the handler
	String() string
}

// RegisterRowsEventHandler adds a new event handler to the internal list.
func (c *Canal) RegisterRowsEventHandler(h RowsEventHandler) {
	c.rsMu.Lock()
	defer c.rsMu.Unlock()

	if c.rsHandlers == nil {
		c.rsHandlers = make([]RowsEventHandler, 0, 4)
	}
	c.rsHandlers = append(c.rsHandlers, h)
}

func (c *Canal) travelRowsEventHandler(ctx context.Context, action string, table ddl.Table, rows [][]interface{}) error {
	c.rsMu.RLock()
	defer c.rsMu.RUnlock()

	erg, ctx := errgroup.WithContext(ctx)

	for _, h := range c.rsHandlers {
		h := h
		erg.Go(func() error {
			err := h.Do(ctx, action, table, rows)
			isInterr := errors.IsInterrupted(err)
			if err != nil && !isInterr {
				c.Log.Info("[binlogsync] Handler.Do error", log.Err(err), log.Stringer("handler_name", h),
					log.String("action", action), log.String("schema", c.DSN.DBName), log.String("table", table.Name))
			} else if isInterr {
				c.Log.Info("[binlogsync] Handler.Do Interrupt", log.Err(err), log.Stringer("handler_name", h),
					log.String("action", action), log.String("schema", c.DSN.DBName), log.String("table", table.Name))
				return errors.Wrap(err, "[binlogsync] travelRowsEventHandler interrupted")
			}
			return nil
		})
	}
	return errors.Wrap(erg.Wait(), "[binlogsync] travelRowsEventHandler errgroup Wait")
}

func (c *Canal) flushEventHandlers(ctx context.Context) error {
	c.rsMu.RLock()
	defer c.rsMu.RUnlock()

	erg, ctx := errgroup.WithContext(ctx)

	for _, h := range c.rsHandlers {
		h := h
		erg.Go(func() error {
			err := h.Complete(ctx)
			isInterr := errors.IsInterrupted(err)
			if err != nil && !isInterr {
				c.Log.Info("[binlogsync] flushEventHandlers.Handler.Complete error", log.Err(err), log.Stringer("handler_name", h))
			} else if isInterr {
				c.Log.Info("[binlogsync] flushEventHandlers.Handler.Complete interrupted", log.Err(err), log.Stringer("handler_name", h))
				return errors.Wrap(err, "[binlogsync] flushEventHandlers interrupted")
			}
			return nil
		})
	}
	return errors.Wrap(erg.Wait(), "[binlogsync] flushEventHandlers errgroup Wait")
}
