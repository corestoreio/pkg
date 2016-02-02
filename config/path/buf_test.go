// Copyright 2015-2016, Cyrill @ Schumacher.fm and the CoreStore contributors
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

package path

import (
	"bytes"
	"testing"

	"strings"

	"github.com/stretchr/testify/assert"
)

func TestPathBuffer(t *testing.T) {
	t.Parallel()
	tests := []struct {
		haveB   byte
		haveBs  []byte
		haveS   string
		wantErr error
		want    string
	}{
		{'a', nil, "", nil, "a"},
		{0, []byte(`aa/bb`), "", nil, "aa/bb"},
		{0, nil, "cc/dd", nil, "cc/dd"},
		{'/', []byte(`rr/ss`), "/tt/uu", nil, "/rr/ss/tt/uu"},
		{0, bytes.Repeat([]byte(`rr/ss`), 150), "/tt/uu", ErrPathTooLong, "/rr/ss/tt/uu"},
		{0, bytes.Repeat([]byte(`rr/ss`), 2), strings.Repeat("/tt/uu", 200), ErrPathTooLong, "/rr/ss/tt/uu"},
	}
	for i, test := range tests {
		buf := newPathBuf()

		haveErr := buf.WriteByte(test.haveB)
		if pathBufferCheckError(t, haveErr, test.wantErr, i) {
			continue
		}
		_, haveErr = buf.Write(test.haveBs)
		if pathBufferCheckError(t, haveErr, test.wantErr, i) {
			continue
		}
		_, haveErr = buf.WriteString(test.haveS)
		if pathBufferCheckError(t, haveErr, test.wantErr, i) {
			continue
		}

		assert.Exactly(t, test.want, string(buf.Bytes()), "Index %d", i)
	}

	buf := newPathBuf()
	_, err := buf.WriteString(strings.Repeat("s", 256))
	assert.NoError(t, err)
	err = buf.WriteByte('t')
	assert.EqualError(t, err, ErrPathTooLong.Error())
	assert.Exactly(t, strings.Repeat("s", 256), string(buf.Bytes()))

}

func pathBufferCheckError(t *testing.T, haveErr, wantErr error, idx int) bool {
	if haveErr != nil && wantErr != nil {
		assert.EqualError(t, haveErr, wantErr.Error(), "Index %d", idx)
		return true
	}
	assert.NoError(t, haveErr, "Index %d", idx)
	return false
}
