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
	"database/sql"
	"database/sql/driver"
	"encoding"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/corestoreio/errors"
	"github.com/gogo/protobuf/proto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	_ fmt.GoStringer             = (*NullUint64)(nil)
	_ fmt.Stringer               = (*NullUint64)(nil)
	_ json.Marshaler             = (*NullUint64)(nil)
	_ json.Unmarshaler           = (*NullUint64)(nil)
	_ encoding.BinaryMarshaler   = (*NullUint64)(nil)
	_ encoding.BinaryUnmarshaler = (*NullUint64)(nil)
	_ encoding.TextMarshaler     = (*NullUint64)(nil)
	_ encoding.TextUnmarshaler   = (*NullUint64)(nil)
	_ gob.GobEncoder             = (*NullUint64)(nil)
	_ gob.GobDecoder             = (*NullUint64)(nil)
	_ driver.Valuer              = (*NullUint64)(nil)
	_ proto.Marshaler            = (*NullUint64)(nil)
	_ proto.Unmarshaler          = (*NullUint64)(nil)
	_ proto.Sizer                = (*NullUint64)(nil)
	_ protoMarshalToer           = (*NullUint64)(nil)
	_ sql.Scanner                = (*NullUint64)(nil)
)

func TestMakeNullUint64(t *testing.T) {

}

func TestNullUint64_Scan(t *testing.T) {
	t.Parallel()

	t.Run("nil", func(t *testing.T) {
		var nv NullUint64
		require.NoError(t, nv.Scan(nil))
		assert.Exactly(t, NullUint64{}, nv)
	})
	t.Run("[]byte", func(t *testing.T) {
		var nv NullUint64
		require.NoError(t, nv.Scan([]byte(`12345678910`)))
		assert.Exactly(t, MakeNullUint64(12345678910), nv)
	})
	t.Run("string unsupported", func(t *testing.T) {
		var nv NullUint64
		err := nv.Scan(`1234567`)
		assert.True(t, errors.Is(err, errors.NotSupported), "Error behaviour should be errors.NotSupported")
		assert.Exactly(t, NullUint64{}, nv)
	})
	t.Run("parse error negative", func(t *testing.T) {
		var nv NullUint64
		err := nv.Scan([]byte(`-1234567`))
		assert.EqualError(t, err, `strconv.ParseUint: parsing "-1234567": invalid syntax`)
		assert.Exactly(t, NullUint64{}, nv)
	})
}
