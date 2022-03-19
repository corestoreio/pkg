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

package dml

import (
	"context"
	"database/sql"
)

var (
	_ QueryExecPreparer = (*dbMock)(nil)
	_ Execer            = (*dbMock)(nil)
)

type dbMock struct {
	error
	prepareFn func(query string) (*sql.Stmt, error)
}

func (pm dbMock) PrepareContext(ctx context.Context, query string) (*sql.Stmt, error) {
	if pm.error != nil {
		return nil, pm.error
	}
	if pm.prepareFn != nil {
		return pm.prepareFn(query)
	}
	return new(sql.Stmt), nil
}

func (pm dbMock) QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error) {
	if pm.error != nil {
		return nil, pm.error
	}
	return new(sql.Rows), nil
}

func (pm dbMock) ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error) {
	if pm.error != nil {
		return nil, pm.error
	}
	return nil, nil
}

func (pm dbMock) QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row {
	return new(sql.Row)
}
