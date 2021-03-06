package myreplicator

import (
	"context"
	"time"

	"github.com/corestoreio/errors"
	"github.com/corestoreio/log"
)

// BinlogStream gets the streaming event.
type BinlogStream struct {
	Log     log.Logger
	bleChan chan *BinlogEvent
	errChan chan error
	err     error
}

// GetEvent gets the binlog event one by one, it will block until Syncer
// receives any events from MySQL or meets a sync error. You can pass a context
// (like Cancel or Timeout) to break the block. May return a temporary error
// behaviour.
func (s *BinlogStream) GetEvent(ctx context.Context) (*BinlogEvent, error) {
	if s.err != nil {
		return nil, errors.Temporary.New(s.err, "[myreplicator] Last sync error or closed, try sync and get event again")
	}

	select {
	case ble := <-s.bleChan:
		return ble, nil
	case s.err = <-s.errChan:
		return nil, errors.WithStack(s.err)
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

// GetEventWithStartTime gets the binlog event with a start time, if current
// binlog event timestamp smaller than specify start time returns a nil event.
func (s *BinlogStream) GetEventWithStartTime(ctx context.Context, startTime time.Time) (*BinlogEvent, error) {
	if s.err != nil {
		return nil, errors.Temporary.New(s.err, "[myreplicator] Last sync error or closed, try sync and get event again")
	}
	startUnix := startTime.Unix()
	select {
	case c := <-s.bleChan:
		if int64(c.Header.Timestamp) >= startUnix {
			return c, nil
		}
		return nil, nil
	case s.err = <-s.errChan:
		return nil, s.err
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

// DumpEvents dumps all left events
func (s *BinlogStream) DumpEvents() []*BinlogEvent {
	count := len(s.bleChan)
	events := make([]*BinlogEvent, 0, count)
	for i := 0; i < count; i++ {
		events = append(events, <-s.bleChan)
	}
	return events
}

func (s *BinlogStream) closeWithError(err error) {
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

func newBinlogStreamer(l log.Logger) *BinlogStream {
	s := new(BinlogStream)
	s.Log = l
	s.bleChan = make(chan *BinlogEvent, 10240)
	s.errChan = make(chan error, 4)

	return s
}
