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
	"database/sql"
	"fmt"
	"sync/atomic"
	"testing"

	"github.com/corestoreio/pkg/util/conv"

	"github.com/corestoreio/errors"

	"github.com/corestoreio/log/logw"
	"github.com/corestoreio/pkg/sql/dml"
	"github.com/corestoreio/pkg/sql/dmltest"
	"github.com/corestoreio/pkg/util/assert"
)

func TestWithLogger_Insert(t *testing.T) {
	uniID := new(int32)
	uniqueIDFunc := func() string {
		return fmt.Sprintf("UNIQ%02d", atomic.AddInt32(uniID, 4))
	}

	buf := new(bytes.Buffer)
	lg := logw.NewLog(
		logw.WithLevel(logw.LevelDebug),
		logw.WithWriter(buf),
		logw.WithFlag(0), // no flags at all
	)

	rConn := dmltest.MustConnectDB(t, dml.WithLogger(lg, uniqueIDFunc))
	defer dmltest.Close(t, rConn)
	installFixtures(t, rConn.DB)

	err := rConn.RegisterByQueryBuilder(map[string]dml.QueryBuilder{
		"insert001":  dml.NewInsert("dml_people").Replace().AddColumns("email", "name").BuildValues(),
		"replace001": dml.NewInsert("dml_people").Replace().AddColumns("email", "name").BuildValues(),
	})
	assert.NoError(t, err)

	t.Run("Conn1Pool", func(t *testing.T) {
		t.Run("Prepare", func(t *testing.T) {
			stmt := rConn.WithPrepareCacheKey(context.Background(), "insert001")
			defer dmltest.Close(t, stmt)
		})

		t.Run("Exec", func(t *testing.T) {
			_, err := rConn.WithCacheKey("insert001").Interpolate().ExecContext(context.Background(), "a@b.c", "John")
			assert.NoError(t, err)
		})

		t.Run("Tx Commit", func(t *testing.T) {
			err := rConn.Transaction(context.Background(), nil, func(tx *dml.Tx) error {
				_, err := tx.WithCacheKey("replace001").ExecContext(context.Background(), "a@b.c", "John")
				return err
			})
			assert.NoError(t, err)
		})
	})

	t.Run("Conn2", func(t *testing.T) {
		conn, err := rConn.Conn(context.Background())
		assert.NoError(t, err)

		t.Run("Exec", func(t *testing.T) {
			_, err := conn.WithCacheKey("insert001").Interpolate().ExecContext(context.Background(), "a@b.zeh", "J0hn")
			assert.NoError(t, err)
		})

		t.Run("Prepare", func(t *testing.T) {
			stmt, err := conn.WithCacheKey("insert001").Prepare(context.Background())
			// oIns.IsBuildValues = false
			assert.NoError(t, err)
			defer dmltest.Close(t, stmt)
		})

		t.Run("Prepare Exec", func(t *testing.T) {
			stmt, err := conn.WithCacheKey("insert001").Prepare(context.Background())
			assert.NoError(t, err)
			defer dmltest.Close(t, stmt)

			_, err = stmt.ExecContext(context.Background(), "mail@e.de", "Hans")
			assert.NoError(t, err)
		})

		t.Run("Tx Commit", func(t *testing.T) {
			err := conn.Transaction(context.Background(), nil, func(tx *dml.Tx) error {
				_, err := tx.WithCacheKey("replace001").ExecContext(context.Background(), "a@b.c", "John")
				return err
			})
			assert.NoError(t, err)
		})

		t.Run("Tx Rollback", func(t *testing.T) {
			assert.Error(t, conn.Transaction(context.Background(), nil, func(tx *dml.Tx) error {
				_, err := tx.WithCacheKey("replace001").Interpolate().ExecContext(context.Background(), "only one arg provided")
				return err
			}))
		})

		t.Run("Tx WithDBR", func(t *testing.T) {
			err := conn.Transaction(context.Background(), nil, func(tx *dml.Tx) error {
				_, err := tx.WithCacheKey("replace001").ExecContext(context.Background(), "a@b.c", "John")
				// more queries
				return err
			})
			assert.NoError(t, err)
		})
	})

	assert.Exactly(t, []string{
		"insert001", "/*$ID$UNIQ08*/REPLACE INTO `dml_people` (`email`,`name`) VALUES (?,?)",
		"replace001", "/*$ID$UNIQ12*/REPLACE INTO `dml_people` (`email`,`name`) VALUES (?,?)",
	}, conv.ToStringSlice(rConn.CachedQueries()))

	assert.MatchesGolden(t, "testdata/TestWithLogger_Insert.want.txt", buf.Bytes(), false)
}

