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
	"database/sql/driver"
	"fmt"
	"math"
	"testing"
	"time"

	"github.com/corestoreio/errors"
	"github.com/corestoreio/pkg/storage/null"
	"github.com/corestoreio/pkg/util/assert"
)

type driverValueBytes []byte

// Value implements the driver.Valuer interface.
func (a driverValueBytes) Value() (driver.Value, error) {
	return []byte(a), nil
}

type driverValueNotSupported uint8

// Value implements the driver.Valuer interface.
func (a driverValueNotSupported) Value() (driver.Value, error) {
	return uint8(a), nil
}

type driverValueNil uint8

// Value implements the driver.Valuer interface.
func (a driverValueNil) Value() (driver.Value, error) {
	return nil, nil
}

type driverValueError uint8

// Value implements the driver.Valuer interface.
func (a driverValueError) Value() (driver.Value, error) {
	return nil, errors.Aborted.Newf("WE've aborted something")
}

func newDBRArgs() *DBR { return &DBR{arguments: make(arguments, 0, 10)} }

func TestArguments_Length_and_Stringer(t *testing.T) {
	t.Parallel()
	nt := now().In(time.UTC)

	t.Run("no slices, nulls valid", func(t *testing.T) {
		args := newDBRArgs().
			Null().Int(-1).Int64(1).Uint(9898).Uint64(2).Float64(3.1).Bool(true).String("eCom1").Bytes([]byte(`eCom2`)).Time(nt).
			NullString(null.MakeString("eCom3")).NullInt64(null.MakeInt64(4)).NullFloat64(null.MakeFloat64(2.7)).
			NullBool(null.MakeBool(true)).NullTime(null.MakeTime(nt))
		assert.Exactly(t, 15, args.Len(), "Length mismatch")

		assert.Exactly(t,
			fmt.Sprint([]interface{}{nil, -1, 1, 9898, 2, 3.1, true, "eCom1", []byte("eCom2"), nt, "eCom3", 4, 2.7, true, nt}),
			fmt.Sprint(args.toInterfaces()))
	})

	t.Run("no slices, nulls invalid", func(t *testing.T) {
		args := newDBRArgs().
			Null().Int(-1).Int64(1).Uint64(2).Float64(3.1).Bool(true).String("eCom1").Bytes([]byte(`eCom2`)).Time(nt).
			NullString(null.String{}).NullInt64(null.Int64{}).NullFloat64(null.Float64{}).
			NullBool(null.Bool{}).NullTime(null.Time{})
		assert.Exactly(t, 14, args.Len(), "Length mismatch")
		assert.Exactly(t,
			fmt.Sprint([]interface{}{
				nil, -1, int64(1), uint64(2), 3.1, true, "eCom1", []byte("eCom2"), nt, nil, nil, nil, nil, nil,
			}),
			fmt.Sprint(args.toInterfaces()),
		)
	})

	t.Run("slices, nulls valid", func(t *testing.T) {
		args := newDBRArgs().
			Null().Int(-1).Int64s(1, 2).Uints(567, 765).Uint64s(2).Float64s(1.2, 3.1).Bools(false, true).Strings("eCom1", "eCom11").BytesSlice(nil, []byte(`eCom2`)).Times(nt, nt).
			NullStrings(null.MakeString("eCom3"), null.MakeString("eCom3")).NullInt64s(null.MakeInt64(4), null.MakeInt64(4)).NullFloat64s(null.MakeFloat64(2.7), null.MakeFloat64(2.7)).
			NullBools(null.MakeBool(true)).NullTimes(null.MakeTime(nt), null.MakeTime(nt))
		assert.Exactly(t, 26, args.Len(), "Length mismatch")
		assert.Exactly(t,
			fmt.Sprint([]interface{}{
				nil, -1, 1, 2, 567, 765, 2, 1.2, 3.1, false, true, "eCom1", "eCom11", []byte(nil), []byte("eCom2"),
				nt,
				nt,
				"eCom3", "eCom3", 4, 4, 2.7, 2.7, true,
				nt,
				nt,
			}),
			fmt.Sprint(args.toInterfaces()),
		)
	})
}

func TestIFaceToArgs(t *testing.T) {
	t.Parallel()
	t.Run("not supported", func(t *testing.T) {
		_, err := iFaceToArgs(arguments{}, time.Minute)
		assert.ErrorIsKind(t, errors.NotSupported, err)
	})
	t.Run("all types", func(t *testing.T) {
		nt := now()
		args, err := iFaceToArgs(arguments{},
			float32(2.3), float64(2.2),
			int64(5), int(6), int32(7), int16(8), int8(9),
			uint32(math.MaxUint32), uint16(math.MaxUint16), uint8(math.MaxUint8),
			true, "Gopher", []byte(`Hello`),
			now(), &nt, nil,
		)
		assert.NoError(t, err)

		assert.Exactly(t, []interface{}{
			float64(2.299999952316284), float64(2.2),
			int64(5), int64(6), int64(7), int64(8), int64(9),
			int64(math.MaxUint32), int64(math.MaxUint16), int64(math.MaxUint8),
			true, "Gopher", []uint8{0x48, 0x65, 0x6c, 0x6c, 0x6f},
			now(), now(), nil,
		}, args.toInterfaces())
	})
}

func TestArguments_Clone(t *testing.T) {
	t.Parallel()

	args := newDBRArgs().Int64(1).String("S1").arguments
	args2 := args.clone()
	args2[0].value = int(2)
	args2[1].value = "S1a"

	assert.Exactly(t, []interface{}{int64(1), "S1"}, args.toInterfaces())
	assert.Exactly(t, []interface{}{int64(2), "S1a"}, args2.toInterfaces())
}
