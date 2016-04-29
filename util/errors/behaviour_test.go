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
	"strings"
	"testing"
)

func TestEB(t *testing.T) {

	teb := eb{
		message: "CoreStore",
	}
	if have, want := teb.Error(), "CoreStore"; have != want {
		t.Errorf("Have %q Want %q", have, want)
	}
}

type testBehave struct{ ret bool }

func (nf testBehave) Fatal() bool {
	return nf.ret
}
func (nf testBehave) NotImplemented() bool {
	return nf.ret
}
func (nf testBehave) Empty() bool {
	return nf.ret
}
func (nf testBehave) WriteFailed() bool {
	return nf.ret
}
func (nf testBehave) NotFound() bool {
	return nf.ret
}
func (nf testBehave) UserNotFound() bool {
	return nf.ret
}
func (nf testBehave) Unauthorized() bool {
	return nf.ret
}
func (nf testBehave) AlreadyExists() bool {
	return nf.ret
}
func (nf testBehave) AlreadyClosed() bool {
	return nf.ret
}
func (nf testBehave) NotSupported() bool {
	return nf.ret
}
func (nf testBehave) NotValid() bool {
	return nf.ret
}
func (nf testBehave) Temporary() bool {
	return nf.ret
}
func (nf testBehave) Timeout() bool {
	return nf.ret
}
func (nf testBehave) Error() string {
	return ""
}

