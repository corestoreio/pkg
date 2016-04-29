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
	"fmt"
	"runtime"
)

// BehaviourFunc defines the signature needed for a function to check
// if an error has a specific behaviour attached.
type BehaviourFunc func(error) bool

func newEB(err error, msg string, pc uintptr) eb {
	return eb{
		err:      err,
		message:  msg,
		location: location(pc),
	}
}

type eb struct {
	err     error
	message string
	location
}

func (e *eb) Error() string {
	if e.err == nil {
		return e.message
	}
	return e.message + ": " + e.err.Error()
}

func ebWrapf(err error, format string, args ...interface{}) eb {
	pc, _, _, _ := runtime.Caller(2)
	return newEB(err, fmt.Sprintf(format, args...), pc)
}

func ebWrap(err error, msg string) eb {
	pc, _, _, _ := runtime.Caller(2)
	return newEB(err, msg, pc)
}

// TODO(cs): add notProvisioned,badRequest,methodNotAllowed,notAssigned,...

type notImplemented struct{ eb }

const notImplementedTxt Error = "Not implemented"

// NewNotImplemented returns an error which wraps err that satisfies
// IsNotImplemented().
func NewNotImplemented(err error, msg string) error {
	if err == nil {
		return nil
	}
	return &notImplemented{ebWrap(err, msg)}
}

// NewNotImplementedf returns an formatted error that satisfies IsNotImplemented().
func NewNotImplementedf(format string, args ...interface{}) error {
	return &notImplemented{ebWrapf(notImplementedTxt, format, args...)}
}

// IsNotImplemented reports whether err was created with NewNotImplemented() or
// has a method receiver "NotImplemented() bool".
func IsNotImplemented(err error) bool {
	type iFace interface {
		NotImplemented() bool
	}
	err = Cause(err)
	var ok bool
	switch et := err.(type) {
	case *notImplemented:
		ok = true
	case iFace:
		ok = et.NotImplemented()
	}
	return ok
}

type empty struct{ eb }

const emptyTxt Error = "Empty value"

// NewEmpty returns an error which wraps err that satisfies
// IsEmpty().
func NewEmpty(err error, msg string) error {
	if err == nil {
		return nil
	}
	return &empty{ebWrap(err, msg)}
}

// NewEmptyf returns an formatted error that satisfies IsEmpty().
func NewEmptyf(format string, args ...interface{}) error {
	return &empty{ebWrapf(emptyTxt, format, args...)}
}

// IsEmpty reports whether err was created with NewEmpty() or
// has a method receiver "Empty() bool".
func IsEmpty(err error) bool {
	type iFace interface {
		Empty() bool
	}
	err = Cause(err)
	var ok bool
	switch et := err.(type) {
	case *empty:
		ok = true
	case iFace:
		ok = et.Empty()
	}
	return ok
}

type writeFailed struct{ eb }

const writeFailedTxt Error = "WriteFailed value"

// NewWriteFailed returns an error which wraps err that satisfies
// IsWriteFailed().
func NewWriteFailed(err error, msg string) error {
	if err == nil {
		return nil
	}
	return &writeFailed{ebWrap(err, msg)}
}

// NewWriteFailedf returns an formatted error that satisfies IsWriteFailed().
func NewWriteFailedf(format string, args ...interface{}) error {
	return &writeFailed{ebWrapf(writeFailedTxt, format, args...)}
}

// IsWriteFailed reports whether err was created with NewWriteFailed() or
// has a method receiver "WriteFailed() bool".
func IsWriteFailed(err error) bool {
	type iFace interface {
		WriteFailed() bool
	}
	err = Cause(err)
	var ok bool
	switch et := err.(type) {
	case *writeFailed:
		ok = true
	case iFace:
		ok = et.WriteFailed()
	}
	return ok
}

type fatal struct{ eb }

const fatalTxt Error = "Fatal"

// NewFatal returns an error which wraps err that satisfies IsFatal().
func NewFatal(err error, msg string) error {
	if err == nil {
		return nil
	}
	return &fatal{ebWrap(err, msg)}
}

// NewFatalf returns an formatted error that satisfies IsFatal().
func NewFatalf(format string, args ...interface{}) error {
	return &fatal{ebWrapf(fatalTxt, format, args...)}
}

