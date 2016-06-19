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

package errors

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestError_Error(t *testing.T) {
	const e1 Error = "e1"
	assert.EqualError(t, e1, "e1")
}

func TestWrapf2(t *testing.T) {
	var e = Wrapf(nil, "Error %d")
	assert.Nil(t, e)
}

func TestErrorContainsAny(t *testing.T) {
	tests := []struct {
		me   error
		vf   []BehaviourFunc
		want bool
	}{
		{NotFound("e0"), []BehaviourFunc{IsNotFound}, true},
		{NotFound("e1"), []BehaviourFunc{IsNotValid}, false},
		{NotFound("e2"), []BehaviourFunc{IsNotValid, IsNotFound}, true},
		{NewNotFound(NewNotValidf("NotValid inner"), "NotFound outer"), []BehaviourFunc{IsNotValid, IsNotFound}, true},
		// once ErrorContainsAny acts recursive the next line will switch to true
		{NewNotFound(NewNotValidf("NotValid inner"), "NotFound outer"), []BehaviourFunc{IsNotValid}, false},
		{nil, []BehaviourFunc{IsNotValid}, false},
		{nil, nil, false},
	}

	for i, test := range tests {
		if have, want := ErrorContainsAny(test.me, test.vf...), test.want; have != want {
			t.Errorf("Index %d: Have %t Want %t", i, have, want)
		}
	}
}
