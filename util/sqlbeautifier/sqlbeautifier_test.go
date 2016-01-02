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

package sqlbeautifier_test

import (
	"testing"

	"bytes"
	"fmt"
	"io/ioutil"

	"go/build"
	"path/filepath"

	"github.com/corestoreio/csfw/util/sqlbeautifier"
	"github.com/stretchr/testify/assert"
)

func TestFromReader(t *testing.T) {
	tests := []struct {
		idx     int
		wantErr error
	}{
		{0, nil},
		{1, nil},
		{2, nil},
		{3, nil},
		{4, nil},
	}
	for _, test := range tests {

		haveF := getTestFile(t, test.idx, "have")
		wantF := getTestFile(t, test.idx, "want")

		bufFormat, haveErr := sqlbeautifier.FromReader(bytes.NewReader(haveF))
		if test.wantErr != nil {
			assert.Error(t, haveErr, "Index %d", test.idx)
			assert.EqualError(t, haveErr, test.wantErr.Error(), "Index %d", test.idx)
			continue
		}

		if bytes.Compare(wantF, bufFormat.Bytes()) != 0 {
			assert.NoError(t, ioutil.WriteFile(fmt.Sprintf("%sformat_%02d.sql", pathPrefix, test.idx), bufFormat.Bytes(), 0644))
			t.Errorf("Error @ Index %d\n%s\n", test.idx, bufFormat.String())
		}
	}
}

var pathPrefix = filepath.Join(build.Default.GOPATH, "src", "github.com", "corestoreio", "csfw", "util", "sqlbeautifier", "test_")

func getTestFile(t *testing.T, i int, test string) []byte {
	n := fmt.Sprintf("%s%02d_%s.sql", pathPrefix, i, test)
	b, err := ioutil.ReadFile(n)
	if err != nil {
		t.Fatal(err)
	}
	return b
}
