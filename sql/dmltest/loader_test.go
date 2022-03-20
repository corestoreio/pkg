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

package dmltest

import (
	"context"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/corestoreio/errors"
	"github.com/corestoreio/pkg/util/assert"
)

func TestSQLDumpLoad(t *testing.T) {
	t.Parallel()

	t.Run("load sql files", func(t *testing.T) {

		exexCmd := func(ctx context.Context, r io.ReadCloser, cmd string, arg ...string) error {
			assert.Exactly(t, `bash`, cmd)

			args := strings.Join(arg, " ")
			assert.Contains(t, args, "mydefaults-")
			assert.Contains(t, args, "-c mysql --defaults-file=")
			assert.Contains(t, args, os.TempDir())
			return nil
		}
		SQLDumpLoad(t, "testdata/*.sql", &SQLDumpOptions{
			DSN:                `cs2:cs2@tcp(localhost:3306)/testDB?parseTime=true&loc=UTC`,
			execCommandContext: exexCmd,
		})
	})

	t.Run("not files found", func(t *testing.T) {
		SQLDumpLoad(testingMock{T: t, wantErr: errors.NotFound}, "testdata/not_files_found", &SQLDumpOptions{
			DSN: `cs2:cs2@tcp(localhost:3306)/testDB?parseTime=true&loc=UTC`,
			execCommandContext: func(ctx context.Context, r io.ReadCloser, cmd string, arg ...string) error {
				return nil
			},
		})
	})

	t.Run("mysql fails", func(t *testing.T) {
		exexCmd := func(ctx context.Context, r io.ReadCloser, cmd string, arg ...string) error {
			return errors.NotImplemented.Newf("Cant handle it")
		}

		SQLDumpLoad(testingMock{T: t, wantErr: errors.NotImplemented}, "testdata/", &SQLDumpOptions{
			DSN:                `cs2:cs2@tcp(localhost:3306)/testDB?parseTime=true&loc=UTC`,
			execCommandContext: exexCmd,
		})
	})
}

type testingMock struct {
	*testing.T
	wantErr errors.Kind
}

func (tm testingMock) Fatalf(format string, args ...any) {
	assert.True(tm.T, tm.wantErr.Match(args[0].(error)))
}
