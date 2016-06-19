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

import "fmt"

// BehaviourFunc defines the signature needed for a function to check
// if an error has a specific behaviour attached.
type BehaviourFunc func(error) bool

// Behaviour constants are returned by function IsBehaviour() to detect which
// behaviour an error has. The order of the constants follows the alphabet.
// Maybe more constants will be added.
const (
	BehaviourAlreadyClosed int = iota + 1
	BehaviourAlreadyExists
	BehaviourEmpty
	BehaviourFatal
	BehaviourNotFound
	BehaviourNotImplemented
	BehaviourNotSupported
	BehaviourNotValid
	BehaviourTemporary
	BehaviourTimeout
	BehaviourUnauthorized
	BehaviourUserNotFound
	BehaviourWriteFailed
)

// HasBehaviour detects which behaviour an error has. It returns 0 when the
// behaviour is not defined.
func HasBehaviour(err error) int {
	var ret int
	switch {
	case IsAlreadyClosed(err):
		ret = BehaviourAlreadyClosed
	case IsAlreadyExists(err):
		ret = BehaviourAlreadyExists
	case IsFatal(err):
		ret = BehaviourFatal
	case IsEmpty(err):
		ret = BehaviourEmpty
	case IsNotFound(err):
		ret = BehaviourNotFound
	case IsNotImplemented(err):
		ret = BehaviourNotImplemented
	case IsNotSupported(err):
		ret = BehaviourNotSupported
	case IsNotValid(err):
		ret = BehaviourNotValid
	case IsTemporary(err):
		ret = BehaviourTemporary
	case IsTimeout(err):
		ret = BehaviourTimeout
	case IsUnauthorized(err):
		ret = BehaviourUnauthorized
	case IsUserNotFound(err):
		ret = BehaviourUserNotFound
	case IsWriteFailed(err):
		ret = BehaviourWriteFailed
	}
	return ret
}

func errWrapf(err error, format string, args ...interface{}) wrapper {
	ret := wrapper{
		cause: cause{
			cause: err,
			msg:   format,
		},
		stack: callers(),
	}
	if len(args) > 0 {
		ret.cause.msg = fmt.Sprintf(format, args...)
	}
	return ret
}

func errNewf(format string, args ...interface{}) (ret _error) {
	ret.msg = format
	ret.stack = callers()
	if len(args) > 0 {
		ret.msg = fmt.Sprintf(format, args...)
	}
	return
}

// TODO(cs): add notProvisioned,badRequest,methodNotAllowed,notAssigned,...

type (
	notImplemented  struct{ wrapper }
	notImplementedf struct{ _error }
)

// const notImplementedTxt NotImplemented = "Not implemented"

// NewNotImplemented returns an error which wraps err that satisfies
// IsNotImplemented().
func NewNotImplemented(err error, msg string) error {
	if err == nil {
		return nil
	}
	return &notImplemented{errWrapf(err, msg)}
}

// NewNotImplementedf returns an formatted error that satisfies IsNotImplemented().
func NewNotImplementedf(format string, args ...interface{}) error {
	return &notImplementedf{errNewf(format, args...)}
}

func isNotImplemented(err error) (ok bool) {
	type iFace interface {
		NotImplemented() bool
	}
	switch et := err.(type) {
	case *notImplemented:
		ok = true
	case *notImplementedf:
		ok = true
	case iFace:
		ok = et.NotImplemented()
	}
	return
}

// IsNotImplemented reports whether err was created with NewNotImplemented() or
// has a method receiver "NotImplemented() bool".
func IsNotImplemented(err error) bool {
	// check if direct hit that err implements the behaviour.
	if isNotImplemented(err) {
		return true
	}

	// unwrap until we get the root cause which might also implement the
	// behaviour.
	return isNotImplemented(Cause(err))
}

type (
	empty  struct{ wrapper }
	emptyf struct{ _error }
)

// NewEmpty returns an error which wraps err that satisfies
// IsEmpty().
func NewEmpty(err error, msg string) error {
	if err == nil {
		return nil
	}
	return &empty{errWrapf(err, msg)}
}

// NewEmptyf returns an formatted error that satisfies IsEmpty().
func NewEmptyf(format string, args ...interface{}) error {
	return &emptyf{errNewf(format, args...)}
}

func isEmpty(err error) (ok bool) {
	type iFace interface {
		Empty() bool
	}
	switch et := err.(type) {
	case *empty:
		ok = true
	case *emptyf:
		ok = true
	case iFace:
		ok = et.Empty()
	}
	return
}

