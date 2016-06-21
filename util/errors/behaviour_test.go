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
			err:  Wrap(NewEmptyf("Err88"), "Wrap88"),
			is:   IsEmpty,
			want: true,
		}, {
			err:  NewNotImplemented(Wrap(NewEmptyf("Err92"), "Wrap92"), ""),
			is:   IsEmpty,
			want: true,
		}, {
			err:  Wrap(NewEmpty(Wrap(NewNotImplementedf("Err92"), "Wrap92"), ""), ""),
			is:   IsEmpty,
			want: false,
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
			err:  Wrap(NewWriteFailedf("Error122"), "Wrap122"),
			is:   IsWriteFailed,
			want: true,
		}, {
			err:  NewNotImplemented(Wrap(NewWriteFailedf("Error130"), "Wrap130"), ""),
			is:   IsWriteFailed,
			want: true,
		}, {
			err:  Wrap(NewWriteFailed(Wrap(NewNotImplementedf("4"), "3"), "2"), "1"),
			is:   IsWriteFailed,
			want: false,
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
			err:  Wrap(NewNotImplementedf("err156"), "Wrap156"),
			is:   IsNotImplemented,
			want: true,
		}, {
			err:  NewAlreadyClosed(Wrap(NewNotImplementedf("err160"), "Wrap160"), ""),
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
			err:  Wrap(NewFatalf("Err182"), "Wrap182"),
			is:   IsFatal,
			want: true,
		}, {
			err:  NewNotImplemented(Wrap(NewFatalf("Err198"), "Wrap198"), ""),
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
			err:  Wrap(NewNotFoundf("Err212"), "Wrap212"),
			is:   IsNotFound,
			want: true,
		}, {
			err:  NewNotImplemented(Wrap(NewNotFoundf("Err232"), "Wrap232"), ""),
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
			err:  Wrap(NewUserNotFoundf("Err246"), "Wrap246"),
			is:   IsUserNotFound,
			want: true,
		}, {
			err:  NewNotImplemented(Wrap(NewUserNotFoundf("Err270"), "Wrap270"), ""),
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
			err:  Wrap(NewUnauthorizedf("Err284"), "Wrap284"),
			is:   IsUnauthorized,
			want: true,
		}, {
			err:  NewNotImplemented(Wrap(NewUnauthorizedf("Err312"), "Wrap312"), ""),
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
			err:  Wrap(NewAlreadyExistsf("Err322"), "Wrap322"),
			is:   IsAlreadyExists,
			want: true,
		}, {
			err:  NewNotImplemented(Wrap(NewAlreadyExistsf("Err354"), "Wrap354"), ""),
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
			err:  Wrap(NewAlreadyClosedf("Err360"), "Wrap360"),
			is:   IsAlreadyClosed,
			want: true,
		}, {
			err:  NewNotImplemented(Wrap(NewAlreadyClosedf("Err396"), "Wrap396"), ""),
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
			err:  Wrap(NewNotSupportedf("Err398"), "Wrap398"),
			is:   IsNotSupported,
			want: true,
		}, {
			err:  NewNotImplemented(Wrap(NewNotSupportedf("Err438"), "Wrap438"), ""),
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
			err:  Wrap(NewNotValidf("Err436"), "Wrap436"),
			is:   IsNotValid,
			want: true,
		}, {
			err:  NewNotImplemented(Wrap(NewNotValidf("Err480"), "Wrap480"), ""),
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
			err:  Wrap(NewTemporaryf("Err474"), "Wrap474"),
			is:   IsTemporary,
			want: true,
		}, {
			err:  NewNotImplemented(Wrap(NewTemporaryf("Err522"), "Wrap522"), ""),
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
			err:  Wrap(NewTimeoutf("Err512"), "Wrap512"),
			is:   IsTimeout,
			want: true,
		}, {
			err:  NewNotImplemented(Wrap(NewTimeoutf("Err564"), "Wrap564"), "ni562"),
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

func TestErrWrapf(t *testing.T) {
	const e Error = "Error1"
	if haveEB, want := errWrapf(e, "Hello World %#v"), "Hello World %#v"; haveEB.msg != want {
		t.Errorf("have %q want %q", haveEB.msg, want)
	}
	if haveEB, want := errWrapf(e, "Hello World %d", 123), "Hello World 123"; haveEB.msg != want {
		t.Errorf("have %q want %q", haveEB.msg, want)
	}
}

func TestErrNewf(t *testing.T) {
	if have, want := errNewf("Hello World %d", 633), "Hello World 633"; have.msg != want {
		t.Errorf("have %q want %q", have.msg, want)
	}
	if have, want := errNewf("Hello World %d"), "Hello World %d"; have.msg != want {
		t.Errorf("have %q want %q", have.msg, want)
	}
}

func TestHasBehaviour(t *testing.T) {
	tests := []struct {
		err           error
		wantBehaviour int
	}{
		{Error("err27"), 0},
		{NewAlreadyClosedf("err28"), BehaviourAlreadyClosed},
		{NewAlreadyExistsf("err29"), BehaviourAlreadyExists},
		{NewEmptyf("err50"), BehaviourEmpty},
		{NewFatalf("err31"), BehaviourFatal},
		{NewNotFoundf("err32"), BehaviourNotFound},
		{NewNotImplementedf("err33"), BehaviourNotImplemented},
		{NewNotSupportedf("err34"), BehaviourNotSupported},
		{NewNotValidf("err35"), BehaviourNotValid},
		{NewTemporaryf("err36"), BehaviourTemporary},
		{NewTimeoutf("err37"), BehaviourTimeout},
		{NewUnauthorizedf("err38"), BehaviourUnauthorized},
		{NewUserNotFoundf("err39"), BehaviourUserNotFound},
		{NewWriteFailedf("err40"), BehaviourWriteFailed},
	}
	for _, test := range tests {
		if have, want := HasBehaviour(test.err), test.wantBehaviour; have != want {
			t.Errorf("%s: Have: %d Want: %d", test.err, have, want)
		}
	}
}
