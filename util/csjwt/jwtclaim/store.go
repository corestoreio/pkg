// Copyright 2015-present, Cyrill @ Schumacher.fm and the CoreStore contributors
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
	"encoding/json"

	"github.com/corestoreio/errors"
	"github.com/corestoreio/pkg/util/conv"
)

// Key.... are available claims in struct Store.
// KeyStore is equal to github.com/corestoreio/pkg/store/storenet.ParamName
const (
	KeyStore  = "store"
	KeyUserID = "userid"
)

// NewStore creates a new Store pointer and makes sure that all sub pointer
// types will also be created.
func NewStore() *Store {
	return &Store{}
}

// Store extends the StandardClaim with important fields for requesting the
// correct store view, user ID and maybe some more useful fields. This struct is
// for your convenience.
//easyjson:json
type Store struct {
	Standard
	Store string `json:"store,omitempty"`
	// UserID add here any user ID you might will be but always bear in mind
	// that when adding a numeric auto increment ID, like customer_id from the
	// MySQL table customer_entity or admin_user you might leak sensitive
	// information.
	UserID string `json:"userid,omitempty"`
}

// TODO(cs) extend Store type with more useful fields

// Set allows to set StoreClaim specific fields and then falls back to the set
// function in StandardClaims
func (s *Store) Set(key string, value interface{}) (err error) {
	switch key {
	case KeyStore:
		s.Store, err = conv.ToStringE(value)
		return errors.Wrap(err, "[jwtclaim] Store.ToString")
	case KeyUserID:
		s.UserID, err = conv.ToStringE(value)
		return errors.Wrap(err, "[jwtclaim] UserID.ToString")
	}

	return s.Standard.Set(key, value)
}

// Get retrieves StoreClaim specific fields and then falls back to the
// StandardClaims Get function.
func (s *Store) Get(key string) (interface{}, error) {
	switch key {
	case KeyStore:
		return s.Store, nil
	case KeyUserID:
		return s.UserID, nil
	}
	return s.Standard.Get(key)
}

// Keys returns all available keys which this type supports.
func (s *Store) Keys() []string {
	return allKeys[:]
}

// String human readable output via JSON, slow.
func (s *Store) String() string {
	b, err := json.Marshal(s)
	if err != nil {
		return errors.Fatal.Newf("[jwtclaim] Store.String(): json.Marshal Error: %s", err).Error()
	}
	return string(b)
}