func TestWithLogger_Delete(t *testing.T) {
	uniID := new(int32)
	uniqueIDFunc := func() string {
		return fmt.Sprintf("UNIQUEID%02d", atomic.AddInt32(uniID, 1))
	}

	buf := new(bytes.Buffer)
	lg := logw.NewLog(
		logw.WithLevel(logw.LevelDebug),
		logw.WithWriter(buf),
		logw.WithFlag(0), // no flags at all
	)

	rConn := dmltest.MustConnectDB(t, dml.WithLogger(lg, uniqueIDFunc))
	defer dmltest.Close(t, rConn)
	installFixtures(t, rConn.DB)

	err := rConn.RegisterByQueryBuilder(map[string]dml.QueryBuilder{
		"delete01": dml.NewDelete("dml_people").Where(dml.Column("id").GreaterOrEqual().Float64(34.56)),
		"delete02": dml.NewDelete("dml_people").Where(dml.Column("id").GreaterOrEqual().PlaceHolder()),
	})
	assert.NoError(t, err)

	t.Run("Conn01Pool", func(t *testing.T) {
		deleteDBR := rConn.WithCacheKey("delete01")

		t.Run("Exec", func(t *testing.T) {
			_, err := deleteDBR.ExecContext(context.Background())
			assert.NoError(t, err)
		})

		t.Run("Prepare", func(t *testing.T) {
			stmt, err := deleteDBR.Prepare(context.Background())
			assert.NoError(t, err)
			defer dmltest.Close(t, stmt)
		})

		t.Run("Tx Commit", func(t *testing.T) {
			assert.NoError(t, rConn.Transaction(context.Background(), nil, func(tx *dml.Tx) error {
				_, err := tx.WithCacheKey("delete01").ExecContext(context.Background())
				return err
			}))
		})
	})

	t.Run("Conn02", func(t *testing.T) {
		conn, err := rConn.Conn(context.Background())
		assert.NoError(t, err)

		deleteDBR := rConn.WithCacheKey("delete02")

		t.Run("Exec", func(t *testing.T) {
			_, err := deleteDBR.Interpolate().ExecContext(context.Background(), 39.56)
			assert.NoError(t, err)
		})

		t.Run("Prepare Exec", func(t *testing.T) {
			stmt, err := deleteDBR.Prepare(context.Background())
			assert.NoError(t, err)
			defer dmltest.Close(t, stmt)

			_, err = stmt.ExecContext(context.Background(), 41.57)
			assert.NoError(t, err)
		})

		t.Run("Tx Commit", func(t *testing.T) {
			assert.NoError(t, conn.Transaction(context.Background(), nil, func(tx *dml.Tx) error {
				_, err := tx.WithCacheKey("delete01").Interpolate().ExecContext(context.Background())
				return err
			}))
		})

		t.Run("Tx Rollback", func(t *testing.T) {
			assert.Error(t, conn.Transaction(context.Background(), nil, func(tx *dml.Tx) error {
				_, err := tx.WithCacheKey("delete02").Interpolate().ExecContext(context.Background())
				return err
			}))
		})
	})

	assert.Exactly(t, []string{
		"delete01", "/*$ID$UNIQUEID02*/DELETE FROM `dml_people` WHERE (`id` >= 34.56)",
		"delete02", "/*$ID$UNIQUEID03*/DELETE FROM `dml_people` WHERE (`id` >= ?)",
	}, conv.ToStringSlice(rConn.CachedQueries()))

	assert.MatchesGolden(t, "testdata/TestWithLogger_Delete.want.txt", buf.Bytes(), false)
}

