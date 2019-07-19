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
	"io/ioutil"
	"sync/atomic"
	"testing"

	"github.com/corestoreio/pkg/sql/dml"
	"github.com/corestoreio/pkg/sql/dmltest"
	"github.com/corestoreio/pkg/util/assert"
)

func TestDriverCallBack(t *testing.T) {
	// Test assumes that the table dml_people does still exists.
	counter := new(int32)

	buf := new(bytes.Buffer)
	db := dmltest.MustConnectDB(t,
		dml.ConnPoolOption{
			UniqueIDFn: func() string { return fmt.Sprintf("RANJID%d", atomic.AddInt32(counter, 1)) },
		},
		dml.WithDSN(
			dmltest.MustGetDSN(t),
			func(fnName string) func(error, string, []driver.NamedValue) error {
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
			},
		))

	ctx := context.TODO()
	sel := db.SelectFrom("dml_people").Star().Where(dml.Column("name").PlaceHolder())
	var ppl dmlPerson
	_, err := sel.WithArgs().String("Bernd").Load(ctx, &ppl)
	assert.NoError(t, err)

	_, err = sel.WithCacheKey("NoCache").SQLNoCache().WithArgs().Interpolate().String("Das Brot").Load(context.Background(), &ppl)
	assert.NoError(t, err)

	con, err := db.Conn(ctx)
	assert.NoError(t, err)

	upd := con.Update("dml_people").AddClauses(dml.Column("name").PlaceHolder())
	_, err = upd.WithArgs().ExecContext(ctx, "Hugo")
	assert.NoError(t, err)

	_, err = upd.WithArgs().String("Bernie").Interpolate().ExecContext(ctx)
	assert.NoError(t, err)

	assert.NoError(t, con.Close())

	dmltest.Close(t, db)
	// t.Log(buf.String())
	// ioutil.WriteFile("testdata/TestDriverCallBack.want2.txt", buf.Bytes(), 0644)
	wantLog, err := ioutil.ReadFile("testdata/TestDriverCallBack.want.txt")
	assert.NoError(t, err)
	if !bytes.Equal(wantLog, buf.Bytes()) {
		t.Error("testdata/TestDriverCallBack.want.txt does not match with `have`.")
	}
}
