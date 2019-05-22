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

// +build csall db

package store

import (
	"context"

	"github.com/corestoreio/errors"
	"github.com/corestoreio/pkg/sql/ddl"
)

// WithLoadFromDB loads the store,group and website data from the database.
// Before loading it clears the cache. This function panics if the tables do not
// exists in the ddl.Tables object.
func WithLoadFromDB(ctx context.Context, tbls *ddl.Tables) Option {
	stmtStore := tbls.MustTable(TableNameStore).Select("*").WithArgs()
	stmtGroup := tbls.MustTable(TableNameStoreGroup).Select("*").WithArgs()
	stmtWebsite := tbls.MustTable(TableNameStoreWebsite).Select("*").WithArgs()

	return Option{
		sortOrder: 199,
		fn: func(s *Service) error {
			s.ClearCache()
			s.mu.Lock()
			defer s.mu.Unlock()

			if _, err := stmtStore.Load(ctx, &s.stores); err != nil {
				return errors.WithStack(err)
			}
			if _, err := stmtGroup.Load(ctx, &s.groups); err != nil {
				return errors.WithStack(err)
			}
			if _, err := stmtWebsite.Load(ctx, &s.websites); err != nil {
				return errors.WithStack(err)
			}
			if err := s.stores.Validate(); err != nil {
				return errors.WithStack(err)
			}
			if err := s.groups.Validate(); err != nil {
				return errors.WithStack(err)
			}
			if err := s.websites.Validate(); err != nil {
				return errors.WithStack(err)
			}
			return nil
		},
	}
}
