// Copyright 2015-present, Cyrill @ Schumacher.fm and the CoreStore contributors
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

package dml

import (
	"bytes"
	"fmt"
	"strconv"
	"unicode/utf8"

	"github.com/corestoreio/errors"
	"github.com/corestoreio/log"
	"github.com/corestoreio/pkg/util/bufferpool"
)

const (
	dmlTypeSelect = 's'
	dmlTypeInsert = 'i'
	dmlTypeUpdate = 'u'
	dmlTypeDelete = 'd'
	dmlTypeWith   = 'w'
	dmlTypeUnion  = 'n'
	// dmlTypeShow   = 'h'
)

// QueryBuilder assembles a query and returns the raw SQL without parameter
// substitution and the arguments.
type QueryBuilder interface {
	ToSQL() (string, []interface{}, error)
}

type QuerySQL string

func (s QuerySQL) ToSQL() (string, []interface{}, error) {
	return string(s), nil, nil
}

type queryBuilder interface {
	toSQL(w *bytes.Buffer, placeHolders []string) ([]string, error)
	writeBuildCache(sql []byte)
	// readBuildCache returns the cached SQL string
	readBuildCache() (sql []byte)
}

// multiplyArguments is only applicable when using *Union as a template.
// multiplyArguments repeats the `args` variable n-times to match the number of
// generated SELECT queries in the final UNION statement. It should be called
// after all calls to `StringReplace` have been made.
func multiplyArguments(templateStmtCount int, args Arguments) Arguments {
	if templateStmtCount == 1 {
		return args
	}
	ret := make(Arguments, len(args)*templateStmtCount)
	lArgs := len(args)
	for i := 0; i < templateStmtCount; i++ {
		copy(ret[i*lArgs:], args)
	}
	return ret
}

// builderCommon
type builderCommon struct {
	// ID of a statement. Used in logging. The ID gets generated with function
	// signature `func() string`. This func gets applied to the logger when
	// setting up a logger.
	id  string     // tracing ID
	Log log.Logger // Log optional logger

	argsRecords []QualifiedRecord
	argsArgs    Arguments
	argsRaw     []interface{}
	// ärgErr represents an argument error caused in one of the three With
	// functions.
	ärgErr error // Sorry Germans for that terrible pun #notSorry

	defaultQualifier string
	// isWithInterfaces will be set to true if the raw interface arguments are
	// getting applied.
	isWithInterfaces bool
	// qualifiedColumns gets collected before calling ToSQL, and clearing the all
	// pointers, to know which columns need values from the QualifiedRecords
	qualifiedColumns []string
	// templateStmtCount only used in case a UNION statement acts as a template.
	// Create one SELECT statement and by setting the data for
	// Union.StringReplace function additional SELECT statements are getting
	// created. Now the arguments must be multiplied by the number of new
	// created SELECT statements. This value  gets stored in templateStmtCount.
	// An example exists in TestUnionTemplate_ReuseArgs.
	templateStmtCount int
	// ColumnMap internal intermediate typ which scans into its own optimized
	// types to avoid lots of allocations. It can also verify UTF8 strings.
	ColumnMap ColumnMap
}

func (bc builderCommon) convertRecordsToArguments() (Arguments, error) {
	if bc.templateStmtCount == 0 {
		bc.templateStmtCount = 1
	}
	if len(bc.argsArgs) == 0 && len(bc.argsRecords) == 0 {
		return bc.argsArgs, nil
	}

	if len(bc.argsArgs) > 0 && len(bc.argsRecords) == 0 && !bc.argsArgs.hasNamedArgs() {
		return multiplyArguments(bc.templateStmtCount, bc.argsArgs), nil
	}

	cm := newColumnMap(make(Arguments, 0, len(bc.argsArgs)+len(bc.argsRecords)), "")
	var unnamedCounter int
	for tsc := 0; tsc < bc.templateStmtCount; tsc++ { // only in case of UNION statements in combination with a template SELECT, can be optimized later
		for _, identifier := range bc.qualifiedColumns { // contains the correct order as the place holders appear in the SQL string
			qualifier, column := splitColumn(identifier)
			if qualifier == "" {
				qualifier = bc.defaultQualifier
			}
			var cut bool
			column, cut = cutPrefix(column, namedArgStartStr)
			cm.columns[0] = column // length is always one!

			if !cut { // if the colon : cannot be found then a simple place holder ? has been detected
				if pArg, ok := bc.argsArgs.unnamedArgByPos(unnamedCounter); ok {
					cm.Args = append(cm.Args, pArg)
				}
				unnamedCounter++
				//continue
			}
			for _, qRec := range bc.argsRecords {
				if qRec.Qualifier == "" {
					qRec.Qualifier = bc.defaultQualifier
				}
				if qRec.Qualifier == qualifier {
					if err := qRec.Record.MapColumns(cm); err != nil {
						return nil, errors.WithStack(err)
					}
				}
			}

			if err := bc.argsArgs.MapColumns(cm); err != nil {
				return nil, errors.WithStack(err)
			}
		}
	}
	if len(cm.Args) == 0 {
		return append(cm.Args, bc.argsArgs...), nil
	}
	return cm.Args, nil
}

