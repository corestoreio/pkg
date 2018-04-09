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

package ccd

import (
	"context"
	"fmt"
	"time"

	"github.com/corestoreio/errors"
	"github.com/corestoreio/log"
	"github.com/corestoreio/pkg/config"
	"github.com/corestoreio/pkg/sql/dml"
	"github.com/corestoreio/pkg/store/scope"
	"github.com/corestoreio/pkg/util/conv"
)

type Options struct {
	TableName           string
	Log                 log.Logger
	QueryContextTimeout time.Duration
	Idle                struct {
		AllKeys time.Duration
		Read    time.Duration
		Write   time.Duration
	}
}

const (
	stateOpen = iota
	stateClosed
)

// DBStorage connects the MySQL DB with the config.Service type. Implements
// interface config.Storager.
type DBStorage struct {
	config   Options
	preparer dml.Preparer

	stmtAll dml.StmtQuerier
	sqlAll  string

	stmtRead dml.StmtQuerier
	sqlRead  string

	stmtWrite dml.StmtQuerier
	sqlWrite  string
}

// NewDBStorage creates a new pointer with resurrecting prepared SQL statements.
// Default logger for the three underlying ResurrectStmt type sports to black
// hole.
//
// All has an idle time of 15s. Read an idle time of 10s. Write an idle time of
// 30s. Implements interface config.Storager.
func NewDBStorage(p dml.Preparer, o Options) (*DBStorage, error) {

	dbs := &DBStorage{
		config:   o,
		preparer: p,
		sqlAll: fmt.Sprintf(
			"SELECT scope,scope_id,path FROM `%s` ORDER BY scope,scope_id,path",
			TableCollection.Name(TableIndexCoreConfigData),
		),

		//Read: csdb.NewResurrectStmt(p, fmt.Sprintf(
		//	"SELECT `value` FROM `%s` WHERE `scope`=? AND `scope_id`=? AND `path`=?",
		//	TableCollection.Name(TableIndexCoreConfigData),
		//)),
		//
		//Write: csdb.NewResurrectStmt(p, fmt.Sprintf(
		//	"INSERT INTO `%s` (`scope`,`scope_id`,`path`,`value`) VALUES (?,?,?,?) ON DUPLICATE KEY UPDATE `value`=?",
		//	TableCollection.Name(TableIndexCoreConfigData),
		//)),
	}
	if dbs.config.Idle.AllKeys == 0 {
		dbs.config.Idle.AllKeys = time.Second * 5
	}
	if dbs.config.Idle.Read == 0 {
		dbs.config.Idle.Read = time.Second * 20
	}
	if dbs.config.Idle.Write == 0 {
		dbs.config.Idle.Write = time.Second * 5
	}
	if dbs.config.TableName == "" {
		dbs.config.TableName = "core_config_data"
	}

	return dbs, nil
}

// MustNewDBStorage same as NewDBStorage but panics on error. Implements
// interface config.Storager.
func MustNewDBStorage(p dml.Preparer, o Options) *DBStorage {
	s, err := NewDBStorage(p, o)
	if err != nil {
		panic(err)
	}
	return s
}

// Set sets a value with its key. Database errors get logged as Info message.
// Enabled debug level logs the insert ID or rows affected.
func (dbs *DBStorage) Set(scp scope.TypeID, path string, value []byte) error {
	// update lastUsed at the end because there might be the slight chance that
	// a statement gets closed despite we're waiting for the result from the
	// server.
	dbs.Write.StartStmtUse()
	defer dbs.Write.StopStmtUse()

	valStr, err := conv.ToStringE(value)
	if err != nil {
		return errors.Wrapf(err, "[ccd] Set.conv.ToStringE. SQL: %q Key: %q Value: %v", dbs.Write.sqlRaw, key, value)
	}

	stmt, err := dbs.Write.Stmt(context.TODO())
	if err != nil {
		return errors.Wrapf(err, "[ccd] Set.Write.Stmt. SQL: %q Key: %q", dbs.Write.sqlRaw, key)
	}

	pathLeveled, err := key.Level(-1)
	if err != nil {
		return errors.Wrapf(err, "[ccd] Set.key.Level. SQL: %q Key: %q", dbs.Write.sqlRaw, key)
	}

	scp, id := key.ScopeID.Unpack()
	result, err := stmt.Exec(scp.StrType(), id, pathLeveled, valStr, valStr)
	if err != nil {
		return errors.Wrapf(err, "[ccd] Set.stmt.Exec. SQL: %q KeyID: %d Scope: %q Path: %q Value: %q", dbs.Write.sqlRaw, id, scp, pathLeveled, valStr)
	}
	if dbs.log.IsDebug() {
		li, err1 := result.LastInsertId()
		ra, err2 := result.RowsAffected()
		dbs.log.Debug(
			"config.DBStorage.Set.Write.Result",
			log.Int64("lastInsertID", li),
			log.ErrWithKey("lastInsertIDErr", err1),
			log.Int64("rowsAffected", ra),
			log.ErrWithKey("rowsAffectedErr", err2),
			log.String("SQL", dbs.Write.sqlRaw),
			log.Stringer("key", key),
			log.Object("value", value),
		)
	}
	return nil
}

