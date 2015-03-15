// Copyright 2015 CoreStore Authors
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
    "testing"
    "errors"
)

func TestGetDSN(t *testing.T) {

    tests := []struct {
        env string
        envContent string
        err error
    }{
        {"TEST_CS_1", "Hello",errors.New("World")},
    }

    for _, test := range tests {
        s,err := getDSN(test.env, test.err)
    }
