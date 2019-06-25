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

package store

import "github.com/corestoreio/errors"

// DefaultStoreID is always 0.
const DefaultStoreID int64 = 0

// WithStores upserts new stores to the Store service. It clears the internal
// cache store and sorts the stores by StoreID.
func WithStores(stores ...*Store) Option {
	return Option{
		sortOrder: 302,
		fn: func(s *Service) error {
			s.mu.Lock()
			defer s.mu.Unlock()

			for _, s1 := range stores {
				if err := s1.Validate(); err != nil {
					return errors.WithStack(err)
				}

				var containsS1 bool
				for idx, s2 := range s.stores.Data {
					if s2.StoreID == s1.StoreID {
						s.stores.Data[idx] = s1
						containsS1 = true
						break
					}
				}
				if !containsS1 {
					s.stores.Data = append(s.stores.Data, s1)
				}
			}
			if err := s.stores.Validate(); err != nil {
				return errors.WithStack(err)
			}
			s.cacheStore = make(map[uint32]*Store, len(s.stores.Data))
			return nil
		},
	}
}

// WithGroups upserts new groups to the Store service. It clears the internal
// cache group and sorts the groups by GroupID.
func WithGroups(groups ...*StoreGroup) Option {
	return Option{
		sortOrder: 301,
		fn: func(s *Service) error {
			s.mu.Lock()
			defer s.mu.Unlock()

			for _, s1 := range groups {
				if err := s1.Validate(); err != nil {
					return errors.WithStack(err)
				}

				var containsS1 bool
				for idx, s2 := range s.groups.Data {
					if s2.GroupID == s1.GroupID {
						s.groups.Data[idx] = s1
						containsS1 = true
						break
					}
				}
				if !containsS1 {
					s.groups.Data = append(s.groups.Data, s1)
				}
			}
			if err := s.groups.Validate(); err != nil {
				return errors.WithStack(err)
			}
			s.cacheGroup = make(map[uint32]*StoreGroup, len(s.groups.Data))
			return nil
		},
	}
}

// WithWebsites upserts new websites to the Store service. It clears the internal
// cache website and sorts the websites by WebsiteID.
func WithWebsites(websites ...*StoreWebsite) Option {
	return Option{
		sortOrder: 300,
		fn: func(s *Service) error {
			s.mu.Lock()
			defer s.mu.Unlock()

			for _, s1 := range websites {
				if err := s1.Validate(); err != nil {
					return errors.WithStack(err)
				}

				var containsS1 bool
				for idx, s2 := range s.websites.Data {
					if s2.WebsiteID == s1.WebsiteID {
						s.websites.Data[idx] = s1
						containsS1 = true
						break
					}
				}
				if !containsS1 {
					s.websites.Data = append(s.websites.Data, s1)
				}
			}
			if err := s.websites.Validate(); err != nil {
				return errors.WithStack(err)
			}
			s.cacheWebsite = make(map[uint32]*StoreWebsite, len(s.websites.Data))
			return nil
		},
	}
}