func TestBehaviour(t *testing.T) {
	tests := []struct {
		err  error
		is   BehaviourFunc
		want bool
	}{
		{
			err:  errors.New("Error1"),
			is:   IsEmpty,
			want: false,
		}, {
			err:  NewEmpty(nil, "Error2"),
			is:   IsEmpty,
			want: false,
		}, {
			err:  NewEmpty(Error("Error2a"), "Error2"),
			is:   IsEmpty,
			want: true,
		}, {
			err:  nil,
			is:   IsEmpty,
			want: false,
		}, {
			err:  testBehave{},
			is:   IsEmpty,
			want: false,
		},

		{
			err:  errors.New("Error1"),
			is:   IsWriteFailed,
			want: false,
		}, {
			err:  NewWriteFailed(nil, "Error2"),
			is:   IsWriteFailed,
			want: false,
		}, {
			err:  NewWriteFailed(Error("Error2a"), "Error2"),
			is:   IsWriteFailed,
			want: true,
		}, {
			err:  nil,
			is:   IsWriteFailed,
			want: false,
		}, {
			err:  testBehave{},
			is:   IsWriteFailed,
			want: false,
		},

		{
			err:  errors.New("Error1"),
			is:   IsNotImplemented,
			want: false,
		}, {
			err:  NewNotImplemented(nil, "Error2"),
			is:   IsNotImplemented,
			want: false,
		}, {
			err:  NewNotImplemented(Error("Error2a"), "Error2"),
			is:   IsNotImplemented,
			want: true,
		}, {
			err:  nil,
			is:   IsNotImplemented,
			want: false,
		}, {
			err:  testBehave{},
			is:   IsNotImplemented,
			want: false,
		},

		{
			err:  errors.New("Error1"),
			is:   IsFatal,
			want: false,
		}, {
			err:  NewFatal(nil, "Error2"),
			is:   IsFatal,
			want: false,
		}, {
			err:  NewFatal(Error("Error2a"), "Error2"),
			is:   IsFatal,
			want: true,
		}, {
			err:  nil,
			is:   IsFatal,
			want: false,
		}, {
			err:  testBehave{},
			is:   IsFatal,
			want: false,
		},

		{
			err:  errors.New("Error1"),
			is:   IsNotFound,
			want: false,
		}, {
			err:  NewNotFound(nil, "Error2"),
			is:   IsNotFound,
			want: false,
		}, {
			err:  NewNotFound(Error("Error2a"), "Error2"),
			is:   IsNotFound,
			want: true,
		}, {
			err:  nil,
			is:   IsNotFound,
			want: false,
		}, {
			err:  testBehave{},
			is:   IsNotFound,
			want: false,
		},

		{
			err:  testBehave{true},
			is:   IsUserNotFound,
			want: true,
		}, {
			err:  errors.New("Error1"),
			is:   IsUserNotFound,
			want: false,
		}, {
			err:  NewUserNotFound(nil, "Error2"),
			is:   IsUserNotFound,
			want: false,
		}, {
			err:  NewUserNotFound(Error("Error2a"), "Error2"),
			is:   IsUserNotFound,
			want: true,
		}, {
			err:  nil,
			is:   IsUserNotFound,
			want: false,
		}, {
			err:  testBehave{},
			is:   IsUserNotFound,
			want: false,
		}, {
			err:  testBehave{true},
			is:   IsUserNotFound,
			want: true,
		},

		{
			err:  testBehave{true},
			is:   IsUnauthorized,
			want: true,
		}, {
			err:  errors.New("Error1"),
			is:   IsUnauthorized,
			want: false,
		}, {
			err:  NewUnauthorized(nil, "Error2"),
			is:   IsUnauthorized,
			want: false,
		}, {
			err:  NewUnauthorized(Error("Error2a"), "Error2"),
			is:   IsUnauthorized,
			want: true,
		}, {
			err:  nil,
			is:   IsUnauthorized,
			want: false,
		}, {
			err:  testBehave{},
			is:   IsUnauthorized,
			want: false,
		}, {
			err:  testBehave{true},
			is:   IsUnauthorized,
			want: true,
		},

		{
			err:  testBehave{true},
			is:   IsAlreadyExists,
			want: true,
		}, {
			err:  errors.New("Error1"),
			is:   IsAlreadyExists,
			want: false,
		}, {
			err:  NewAlreadyExists(nil, "Error2"),
			is:   IsAlreadyExists,
			want: false,
		}, {
			err:  NewAlreadyExists(Error("Error2a"), "Error2"),
			is:   IsAlreadyExists,
			want: true,
		}, {
			err:  nil,
			is:   IsAlreadyExists,
			want: false,
		}, {
			err:  testBehave{},
			is:   IsAlreadyExists,
			want: false,
		}, {
			err:  testBehave{true},
			is:   IsAlreadyExists,
			want: true,
		},

		{
			err:  testBehave{true},
			is:   IsAlreadyClosed,
			want: true,
		}, {
			err:  errors.New("Error1"),
			is:   IsAlreadyClosed,
			want: false,
		}, {
			err:  NewAlreadyClosed(nil, "Error2"),
			is:   IsAlreadyClosed,
			want: false,
		}, {
			err:  NewAlreadyClosed(Error("Error2a"), "Error2"),
			is:   IsAlreadyClosed,
			want: true,
		}, {
			err:  nil,
			is:   IsAlreadyClosed,
			want: false,
		}, {
			err:  testBehave{},
			is:   IsAlreadyClosed,
			want: false,
		}, {
			err:  testBehave{true},
			is:   IsAlreadyClosed,
			want: true,
		},

		{
			err:  testBehave{true},
			is:   IsNotSupported,
			want: true,
		}, {
			err:  errors.New("Error1"),
			is:   IsNotSupported,
			want: false,
		}, {
			err:  NewNotSupported(nil, "Error2"),
			is:   IsNotSupported,
			want: false,
		}, {
			err:  NewNotSupported(Error("Error2a"), "Error2"),
			is:   IsNotSupported,
			want: true,
		}, {
			err:  nil,
			is:   IsNotSupported,
			want: false,
		}, {
			err:  testBehave{},
			is:   IsNotSupported,
			want: false,
		}, {
			err:  testBehave{true},
			is:   IsNotSupported,
			want: true,
		},

		{
			err:  testBehave{true},
			is:   IsNotValid,
			want: true,
		}, {
			err:  errors.New("Error1"),
			is:   IsNotValid,
			want: false,
		}, {
			err:  NewNotValid(nil, "Error2"),
			is:   IsNotValid,
			want: false,
		}, {
			err:  NewNotValid(Error("Error2a"), "Error2"),
			is:   IsNotValid,
			want: true,
		}, {
			err:  nil,
			is:   IsNotValid,
			want: false,
		}, {
			err:  testBehave{},
			is:   IsNotValid,
			want: false,
		}, {
			err:  testBehave{true},
			is:   IsNotValid,
			want: true,
		},

		{
			err:  testBehave{true},
			is:   IsTemporary,
			want: true,
		}, {
			err:  errors.New("Error1"),
			is:   IsTemporary,
			want: false,
		}, {
			err:  NewTemporary(nil, "Error2"),
			is:   IsTemporary,
			want: false,
		}, {
			err:  NewTemporary(Error("Error2a"), "Error2"),
			is:   IsTemporary,
			want: true,
		}, {
			err:  nil,
			is:   IsTemporary,
			want: false,
		}, {
			err:  testBehave{},
			is:   IsTemporary,
			want: false,
		}, {
			err:  testBehave{true},
			is:   IsTemporary,
			want: true,
		},

		{
			err:  testBehave{true},
			is:   IsTimeout,
			want: true,
		}, {
			err:  errors.New("Error1"),
			is:   IsTimeout,
			want: false,
		}, {
			err:  NewTimeout(nil, "Error2"),
			is:   IsTimeout,
			want: false,
		}, {
			err:  NewTimeout(Error("Error2a"), "Error2"),
			is:   IsTimeout,
			want: true,
		}, {
			err:  nil,
			is:   IsTimeout,
			want: false,
		}, {
			err:  testBehave{},
			is:   IsTimeout,
			want: false,
		}, {
			err:  testBehave{true},
			is:   IsTimeout,
			want: true,
		},
	}
	for i, test := range tests {
		if want, have := test.want, test.is(test.err); want != have {
			t.Errorf("Index %d: Want %t Have %t", i, want, have)
		}
	}
}

