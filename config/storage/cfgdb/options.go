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

package cfgdb

import (
	"context"
	"time"

	"github.com/corestoreio/errors"
	"github.com/corestoreio/pkg/config"
	"github.com/corestoreio/pkg/sql/ddl"
	"github.com/corestoreio/pkg/sql/dml"
	"github.com/corestoreio/pkg/store/scope"
)

// WithCoreConfigData reads the table core_config_data into the Service and
// overrides existing values. Stops on errors.
func WithCoreConfigData(tbls *ddl.Tables, o Options) config.Option {
	return func(s *config.Service) error {

		tn := o.TableName
		if tn == "" {
			tn = TableNameCoreConfigData
		}

		tbl, err := tbls.Table(tn)
		if err != nil {
			return errors.WithStack(err)
		}

		if o.ContextTimeoutRead == 0 {
			o.ContextTimeoutRead = time.Second * 10 // just a guess
		}

		ctx, cancel := context.WithTimeout(context.Background(), o.ContextTimeoutRead)
		defer cancel()

		return tbl.SelectAll().WithArgs().IterateSerial(ctx, func(cm *dml.ColumnMap) error {
			var ccd TableCoreConfigData
			if err := ccd.MapColumns(cm); err != nil {
				return errors.Wrapf(err, "[ccd] dbs.stmtAll.IterateSerial at row %d", cm.Count)
			}

			var v []byte
			if ccd.Value.Valid {
				v = []byte(ccd.Value.String)
			}
			scp := scope.FromString(ccd.Scope).Pack(ccd.ScopeID)
			p, err := config.NewPathWithScope(scp, ccd.Path)
			if err != nil {
				return errors.Wrapf(err, "[ccd] WithCoreConfigData.config.NewPathWithScope Path %q Scope: %q ID: %d", ccd.Path, scp, ccd.ConfigID)
			}
			if err = s.Set(p, v); err != nil {
				return errors.Wrapf(err, "[ccd] WithCoreConfigData.Service.Write Path %q Scope: %q ID: %d", ccd.Path, scp, ccd.ConfigID)
			}

			return nil
		})
	}
}
