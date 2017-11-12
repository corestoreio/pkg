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

package util_test

import (
	"errors"
	"testing"

	"github.com/corestoreio/pkg/util"
	"github.com/stretchr/testify/assert"
)

func TestParseTimeStrict(t *testing.T) {
	t.Parallel()
	tests := []struct {
		layout  string
		value   string
		wantErr error
		want    string
	}{
		{"1/2/06", "11/31/15", errors.New("invalid time: \"11/31/15\""), ""},
		{"1/2/06", "11/30/15", nil, "2015-11-30 00:00:00 +0000 UTC"},
	}
	for _, test := range tests {

		tt, err := util.ParseTimeStrict(test.layout, test.value)
		if test.wantErr != nil {
			assert.Error(t, err, "Test %v", test)
			continue
		}
		assert.NoError(t, err, "Test %v", test)
		assert.Equal(t, test.want, tt.String(), "Test %v", test)
	}
}
