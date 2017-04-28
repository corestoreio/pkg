package dbr

import (
	"strings"
	"unicode/utf8"

	"github.com/corestoreio/csfw/util/bufferpool"
	"github.com/corestoreio/errors"
)

// Repeat takes a SQL string and repeats the question marks with the provided
// arguments. If the amount of arguments does not match the number of questions
// marks, a Mismatch error gets returned. The arguments are getting converted to
// an interface slice to easy passing into the db.Query/db.Exec/etc functions at
// an argument.
//		Repeat("SELECT * FROM table WHERE id IN (?) AND status IN (?)", ArgInt(myIntSlice...), ArgString(myStrSlice...))
// Gets converted to:
//		SELECT * FROM table WHERE id IN (?,?) AND status IN (?,?,?)
// The questions marks are of course depending on the values in the Arg*
// functions. This function should be generally used when dealing with prepared
// statements.
func Repeat(sql string, args ...Argument) (string, []interface{}, error) {
	const qMarkStr = `?`
	const qMarkRne = '?'

	markCount := strings.Count(sql, qMarkStr)
	if want := len(args); markCount != want || want == 0 {
		return "", nil, errors.NewMismatchf("[dbr] Repeat: Number of %s:%d do not match the number of repetitions: %d", qMarkStr, markCount, want)
	}

	retArgs := make([]interface{}, 0, len(args)*2)

	buf := bufferpool.Get()
	defer bufferpool.Put(buf)

	n := markCount
	i := 0
	for i < n {
		m := strings.IndexByte(sql, qMarkRne)
		if m < 0 {
			break
		}
		buf.WriteString(sql[:m])

		if i < len(args) {
			prevLen := len(retArgs)
			retArgs = args[i].toIFace(retArgs)
			reps := len(retArgs) - prevLen
			for r := 0; r < reps; r++ {
				buf.WriteByte(qMarkRne)
				if r < reps-1 {
					buf.WriteByte(',')
				}
			}
		}
		sql = sql[m+len(qMarkStr):]
		i++
	}
	buf.WriteString(sql)
	return buf.String(), retArgs, nil
}

// Interpolate takes an SQL string with placeholders and a list of arguments to
// replace them with. It returns a blank string and error if the number of placeholders
// does not match the number of arguments.
func Interpolate(sql string, args ...Argument) (string, error) {
	// Get the number of arguments to add to this query
	if sql == "" {
		if len(args) != 0 {
			return "", errors.NewNotValidf("[dbr] Arguments are imbalanced")
		}
		return "", nil
	}

	var buf = bufferpool.Get()
	defer bufferpool.Put(buf)

	qCountTotal := 0
	qCount := -1
	argIndex := 0
	argLength := 0
	if len(args) > 0 {
		argLength = args[0].len()
	}
	pos := 0
	for pos < len(sql) {
		r, w := utf8.DecodeRuneInString(sql[pos:])
		pos += w

		switch {
		case r == '?':
			if qCount >= argLength {
				return "", errors.NewNotValidf("[dbr] Arguments are imbalanced. Placeholder count: %d Current argument count: %d", qCount, args[argIndex].len())
			}

			if qCount < argLength-1 {
				qCount++
			} else {
				qCount = 0 // next argument set starts
				argIndex++
				if argIndex >= len(args) {
					return "", errors.NewNotValidf("[dbr] Arguments are imbalanced. Argument Index %d but argument count was %d", argIndex, len(args)-1)
				}
				argLength = args[argIndex].len()
			}
			if argLength == 0 {
				return "", errors.NewEmptyf("[dbr] Empty Argument for position %d", qCountTotal+1)
			}

			if err := args[argIndex].writeTo(buf, qCount); err != nil {
				return "", errors.Wrap(err, "[dbr] Interpolate writeTo arguments")
			}

			qCountTotal++
		case r == '`', r == '\'', r == '"':
			p := strings.IndexRune(sql[pos:], r)
			if p == -1 {
				return "", errors.NewNotValidf("[dbr] Interpolate: Invalid syntax")
			}
			if r == '"' {
				r = '\''
			}
			buf.WriteRune(r)
			buf.WriteString(sql[pos : pos+p])
			buf.WriteRune(r)
			pos += p + 1
		case r == '[':
			w := strings.IndexRune(sql[pos:], ']')
			col := sql[pos : pos+w]
			dialect.EscapeIdent(buf, col)
			pos += w + 1 // size of ']'
		default:
			buf.WriteRune(r)
		}
	}

	if al := Arguments(args).len(); qCountTotal != al {
		return "", errors.NewNotValidf("[dbr] Arguments are imbalanced. Placeholders: %d Current argument count: %d or %d", qCountTotal, al, len(args))
	}
	return buf.String(), nil
}