// IsFatal reports whether err was created with NewFatal() or
// has a method receiver "Fatal() bool".
func IsFatal(err error) bool {
	type iFace interface {
		Fatal() bool
	}
	err = Cause(err)
	var ok bool
	switch et := err.(type) {
	case *fatal:
		ok = true
	case iFace:
		ok = et.Fatal()
	}
	return ok
}

type notFound struct{ eb }

const notFoundTxt Error = "Not found"

// NewNotFound returns an error which wraps err that satisfies
// IsNotFound().
func NewNotFound(err error, msg string) error {
	if err == nil {
		return nil
	}
	return &notFound{ebWrap(err, msg)}
}

// NewNotFoundf returns an formatted error that satisfies IsNotFound().
func NewNotFoundf(format string, args ...interface{}) error {
	return &notFound{ebWrapf(notFoundTxt, format, args...)}
}

// IsNotFound reports whether err was created with NewNotFound() or
// has a method receiver "NotFound() bool".
func IsNotFound(err error) bool {
	type iFace interface {
		NotFound() bool
	}
	err = Cause(err)
	var ok bool
	switch et := err.(type) {
	case *notFound:
		ok = true
	case iFace:
		ok = et.NotFound()
	}
	return ok
}

type userNotFound struct{ eb }

const userNotFoundTxt Error = "User not found"

// NewUserNotFound returns an error which wraps err and satisfies
// IsUserNotFound().
func NewUserNotFound(err error, msg string) error {
	if err == nil {
		return nil
	}
	return &userNotFound{ebWrap(err, msg)}
}

// NewUserNotFoundf returns an formatted error that satisfies IsUserNotFound().
func NewUserNotFoundf(format string, args ...interface{}) error {
	return &userNotFound{ebWrapf(userNotFoundTxt, format, args...)}
}

// IsUserNotFound reports whether err was created with NewUserNotFound() or
// has a method receiver "UserNotFound() bool".
func IsUserNotFound(err error) bool {
	type iFace interface {
		UserNotFound() bool
	}
	err = Cause(err)
	var ok bool
	switch et := err.(type) {
	case *userNotFound:
		ok = true
	case iFace:
		ok = et.UserNotFound()
	}
	return ok
}

type unauthorized struct{ eb }

const unauthorizedTxt Error = "Unauthorized"

// NewUnauthorized returns an error which wraps err and satisfies
// IsUnauthorized().
func NewUnauthorized(err error, msg string) error {
	if err == nil {
		return nil
	}
	return &unauthorized{ebWrap(err, msg)}
}

// NewUnauthorizedf returns an formatted error that satisfies IsUnauthorized().
func NewUnauthorizedf(format string, args ...interface{}) error {
	return &unauthorized{ebWrapf(unauthorizedTxt, format, args...)}
}

// IsUnauthorized reports whether err was created with NewUnauthorized() or
// has a method receiver "Unauthorized() bool".
func IsUnauthorized(err error) bool {
	type iFace interface {
		Unauthorized() bool
	}
	err = Cause(err)
	var ok bool
	switch et := err.(type) {
	case *unauthorized:
		ok = true
	case iFace:
		ok = et.Unauthorized()
	}
	return ok
}

type alreadyExists struct{ eb }

const alreadyExistsTxt Error = "Already exists"

// NewAlreadyExists returns an error which wraps err and satisfies
// IsAlreadyExists().
func NewAlreadyExists(err error, msg string) error {
	if err == nil {
		return nil
	}
	return &alreadyExists{ebWrap(err, msg)}
}

// NewAlreadyExistsf returns an formatted error that satisfies IsAlreadyExists().
func NewAlreadyExistsf(format string, args ...interface{}) error {
	return &alreadyExists{ebWrapf(alreadyExistsTxt, format, args...)}
}

// IsAlreadyExists reports whether err was created with NewAlreadyExists() or
// has a method receiver "AlreadyExists() bool".
func IsAlreadyExists(err error) bool {
	type iFace interface {
		AlreadyExists() bool
	}
	err = Cause(err)
	var ok bool
	switch et := err.(type) {
	case *alreadyExists:
		ok = true
	case iFace:
		ok = et.AlreadyExists()
	}
	return ok
}

type alreadyClosed struct{ eb }

const alreadyClosedTxt Error = "Already closed"

// NewAlreadyClosed returns an error which wraps err and satisfies
// IsAlreadyClosed().
func NewAlreadyClosed(err error, msg string) error {
	if err == nil {
		return nil
	}
	return &alreadyClosed{ebWrap(err, msg)}
}