func TestWithLogger_Select(t *testing.T) {
	uniID := new(int32)
	uniqueIDFunc := func() string {
		return fmt.Sprintf("UNIQ%02d", atomic.AddInt32(uniID, 1))
	}

	buf := new(bytes.Buffer)
	lg := logw.NewLog(
		logw.WithLevel(logw.LevelDebug),
		logw.WithWriter(buf),
		logw.WithFlag(0), // no flags at all
	)
	rConn := dmltest.MustConnectDB(t, dml.WithLogger(lg, uniqueIDFunc))
	defer dmltest.Close(t, rConn)
	installFixtures(t, rConn.DB)

	err := rConn.RegisterByQueryBuilder(map[string]dml.QueryBuilder{
		"SELECT_mail_people_id_gt":           dml.NewSelect("email").From("dml_people").Where(dml.Column("id").Greater().PlaceHolder()),
		"SELECT_id_people_id_lt":             dml.NewSelect("id").From("dml_people").Where(dml.Column("id").LessOrEqual().PlaceHolder()),
		"SELECT_total_income_people_id_gt":   dml.NewSelect("total_income").From("dml_people").Where(dml.Column("id").Greater().PlaceHolder()),
		"SELECT_email_name_people_id_in79":   dml.NewSelect("name", "email").From("dml_people").Where(dml.Column("id").In().Int64s(7, 9)),
		"SELECT_email_name_people_id_in7191": dml.NewSelect("name", "email").From("dml_people").Where(dml.Column("id").In().Int64s(71, 91)),
		"SELECT_email_name_people_id_ph":     dml.NewSelect("name", "email").From("dml_people").Where(dml.Column("id").In().PlaceHolder()),
		"SELECT_email_name_dp2_id_lt":        dml.NewSelect("name", "email").FromAlias("dml_people", "dp2").Where(dml.Column("id").Less().PlaceHolder()),
	})
	assert.NoError(t, err)

	t.Run("ConnPool", func(t *testing.T) {
		pplSel := rConn.WithCacheKey("SELECT_mail_people_id_gt")

		t.Run("Query Error interpolation with iFace slice", func(t *testing.T) {
			rows, err := pplSel.Interpolate().QueryContext(context.Background(), 67896543123)
			assert.NotNil(t, rows)
			assert.NoError(t, err)
		})
		t.Run("Query Correct", func(t *testing.T) {
			rows, err := pplSel.QueryContext(context.Background(), 67896543123)
			assert.NoError(t, err)
			assert.NoError(t, rows.Close())
		})

		t.Run("Load", func(t *testing.T) {
			p := &dmlPerson{}
			_, err := pplSel.Load(context.Background(), p, 67896543113)
			assert.NoError(t, err)
		})

		pplSel = rConn.WithCacheKey("SELECT_id_people_id_lt")

		t.Run("LoadInt64", func(t *testing.T) {
			_, _, err := pplSel.LoadNullInt64(context.Background(), 67896543124)
			if !errors.NotFound.Match(err) {
				assert.NoError(t, err)
			}
		})

		t.Run("LoadInt64s", func(t *testing.T) {
			_, err := pplSel.LoadInt64s(context.Background(), nil, 67896543125)
			assert.NoError(t, err)
		})

		t.Run("LoadUint64", func(t *testing.T) {
			_, _, err := pplSel.LoadNullUint64(context.Background(), 67896543126)
			if !errors.NotFound.Match(err) {
				assert.NoError(t, err)
			}
		})

		t.Run("LoadUint64s", func(t *testing.T) {
			_, err := pplSel.LoadUint64s(context.Background(), nil, 67896543127)
			assert.NoError(t, err)
		})

		pplSel = rConn.WithCacheKey("SELECT_total_income_people_id_gt")

		t.Run("LoadFloat64", func(t *testing.T) {
			_, _, err := pplSel.LoadNullFloat64(context.Background(), 678965.43125)
			if !errors.NotFound.Match(err) {
				assert.NoError(t, err)
			}
		})

		t.Run("LoadFloat64s", func(t *testing.T) {
			_, err := pplSel.LoadFloat64s(context.Background(), nil, 6789654.3125)
			assert.NoError(t, err)
		})

		pplSel = rConn.WithCacheKey("SELECT_mail_people_id_gt")

		t.Run("LoadString", func(t *testing.T) {
			_, _, err := pplSel.LoadNullString(context.Background(), "hello")
			if !errors.NotFound.Match(err) {
				assert.NoError(t, err)
			}
		})

		t.Run("LoadStrings", func(t *testing.T) {
			_, err := pplSel.LoadStrings(context.Background(), nil, 99987)
			assert.NoError(t, err)
		})

		t.Run("Prepare", func(t *testing.T) {
			stmt, err := pplSel.Prepare(context.Background())
			assert.NoError(t, err)
			defer dmltest.Close(t, stmt)
		})

		t.Run("Tx Commit", func(t *testing.T) {
			assert.NoError(t, rConn.Transaction(context.Background(), nil, func(tx *dml.Tx) error {
				rows, err := tx.WithCacheKey("SELECT_email_name_people_id_in79").QueryContext(context.Background())
				assert.NoError(t, err)
				return rows.Close()
			}))
		})
	})

	t.Run("ConnSingle", func(t *testing.T) {
		conn, err := rConn.Conn(context.Background())
		defer dmltest.Close(t, conn)
		assert.NoError(t, err)

		pplSel := conn.WithCacheKey("SELECT_email_name_dp2_id_lt")

		t.Run("Query", func(t *testing.T) {
			rows, err := pplSel.QueryContext(context.Background(), -3)
			assert.NoError(t, err)
			dmltest.Close(t, rows)
		})

		t.Run("Load", func(t *testing.T) {
			p := &dmlPerson{}
			_, err := pplSel.Load(context.Background(), p, -2)
			assert.NoError(t, err)
		})

		t.Run("Prepare", func(t *testing.T) {
			stmt, err := pplSel.Prepare(context.Background())
			assert.NoError(t, err)
			defer dmltest.Close(t, stmt)

			t.Run("QueryRow", func(t *testing.T) {
				rows := stmt.QueryRowContext(context.Background(), -8)
				var x string
				err := rows.Scan(&x)
				assert.True(t, errors.Cause(err) == sql.ErrNoRows, "but got this error: %#v", err)
				_ = x
			})

			t.Run("Query", func(t *testing.T) {
				rows, err := stmt.QueryContext(context.Background(), -4)
				assert.NoError(t, err)
				dmltest.Close(t, rows)
			})

			t.Run("Load", func(t *testing.T) {
				p := &dmlPerson{}
				_, err := stmt.Load(context.Background(), p, -6)
				assert.NoError(t, err)
			})

			t.Run("LoadInt64", func(t *testing.T) {
				_, _, err := stmt.LoadNullInt64(context.Background(), -7)
				if !errors.NotFound.Match(err) {
					assert.NoError(t, err)
				}
			})

			t.Run("LoadInt64s", func(t *testing.T) {
				iSl, err := stmt.LoadInt64s(context.Background(), nil, -7)
				assert.NoError(t, err)
				assert.Nil(t, iSl)
			})
		})

		t.Run("Tx Commit", func(t *testing.T) {
			assert.NoError(t, conn.Transaction(context.Background(), nil, func(tx *dml.Tx) error {
				rows, err := tx.WithCacheKey("SELECT_email_name_people_id_in7191").QueryContext(context.Background())
				if err != nil {
					return err
				}
				return rows.Close()
			}))
		})

		t.Run("Tx Rollback", func(t *testing.T) {
			assert.Error(t, conn.Transaction(context.Background(), nil, func(tx *dml.Tx) error {
				rows, err := tx.WithCacheKey("SELECT_email_name_people_id_ph").QueryContext(context.Background())
				if err != nil {
					return err
				}
				return rows.Close()
			}))
		})
	})

	assert.Exactly(t, []string{
		"SELECT_email_name_dp2_id_lt", "/*$ID$UNIQ02*/SELECT `name`, `email` FROM `dml_people` AS `dp2` WHERE (`id` < ?)",
		"SELECT_email_name_people_id_in7191", "/*$ID$UNIQ03*/SELECT `name`, `email` FROM `dml_people` WHERE (`id` IN (71,91))",
		"SELECT_email_name_people_id_in79", "/*$ID$UNIQ04*/SELECT `name`, `email` FROM `dml_people` WHERE (`id` IN (7,9))",
		"SELECT_email_name_people_id_ph", "/*$ID$UNIQ05*/SELECT `name`, `email` FROM `dml_people` WHERE (`id` IN ?)",
		"SELECT_id_people_id_lt", "/*$ID$UNIQ06*/SELECT `id` FROM `dml_people` WHERE (`id` <= ?)",
		"SELECT_mail_people_id_gt", "/*$ID$UNIQ07*/SELECT `email` FROM `dml_people` WHERE (`id` > ?)",
		"SELECT_total_income_people_id_gt", "/*$ID$UNIQ08*/SELECT `total_income` FROM `dml_people` WHERE (`id` > ?)",
	}, conv.ToStringSlice(rConn.CachedQueries()))

	assert.MatchesGolden(t, "testdata/TestWithLogger_Select.want.txt", buf.Bytes(), false)
}

