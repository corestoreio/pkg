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

package jwtclaim

import (
	"github.com/corestoreio/csfw/util/conv"
	"github.com/juju/errors"
)

//go:generate ffjson $GOFILE

// Available claims in struct StoreClaim.
// ClaimStore is equal to github.com/corestoreio/csfw/store/storenet.ParamName
const (
	ClaimStore  = "store"
	ClaimUserID = "userid"
)

// NewStore creates a new Store pointer and makes sure that all sub pointer
// types will also be created.
func NewStore() *Store {
	return &Store{
		Standard: new(Standard),
	}
}

// Store extends the StandardClaim with important fields for requesting the
// correct store view, user ID and maybe some more useful fields.
// This struct is for your convenience.
// ffjson: noencoder
type Store struct {
	*Standard
	Store string `json:"store,omitempty"`
	// UserID add here any user ID you might will be but always bear in mind that
	// when adding a numeric auto increment ID, like customer_id from the MySQL
	// table customer_entity or admin_user you might leak sensitive information.
	UserID string `json:"userid,omitempty"`
	// todo extend with more useful fields
}

// Set allows to set StoreClaim specific fields and then falls back to the set
// function in StandardClaims
func (s *Store) Set(key string, value interface{}) (err error) {

	switch key {
	case ClaimStore:
		s.Store, err = conv.ToStringE(value)
		err = errors.Mask(err)
	case ClaimUserID:
		s.UserID, err = conv.ToStringE(value)
		err = errors.Mask(err)
	}
	if err != nil {
		return err
	}

	return s.Standard.Set(key, value)
}

// Get retrieves StoreClaim specific fields and then falls back to the
// StandardClaims Get function.
func (s *Store) Get(key string) (value interface{}, err error) {
	switch key {
	case ClaimStore:
		return s.Store, nil
	case ClaimUserID:
		return s.UserID, nil
	}
	return s.Standard.Get(key)
}
