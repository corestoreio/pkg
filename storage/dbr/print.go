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
	ToSQL() (string, Arguments, error)
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
}

var _ queryWriter = (*backHole)(nil)

type backHole struct{} // TODO(CyS) just a temporary implementation. should get removed later

func (backHole) WriteString(s string) (n int, err error) { return }
func (backHole) WriteRune(r rune) (n int, err error)     { return }
func (backHole) WriteByte(c byte) error                  { return nil }
func (backHole) Write(p []byte) (n int, err error)       { return }

// toSQL generates the SQL string and its place holders. Takes care of caching
// and interpolation. It returns the string with placeholders and a slice of
// query arguments. With switched on interpolation, it only returns a string
// including the stringyfied arguments. With an enabled cache, the arguments
// gets regenerated each time a call to ToSQL happens.
func toSQL(b queryBuilder, isInterpolate bool) (string, Arguments, error) {
	var ipBuf *bytes.Buffer // ip = interpolate buffer
	if isInterpolate {
		ipBuf = bufferpool.Get()
		defer bufferpool.Put(ipBuf)
	}

	useCache := b.hasBuildCache()
	if useCache {
		sql, args, err := b.readBuildCache()
		if err != nil {
			return "", nil, errors.Wrap(err, "[dbr] toSQL.readBuildCache")
		}
		if sql != nil {
			if isInterpolate {
				err := interpolate(ipBuf, sql, args...)
				return ipBuf.String(), nil, errors.Wrap(err, "[dbr] toSQL.Interpolate")
			}
			return string(sql), args, nil
		}
	}

	buf := bufferpool.Get()
	defer bufferpool.Put(buf)

	if err := b.toSQL(buf); err != nil {
		return "", nil, errors.Wrap(err, "[dbr] toSQL.toSQL")
	}
	// capacity of Arguments gets handled in the concret implementation of `b`
	args, err := b.appendArgs(Arguments{})
	if err != nil {
		return "", nil, errors.Wrap(err, "[dbr] toSQL.appendArgs")
	}
	if useCache {
		sqlCopy := make([]byte, buf.Len())
		copy(sqlCopy, buf.Bytes())
		b.writeBuildCache(sqlCopy)
	}

	if isInterpolate {
		err := interpolate(ipBuf, buf.Bytes(), args...)
		return ipBuf.String(), nil, errors.Wrap(err, "[dbr] toSQL.Interpolate")
	}
	return buf.String(), args, nil
}

func toSQLPrepared(b queryBuilder) (string, error) {
	// TODO(CyS) implement build cache like the toSQL function. see above.
	buf := bufferpool.Get()
	defer bufferpool.Put(buf)
	err := b.toSQL(buf)
	return buf.String(), errors.Wrap(err, "[dbr] toSQLPrepared.toSQL")
}

func makeSQL(b queryBuilder, isInterpolate bool) string {
	sRaw, _, err := toSQL(b, isInterpolate)
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

func sqlWriteUnionAll(w queryWriter, isAll bool) {
	w.WriteString("\nUNION")
	if isAll {
		w.WriteString(" ALL")
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
	for i, c := range orderBys {
		if i > 0 {
			w.WriteString(", ")
		}
		c.FquoteAs(w)
		// TODO append arguments
	}
}

func sqlWriteLimitOffset(w queryWriter, limitValid bool, limitCount uint64, offsetValid bool, offsetCount uint64) {
	if limitValid {
		w.WriteString(" LIMIT ")
		w.WriteString(strconv.FormatUint(limitCount, 10))
	}
	if offsetValid {
		w.WriteString(" OFFSET ")
		w.WriteString(strconv.FormatUint(offsetCount, 10))
	}
}
