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
	"github.com/corestoreio/csfw/config/cfgpath"
	"github.com/corestoreio/csfw/storage/dbr"
	"github.com/corestoreio/csfw/store/scope"
	"github.com/corestoreio/errors"
	"github.com/corestoreio/log"
)

// WithCoreConfigData reads the table core_config_data into the Service and
// overrides existing values. If the column `value` is NULL entry will be
// ignored. Stops on errors.
func WithCoreConfigData(dbrSess dbr.SessionRunner) config.Option {
	return func(s *config.Service) error {

		var ccd TableCoreConfigDataSlice
		loadedRows, err := TableCollection.MustTable(TableIndexCoreConfigData).LoadSlice(dbrSess, &ccd)
		if s.Log.IsDebug() {
			s.Log.Debug("ccd.WithCoreConfigData.LoadSlice", log.Int("rows", loadedRows), log.Err(err))
		}
		if err != nil {
			return errors.Wrap(err, "[ccd] WithCoreConfigData.csdb.LoadSlice")
		}

		var writtenRows int
		for _, cd := range ccd {
			if cd.Value.Valid {
				var p cfgpath.Path
				p, err = cfgpath.NewByParts(cd.Path)
				if err != nil {
					return errors.Wrapf(err, "[ccd] cfgpath.NewByParts Path %q", cd.Path)
				}

				if err = s.Write(p.Bind(scope.FromString(cd.Scope).Pack(cd.ScopeID)), cd.Value.String); err != nil {
					return errors.Wrapf(err, "[ccd] cfgpath.NewByParts Path %q Scope: %q ID: %d Value: %q", cd.Path, cd.Scope, cd.ScopeID, cd.Value.String)
				}
				writtenRows++
			}
		}
		if s.Log.IsDebug() {
			s.Log.Debug("ccd.WithCoreConfigData.Written", log.Int("loadedRows", loadedRows), log.Int("writtenRows", writtenRows))
		}
		return nil
	}
}