func TestBehaviourF(t *testing.T) {
	tests := []struct {
		constErr error
		errf     func(format string, args ...interface{}) error
		is       BehaviourFunc
		want     bool
	}{
		{notImplementedTxt, NewNotImplementedf, IsNotImplemented, true},
		{fatalTxt, NewFatalf, IsFatal, true},
		{notFoundTxt, NewNotFoundf, IsNotFound, true},
		{userNotFoundTxt, NewUserNotFoundf, IsUserNotFound, true},
		{unauthorizedTxt, NewUnauthorizedf, IsUnauthorized, true},
		{alreadyExistsTxt, NewAlreadyExistsf, IsAlreadyExists, true},
		{alreadyClosedTxt, NewAlreadyClosedf, IsAlreadyClosed, true},
		{notSupportedTxt, NewNotSupportedf, IsNotSupported, true},
		{notValidTxt, NewNotValidf, IsNotValid, true},
		{temporaryTxt, NewTemporaryf, IsTemporary, true},
		{timeoutTxt, NewTimeoutf, IsTimeout, true},
	}
	const substrLocation = `github.com/corestoreio/csfw/util/errors/behaviour_test.go`
	for i, test := range tests {
		haveErr := test.errf("Gopher %d", i)
		if want, have := test.want, test.is(haveErr); want != have {
			t.Errorf("Index %d: Want %t Have %t", i, want, have)
		}
		loca := PrintLoc(haveErr)
		if !strings.Contains(loca, substrLocation) {
			t.Errorf("Index %d: Cannot find %q in %q", i, substrLocation, loca)
		}
		if substr := test.constErr.Error(); !strings.Contains(loca, substr) {
			t.Errorf("Index %d: Cannot find %q in %q", i, substr, loca)
		}
	}
}

var benchmarkAsserted bool

func BenchmarkAssertBehaviourInterface(b *testing.B) {
	const hell AlreadyExists = "Hell"
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		benchmarkAsserted = IsAlreadyExists(hell)
		if !benchmarkAsserted {
			b.Error("Hell should already exists.")
		}
	}
}

func BenchmarkAssertBehaviourPointer(b *testing.B) {
	var hell = NewAlreadyExists(errors.New("Hell"), "There is already a place for you")
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		benchmarkAsserted = IsAlreadyExists(hell)
		if !benchmarkAsserted {
			b.Error("Hell should already exists.")
		}
	}
}

func BenchmarkAssertBehaviourNoMatch(b *testing.B) {
	var hell = NewAlreadyExists(errors.New("Hell"), "There is already a place for you")
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		benchmarkAsserted = IsAlreadyClosed(hell)
		if benchmarkAsserted {
			b.Error("Hell should already be clsoed.")
		}
	}
}
