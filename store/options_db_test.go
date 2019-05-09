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

// +build csall db

package store

import (
	"context"
	"flag"
	"testing"

	"github.com/alecthomas/repr"
	"github.com/corestoreio/pkg/sql/dmltest"

	"github.com/corestoreio/pkg/util/assert"
)

var (
	runIntegration = flag.Bool("integration", false, "Enables MySQL/MariaDB integration tests, env var CS_DSN must be set")
)

func TestWithLoadFromDB(t *testing.T) {
	if !*runIntegration {
		t.Skip("Skipping integration tests. You can enable them with via CLI option `-integration`")
	}

	dbc := dmltest.MustConnectDB(t)
	defer dmltest.Close(t, dbc)

	srv, err := NewService(WithLoadFromDB(context.TODO(), dbc.DB))
	assert.NoError(t, err)

	repr.Println(srv.stores, srv.groups, srv.websites)

	st, err := srv.DefaultStoreView()
	assert.NoError(t, err)
	t.Logf("%q", st.Code)
}