func TestWithLogger_Union(t *testing.T) {
	uniID := new(int32)
	uniqueIDFunc := func() string {
		return fmt.Sprintf("UNIQ%02d", atomic.AddInt32(uniID, 1))
	}

	buf := new(bytes.Buffer)
	lg := logw.NewLog(
		logw.WithLevel(logw.LevelDebug),
		logw.WithWriter(buf),
		logw.WithFlag(0), // no flags at all
	)

	rConn := dmltest.MustConnectDB(t, dml.WithLogger(lg, uniqueIDFunc))
	defer dmltest.Close(t, rConn)
	installFixtures(t, rConn.DB)

	t.Run("ConnPool", func(t *testing.T) {
		u := rConn.WithQueryBuilder(dml.NewUnion(
			dml.NewSelect("name").AddColumnsAliases("email", "email").From("dml_people"),
			dml.NewSelect("name", "email").FromAlias("dml_people", "dp2").Where(dml.Column("id").In().Int64s(6, 8)),
		))

		t.Run("Query", func(t *testing.T) {
			rows, err := u.QueryContext(context.Background())
			assert.NoError(t, err)
			assert.NoError(t, rows.Close())
		})

		t.Run("Load", func(t *testing.T) {
			p := &dmlPerson{}
			_, err := u.Interpolate().Load(context.Background(), p)
			assert.NoError(t, err)
		})

		t.Run("Prepare", func(t *testing.T) {
			stmt, err := u.Prepare(context.Background())
			assert.NoError(t, err)
			defer dmltest.Close(t, stmt)
		})

		t.Run("Tx Commit", func(t *testing.T) {
			assert.NoError(t, rConn.Transaction(context.Background(), nil, func(tx *dml.Tx) error {
				rows, err := tx.WithQueryBuilder(dml.NewUnion(
					dml.NewSelect("name").AddColumnsAliases("email", "email").From("dml_people"),
					dml.NewSelect("name", "email").FromAlias("dml_people", "dp2").Where(dml.Column("id").In().Int64s(7, 9)),
				)).Interpolate().QueryContext(context.Background())
				assert.NoError(t, err)
				assert.NoError(t, rows.Close())
				return err
			}))
		})
	})

	t.Run("Conn", func(t *testing.T) {
		conn, err := rConn.Conn(context.Background())
		assert.NoError(t, err)

		u := conn.WithQueryBuilder(dml.NewUnion(
			dml.NewSelect("name").AddColumnsAliases("email", "email").From("dml_people"),
			dml.NewSelect("name", "email").FromAlias("dml_people", "dp2").Where(dml.Column("id").In().Int64s(61, 81)),
		))
		t.Run("Query", func(t *testing.T) {
			rows, err := u.Interpolate().QueryContext(context.Background())
			assert.NoError(t, err)
			assert.NoError(t, rows.Close())
		})

		t.Run("Load", func(t *testing.T) {
			p := &dmlPerson{}
			_, err := u.Load(context.Background(), p)
			assert.NoError(t, err)
		})

		t.Run("Prepare", func(t *testing.T) {
			stmt, err := u.Prepare(context.Background())
			assert.NoError(t, err)
			defer dmltest.Close(t, stmt)
		})

		t.Run("Tx Commit", func(t *testing.T) {
			assert.NoError(t, conn.Transaction(context.Background(), nil, func(tx *dml.Tx) error {
				rows, err := tx.WithQueryBuilder(dml.NewUnion(
					dml.NewSelect("name").AddColumnsAliases("email", "email").From("dml_people"),
					dml.NewSelect("name", "email").FromAlias("dml_people", "dp2").Where(dml.Column("id").In().Int64s(71, 91)),
				)).Interpolate().QueryContext(context.Background())
				if err != nil {
					return err
				}
				return rows.Close()
			}))
		})

		t.Run("Tx Rollback", func(t *testing.T) {
			assert.Error(t, conn.Transaction(context.Background(), nil, func(tx *dml.Tx) error {
				rows, err := tx.WithQueryBuilder(dml.NewUnion(
					dml.NewSelect("name").AddColumnsAliases("email", "email").From("dml_people"),
					dml.NewSelect("name", "email").FromAlias("dml_people", "dp2").Where(dml.Column("id").In().PlaceHolder()),
				)).Interpolate().QueryContext(context.Background())
				if err != nil {
					return err
				}
				return rows.Close()
			}))
		})
	})

	assert.Exactly(t, []string{
		"(SELECT448f828630cbf059", "/*$ID$UNIQ04*/(SELECT `name`, `email` AS `email` FROM `dml_people`)\nUNION\n(SELECT `name`, `email` FROM `dml_people` AS `dp2` WHERE (`id` IN (7,9)))",
		"(SELECT51bba5323b945cd9", "/*$ID$UNIQ08*/(SELECT `name`, `email` AS `email` FROM `dml_people`)\nUNION\n(SELECT `name`, `email` FROM `dml_people` AS `dp2` WHERE (`id` IN (71,91)))",
		"(SELECT7d4c08cb1f95bf31", "/*$ID$UNIQ02*/(SELECT `name`, `email` AS `email` FROM `dml_people`)\nUNION\n(SELECT `name`, `email` FROM `dml_people` AS `dp2` WHERE (`id` IN (6,8)))",
		"(SELECT939348261b065259", "/*$ID$UNIQ10*/(SELECT `name`, `email` AS `email` FROM `dml_people`)\nUNION\n(SELECT `name`, `email` FROM `dml_people` AS `dp2` WHERE (`id` IN ?))",
		"(SELECTc4278b06f5c4a833", "/*$ID$UNIQ06*/(SELECT `name`, `email` AS `email` FROM `dml_people`)\nUNION\n(SELECT `name`, `email` FROM `dml_people` AS `dp2` WHERE (`id` IN (61,81)))",
	}, conv.ToStringSlice(rConn.CachedQueries()))

	assert.MatchesGolden(t, "testdata/TestWithLogger_Union.want.txt", buf.Bytes(), false)
}