// Get returns a value from the database by its key. It is guaranteed that the
// type in the empty interface is a string. It returns nil on error but errors
// get logged as info message. Error behaviour: NotFound
func (dbs *DBStorage) Value(scp scope.TypeID, path string) (v []byte, ok bool, err error) {
	// update lastUsed at the end because there might be the slight chance that
	// a statement gets closed despite we're waiting for the result from the
	// server.
	dbs.Read.StartStmtUse()
	defer dbs.Read.StopStmtUse()

	stmt, err := dbs.Read.Stmt(context.TODO())
	if err != nil {
		return nil, errors.Wrapf(err, "[ccd] Get.Read.Stmt. SQL: %q Key: %q", dbs.Read.sqlRaw, key)
	}

	pl, err := key.Level(-1)
	if err != nil {
		return nil, errors.Wrapf(err, "[ccd] Get.key.Level. SQL: %q Key: %q", dbs.Read.sqlRaw, key)
	}

	var data null.String
	scp, id := key.ScopeID.Unpack()
	err = stmt.QueryRow(scp.StrType(), id, pl).Scan(&data)
	if err != nil {
		return nil, errors.Wrapf(err, "[ccd] Get.QueryRow. SQL: %q Key: %q PathLevel: %q", dbs.Read.sqlRaw, key, pl)
	}
	if data.Valid {
		return data.String, nil
	}
	return nil, errKeyNotFound
}

// AllKeys returns all available keys. Database errors get logged as info message.
func (dbs *DBStorage) AllKeys() (scps scope.TypeIDs, paths []string, err error) {
	// update lastUsed at the end because there might be the slight chance
	// that a statement gets closed despite we're waiting for the result
	// from the server.
	dbs.All.StartStmtUse()
	defer dbs.All.StopStmtUse()

	stmt, err := dbs.All.Stmt(context.TODO())
	if err != nil {
		return nil, errors.Wrapf(err, "[ccd] AllKeys.All.Stmt. SQL: %q", dbs.All.sqlRaw)
	}

	rows, err := stmt.Query()
	if err != nil {
		return nil, errors.Wrapf(err, "[ccd] AllKeys.All.Query. SQL: %q", dbs.All.sqlRaw)
	}
	defer rows.Close()

	const maxCap = 750 // Just a guess the 750
	var ret = make(config.PathSlice, 0, maxCap)
	var sqlScope null.String
	var sqlScopeID null.Int64
	var sqlPath null.String

	for rows.Next() {
		if err := rows.Scan(&sqlScope, &sqlScopeID, &sqlPath); err != nil {
			return nil, errors.Wrapf(err, "[ccd] AllKeys.rows.Scan. SQL: %q", dbs.All.sqlRaw)
		}
		if sqlPath.Valid {
			p, err := config.MakeByString(sqlPath.String)
			if err != nil {
				return ret, errors.Wrapf(err, "[ccd] AllKeys.rows.config.MakeByString. SQL: %q: Path: %q", dbs.All.sqlRaw, sqlPath.String)
			}
			ret = append(ret, p.Bind(scope.FromString(sqlScope.String).Pack(sqlScopeID.Int64)))
		}
		sqlScope.String = ""
		sqlScope.Valid = false
		sqlScopeID.Int64 = 0
		sqlScopeID.Valid = false
		sqlPath.String = ""
		sqlPath.Valid = false
	}
	return ret, nil
}
