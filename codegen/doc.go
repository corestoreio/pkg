// Copyright 2015-2017, Cyrill @ Schumacher.fm and the CoreStore contributors
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

// Package codegen contains functions for writing database tables into Go structs.
//
// It supports package storage/dbr and its full functionality
//
// 3rd party packages which can support in code generation:
// - https://github.com/fatih/astrewrite
// - https://github.com/matroskin13/grizzly Collections generator for Golang
// - https://github.com/awalterschulze/goderive

// Idea when generating the code, but in the above case column email has a unique index and hence
// we can throw an error NewDuplicated
//// https://github.com/gobuffalo/authrecipe/blob/master/models/user.go#53
//
//// Validate gets run every time you call a "pop.Validate" method.
//func (u *User) Validate(tx *pop.Connection) (*validate.Errors, error) {
//	var err error
//	return validate.Validate(
//		&validators.StringIsPresent{Field: u.Email, Name: "Email"},
//		&validators.StringIsPresent{Field: u.PasswordHash, Name: "PasswordHash"},
//		// check to see if the email address is already taken:
//		&validators.FuncValidator{
//			Field:   u.Email,
//			Name:    "Email",
//			Message: "%s is already taken",
//			Fn: func() bool {
//				var b bool
//				q := tx.Where("email = ?", u.Email)
//				if u.ID != uuid.Nil {
//					q = q.Where("id != ?", u.ID)
//				}
//				b, err = q.Exists(u)
//				if err != nil {
//					return false
//				}
//				return !b
//			},
//		},
//	), err
//}

package codegen
