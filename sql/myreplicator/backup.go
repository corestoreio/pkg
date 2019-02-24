package myreplicator

import (
	"context"
	"io"
	"os"
	"path"
	"time"

	"github.com/corestoreio/errors"
	"github.com/corestoreio/pkg/sql/ddl"
)

// StartBackup enables, like the mysqlbinlog command line tool, a remote raw
// backup. Backup remote binlog from position (filename, offset) and write in
// backupDir.
func (b *BinlogSyncer) StartBackup(backupDir string, perm os.FileMode, p ddl.MasterStatus, timeout time.Duration) error {
	if timeout == 0 {
		// a very long timeout here
		timeout = 30 * 3600 * 24 * time.Second
	}

	// Force use raw mode
	b.parser.SetRawMode(true)

	if err := os.MkdirAll(backupDir, perm); err != nil {
		return errors.WithStack(err)
	}

	s, err := b.StartSync(p)
	if err != nil {
		return errors.WithStack(err)
	}

	var filename string
	var offset uint32

	var f *os.File
	defer func() {
		if f != nil {
			if err2 := f.Close(); err2 != nil && err == nil {
				err = errors.WithStack(err2)
			}
		}
	}()

	for {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		e, err := s.GetEvent(ctx)
		cancel()

		if err == context.DeadlineExceeded {
			return err
		}
		if err != nil {
			return errors.WithStack(err)
		}

		offset = e.Header.LogPos

		if e.Header.EventType == ROTATE_EVENT {
			rotateEvent := e.Event.(*RotateEvent)
			filename = string(rotateEvent.NextLogName)

			if e.Header.Timestamp == 0 || offset == 0 {
				// fake rotate event
				continue
			}
		} else if e.Header.EventType == FORMAT_DESCRIPTION_EVENT {
			// FormateDescriptionEvent is the first event in binlog, we will close old one and create a new

			if f != nil {
				if err2 := f.Close(); err2 != nil && err == nil {
					return errors.WithStack(err2)
				}
			}

			if len(filename) == 0 {
				return errors.Empty.Newf("[myreplicator] empty binlog filename for FormateDescriptionEvent")
			}

			f, err = os.OpenFile(path.Join(backupDir, filename), os.O_CREATE|os.O_WRONLY, 0644)
			if err != nil {
				return errors.WithStack(err)
			}

			// write binlog header fe'bin'
			if _, err = f.Write(BinLogFileHeader); err != nil {
				return errors.WithStack(err)
			}

		}

		if n, err := f.Write(e.RawData); err != nil {
			return errors.WithStack(err)
		} else if n != len(e.RawData) {
			return errors.WithStack(io.ErrShortWrite)
		}
	}
}
