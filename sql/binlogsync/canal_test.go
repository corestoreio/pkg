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

package binlogsync

import (
	"testing"

	"github.com/corestoreio/errors"
	"github.com/corestoreio/pkg/config"
	"github.com/corestoreio/pkg/config/storage"
	"github.com/corestoreio/pkg/util/assert"
)

var (
	_ errors.Kinder = (*errTableNotAllowed)(nil)
	_ error         = (*errTableNotAllowed)(nil)
)

func TestErrTableNotAllowed(t *testing.T) {
	t.Parallel()
	err := errTableNotAllowed("Errr")
	assert.Exactly(t, errors.NotAllowed, err.ErrorKind())
	assert.Exactly(t, "[binlogsync] Table \"Errr\" is not allowed", err.Error(), "%q", err.Error())
}

func TestOptions_LoadFromConfigService(t *testing.T) {
	t.Parallel()

	cfgScp := config.NewFakeService(storage.NewMap(
		`default/0/sql/binlogsync/include_table_regex`, "^sales_order$,^catalog_[a-z]+$",
		`default/0/sql/binlogsync/exclude_table_regex`, "wishlist.+,core.+",
		`default/0/sql/binlogsync/binlog_start_file`, "my.bin.x",
		`default/0/sql/binlogsync/binlog_start_position`, "123456",
		`default/0/sql/binlogsync/binlog_slave_id`, "4711",
		`default/0/sql/binlogsync/server_flavor`, "mysql",
	)).Scoped(1, 1)

	o := &Options{
		ConfigScoped: cfgScp,
	}
	err := o.loadFromConfigService()
	assert.NoError(t, err, "\n%+v", err)

	assert.Exactly(t, []string{"^sales_order$", "^catalog_[a-z]+$"}, o.IncludeTableRegex, "IncludeTableRegex")
	assert.Exactly(t, []string{"wishlist.+", "core.+"}, o.ExcludeTableRegex, "ExcludeTableRegex")
	assert.Exactly(t, "my.bin.x", o.BinlogStartFile, "BinlogStartFile")
	assert.Exactly(t, uint64(123456), o.BinlogStartPosition, "BinlogStartPosition")
	assert.Exactly(t, uint64(4711), o.BinlogSlaveId, "BinlogSlaveId")
	assert.Exactly(t, "mysql", o.Flavor, "Flavor")
}