// NewAlreadyClosedf returns an formatted error that satisfies IsAlreadyClosed().
func NewAlreadyClosedf(format string, args ...interface{}) error {
	return &alreadyClosed{ebWrapf(alreadyClosedTxt, format, args...)}
}

// IsAlreadyClosed reports whether err was created with NewAlreadyClosed() or
// has a method receiver "AlreadyClosed() bool".
func IsAlreadyClosed(err error) bool {
	type iFace interface {
		AlreadyClosed() bool
	}
	err = Cause(err)
	var ok bool
	switch et := err.(type) {
	case *alreadyClosed:
		ok = true
	case iFace:
		ok = et.AlreadyClosed()
	}
	return ok
}

type notSupported struct{ eb }

const notSupportedTxt Error = "Not supported"

// NewNotSupported returns an error which wraps err and satisfies
// IsNotSupported().
func NewNotSupported(err error, msg string) error {
	if err == nil {
		return nil
	}
	return &notSupported{ebWrap(err, msg)}
}

// NewNotSupportedf returns an formatted error that satisfies IsNotSupported().
func NewNotSupportedf(format string, args ...interface{}) error {
	return &notSupported{ebWrapf(notSupportedTxt, format, args...)}
}

// IsNotSupported reports whether err was created with NewNotSupported() or
// has a method receiver "NotSupported() bool".
func IsNotSupported(err error) bool {
	type iFace interface {
		NotSupported() bool
	}
	err = Cause(err)
	var ok bool
	switch et := err.(type) {
	case *notSupported:
		ok = true
	case iFace:
		ok = et.NotSupported()
	}
	return ok
}

type notValid struct{ eb }

const notValidTxt Error = "Not valid"

// NewNotValid returns an error which wraps err and satisfies
// IsNotValid().
func NewNotValid(err error, msg string) error {
	if err == nil {
		return nil
	}
	return &notValid{ebWrap(err, msg)}
}

// NewNotValidf returns an formatted error that satisfies IsNotValid().
func NewNotValidf(format string, args ...interface{}) error {
	return &notValid{ebWrapf(notValidTxt, format, args...)}
}

// IsNotValid reports whether err was created with NewNotValid() or
// has a method receiver "NotValid() bool".
func IsNotValid(err error) bool {
	type iFace interface {
		NotValid() bool
	}
	err = Cause(err)
	var ok bool
	switch et := err.(type) {
	case *notValid:
		ok = true
	case iFace:
		ok = et.NotValid()
	}
	return ok
}

type temporary struct{ eb }

const temporaryTxt Error = "Temporary"

// NewTemporary returns an error which wraps err and satisfies
// IsTemporary().
func NewTemporary(err error, msg string) error {
	if err == nil {
		return nil
	}
	return &temporary{ebWrap(err, msg)}
}

// NewTemporaryf returns an formatted error that satisfies IsTemporary().
func NewTemporaryf(format string, args ...interface{}) error {
	return &temporary{ebWrapf(temporaryTxt, format, args...)}
}

// IsTemporary reports whether err was created with NewTemporary() or
// has a method receiver "Temporary() bool".
func IsTemporary(err error) bool {
	type iFace interface {
		Temporary() bool
	}
	err = Cause(err)
	var ok bool
	switch et := err.(type) {
	case *temporary:
		ok = true
	case iFace:
		ok = et.Temporary()
	}
	return ok
}

type timeout struct{ eb }

const timeoutTxt Error = "Timeout"

// NewTimeout returns an error which wraps err and satisfies
// IsTimeout().
func NewTimeout(err error, msg string) error {
	if err == nil {
		return nil
	}
	return &timeout{ebWrap(err, msg)}
}

// NewTimeoutf returns an formatted error that satisfies IsTimeout().
func NewTimeoutf(format string, args ...interface{}) error {
	return &timeout{ebWrapf(timeoutTxt, format, args...)}
}

// IsTimeout reports whether err was created with NewTimeout() or
// has a method receiver "Timeout() bool".
func IsTimeout(err error) bool {
	type iFace interface {
		Timeout() bool
	}
	err = Cause(err)
	var ok bool
	switch et := err.(type) {
	case *timeout:
		ok = true
	case iFace:
		ok = et.Timeout()
	}
	return ok
}