func TestWithLogger_Update(t *testing.T) {
	uniID := new(int32)
	uniqueIDFunc := func() string {
		return fmt.Sprintf("UNIQ%02d", atomic.AddInt32(uniID, 3))
	}

	buf := new(bytes.Buffer)
	lg := logw.NewLog(
		logw.WithLevel(logw.LevelDebug),
		logw.WithWriter(buf),
		logw.WithFlag(0), // no flags at all
	)

	rConn := dmltest.MustConnectDB(t, dml.WithLogger(lg, uniqueIDFunc))
	defer dmltest.Close(t, rConn)
	installFixtures(t, rConn.DB)

	t.Run("ConnPool", func(t *testing.T) {
		d := rConn.WithQueryBuilder(dml.NewUpdate("dml_people").AddClauses(
			dml.Column("email").Str("new@email.com"),
		).Where(dml.Column("id").GreaterOrEqual().Float64(78.31)))

		t.Run("Exec", func(t *testing.T) {
			_, err := d.ExecContext(context.Background())
			assert.NoError(t, err)
		})

		t.Run("Prepare", func(t *testing.T) {
			stmt, err := d.Prepare(context.Background())
			assert.NoError(t, err)
			defer dmltest.Close(t, stmt)
		})

		t.Run("Tx Commit", func(t *testing.T) {
			assert.NoError(t, rConn.Transaction(context.Background(), nil, func(tx *dml.Tx) error {
				_, err := tx.WithQueryBuilder(dml.NewUpdate("dml_people").AddClauses(
					dml.Column("email").Str("new@email.com"),
				).Where(dml.Column("id").GreaterOrEqual().Float64(36.56))).ExecContext(context.Background())
				return err
			}))
		})
	})

	t.Run("Conn", func(t *testing.T) {
		conn, err := rConn.Conn(context.Background())
		assert.NoError(t, err)

		d := conn.WithQueryBuilder(dml.NewUpdate("dml_people").AddClauses(
			dml.Column("email").Str("new@email.com"),
		).Where(dml.Column("id").GreaterOrEqual().Float64(21.56)))

		t.Run("Exec", func(t *testing.T) {
			_, err := d.ExecContext(context.Background())
			assert.NoError(t, err)
		})

		t.Run("Prepare Exec", func(t *testing.T) {
			stmt, err := d.Prepare(context.Background())
			assert.NoError(t, err)
			defer dmltest.Close(t, stmt)

			_, err = stmt.ExecContext(context.Background())
			assert.NoError(t, err)
		})

		t.Run("Tx Commit", func(t *testing.T) {
			assert.NoError(t, conn.Transaction(context.Background(), nil, func(tx *dml.Tx) error {
				_, err := tx.WithQueryBuilder(dml.NewUpdate("dml_people").AddClauses(
					dml.Column("email").Str("new@email.com"),
				).Where(dml.Column("id").GreaterOrEqual().Float64(39.56))).ExecContext(context.Background())
				return err
			}))
		})

		t.Run("Tx Rollback", func(t *testing.T) {
			assert.Error(t, conn.Transaction(context.Background(), nil, func(tx *dml.Tx) error {
				_, err := tx.WithQueryBuilder(dml.NewUpdate("dml_people").AddClauses(
					dml.Column("email").Str("new@email.com"),
				).Where(dml.Column("id").GreaterOrEqual().PlaceHolder())).ExecContext(context.Background())
				return err
			}))
		})
	})

	assert.Exactly(t, []string{
		"UPDATE0e53a96824343e63", "/*$ID$UNIQ18*/UPDATE `dml_people` SET `email`='new@email.com' WHERE (`id` >= 21.56)",
		"UPDATE43a36f841d14149e", "/*$ID$UNIQ24*/UPDATE `dml_people` SET `email`='new@email.com' WHERE (`id` >= 39.56)",
		"UPDATE81232c5a0bb3c1ac", "/*$ID$UNIQ06*/UPDATE `dml_people` SET `email`='new@email.com' WHERE (`id` >= 78.31)",
		"UPDATEbca8ee5774af2a17", "/*$ID$UNIQ12*/UPDATE `dml_people` SET `email`='new@email.com' WHERE (`id` >= 36.56)",
		"UPDATEc2e57118892eb042", "/*$ID$UNIQ30*/UPDATE `dml_people` SET `email`='new@email.com' WHERE (`id` >= ?)",
	}, conv.ToStringSlice(rConn.CachedQueries()))

	assert.MatchesGolden(t, "testdata/TestWithLogger_Update.want.txt", buf.Bytes(), false)
}

