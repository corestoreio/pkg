// Copyright 2015, Cyrill @ Schumacher.fm and the CoreStore contributors
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

package user

import (
	"errors"

	"github.com/corestoreio/csfw/storage/csdb"
)

var (
	// TableCollection handles all tables and its columns. init() in generated Go file will set the value.
	TableCollection csdb.Manager

	ErrUserNotFound = errors.New("Admin user not found")
)

// FindByID returns a TableAdminUser if found by id or an error
func (s TableAdminUserSlice) FindByID(id int64) (*TableAdminUser, error) {
	for _, u := range s {
		if u != nil && u.UserID == id {
			return u, nil
		}
	}
	return nil, ErrUserNotFound
}

// FindByUsername returns a TableAdminUser if found by code or an error
func (s TableAdminUserSlice) FindByUsername(username string) (*TableAdminUser, error) {
	for _, u := range s {
		if u != nil && u.Username.Valid && u.Username.String == username {
			return u, nil
		}
	}
	return nil, ErrUserNotFound
}

// Filter returns a new slice filtered by predicate f
func (s TableAdminUserSlice) Filter(f func(*TableAdminUser) bool) TableAdminUserSlice {
	var tws TableAdminUserSlice
	for _, w := range s {
		if w != nil && f(w) {
			tws = append(tws, w)
		}
	}
	return tws
}
