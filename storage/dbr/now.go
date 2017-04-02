package dbr

import (
	"database/sql/driver"
	"time"
)

type nowSentinel struct{}

var now = time.Now

// Now is a value that serializes to the current time
var Now = nowSentinel{}

const timeFormat = "2006-01-02 15:04:05"

// Value implements a valuer for compatibility
func (n nowSentinel) Value() (driver.Value, error) {
	fnow := n.UTC().Format(timeFormat)
	return fnow, nil
}

// String returns the time string in format "2006-01-02 15:04:05"
func (n nowSentinel) String() string {
	return n.UTC().Format(timeFormat)
}

// UTC returns the UTC time
func (n nowSentinel) UTC() time.Time {
	return now().UTC()
}
