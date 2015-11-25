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
	"github.com/corestoreio/csfw/config/scope"
	"github.com/corestoreio/csfw/utils"
	"github.com/corestoreio/csfw/utils/cast"
	"github.com/juju/errgo"
	"time"
)

var _ Storager = (*DBStorage)(nil)

type StmtUsage struct {
	SQL        string
	stmt       *sql.Stmt
	closed     bool
	Idle       time.Duration
	lastUsed   time.Time
	inUse      bool
	stopTicker chan struct{}
	closeError error
}

type DBStorage struct {
	db    *sql.DB
	Read  *StmtUsage
	Write *StmtUsage
}

func NewDBStorage(db *sql.DB) *DBStorage {
	// idea: as this is a long running service we should have here
	// two prepared statements one for select and one for insert.
	// after time x in which nothing happens neither select nor
	// insert not update the preprared statement gets closed
	// and once there is a new action then we recreate a prepared
	// statement.
	dbs := &DBStorage{
		db: db,
		Read: &StmtUsage{
			SQL: fmt.Sprintf(
				"SELECT `value` FROM `%s` WHERE `scope`=? AND `scope_id`=? AND `path`=?",
				TableCollection.Name(TableIndexCoreConfigData),
			),
			Idle:       time.Second * 10,
			stopTicker: make(chan struct{}),
			closed:     true,
		},
		Write: &StmtUsage{
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

func (su *StmtUsage) close() {
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

func (su *StmtUsage) checkUsage() {
	ticker := time.NewTicker(su.Idle)
	for {
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

func (su *StmtUsage) getStmt(db *sql.DB) (*sql.Stmt, error) {
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

func (su *StmtUsage) startUse() {
	su.lastUsed = time.Now()
	su.inUse = true
}

func (su *StmtUsage) stopUse() {
	su.lastUsed = time.Now()
	su.inUse = false
}

func (dbs *DBStorage) Start() *DBStorage {
	go dbs.Read.checkUsage()
	go dbs.Write.checkUsage()
	return dbs
}

func (dbs *DBStorage) Stop() (err error) {
	dbs.Read.stopTicker <- struct{}{}
	dbs.Write.stopTicker <- struct{}{}
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

	scope, scopeID, path, err := scope.ReverseFQPath(key)
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
		PkgLog.Info("config.DBStorage.Set.Write.Result", "lastInderID", li, "lastInderIDErr", err1, "rowsAffected", ra, "rowsAffectedErr", err2, "SQL", dbs.Write.SQL, "key", key, "value", value)
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

	scope, scopeID, path, err := scope.ReverseFQPath(key)
	if err != nil {
		PkgLog.Info("config.DBStorage.Get.ReverseFQPath", "err", err, "key", key)
		return nil
	}

	var data TableCoreConfigData
	err = stmt.QueryRow(scope, scopeID, path).Scan(&data.Value)

	return data.Value.String
}
func (sp *DBStorage) AllKeys() []string {
	var ret = make(utils.StringSlice, 0)
	return ret.Sort()
}
