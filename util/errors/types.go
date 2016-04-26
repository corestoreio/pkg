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

// Error type can be used for constant errors and says nothing about its
// behaviour.
// http://dave.cheney.net/2016/04/07/constant-errors
type Error string

// Error implements the error interface
func (e Error) Error() string { return string(e) }

// Empty represents an error with the behaviour that an entity has no value.
type Empty string

// Error implements the error interface
func (e Empty) Error() string { return string(e) }

// Empty satisfies the function IsEmpty()
func (e Empty) Empty() bool { return true }

// NotImplemented represents an error with the behaviour that an entity hasn't
// implemented a feature.
type NotImplemented string

// Error implements the error interface
func (e NotImplemented) Error() string { return string(e) }

// NotImplemented satisfies the function IsNotImplemented()
func (e NotImplemented) NotImplemented() bool { return true }

// Fatal represents an error with the behaviour that the function should
// terminate.
type Fatal string

// Error implements the error interface
func (e Fatal) Error() string { return string(e) }

// Fatal satisfies the function IsFatal()
func (e Fatal) Fatal() bool { return true }

// AlreadyClosed represents an error with the behaviour that an entity
// already has been closed.
type AlreadyClosed string

// Error implements the error interface
func (e AlreadyClosed) Error() string { return string(e) }

// AlreadyClosed satisfies the function IsAlreadyClosed()
func (e AlreadyClosed) AlreadyClosed() bool { return true }

// AlreadyExists represents an error with the behaviour that an entity
// already exists.
type AlreadyExists string

// Error implements the error interface
func (e AlreadyExists) Error() string { return string(e) }

// AlreadyExists satisfies the function IsAlreadyExists()
func (e AlreadyExists) AlreadyExists() bool { return true }

// NotFound represents an error with the behaviour that an entity
// cannot be found.
type NotFound string

// Error implements the error interface
func (e NotFound) Error() string { return string(e) }

// NotFound satisfies the function IsNotFound()
func (e NotFound) NotFound() bool { return true }

// NotSupported represents an error with the behaviour when something
// cannot be supported.
type NotSupported string

// Error implements the error interface
func (e NotSupported) Error() string { return string(e) }

// NotSupported satisfies the function IsNotSupported()
func (e NotSupported) NotSupported() bool { return true }

// NotValid represents an error with the behaviour that an entity
// is not valid.
type NotValid string

// Error implements the error interface
func (e NotValid) Error() string { return string(e) }

// NotValid satisfies the function IsNotValid()
func (e NotValid) NotValid() bool { return true }

// Temporary represents an error with the behaviour the current error
// has only a temporary duration.
type Temporary string

// Error implements the error interface
func (e Temporary) Error() string { return string(e) }

// Temporary satisfies the function IsTemporary()
func (e Temporary) Temporary() bool { return true }

// Timeout represents an error with the behaviour that an entity
// has timed out.
type Timeout string

// Error implements the error interface
func (e Timeout) Error() string { return string(e) }

// Timeout satisfies the function IsTimeout()
func (e Timeout) Timeout() bool { return true }

// Unauthorized represents an error with the behaviour that an entity
// cannot be access granted.
type Unauthorized string

// Error implements the error interface
func (e Unauthorized) Error() string { return string(e) }

// Unauthorized satisfies the function IsUnauthorized()
func (e Unauthorized) Unauthorized() bool { return true }

// UserNotFound represents an error with the behaviour that an user
// cannot be found.
type UserNotFound string

// Error implements the error interface
func (e UserNotFound) Error() string { return string(e) }

// UserNotFound satisfies the function IsUserNotFound()
func (e UserNotFound) UserNotFound() bool { return true }
