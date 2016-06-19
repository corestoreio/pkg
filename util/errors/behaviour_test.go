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
	"fmt"
	"strings"
	"testing"
)

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

func TestBehaviourPlain(t *testing.T) {
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
			err:  Wrap(Empty("Err88"), "Wrap88"),
			is:   IsEmpty,
			want: true,
		}, {
			err:  NewEmptyf("Error3"),
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

		{ // 8
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
			err:  NewWriteFailedf("Error118"),
			is:   IsWriteFailed,
			want: true,
		}, {
			err:  Wrap(WriteFailed("Error122"), "Wrap122"),
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

		{ // 15
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
			err:  NewNotImplementedf("err152"),
			is:   IsNotImplemented,
			want: true,
		}, {
			err:  Wrap(NotImplemented("err156"), "Wrap156"),
			is:   IsNotImplemented,
			want: true,
		}, {
			err:  testBehave{},
			is:   IsNotImplemented,
			want: false,
		},

		{ // 22
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
			err:  NewFatalf("Err178"),
			is:   IsFatal,
			want: true,
		}, {
			err:  Wrap(Fatal("Err182"), "Wrap182"),
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

		{ // 29
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
			err:  NewNotFoundf("Err208"),
			is:   IsNotFound,
			want: true,
		}, {
			err:  Wrap(NotFound("Err212"), "Wrap212"),
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

		{ // 35
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
			err:  NewUserNotFoundf("Err242"),
			is:   IsUserNotFound,
			want: true,
		}, {
			err:  Wrap(UserNotFound("Err246"), "Wrap246"),
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

		{ // 44
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
			err:  NewUnauthorizedf("Err280"),
			is:   IsUnauthorized,
			want: true,
		}, {
			err:  Wrap(Unauthorized("Err284"), "Wrap284"),
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
			err:  NewAlreadyExistsf("Err318"),
			is:   IsAlreadyExists,
			want: true,
		}, {
			err:  Wrap(AlreadyExists("Err322"), "Wrap322"),
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
			err:  NewAlreadyClosedf("Err356"),
			is:   IsAlreadyClosed,
			want: true,
		}, {
			err:  Wrap(AlreadyClosed("Err360"), "Wrap360"),
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
			err:  NewNotSupportedf("Err394"),
			is:   IsNotSupported,
			want: true,
		}, {
			err:  Wrap(NotSupported("Err398"), "Wrap398"),
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
			err:  NewNotValidf("Err432"),
			is:   IsNotValid,
			want: true,
		}, {
			err:  Wrap(NotValid("Err436"), "Wrap436"),
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
			err:  NewTemporaryf("Err470"),
			is:   IsTemporary,
			want: true,
		}, {
			err:  Wrap(Temporary("Err474"), "Wrap474"),
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
			err:  NewTimeoutf("Err508"),
			is:   IsTimeout,
			want: true,
		}, {
			err:  Wrap(Timeout("Err512"), "Wrap512"),
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

func TestBehaviourFormat(t *testing.T) {
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
	// const substrLocation = `github.com/corestoreio/csfw/util/errors/behaviour.go`
	const substrLocation = `/behaviour.go`
	for i, test := range tests {
		haveErr := test.errf("Gopher %d", i)
		if want, have := test.want, test.is(haveErr); want != have {
			t.Errorf("Index %d: Want %t Have %t", i, want, have)
		}
		loca := fmt.Sprintf("%+v", haveErr)
		if !strings.Contains(loca, substrLocation) {
			t.Errorf("Index %d: Cannot find %q in %q", i, substrLocation, loca)
		}
		if substr := test.constErr.Error(); !strings.Contains(loca, substr) {
			t.Errorf("Index %d: Cannot find %q in %q", i, substr, loca)
		}
	}
}

func TestEbWrapf(t *testing.T) {
	const e Error = "Error1"
	if haveEB, want := ebWrapf(e, "Hello World %#v"), "Hello World %#v"; haveEB.msg != want {
		t.Errorf("have %q want %q", haveEB.msg, want)
	}
	if haveEB, want := ebWrapf(e, "Hello World %d", 123), "Hello World 123"; haveEB.msg != want {
		t.Errorf("have %q want %q", haveEB.msg, want)
	}
}