// IsEmpty reports whether err was created with NewEmpty() or
// has a method receiver "Empty() bool".
func IsEmpty(err error) bool {
	// check if direct hit that err implements the behaviour.
	if isEmpty(err) {
		return true
	}
	// unwrap until we get the root cause which might also implement the
	// behaviour.
	return isEmpty(Cause(err))
}

type (
	writeFailed  struct{ wrapper }
	writeFailedf struct{ _error }
)

// NewWriteFailed returns an error which wraps err that satisfies
// IsWriteFailed().
func NewWriteFailed(err error, msg string) error {
	if err == nil {
		return nil
	}
	return &writeFailed{errWrapf(err, msg)}
}

// NewWriteFailedf returns an formatted error that satisfies IsWriteFailed().
func NewWriteFailedf(format string, args ...interface{}) error {
	return &writeFailedf{errNewf(format, args...)}
}

func isWriteFailed(err error) (ok bool) {
	type iFace interface {
		WriteFailed() bool
	}
	switch et := err.(type) {
	case *writeFailed:
		ok = true
	case *writeFailedf:
		ok = true
	case iFace:
		ok = et.WriteFailed()
	}
	return ok
}

// IsWriteFailed reports whether err was created with NewWriteFailed() or
// has a method receiver "WriteFailed() bool".
func IsWriteFailed(err error) bool {
	// check if direct hit that err implements the behaviour.
	if isWriteFailed(err) {
		return true
	}
	// unwrap until we get the root cause which might also implement the
	// behaviour.
	return isWriteFailed(Cause(err))
}

type (
	fatal  struct{ wrapper }
	fatalf struct{ _error }
)

// NewFatal returns an error which wraps err that satisfies IsFatal().
func NewFatal(err error, msg string) error {
	if err == nil {
		return nil
	}
	return &fatal{errWrapf(err, msg)}
}

// NewFatalf returns an formatted error that satisfies IsFatal().
func NewFatalf(format string, args ...interface{}) error {
	return &fatalf{errNewf(format, args...)}
}

func isFatal(err error) (ok bool) {
	type iFace interface {
		Fatal() bool
	}
	switch et := err.(type) {
	case *fatal:
		ok = true
	case *fatalf:
		ok = true
	case iFace:
		ok = et.Fatal()
	}
	return
}

// IsFatal reports whether err was created with NewFatal() or
// has a method receiver "Fatal() bool".
func IsFatal(err error) bool {
	// check if direct hit that err implements the behaviour.
	if isFatal(err) {
		return true
	}
	// unwrap until we get the root cause which might also implement the
	// behaviour.
	return isFatal(Cause(err))
}

type (
	notFound  struct{ wrapper }
	notFoundf struct{ _error }
)

// NewNotFound returns an error which wraps err that satisfies
// IsNotFound().
func NewNotFound(err error, msg string) error {
	if err == nil {
		return nil
	}
	return &notFound{errWrapf(err, msg)}
}

// NewNotFoundf returns an formatted error that satisfies IsNotFound().
func NewNotFoundf(format string, args ...interface{}) error {
	return &notFoundf{errNewf(format, args...)}
}

func isNotFound(err error) (ok bool) {
	type iFace interface {
		NotFound() bool
	}
	switch et := err.(type) {
	case *notFound:
		ok = true
	case *notFoundf:
		ok = true
	case iFace:
		ok = et.NotFound()
	}
	return ok
}

// IsNotFound reports whether err was created with NewNotFound() or
// has a method receiver "NotFound() bool".
func IsNotFound(err error) bool {
	// check if direct hit that err implements the behaviour.
	if isNotFound(err) {
		return true
	}
	// unwrap until we get the root cause which might also implement the
	// behaviour.
	return isNotFound(Cause(err))
}

type (
	userNotFound  struct{ wrapper }
	userNotFoundf struct{ _error }
)

// NewUserNotFound returns an error which wraps err and satisfies
// IsUserNotFound().
func NewUserNotFound(err error, msg string) error {
	if err == nil {
		return nil
	}
	return &userNotFound{errWrapf(err, msg)}
}

// NewUserNotFoundf returns an formatted error that satisfies IsUserNotFound().
func NewUserNotFoundf(format string, args ...interface{}) error {
	return &userNotFoundf{errNewf(format, args...)}
}