// estimatedCachedSQLSize 1024 bytes value got retrieved by analyzing and
// reviewing some M2 SQL queries.
const estimatedCachedSQLSize = 1024

// BuilderBase contains fields which all SQL query builder have in common, the
// same base. Exported for documentation reasons.
type BuilderBase struct {
	builderCommon
	// cachedSQL contains the final SQL string which gets send to the server.
	cachedSQL []byte
	// EstimatedCachedSQLSize specifies the estimated size in bytes of the final
	// SQL string. This value gets used during SQL string building process to
	// reduce the allocations and speed up the process. Default Value is xxxx
	// Bytes.
	EstimatedCachedSQLSize uint16
	RawFullSQL             string
	Table                  id
	// PropagationStopped set to true if you would like to interrupt the
	// listener chain. Once set to true all sub sequent calls of the next
	// listeners will be suppressed.
	PropagationStopped   bool
	IsInterpolate        bool // See Interpolate()
	IsBuildCacheDisabled bool // see DisableBuildCache()
	IsExpandPlaceHolders bool // see ExpandPlaceHolders()
	// IsUnsafe if set to true the functions AddColumn* will turn any
	// non valid identifier (not `{a-z}[a-z0-9$_]+`i) into an expression.
	IsUnsafe bool
	// propagationStoppedAt position in the slice where the stopped propagation
	// has been requested. for every new iteration the propagation must stop at
	// this position.
	propagationStoppedAt int
}

// hasBuildCache satisfies partially interface queryBuilder
func (bb *BuilderBase) hasBuildCache() bool {
	return !bb.IsBuildCacheDisabled
}

func (bb *BuilderBase) resetArgs() {
	bb.argsArgs = bb.argsArgs[:0]
	bb.argsRaw = bb.argsRaw[:0]
	bb.argsRecords = bb.argsRecords[:0]
}

func (bb *BuilderBase) withArgs(args []interface{}) {
	bb.resetArgs()
	bb.argsRaw = args
	bb.isWithInterfaces = true
}

func (bb *BuilderBase) withArguments(args Arguments) {
	bb.resetArgs()
	bb.argsArgs = args
	bb.isWithInterfaces = false
}

func (bb *BuilderBase) withRecords(records []QualifiedRecord) {
	bb.resetArgs()
	bb.argsRecords = records
	bb.isWithInterfaces = false
}

// buildToSQL builds the raw SQL string and caches it as a byte slice. It gets
// called by toSQL.
func (bb *BuilderBase) buildToSQL(qb queryBuilder) ([]byte, error) {
	if bb.ärgErr != nil {
		return nil, errors.WithStack(bb.ärgErr)
	}
	rawSQL := qb.readBuildCache()
	if rawSQL == nil || bb.IsBuildCacheDisabled {
		bb.qualifiedColumns = bb.qualifiedColumns[:0]
		// Pre allocating that with a decent size, can speed up writing due to
		// less re-slicing / buffer.Grow.
		size := bb.EstimatedCachedSQLSize
		if size == 0 {
			size = estimatedCachedSQLSize
		}
		buf := bytes.NewBuffer(make([]byte, 0, size))
		var err error
		bb.qualifiedColumns, err = qb.toSQL(buf, bb.qualifiedColumns)
		if err != nil {
			return nil, errors.WithStack(err)
		}
		if !bb.IsBuildCacheDisabled {
			qb.writeBuildCache(buf.Bytes())
		}
		rawSQL = buf.Bytes()
	}
	return rawSQL, nil
}

// buildArgsAndSQL generates the SQL string and its place holders. Takes care of
// caching and interpolation. It returns the string with placeholders and a
// slice of query arguments. With switched on interpolation, it only returns a
// string including the stringyfied arguments. With an enabled cache, the
// arguments gets regenerated each time a call to ToSQL happens.
func (bb *BuilderBase) buildArgsAndSQL(qb queryBuilder) (string, []interface{}, error) {
	rawSQL, err := bb.buildToSQL(qb)
	if err != nil {
		return "", nil, errors.WithStack(err)
	}

	args, err := bb.convertRecordsToArguments()
	if err != nil {
		return "", nil, errors.WithStack(err)
	}

	if bb.IsExpandPlaceHolders {
		if phCount := bytes.Count(rawSQL, placeHolderBytes); phCount < args.Len() {
			var buf bytes.Buffer
			if err = expandPlaceHolders(&buf, rawSQL, args); err != nil {
				return "", nil, errors.WithStack(err)
			}
			qb.writeBuildCache(buf.Bytes())
			rawSQL = buf.Bytes()
			bb.IsExpandPlaceHolders = false
		}
	}

	if lArgs := len(args); bb.IsInterpolate && lArgs > 0 {
		buf := bufferpool.Get()
		err = writeInterpolate(buf, rawSQL, args)
		s := buf.String()
		bufferpool.Put(buf)
		return s, nil, errors.WithStack(err)
	} else if bb.IsInterpolate && lArgs == 0 && len(bb.argsRaw) > 0 {
		return "", nil, errors.NotAllowed.Newf("[dml] Interpolation does only work with an Arguments slice, but you provided an interface slice: %#v", bb.argsRaw)
	}

	if !bb.isWithInterfaces {
		bb.argsRaw = bb.argsRaw[:0]
	}
	bb.argsRaw = append(bb.argsRaw, args.Interfaces()...) // TODO optimize
	return string(rawSQL), bb.argsRaw, errors.WithStack(err)
}

