package myreplicator

import (
	"context"

	"github.com/corestoreio/csfw/log"
	"github.com/corestoreio/csfw/util/errors"
)

var (
	errSyncAlreadyClosed = errors.NewAlreadyClosedf("[myreplicator] Sync was closed")
)

// BinlogStreamer gets the streaming event.
type BinlogStreamer struct {
	Log     log.Logger
	bleChan chan *BinlogEvent
	errChan chan error
	err     error
}

// GetEvent gets the binlog event one by one, it will block until Syncer receives any events from MySQL
// or meets a sync error. You can pass a context (like Cancel or Timeout) to break the block.
// Returns a temporary error behaviour
func (s *BinlogStreamer) GetEvent(ctx context.Context) (*BinlogEvent, error) {
	if s.err != nil {
		return nil, errors.NewTemporaryf("[myreplicator] Last sync error or closed, try sync and get event again")
	}

	select {
	case ble := <-s.bleChan:
		return ble, nil
	case s.err = <-s.errChan:
		return nil, errors.Wrap(s.err, "[myreplicator] GetEvent error")
	case <-ctx.Done():
		return nil, errors.Wrap(ctx.Err(), "[myreplicator] GetEvent context error")
	}
}

func (s *BinlogStreamer) close() {
	s.closeWithError(errors.Wrap(errSyncAlreadyClosed, "[myreplicator] binlogstreamer close"))
}

func (s *BinlogStreamer) closeWithError(err error) {
	if err == nil {
		err = errors.Wrap(errSyncAlreadyClosed, "")
	}
	// log.Errorf("close sync with err: %v", err)
	select {
	case s.errChan <- err:
		if s.Log.IsInfo() {
			s.Log.Info("[myreplicator] closeWithError", log.Err(err))
		}
	default:
	}
}

func newBinlogStreamer(l log.Logger) *BinlogStreamer {
	s := new(BinlogStreamer)
	s.Log = l
	s.bleChan = make(chan *BinlogEvent, 10240)
	s.errChan = make(chan error, 4)

	return s
}