func isUserNotFound(err error) (ok bool) {
	type iFace interface {
		UserNotFound() bool
	}
	switch et := err.(type) {
	case *userNotFound:
		ok = true
	case *userNotFoundf:
		ok = true
	case iFace:
		ok = et.UserNotFound()
	}
	return ok
}

// IsUserNotFound reports whether err was created with NewUserNotFound() or
// has a method receiver "UserNotFound() bool".
func IsUserNotFound(err error) bool {
	// check if direct hit that err implements the behaviour.
	if isUserNotFound(err) {
		return true
	}
	// unwrap until we get the root cause which might also implement the
	// behaviour.
	return isUserNotFound(Cause(err))
}

type (
	unauthorized  struct{ wrapper }
	unauthorizedf struct{ _error }
)

// NewUnauthorized returns an error which wraps err and satisfies
// IsUnauthorized().
func NewUnauthorized(err error, msg string) error {
	if err == nil {
		return nil
	}
	return &unauthorized{errWrapf(err, msg)}
}

// NewUnauthorizedf returns an formatted error that satisfies IsUnauthorized().
func NewUnauthorizedf(format string, args ...interface{}) error {
	return &unauthorizedf{errNewf(format, args...)}
}

func isUnauthorized(err error) (ok bool) {
	type iFace interface {
		Unauthorized() bool
	}
	switch et := err.(type) {
	case *unauthorized:
		ok = true
	case *unauthorizedf:
		ok = true
	case iFace:
		ok = et.Unauthorized()
	}
	return
}

// IsUnauthorized reports whether err was created with NewUnauthorized() or
// has a method receiver "Unauthorized() bool".
func IsUnauthorized(err error) bool {
	// check if direct hit that err implements the behaviour.
	if isUnauthorized(err) {
		return true
	}
	// unwrap until we get the root cause which might also implement the
	// behaviour.
	return isUnauthorized(Cause(err))
}

type (
	alreadyExists  struct{ wrapper }
	alreadyExistsf struct{ _error }
)

// NewAlreadyExists returns an error which wraps err and satisfies
// IsAlreadyExists().
func NewAlreadyExists(err error, msg string) error {
	if err == nil {
		return nil
	}
	return &alreadyExists{errWrapf(err, msg)}
}

// NewAlreadyExistsf returns an formatted error that satisfies IsAlreadyExists().
func NewAlreadyExistsf(format string, args ...interface{}) error {
	return &alreadyExistsf{errNewf(format, args...)}
}

func isAlreadyExists(err error) (ok bool) {
	type iFace interface {
		AlreadyExists() bool
	}
	switch et := err.(type) {
	case *alreadyExists:
		ok = true
	case *alreadyExistsf:
		ok = true
	case iFace:
		ok = et.AlreadyExists()
	}
	return
}

// IsAlreadyExists reports whether err was created with NewAlreadyExists() or
// has a method receiver "AlreadyExists() bool".
func IsAlreadyExists(err error) bool {
	// check if direct hit that err implements the behaviour.
	if isAlreadyExists(err) {
		return true
	}
	// unwrap until we get the root cause which might also implement the
	// behaviour.
	return isAlreadyExists(Cause(err))
}

type (
	alreadyClosed  struct{ wrapper }
	alreadyClosedf struct{ _error }
)

// NewAlreadyClosed returns an error which wraps err and satisfies
// IsAlreadyClosed().
func NewAlreadyClosed(err error, msg string) error {
	if err == nil {
		return nil
	}
	return &alreadyClosed{errWrapf(err, msg)}
}

// NewAlreadyClosedf returns an formatted error that satisfies IsAlreadyClosed().
func NewAlreadyClosedf(format string, args ...interface{}) error {
	return &alreadyClosedf{errNewf(format, args...)}
}

func isAlreadyClosed(err error) (ok bool) {
	type iFace interface {
		AlreadyClosed() bool
	}
	switch et := err.(type) {
	case *alreadyClosed:
		ok = true
	case *alreadyClosedf:
		ok = true
	case iFace:
		ok = et.AlreadyClosed()
	}
	return ok
}

// IsAlreadyClosed reports whether err was created with NewAlreadyClosed() or
// has a method receiver "AlreadyClosed() bool".
func IsAlreadyClosed(err error) bool {
	// check if direct hit that err implements the behaviour.
	if isAlreadyClosed(err) {
		return true
	}
	// unwrap until we get the root cause which might also implement the
	// behaviour.
	return isAlreadyClosed(Cause(err))
}

type (
	notSupported  struct{ wrapper }
	notSupportedf struct{ _error }
)