// BuilderConditional defines base fields used in statements which can have
// conditional constraints like WHERE, JOIN, ORDER, etc. Exported for
// documentation reasons.
type BuilderConditional struct {
	Joins      Joins
	Wheres     Conditions
	OrderBys   ids
	LimitCount uint64
	LimitValid bool
}

func (b *BuilderConditional) join(j string, t id, on ...*Condition) {
	jf := &join{
		JoinType: j,
		Table:    t,
	}
	jf.On = append(jf.On, on...)
	b.Joins = append(b.Joins, jf)
}

func sqlObjToString(rawSQL []byte, err error) string {
	if err != nil {
		return fmt.Sprintf("[dml] String Error: %+v", err)
	}
	return string(rawSQL)
}

// String returns a string representing a preprocessed, interpolated, query.
// On error, the error gets printed. Fulfills interface fmt.Stringer.
func (b *Delete) String() string {
	return sqlObjToString(b.buildToSQL(b))
}

// String returns a string representing a preprocessed, interpolated, query.
// On error, the error gets printed. Fulfills interface fmt.Stringer.
func (b *Insert) String() string {
	return sqlObjToString(b.buildToSQL(b))
}

// String returns a string representing a preprocessed, interpolated, query.
// On error, the error gets printed. Fulfills interface fmt.Stringer.
func (b *Select) String() string {
	return sqlObjToString(b.buildToSQL(b))
}

// String returns a string representing a preprocessed, interpolated, query.
// On error, the error gets printed. Fulfills interface fmt.Stringer.
func (b *Update) String() string {
	return sqlObjToString(b.buildToSQL(b))
}

// String returns a string representing a preprocessed, interpolated, query.
// On error, the error gets printed. Fulfills interface fmt.Stringer.
func (u *Union) String() string {
	return sqlObjToString(u.buildToSQL(u))
}

// String returns a string representing a preprocessed, interpolated, query.
// On error, the error gets printed. Fulfills interface fmt.Stringer.
func (b *With) String() string {
	return sqlObjToString(b.buildToSQL(b))
}

// String returns a string representing a preprocessed, interpolated, query.
// On error, the error gets printed. Fulfills interface fmt.Stringer.
func (b *Show) String() string {
	return sqlObjToString(b.buildToSQL(b))
}

func sqlWriteUnionAll(w *bytes.Buffer, isAll bool, isIntersect bool, isExcept bool) {
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

func sqlWriteOrderBy(w *bytes.Buffer, orderBys ids, br bool) {
	if len(orderBys) == 0 {
		return
	}
	brS := ' '
	if br {
		brS = '\n'
	}
	w.WriteRune(brS)
	w.WriteString("ORDER BY ")
	orderBys.writeQuoted(w, nil)
}

func sqlWriteLimitOffset(w *bytes.Buffer, limitValid bool, limitCount uint64, offsetValid bool, offsetCount uint64) {
	if limitValid {
		w.WriteString(" LIMIT ")
		writeUint64(w, limitCount)
	}
	if offsetValid {
		w.WriteString(" OFFSET ")
		writeUint64(w, offsetCount)
	}
}

func writeFloat64(w *bytes.Buffer, f float64) (err error) {
	d := w.Bytes()
	w.Reset()
	_, err = w.Write(strconv.AppendFloat(d, f, 'g', -1, 64))
	return err
}

func writeInt64(w *bytes.Buffer, i int64) (err error) {
	d := w.Bytes()
	w.Reset()
	_, err = w.Write(strconv.AppendInt(d, i, 10))
	return err
}

func writeUint64(w *bytes.Buffer, i uint64) (err error) {
	d := w.Bytes()
	w.Reset()
	_, err = w.Write(strconv.AppendUint(d, i, 10))
	return err
}

func writeBytes(w *bytes.Buffer, p []byte) (err error) {
	switch {
	case p == nil:
		_, err = w.WriteString(sqlStrNullUC)
	case !utf8.Valid(p):
		dialect.EscapeBinary(w, p)
	default:
		dialect.EscapeString(w, string(p)) // maybe create an EscapeByteString version to avoid one alloc ;-)
	}
	return
}

func writeStmtID(w *bytes.Buffer, id string) {
	if id != "" {
		w.WriteString("/*ID:")
		w.WriteString(id)
		w.WriteString("*/ ")
	}
}
