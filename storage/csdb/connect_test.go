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

package csdb

import (
	"errors"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetDSN(t *testing.T) {

	tests := []struct {
		env        string
		envContent string
		err        error
		returnErr  bool
	}{
		{
			env:        "TEST_CS_1",
			envContent: "Hello",
			err:        errors.New("World"),
			returnErr:  false,
		},
	}

	for _, test := range tests {
		os.Setenv(test.env, test.envContent)
		s, aErr := getDSN(test.env, test.err)
		assert.Equal(t, test.envContent, s)
		assert.NoError(t, aErr)

		s, aErr = getDSN(test.env+"NOTFOUND", test.err)
		assert.Equal(t, "", s)
		assert.Error(t, aErr)
		assert.Equal(t, test.err, aErr)
	}
}
