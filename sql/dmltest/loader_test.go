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
	"testing"

	"github.com/corestoreio/errors"
	"github.com/corestoreio/pkg/util/assert"
)

func TestSQLDumpLoad(t *testing.T) {
	t.Parallel()

	t.Run("load sql files", func(t *testing.T) {

		exexCmd := func(_ context.Context, name string, arg ...string) error {
			assert.Exactly(t, `mysql`, name)

			fileName := arg[len(arg)-1]
			assert.Contains(t, fileName, "testdata/product_entity_0")

			arg = arg[:len(arg)-1]
			assert.Exactly(t, []string{
				"--host", "localhost:3306", "--user", "cs2", "--password", "cs2", "--database", "testDB", "<",
			}, arg)
			return nil
		}

		err := SQLDumpLoad(context.TODO(), `cs2:cs2@tcp(localhost:3306)/testDB?parseTime=true&loc=UTC`, "testdata/*.sql", SQLDumpOptions{
			execCommandContext: exexCmd,
		})
		assert.NoError(t, err)
	})

	t.Run("mysql fails", func(t *testing.T) {

		exexCmd := func(_ context.Context, name string, arg ...string) error {
			return errors.NotImplemented.Newf("Cant handle it")
		}

		err := SQLDumpLoad(context.TODO(), `cs2:cs2@tcp(localhost:3306)/testDB?parseTime=true&loc=UTC`, "testdata/", SQLDumpOptions{
			execCommandContext: exexCmd,
		})
		assert.True(t, errors.NotImplemented.Match(err), "%+v", err)
	})

}
