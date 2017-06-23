// Copyright 2015-2017, Cyrill @ Schumacher.fm and the CoreStore contributors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package dbr

import (
	"bytes"
	"fmt"
	"strconv"

	"github.com/corestoreio/csfw/util/bufferpool"
	"github.com/corestoreio/errors"
)

// QueryBuilder assembles a query and returns the raw SQL without parameter
// substitution and the arguments.
type QueryBuilder interface {
	ToSQL() (string, []interface{}, error)
}

type queryBuilder interface {
	toSQL(queryWriter) error
	// appendArgs appends the arguments to Arguments and returns them. If
	// argument `Arguments` is nil, allocates new bytes
	appendArgs(Arguments) (Arguments, error)
	hasBuildCache() bool
	writeBuildCache(sql []byte)
	// readBuildCache returns the cached SQL string including its place holders.
	readBuildCache() (sql []byte, args Arguments, err error)
}

// queryWriter at used to generate a query.
type queryWriter interface {
	WriteString(s string) (n int, err error)
	WriteRune(r rune) (n int, err error)
	WriteByte(c byte) error
	Write(p []byte) (n int, err error)
	Reset()
	Bytes() []byte
}

var _ queryWriter = (*backHole)(nil)

type backHole struct{} // TODO(CyS) just a temporary implementation. should get removed later

func (backHole) WriteString(s string) (n int, err error) { return }
func (backHole) WriteRune(r rune) (n int, err error)     { return }
func (backHole) WriteByte(c byte) error                  { return nil }
func (backHole) Write(p []byte) (n int, err error)       { return }
func (backHole) Bytes() []byte                           { return nil }
func (backHole) Reset()                                  {}

// For the sake of readability within the source code, because boolean arguments
// are terrible.
const (
	isNotPrepared    = false
	isPrepared       = true
	isNotInterpolate = false
)

// toSQL generates the SQL string and its place holders. Takes care of caching
// and interpolation. It returns the string with placeholders and a slice of
// query arguments. With switched on interpolation, it only returns a string
// including the stringyfied arguments. With an enabled cache, the arguments
// gets regenerated each time a call to ToSQL happens.
// isPrepared if true skips assembling the arguments.
func toSQL(b queryBuilder, isInterpolate, isPrepared bool) (string, []interface{}, error) {
	var ipBuf *bytes.Buffer // ip = interpolate buffer
	if isInterpolate {
		ipBuf = bufferpool.Get()
		defer bufferpool.Put(ipBuf)
	}

	useCache := b.hasBuildCache()
	if useCache {
		sql, args, err := b.readBuildCache()
		if err != nil {
			return "", nil, errors.WithStack(err)
		}
		if sql != nil {
			if isInterpolate && !isPrepared {
				err := interpolate(ipBuf, sql, args...)
				return ipBuf.String(), nil, errors.WithStack(err)
			}
			return string(sql), args.Interfaces(), nil
		}
	}

	buf := bufferpool.Get()
	defer bufferpool.Put(buf)

	if err := b.toSQL(buf); err != nil {
		return "", nil, errors.WithStack(err)
	}
	if useCache {
		sqlCopy := make([]byte, buf.Len())
		copy(sqlCopy, buf.Bytes())
		b.writeBuildCache(sqlCopy)
	}
	if isPrepared {
		return buf.String(), nil, nil
	}

	// capacity of Arguments gets handled in the concret implementation of `b`
	args, err := b.appendArgs(Arguments{})
	if err != nil {
		return "", nil, errors.WithStack(err)
	}

	if isInterpolate {
		err := interpolate(ipBuf, buf.Bytes(), args...)
		return ipBuf.String(), nil, errors.WithStack(err)
	}
	return buf.String(), args.Interfaces(), nil
}

func makeSQL(b queryBuilder, isInterpolate bool) string {
	sRaw, _, err := toSQL(b, isInterpolate, isNotPrepared)
	if err != nil {
		return fmt.Sprintf("[dbr] ToSQL Error: %+v", err)
	}
	return sRaw
}

// String returns a string representing a preprocessed, interpolated, query.
// On error, the error gets printed. Fulfills interface fmt.Stringer.
func (b *Delete) String() string {
	return makeSQL(b, b.IsInterpolate)
}

// String returns a string representing a preprocessed, interpolated, query.
// On error, the error gets printed. Fulfills interface fmt.Stringer.
func (b *Insert) String() string {
	return makeSQL(b, b.IsInterpolate)
}

// String returns a string representing a preprocessed, interpolated, query.
// On error, the error gets printed. Fulfills interface fmt.Stringer.
func (b *Select) String() string {
	return makeSQL(b, b.IsInterpolate)
}

// String returns a string representing a preprocessed, interpolated, query.
// On error, the error gets printed. Fulfills interface fmt.Stringer.
func (b *Update) String() string {
	return makeSQL(b, b.IsInterpolate)
}

func sqlWriteUnionAll(w queryWriter, isAll bool, isIntersect bool, isExcept bool) {
	w.WriteByte('\n')
	switch {
	case isIntersect:
		w.WriteString("INTERSECT") // MariaDB >= 10.3
	case isExcept:
		w.WriteString("EXCEPT") // MariaDB >= 10.3
	default:
		w.WriteString("UNION")
		if isAll {
			w.WriteString(" ALL")
		}
	}
	w.WriteByte('\n')
}

func sqlWriteOrderBy(w queryWriter, orderBys aliases, br bool) {
	if len(orderBys) == 0 {
		return
	}
	brS := ' '
	if br {
		brS = '\n'
	}
	w.WriteRune(brS)
	w.WriteString("ORDER BY ")
	orderBys.WriteQuoted(w)
}

func sqlWriteLimitOffset(w queryWriter, limitValid bool, limitCount uint64, offsetValid bool, offsetCount uint64) {
	if limitValid {
		w.WriteString(" LIMIT ")
		writeUint64(w, limitCount)
	}
	if offsetValid {
		w.WriteString(" OFFSET ")
		writeUint64(w, offsetCount)
	}
}

func writeFloat64(w queryWriter, f float64) (err error) {
	d := w.Bytes()
	w.Reset()
	_, err = w.Write(strconv.AppendFloat(d, f, 'f', -1, 64))
	return err
}

func writeInt64(w queryWriter, i int64) (err error) {
	d := w.Bytes()
	w.Reset()
	_, err = w.Write(strconv.AppendInt(d, i, 10))
	return err
}

func writeUint64(w queryWriter, i uint64) (err error) {
	d := w.Bytes()
	w.Reset()
	_, err = w.Write(strconv.AppendUint(d, i, 10))
	return err
}
