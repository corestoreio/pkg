// Copyright 2015-2016, Cyrill @ Schumacher.fm and the CoreStore contributors
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
	"fmt"
	"time"

	"github.com/corestoreio/csfw/config/cfgpath"
	"github.com/corestoreio/csfw/log"
	"github.com/corestoreio/csfw/storage/csdb"
	"github.com/corestoreio/csfw/storage/dbr"
	"github.com/corestoreio/csfw/store/scope"
	"github.com/corestoreio/csfw/util/conv"
	"github.com/corestoreio/csfw/util/errors"
)

// TableCollection handles all tables and its columns. init() in generated Go file will set the value.
var TableCollection csdb.TableManager

// DBStorage connects the MySQL DB with the config.Service type.
type DBStorage struct {
	log log.Logger
	// All is a SQL statement for the all keys query
	All *csdb.ResurrectStmt
	// Read is a SQL statement for selecting a value from a path/key
	Read *csdb.ResurrectStmt
	// Write statement inserts or updates a value
	Write *csdb.ResurrectStmt
}

// NewDBStorage creates a new pointer with resurrecting prepared SQL statements.
// Default logger for the three underlying ResurrectStmt type sports to black hole.
//
// All has an idle time of 15s. Read an idle time of 10s. Write an idle time of 30s.
func NewDBStorage(p csdb.Preparer) (*DBStorage, error) {
	// todo: instead of logging the error we may write it into an
	// error channel and the gopher who calls NewDBStorage is responsible
	// for continuously reading from the error channel. or we accept an error channel
	// as argument here and then writing to it ...

	dbs := &DBStorage{
		log: log.BlackHole{}, // skip debug and info level via init with empty fields
		All: csdb.NewResurrectStmt(p, fmt.Sprintf(
			"SELECT scope,scope_id,path FROM `%s` ORDER BY scope,scope_id,path",
			TableCollection.Name(TableIndexCoreConfigData),
		)),
		Read: csdb.NewResurrectStmt(p, fmt.Sprintf(
			"SELECT `value` FROM `%s` WHERE `scope`=? AND `scope_id`=? AND `path`=?",
			TableCollection.Name(TableIndexCoreConfigData),
		)),

		Write: csdb.NewResurrectStmt(p, fmt.Sprintf(
			"INSERT INTO `%s` (`scope`,`scope_id`,`path`,`value`) VALUES (?,?,?,?) ON DUPLICATE KEY UPDATE `value`=?",
			TableCollection.Name(TableIndexCoreConfigData),
		)),
	}
	dbs.All.Idle = time.Second * 15
	dbs.All.Log = dbs.log
	dbs.Read.Idle = time.Second * 10
	dbs.Read.Log = dbs.log
	dbs.Write.Idle = time.Second * 30
	dbs.Write.Log = dbs.log
	// in the future we may add errors ... just to have for now the func signature
	return dbs, nil
}

// MustNewDBStorage same as NewDBStorage but panics on error
func MustNewDBStorage(p csdb.Preparer) *DBStorage {
	s, err := NewDBStorage(p)
	if err != nil {
		panic(err)
	}
	return s
}

// SetLogger applies your custom logger
func (dbs *DBStorage) SetLogger(l log.Logger) *DBStorage {
	dbs.log = l
	dbs.All.Log = l
	dbs.Read.Log = l
	dbs.Write.Log = l
	return dbs
}

// Start starts the internal idle time checker for the resurrecting SQL statements.
func (dbs *DBStorage) Start() *DBStorage {
	dbs.All.StartIdleChecker()
	dbs.Read.StartIdleChecker()
	dbs.Write.StartIdleChecker()
	return dbs
}

// Stop stops the internal goroutines for idle time checking. Returns the
// first occurring sql.Stmt.Close() error.
func (dbs *DBStorage) Stop() error {
	if err := dbs.All.StopIdleChecker(); err != nil {
		return errors.Wrap(err, "[ccd] All.StopIdleChecker")
	}
	if err := dbs.Read.StopIdleChecker(); err != nil {
		return errors.Wrap(err, "[ccd] Read.StopIdleChecker")
	}
	if err := dbs.Write.StopIdleChecker(); err != nil {
		return errors.Wrap(err, "[ccd] Write.StopIdleChecker")
	}
	return nil
}

