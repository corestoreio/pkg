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

package ddl

import (
	"github.com/corestoreio/errors"
	"github.com/corestoreio/pkg/sql/dml"
)

type DBM struct {
	*Tables
	queries map[string]string
}

func NewDBM(opt []TableOption, cbs ...func(*DBM) dml.QueryBuilder) (*DBM, error) {
	t, err := NewTables(opt...)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	dbm := &DBM{
		Tables:  t,
		queries: map[string]string{},
	}
	for _, cb := range cbs {
		sqlStr, _, err := cb(dbm).ToSQL()
		if err != nil {
			return nil, errors.WithStack(err)
		}
		dbm.queries[sqlStr] = sqlStr // does not make sense
	}
	return dbm, nil
}

func (dbm *DBM) CachedQuery(key string) string {
	return dbm.queries[key]
}
