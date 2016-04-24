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

// TODO(cs): add notImplemented,notProvisioned,badRequest,methodNotAllowed,notAssigned,...

type fatal struct {
	error
}

// NewFatal returns an error which wraps err that satisfies
// IsFatal().
func NewFatal(err error, msg string, args ...interface{}) error {
	return &fatal{Wrapf(err, msg, args...)}
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

type notFound struct {
	error
}

// NewNotFound returns an error which wraps err that satisfies
// IsNotFound().
func NewNotFound(err error, msg string, args ...interface{}) error {
	return &notFound{Wrapf(err, msg, args...)}
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

type userNotFound struct {
	error
}

// NewUserNotFound returns an error which wraps err and satisfies
// IsUserNotFound().
func NewUserNotFound(err error, msg string, args ...interface{}) error {
	return &userNotFound{Wrapf(err, msg, args...)}
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

type unauthorized struct {
	error
}

// NewUnauthorized returns an error which wraps err and satisfies
// IsUnauthorized().
func NewUnauthorized(err error, msg string, args ...interface{}) error {
	return &unauthorized{Wrapf(err, msg, args...)}
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

type alreadyExists struct {
	error
}

// NewAlreadyExists returns an error which wraps err and satisfies
// IsAlreadyExists().
func NewAlreadyExists(err error, msg string, args ...interface{}) error {
	return &alreadyExists{Wrapf(err, msg, args...)}
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

type alreadyClosed struct {
	error
}

// NewAlreadyClosed returns an error which wraps err and satisfies
// IsAlreadyClosed().
func NewAlreadyClosed(err error, msg string, args ...interface{}) error {
	return &alreadyClosed{Wrapf(err, msg, args...)}
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

type notSupported struct {
	error
}

// NewNotSupported returns an error which wraps err and satisfies
// IsNotSupported().
func NewNotSupported(err error, msg string, args ...interface{}) error {
	return &notSupported{Wrapf(err, msg, args...)}
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

type notValid struct {
	error
}

// NewNotValid returns an error which wraps err and satisfies
// IsNotValid().
func NewNotValid(err error, msg string, args ...interface{}) error {
	return &notValid{Wrapf(err, msg, args...)}
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

type temporary struct {
	error
}

// NewTemporary returns an error which wraps err and satisfies
// IsTemporary().
func NewTemporary(err error, msg string, args ...interface{}) error {
	return &temporary{Wrapf(err, msg, args...)}
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

type timeout struct {
	error
}

// NewTimeout returns an error which wraps err and satisfies
// IsTimeout().
func NewTimeout(err error, msg string, args ...interface{}) error {
	return &timeout{Wrapf(err, msg, args...)}
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
