// Copyright 2015, Cyrill @ Schumacher.fm and the CoreStore contributors
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

package config

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/corestoreio/csfw/config/scope"
	"github.com/corestoreio/csfw/storage/dbr"
	"github.com/corestoreio/csfw/utils/cast"
	"github.com/juju/errgo"
)

var _ Storager = (*DBStorage)(nil)

type stmtUsage struct {
	SQL        string
	stmt       *sql.Stmt
	closed     bool
	Idle       time.Duration
	lastUsed   time.Time
	inUse      bool
	stopTicker chan struct{}
	closeError error
}

func (su *stmtUsage) close() {
	if su.stmt == nil {
		return
	}
	if su.closeError = errgo.Mask(su.stmt.Close()); su.closeError != nil {
		PkgLog.Info("config.StmtUsage.stmt.Close.error", "err", su.closeError, "SQL", su.SQL)
	} else {
		su.closed = true
	}
	if PkgLog.IsDebug() {
		PkgLog.Debug("config.StmtUsage.stmt.Close", "SQL", su.SQL)
	}
}

func (su *stmtUsage) checkUsage() {
	ticker := time.NewTicker(su.Idle)
	for {
		// maybe squeeze all three go routines into one. for each statement one select case.
		select {
		case t, ok := <-ticker.C:
			if !ok {
				// todo maybe debug log?
				return
			}
			if t.After(su.lastUsed) && !su.closed && !su.inUse {
				// stmt has not been used within the last x seconds.
				// so close the stmt and release the resources in the DB.
				su.close()
			}
		case <-su.stopTicker:
			ticker.Stop()
			su.close()
			return
		}
	}
}

func (su *stmtUsage) getStmt(db *sql.DB) (*sql.Stmt, error) {
	if false == su.closed {
		return su.stmt, nil
	}
	var err error
	su.stmt, err = db.Prepare(su.SQL)
	if err != nil {
		return nil, errgo.Mask(err)
	}
	if PkgLog.IsDebug() {
		PkgLog.Debug("config.StmtUsage.stmt.Prepare", "SQL", su.SQL)
	}
	su.closed = false
	return su.stmt, nil
}

func (su *stmtUsage) startUse() {
	su.lastUsed = time.Now()
	su.inUse = true
}

func (su *stmtUsage) stopUse() {
	su.lastUsed = time.Now()
	su.inUse = false
}

type DBStorage struct {
	db *sql.DB
	// All is a SQL statement for the all keys query
	All *stmtUsage
	// Read is a SQL statement for selecting a value from a path/key
	Read *stmtUsage
	// Write statement inserts or updates a value
	Write *stmtUsage
}

func NewDBStorage(db *sql.DB) *DBStorage {
	// idea: as this is a long running service we should have
	// two prepared statements for select, for insert and for all keys.
	// After time x in which nothing happens neither select nor
	// insert nor an update the prepared statement gets closed
	// and once there is a new action then we recreate a prepared
	// statement.
	dbs := &DBStorage{
		db: db,
		All: &stmtUsage{
			SQL: fmt.Sprintf(
				"SELECT CONCAT(scope,'%s',scope_id,'%s',path) AS `fqpath` FROM `%s` ORDER BY scope,scope_id,path",
				scope.PS,
				scope.PS,
				TableCollection.Name(TableIndexCoreConfigData),
			),
			Idle:       time.Second * 15,
			stopTicker: make(chan struct{}),
			closed:     true,
		},
		Read: &stmtUsage{
			SQL: fmt.Sprintf(
				"SELECT `value` FROM `%s` WHERE `scope`=? AND `scope_id`=? AND `path`=?",
				TableCollection.Name(TableIndexCoreConfigData),
			),
			Idle:       time.Second * 10,
			stopTicker: make(chan struct{}),
			closed:     true,
		},
		Write: &stmtUsage{
			SQL: fmt.Sprintf(
				"INSERT INTO `%s` (`scope`,`scope_id`,`path`,`value`) VALUES (?,?,?,?) ON DUPLICATE KEY UPDATE `value`=?",
				TableCollection.Name(TableIndexCoreConfigData),
			),
			Idle:       time.Second * 30,
			stopTicker: make(chan struct{}),
			closed:     true,
		},
	}
	return dbs
}

