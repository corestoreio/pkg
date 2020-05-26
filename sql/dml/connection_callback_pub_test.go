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

package dml_test

import (
	"bytes"
	"context"
	"database/sql/driver"
	"fmt"
	"sync/atomic"
	"testing"

	"github.com/corestoreio/pkg/util/conv"

	"github.com/corestoreio/pkg/sql/dml"
	"github.com/corestoreio/pkg/sql/dmltest"
	"github.com/corestoreio/pkg/util/assert"
)

func TestDriverCallBack(t *testing.T) {
	// Test assumes that the table dml_people does still exists.
	counter := new(int32)

	buf := new(bytes.Buffer)
	dbc := dmltest.MustConnectDB(t, // dbc == database connection
		dml.ConnPoolOption{
			UniqueIDFn: func() string { return fmt.Sprintf("RANJID%d", atomic.AddInt32(counter, 1)) },
		},
		dml.WithDriverCallBack(func(fnName string) func(error, string, []driver.NamedValue) error {
			start := now()
			return func(err error, query string, namedArgs []driver.NamedValue) error {
				fmt.Fprintf(buf, "%q Took: %s\n", fnName, now().Sub(start))
				if err != nil {
					fmt.Fprintf(buf, "Error: %s\n", err)
				}
				if query != "" {
					fmt.Fprintf(buf, "Query: %q\n", query)
				}
				if len(namedArgs) > 0 {
					fmt.Fprintf(buf, "NamedArgs: %#v\n", namedArgs)
				}
				fmt.Fprint(buf, "\n")
				return err
			}
		}),
	)
	installFixtures(t, dbc.DB)

	ctx := context.Background()

	selSQL := dml.NewSelect("*").From("dml_people").Where(dml.Column("name").PlaceHolder())
	err := dbc.RegisterByQueryBuilder(map[string]dml.QueryBuilder{
		"sel1c":  selSQL,
		"sel1nc": selSQL.Clone().SQLNoCache(),
		"up":     dml.NewUpdate("dml_people").AddClauses(dml.Column("name").PlaceHolder()),
	})
	assert.NoError(t, err)

	assert.Exactly(t, []string{
		"sel1c", "/*$ID$RANJID1*/SELECT * FROM `dml_people` WHERE (`name` = ?)",
		"sel1nc", "/*$ID$RANJID2*/SELECT SQL_NO_CACHE * FROM `dml_people` WHERE (`name` = ?)",
		"up", "/*$ID$RANJID3*/UPDATE `dml_people` SET `name`=?",
	}, conv.ToStringSlice(dbc.CachedQueries()))

	var ppl dmlPerson
	_, err = dbc.WithCacheKey("sel1c").Load(ctx, &ppl, "Bernd")
	assert.NoError(t, err)
	_, err = dbc.WithCacheKey("sel1nc").Interpolate().Load(context.Background(), &ppl, "Das Brot")
	assert.NoError(t, err)

	con, err := dbc.Conn(ctx)
	assert.NoError(t, err)

	up := dbc.WithCacheKey("up")
	_, err = up.ExecContext(ctx, "Hugo")
	assert.NoError(t, err)

	_, err = up.Interpolate().ExecContext(ctx, "Bernie")
	assert.NoError(t, err)

	assert.NoError(t, con.Close())
	dmltest.Close(t, dbc)

	assert.MatchesGolden(t, "testdata/TestDriverCallBack.want.txt", buf.Bytes(), false, func(goldenData []byte) []byte {
		return bytes.ReplaceAll(goldenData, []byte(`{{SCHEMA}}`), []byte(dbc.Schema()))
	})
}