func TestWithLogger_WithCTE(t *testing.T) {
	uniID := new(int32)
	uniqueIDFunc := func() string {
		return fmt.Sprintf("UNIQ%02d", atomic.AddInt32(uniID, 2))
	}

	buf := new(bytes.Buffer)
	lg := logw.NewLog(
		logw.WithLevel(logw.LevelDebug),
		logw.WithWriter(buf),
		logw.WithFlag(0), // no flags at all
	)

	rConn := dmltest.MustConnectDB(t, dml.WithLogger(lg, uniqueIDFunc))
	defer dmltest.Close(t, rConn)
	installFixtures(t, rConn.DB)

	cte := dml.WithCTE{
		Name:    "zehTeEh",
		Columns: []string{"name2", "email2"},
		Union: dml.NewUnion(
			dml.NewSelect("name").AddColumnsAliases("email", "email").From("dml_people"),
			dml.NewSelect("name", "email").FromAlias("dml_people", "dp2").Where(dml.Column("id").In().Int64s(6, 8)),
		).All(),
	}
	cteSel := dml.NewSelect().Star().From("zehTeEh")

	t.Run("ConnPool", func(t *testing.T) {
		wth := rConn.WithQueryBuilder(dml.NewWith(cte).Select(cteSel))

		t.Run("Query", func(t *testing.T) {
			rows, err := wth.Interpolate().QueryContext(context.Background())
			assert.NoError(t, err)
			assert.NoError(t, rows.Close())
		})

		t.Run("Load", func(t *testing.T) {
			p := &dmlPerson{}
			_, err := wth.Interpolate().Load(context.Background(), p)
			assert.NoError(t, err)
		})

		t.Run("Prepare", func(t *testing.T) {
			stmt, err := wth.Prepare(context.Background())
			assert.NoError(t, err)
			defer dmltest.Close(t, stmt)
		})

		t.Run("Tx Commit", func(t *testing.T) {
			assert.NoError(t, rConn.Transaction(context.Background(), nil, func(tx *dml.Tx) error {
				rows, err := tx.WithQueryBuilder(dml.NewWith(
					dml.WithCTE{
						Name:    "zehTeEh",
						Columns: []string{"name2", "email2"},
						Union: dml.NewUnion(
							dml.NewSelect("name").AddColumnsAliases("email", "email").From("dml_people"),
							dml.NewSelect("name", "email").FromAlias("dml_people", "dp2").Where(dml.Column("id").In().Int64s(6, 8)),
						).All(),
					},
				).Recursive().
					Select(dml.NewSelect().Star().From("zehTeEh"))).Interpolate().QueryContext(context.Background())

				assert.NoError(t, err)
				return rows.Close()
			}))
		})
	})

	t.Run("Conn", func(t *testing.T) {
		conn, err := rConn.Conn(context.Background())
		assert.NoError(t, err)

		u := conn.WithQueryBuilder(dml.NewWith(cte).Select(cteSel))

		t.Run("Query", func(t *testing.T) {
			rows, err := u.Interpolate().QueryContext(context.Background())
			assert.NoError(t, err)
			assert.NoError(t, rows.Close())
		})

		t.Run("Load", func(t *testing.T) {
			p := &dmlPerson{}
			_, err := u.Load(context.Background(), p)
			assert.NoError(t, err)
		})

		t.Run("Prepare", func(t *testing.T) {
			stmt, err := u.Prepare(context.Background())
			assert.NoError(t, err)
			defer dmltest.Close(t, stmt)
		})

		t.Run("Tx Commit", func(t *testing.T) {
			assert.NoError(t, conn.Transaction(context.Background(), nil, func(tx *dml.Tx) error {
				rows, err := tx.WithQueryBuilder(dml.NewWith(cte).Select(cteSel)).QueryContext(context.Background())
				if err != nil {
					return err
				}
				return rows.Close()
			}))
		})

		t.Run("Tx Rollback", func(t *testing.T) {
			assert.Error(t, conn.Transaction(context.Background(), nil, func(tx *dml.Tx) error {
				rows, err := tx.WithQueryBuilder(
					dml.NewWith(cte).Select(cteSel.Where(dml.Column("email").In().PlaceHolder())),
				).QueryContext(context.Background())
				if err != nil {
					return err
				}
				return rows.Close()
			}))
		})
	})
	assert.Exactly(t, []string{
		"WITH46e1ea15e42d9f39", "/*$ID$UNIQ16*/WITH `zehTeEh` (`name2`,`email2`) AS ((SELECT `name`, `email` AS `email` FROM `dml_people`)\nUNION ALL\n(SELECT `name`, `email` FROM `dml_people` AS `dp2` WHERE (`id` IN (6,8))))\nSELECT * FROM `zehTeEh` WHERE (`email` IN ?)",
		"WITH9118b6dd9042bd4f", "/*$ID$UNIQ04*/WITH `zehTeEh` (`name2`,`email2`) AS ((SELECT `name`, `email` AS `email` FROM `dml_people`)\nUNION ALL\n(SELECT `name`, `email` FROM `dml_people` AS `dp2` WHERE (`id` IN (6,8))))\nSELECT * FROM `zehTeEh`",
		"WITHf581a58421d8e155", "/*$ID$UNIQ08*/WITH RECURSIVE `zehTeEh` (`name2`,`email2`) AS ((SELECT `name`, `email` AS `email` FROM `dml_people`)\nUNION ALL\n(SELECT `name`, `email` FROM `dml_people` AS `dp2` WHERE (`id` IN (6,8))))\nSELECT * FROM `zehTeEh`",
	}, conv.ToStringSlice(rConn.CachedQueries()))

	assert.MatchesGolden(t, "testdata/TestWithLogger_WithCTE.want.txt", buf.Bytes(), false)
}
