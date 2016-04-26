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
	"errors"
	"testing"
)

func TestTypeErrorBehaviour(t *testing.T) {
	const errHellFa Fatal = "Hell"
	const errHellE AlreadyExists = "Hell"
	const errHellC AlreadyClosed = "Hell"
	const errHellF NotFound = "Hell"
	const errHellS NotSupported = "Hell"
	const errHellV NotValid = "Hell"
	const errHellT Temporary = "Hell"
	const errHellTo Timeout = "Hell"
	const errHellU Unauthorized = "Hell"
	const errHellUnf UserNotFound = "Hell"
	const errHellNi NotImplemented = "Hell"
	const errHellEm Empty = "Hell"
	tests := []struct {
		err   error
		check func(error) bool
		want  bool
	}{
		{errHellEm, IsEmpty, true},
		{errHellEm, IsNotFound, false},
		{errors.New("Paradise"), IsEmpty, false},
		{nil, IsEmpty, false},

		{errHellNi, IsNotImplemented, true},
		{errHellNi, IsNotFound, false},
		{errors.New("Paradise"), IsNotImplemented, false},
		{nil, IsNotImplemented, false},

		{errHellFa, IsFatal, true},
		{errHellFa, IsNotFound, false},
		{errors.New("Paradise"), IsFatal, false},
		{nil, IsFatal, false},

		{errHellE, IsAlreadyExists, true},
		{errHellE, IsNotFound, false},
		{errors.New("Paradise"), IsAlreadyExists, false},
		{nil, IsAlreadyExists, false},

		{errHellC, IsAlreadyClosed, true},
		{errHellC, IsNotFound, false},
		{errors.New("Paradise"), IsAlreadyClosed, false},
		{nil, IsAlreadyClosed, false},

		{errHellF, IsNotFound, true},
		{errHellF, IsAlreadyClosed, false},
		{errors.New("Paradise"), IsNotFound, false},
		{nil, IsNotFound, false},

		{errHellS, IsNotSupported, true},
		{errHellS, IsAlreadyClosed, false},
		{errors.New("Paradise"), IsNotSupported, false},
		{nil, IsNotSupported, false},

		{errHellV, IsNotValid, true},
		{errHellV, IsAlreadyClosed, false},
		{errors.New("Paradise"), IsNotValid, false},
		{nil, IsNotValid, false},

		{errHellT, IsTemporary, true},
		{errHellT, IsAlreadyClosed, false},
		{errors.New("Paradise"), IsTemporary, false},
		{nil, IsTemporary, false},

		{errHellTo, IsTimeout, true},
		{errHellTo, IsAlreadyClosed, false},
		{errors.New("Paradise"), IsTimeout, false},
		{nil, IsTimeout, false},

		{errHellU, IsUnauthorized, true},
		{errHellU, IsAlreadyClosed, false},
		{errors.New("Paradise"), IsUnauthorized, false},
		{nil, IsUnauthorized, false},

		{errHellUnf, IsUserNotFound, true},
		{errHellUnf, IsAlreadyClosed, false},
		{errors.New("Paradise"), IsUserNotFound, false},
		{nil, IsUserNotFound, false},
	}
	for i, test := range tests {
		if have, want := test.check(test.err), test.want; have != want {
			t.Errorf("(%02d) Error: %q => Have %t Want %t", i, test.err, have, want)
		}
		if test.err != nil && test.err.Error() == "" {
			t.Errorf("(%02d) Error: %q => Missing error string", i, test.err)
		}
	}
}