// Set sets a value with its key. Database errors get logged as Info message.
// Enabled debug level logs the insert ID or rows affected.
func (dbs *DBStorage) Set(key cfgpath.Path, value interface{}) error {
	// update lastUsed at the end because there might be the slight chance
	// that a statement gets closed despite we're waiting for the result
	// from the server.
	dbs.Write.StartStmtUse()
	defer dbs.Write.StopStmtUse()

	valStr, err := conv.ToStringE(value)
	if err != nil {
		return errors.Wrapf(err, "[ccd] Set.conv.ToStringE. SQL: %q Key: %q Value: %v", dbs.Write.SQL, key, value)
	}

	stmt, err := dbs.Write.Stmt()
	if err != nil {
		return errors.Wrapf(err, "[ccd] Set.Write.Stmt. SQL: %q Key: %q", dbs.Write.SQL, key)
	}

	pathLeveled, err := key.Level(-1)
	if err != nil {
		return errors.Wrapf(err, "[ccd] Set.key.Level. SQL: %q Key: %q", dbs.Write.SQL, key)
	}

	scp, id := key.ScopeHash.Unpack()
	result, err := stmt.Exec(scp.StrScope(), id, pathLeveled, valStr, valStr)
	if err != nil {
		return errors.Wrapf(err, "[ccd] Set.stmt.Exec. SQL: %q KeyID: %d Scope: %q Path: %q Value: %q", dbs.Write.SQL, id, scp, pathLeveled, valStr)
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
			log.String("SQL", dbs.Write.SQL),
			log.Stringer("key", key),
			log.Object("value", value),
		)
	}
	return nil
}

// Get returns a value from the database by its key. It is guaranteed that the
// type in the empty interface is a string. It returns nil on error but errors
// get logged as info message.
// Error behaviour: NotFound
func (dbs *DBStorage) Get(key cfgpath.Path) (interface{}, error) {
	// update lastUsed at the end because there might be the slight chance
	// that a statement gets closed despite we're waiting for the result
	// from the server.
	dbs.Read.StartStmtUse()
	defer dbs.Read.StopStmtUse()

	stmt, err := dbs.Read.Stmt()
	if err != nil {
		return nil, errors.Wrapf(err, "[ccd] Get.Read.Stmt. SQL: %q Key: %q", dbs.Read.SQL, key)
	}

	pl, err := key.Level(-1)
	if err != nil {
		return nil, errors.Wrapf(err, "[ccd] Get.key.Level. SQL: %q Key: %q", dbs.Read.SQL, key)
	}

	var data dbr.NullString
	scp, id := key.ScopeHash.Unpack()
	err = stmt.QueryRow(scp.StrScope(), id, pl).Scan(&data)
	if err != nil {
		return nil, errors.Wrapf(err, "[ccd] Get.QueryRow. SQL: %q Key: %q PathLevel: %q", dbs.Read.SQL, key, pl)
	}
	if data.Valid {
		return data.String, nil
	}
	return nil, errKeyNotFound
}

var errKeyNotFound = errors.NewNotFoundf(`[ccd] Key not found`) // todo add test

// AllKeys returns all available keys. Database errors get logged as info message.
func (dbs *DBStorage) AllKeys() (cfgpath.PathSlice, error) {
	// update lastUsed at the end because there might be the slight chance
	// that a statement gets closed despite we're waiting for the result
	// from the server.
	dbs.All.StartStmtUse()
	defer dbs.All.StopStmtUse()

	stmt, err := dbs.All.Stmt()
	if err != nil {
		return nil, errors.Wrapf(err, "[ccd] AllKeys.All.Stmt. SQL: %q", dbs.All.SQL)
	}

	rows, err := stmt.Query()
	if err != nil {
		return nil, errors.Wrapf(err, "[ccd] AllKeys.All.Query. SQL: %q", dbs.All.SQL)
	}
	defer rows.Close()

	const maxCap = 750 // Just a guess the 750
	var ret = make(cfgpath.PathSlice, 0, maxCap)
	var sqlScope dbr.NullString
	var sqlScopeID dbr.NullInt64
	var sqlPath dbr.NullString

	for rows.Next() {
		if err := rows.Scan(&sqlScope, &sqlScopeID, &sqlPath); err != nil {
			return nil, errors.Wrapf(err, "[ccd] AllKeys.rows.Scan. SQL: %q", dbs.All.SQL)
		}
		if sqlPath.Valid {
			p, err := cfgpath.NewByParts(sqlPath.String)
			if err != nil {
				return ret, errors.Wrapf(err, "[ccd] AllKeys.rows.cfgpath.NewByParts. SQL: %q: Path: %q", dbs.All.SQL, sqlPath.String)
			}
			ret = append(ret, p.Bind(scope.FromString(sqlScope.String), sqlScopeID.Int64))
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
