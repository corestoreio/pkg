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
	"github.com/corestoreio/csfw/config"
	"github.com/corestoreio/csfw/storage/csdb"
	"github.com/corestoreio/csfw/storage/dbr"
	"github.com/corestoreio/csfw/utils/crypto"
)

// @todo app/code/Magento/User/Model/User.php
// this whole API will change just some brain storming
// instead of returning bool in some functions we return nil (success) or an error

type UserSlice []*User

// UserOption can be used as an argument in NewUser to configure a user.
type UserOption func(*User)

type User struct {
	mc   config.ModelConstructor // @todo see directory pkg
	Data *TableAdminUser
}

func LoadOne(dbrSess dbr.SessionRunner, cbs ...dbr.SelectCb) UserOption {
	return func(u *User) {
		// todo
	}
}

func NewUser(opts ...UserOption) *User {
	u := new(User)
	u.ApplyOptions(opts...)
	return u
}

func (u *User) ApplyOptions(opts ...UserOption) *User {
	for _, o := range opts {
		if o != nil {
			o(u)
		}
	}
	return u
}

func (u *User) Authenticate(cr config.Reader, h crypto.Hasher, username, password string) error {
	isCaseSensitive := cr.GetBool(config.Path("admin/security/use_case_sensitive_login"))

	if !isCaseSensitive {
		// ... hmm
	}

	return nil
}

func (u *User) VerifyIdentity() error {
	// validateHash()
	// getIsActive() ?
	// hasAssigned2Role()
	return nil
}

func (u *User) Login(username, password string) error {
	// u.Authenticate()
	// recordLogin()
	return nil
}

// HasAssigned2Role check if user is assigned to any role
func (u *User) HasAssigned2Role() error {
	// check entries in table authorization_role

	return nil
}

func (u *User) Reload() error {
	// reload Data

	return nil
}