func (dbs *DBStorage) Start() *DBStorage {
	go dbs.All.checkUsage()
	go dbs.Read.checkUsage()
	go dbs.Write.checkUsage()
	return dbs
}

func (dbs *DBStorage) Stop() (err error) {
	dbs.All.stopTicker <- struct{}{}
	dbs.Read.stopTicker <- struct{}{}
	dbs.Write.stopTicker <- struct{}{}
	if dbs.All.closeError != nil {
		return dbs.All.closeError
	}
	if dbs.Read.closeError != nil {
		return dbs.Read.closeError
	}
	if dbs.Write.closeError != nil {
		return dbs.Write.closeError
	}
	return nil
}

func (dbs *DBStorage) Set(key string, value interface{}) {
	// update lastUsed at the end because there might be the slight chance
	// that a statement gets closed despite we're waiting for the result
	// from the server.
	dbs.Write.startUse()
	defer dbs.Write.stopUse()

	valStr, err := cast.ToStringE(value)
	if err != nil {
		PkgLog.Info("config.DBStorage.Set.ToString", "err", err, "SQL", dbs.Write.SQL, "value", value)
		return
	}

	stmt, err := dbs.Write.getStmt(dbs.db)
	if err != nil {
		PkgLog.Info("config.DBStorage.Set.Write.getStmt", "err", err, "SQL", dbs.Write.SQL)
		return
	}

	scope, scopeID, path, err := scope.SplitFQPath(key)
	if err != nil {
		PkgLog.Info("config.DBStorage.Set.ReverseFQPath", "err", err, "key", key)
		return
	}

	result, err := stmt.Exec(scope, scopeID, path, valStr, valStr)
	if err != nil {
		PkgLog.Info("config.DBStorage.Set.Write.Exec", "err", err, "SQL", dbs.Write.SQL, "key", key, "value", value)
		return
	}
	if PkgLog.IsDebug() {
		li, err1 := result.LastInsertId()
		ra, err2 := result.RowsAffected()
		PkgLog.Info("config.DBStorage.Set.Write.Result", "lastInsertID", li, "lastInsertIDErr", err1, "rowsAffected", ra, "rowsAffectedErr", err2, "SQL", dbs.Write.SQL, "key", key, "value", value)
	}
}

func (dbs *DBStorage) Get(key string) interface{} {
	// update lastUsed at the end because there might be the slight chance
	// that a statement gets closed despite we're waiting for the result
	// from the server.
	dbs.Read.startUse()
	defer dbs.Read.stopUse()

	stmt, err := dbs.Read.getStmt(dbs.db)
	if err != nil {
		PkgLog.Info("config.DBStorage.Get.Read.getStmt", "err", err, "SQL", dbs.Read.SQL)
		return nil
	}

	scope, scopeID, path, err := scope.SplitFQPath(key)
	if err != nil {
		PkgLog.Info("config.DBStorage.Get.ReverseFQPath", "err", err, "key", key)
		return nil
	}

	var data dbr.NullString
	err = stmt.QueryRow(scope, scopeID, path).Scan(&data)
	if err != nil {
		PkgLog.Info("config.DBStorage.Get.QueryRow", "err", err, "key", key)
		return nil
	}
	if data.Valid {
		return data.String
	}
	return nil
}

func (dbs *DBStorage) AllKeys() []string {
	// update lastUsed at the end because there might be the slight chance
	// that a statement gets closed despite we're waiting for the result
	// from the server.
	dbs.All.startUse()
	defer dbs.All.stopUse()

	stmt, err := dbs.All.getStmt(dbs.db)
	if err != nil {
		PkgLog.Info("config.DBStorage.AllKeys.All.getStmt", "err", err, "SQL", dbs.All.SQL)
		return nil
	}

	rows, err := stmt.Query()
	if err != nil {
		PkgLog.Info("config.DBStorage.AllKeys.All.Query", "err", err, "SQL", dbs.All.SQL)
		return nil
	}
	defer rows.Close()

	var ret = make([]string, 0, 100)
	var data dbr.NullString
	for rows.Next() {
		if err := rows.Scan(&data); err != nil {
			PkgLog.Info("config.DBStorage.AllKeys.All.Rows.Scan", "err", err, "SQL", dbs.All.SQL)
			return nil
		}
		if data.Valid {
			ret = append(ret, data.String)
		}
		data.String = ""
		data.Valid = false
	}
	return ret
}
