package myreplicator

import (
	"context"

	"github.com/corestoreio/errors"
	"github.com/corestoreio/log"
)

// BinlogStreamer gets the streaming event.
type BinlogStreamer struct {
	Log     log.Logger
	bleChan chan *BinlogEvent
	errChan chan error
	err     error
}

// GetEvent gets the binlog event one by one, it will block until Syncer
// receives any events from MySQL or meets a sync error. You can pass a context
// (like Cancel or Timeout) to break the block. May return a temporary error
// behaviour.
func (s *BinlogStreamer) GetEvent(ctx context.Context) (*BinlogEvent, error) {
	if s.err != nil {
		return nil, errors.Temporary.Newf("[myreplicator] Last sync error or closed, try sync and get event again")
	}

	select {
	case ble := <-s.bleChan:
		return ble, nil
	case s.err = <-s.errChan:
		return nil, errors.WithStack(s.err)
	case <-ctx.Done():
		return nil, errors.WithStack(ctx.Err())
	}
}

// DumpEvents dumps all left events
func (s *BinlogStreamer) DumpEvents() []*BinlogEvent {
	count := len(s.bleChan)
	events := make([]*BinlogEvent, 0, count)
	for i := 0; i < count; i++ {
		events = append(events, <-s.bleChan)
	}
	return events
}

func (s *BinlogStreamer) close() {
	s.closeWithError(errors.AlreadyClosed.Newf("[myreplicator] Sync already closed"))
}

func (s *BinlogStreamer) closeWithError(err error) {
	if err == nil {
		err = errors.AlreadyClosed.Newf("[myreplicator] Sync closed")
	}

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
