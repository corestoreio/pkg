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
	"github.com/corestoreio/csfw/config"
	"github.com/corestoreio/csfw/config/path"
	"github.com/corestoreio/csfw/storage/csdb"
	"github.com/corestoreio/csfw/storage/dbr"
	"github.com/corestoreio/csfw/store/scope"
	"github.com/juju/errgo"
)

// WithDBStorage applies the MySQL storage to a new Service. It
// starts the idle checker of the DBStorage type.
func WithDBStorage(p csdb.Preparer) config.ServiceOption {
	return func(s *config.Service) {
		s.Storage = MustNewDBStorage(p).Start()
	}
}

// WithCoreConfigData reads the table core_config_data into the Service and overrides
// existing values. If the column `value` is NULL entry will be ignored.
// Stops on errors.
func WithCoreConfigData(dbrSess dbr.SessionRunner) config.ServiceOption {

	return func(s *config.Service) {

		var ccd TableCoreConfigDataSlice
		loadedRows, err := csdb.LoadSlice(dbrSess, TableCollection, TableIndexCoreConfigData, &ccd)
		if PkgLog.IsDebug() {
			PkgLog.Debug("ccd.WithCoreConfigData.LoadSlice", "rows", loadedRows)
		}
		if err != nil {
			if PkgLog.IsDebug() {
				PkgLog.Debug("ccd.WithCoreConfigData.LoadSlice.err", "err", err)
			}
			s.Errors = append(s.Errors, errgo.Mask(err))
			return
		}

		var writtenRows int
		for _, cd := range ccd {
			if cd.Value.Valid {
				var p path.Path
				p, err = path.NewByParts(cd.Path)
				if err != nil {
					s.Errors = append(s.Errors, errgo.Mask(err))
					return
				}

				if err = s.Write(p.Bind(scope.FromString(cd.Scope), cd.ScopeID), cd.Value.String); err != nil {
					s.Errors = append(s.Errors, errgo.Mask(err))
					return
				}
				writtenRows++
			}
		}
		if PkgLog.IsDebug() {
			PkgLog.Debug("ccd.WithCoreConfigData.Written", "loadedRows", loadedRows, "writtenRows", writtenRows)
		}
	}
}
