// Go MySQL Driver - A MySQL-Driver for Go's database/sql package
//
// Copyright 2012 The Go-MySQL-Driver Authors. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this file,
// You can obtain one at http://mozilla.org/MPL/2.0/.

package null

import (
	"database/sql"
	"strings"
	"time"

	"github.com/corestoreio/errors"
)

/******************************************************************************
*                           Time related utils                                *
******************************************************************************/

// Time represents a time.Time that may be NULL.
// Time implements the Scanner interface so
// it can be used as a scan destination:
//
//  var nt Time
//  err := db.QueryRow("SELECT time FROM foo WHERE id=?", id).Scan(&nt)
//  ...
//  if nt.Valid {
//     // use nt.Time
//  } else {
//     // NULL value
//  }
//
// This Time implementation is not driver-specific
type Time struct {
	sql.NullTime
}

// Scan implements the Scanner interface. The value type must be time.Time or
// string / []byte (formatted time-string), otherwise Scan fails. It supports
// more input data than database/sql.NullTime.Scan
func (nt *Time) Scan(value interface{}) (err error) {
	nt.Time, nt.Valid = time.Time{}, false
	if value == nil {
		return
	}

	switch v := value.(type) {
	case time.Time:
		nt.Time = v
	case []byte:
		if v == nil {
			return
		}
		*nt, err = ParseDateTime(string(v), time.UTC)
	case string:
		if v == "" {
			return
		}
		*nt, err = ParseDateTime(v, time.UTC)
	default:
		err = errors.NotValid.Newf("[dml] Can't convert %T to time.Time. Maybe not yet implemented.", value)
	}
	nt.Valid = err == nil
	return
}

// ParseDateTime parses a string into a Time type. Empty string is considered NULL.
func ParseDateTime(str string, loc *time.Location) (t Time, err error) {
	if str == "" {
		return
	}
	zeroBase := "0000-00-00 00:00:00.000000000+00:00"
	base := "2006-01-02 15:04:05.999999999 07:00"
	if strings.IndexByte(str, 'T') > 0 {
		base = time.RFC3339Nano
	}

	switch lStr := len(str); lStr {
	case 10, 19, 21, 22, 23, 24, 25, 26, 27, 28, 29, 35: // up to "YYYY-MM-DD HH:MM:SS.MMMMMMM+HH:II"
		if str == zeroBase[:lStr] {
			return
		}
		t.Time, err = time.Parse(base[:lStr], str) // time.RFC3339Nano cannot be used due to the T
		t.Valid = err == nil
	default:
		err = errors.NotValid.Newf("invalid time string: %q", str)
		return
	}

	// Adjust location
	if err == nil && loc != time.UTC {
		y, mo, d := t.Time.Date()
		h, mi, s := t.Time.Clock()
		t.Time, err = time.Date(y, mo, d, h, mi, s, t.Time.Nanosecond(), loc), nil
	}

	return
}