// NewNotSupported returns an error which wraps err and satisfies
// IsNotSupported().
func NewNotSupported(err error, msg string) error {
	if err == nil {
		return nil
	}
	return &notSupported{errWrapf(err, msg)}
}

// NewNotSupportedf returns an formatted error that satisfies IsNotSupported().
func NewNotSupportedf(format string, args ...interface{}) error {
	return &notSupportedf{errNewf(format, args...)}
}

func isNotSupported(err error) (ok bool) {
	type iFace interface {
		NotSupported() bool
	}
	switch et := err.(type) {
	case *notSupported:
		ok = true
	case *notSupportedf:
		ok = true
	case iFace:
		ok = et.NotSupported()
	}
	return
}

// IsNotSupported reports whether err was created with NewNotSupported() or
// has a method receiver "NotSupported() bool".
func IsNotSupported(err error) bool {
	// check if direct hit that err implements the behaviour.
	if isNotSupported(err) {
		return true
	}
	// unwrap until we get the root cause which might also implement the
	// behaviour.
	return isNotSupported(Cause(err))
}

type (
	notValid  struct{ wrapper }
	notValidf struct{ _error }
)

// NewNotValid returns an error which wraps err and satisfies
// IsNotValid().
func NewNotValid(err error, msg string) error {
	if err == nil {
		return nil
	}
	return &notValid{errWrapf(err, msg)}
}

// NewNotValidf returns an formatted error that satisfies IsNotValid().
func NewNotValidf(format string, args ...interface{}) error {
	return &notValidf{errNewf(format, args...)}
}

func isNotValid(err error) (ok bool) {
	type iFace interface {
		NotValid() bool
	}
	switch et := err.(type) {
	case *notValid:
		ok = true
	case *notValidf:
		ok = true
	case iFace:
		ok = et.NotValid()
	}
	return
}

// IsNotValid reports whether err was created with NewNotValid() or
// has a method receiver "NotValid() bool".
func IsNotValid(err error) bool {
	// check if direct hit that err implements the behaviour.
	if isNotValid(err) {
		return true
	}
	// unwrap until we get the root cause which might also implement the
	// behaviour.
	return isNotValid(Cause(err))
}

type (
	temporary  struct{ wrapper }
	temporaryf struct{ _error }
)

// NewTemporary returns an error which wraps err and satisfies
// IsTemporary().
func NewTemporary(err error, msg string) error {
	if err == nil {
		return nil
	}
	return &temporary{errWrapf(err, msg)}
}

// NewTemporaryf returns an formatted error that satisfies IsTemporary().
func NewTemporaryf(format string, args ...interface{}) error {
	return &temporaryf{errNewf(format, args...)}
}

func isTemporary(err error) (ok bool) {
	type iFace interface {
		Temporary() bool
	}
	switch et := err.(type) {
	case *temporary:
		ok = true
	case *temporaryf:
		ok = true
	case iFace:
		ok = et.Temporary()
	}
	return
}

// IsTemporary reports whether err was created with NewTemporary() or
// has a method receiver "Temporary() bool".
func IsTemporary(err error) bool {
	// check if direct hit that err implements the behaviour.
	if isTemporary(err) {
		return true
	}
	// unwrap until we get the root cause which might also implement the
	// behaviour.
	return isTemporary(Cause(err))
}

type (
	timeout  struct{ wrapper }
	timeoutf struct{ _error }
)

// NewTimeout returns an error which wraps err and satisfies
// IsTimeout().
func NewTimeout(err error, msg string) error {
	if err == nil {
		return nil
	}
	return &timeout{errWrapf(err, msg)}
}

// NewTimeoutf returns an formatted error that satisfies IsTimeout().
func NewTimeoutf(format string, args ...interface{}) error {
	return &timeoutf{errNewf(format, args...)}
}

func isTimeout(err error) (ok bool) {
	type iFace interface {
		Timeout() bool
	}
	switch et := err.(type) {
	case *timeout:
		ok = true
	case *timeoutf:
		ok = true
	case iFace:
		ok = et.Timeout()
	}
	return
}

// IsTimeout reports whether err was created with NewTimeout() or
// has a method receiver "Timeout() bool".
func IsTimeout(err error) bool {
	// check if direct hit that err implements the behaviour.
	if isTimeout(err) {
		return true
	}
	// unwrap until we get the root cause which might also implement the
	// behaviour.
	return isTimeout(Cause(err))
}
