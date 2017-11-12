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

package binlogsync_test

import (
	"context"
	"math/rand"
	"testing"
	"time"

	"github.com/corestoreio/pkg/sql/binlogsync"
	"github.com/corestoreio/pkg/sql/ddl"
)

// TODO(CyS): Add more tests

func TestIntegrationNewCanal(t *testing.T) {
	t.Parallel()
	dsn, err := ddl.GetParsedDSN()
	if err != nil {
		t.Skipf("Failed to get DSN from env %q with %s", ddl.EnvDSN, err)
	}
	c, err := binlogsync.NewCanal(dsn, binlogsync.WithMySQL())
	if err != nil {
		t.Fatalf("%+v", err)
	}

	c.RegisterRowsEventHandler(catalogProductEvent{idx: 1001, t: t})
	//c.RegisterRowsEventHandler(catalogProductEvent{idx: 1002, t: t})

	if err := c.Start(context.Background()); err != nil {
		t.Fatalf("%+v", err)
	}
	time.Sleep(time.Second * 10)
	c.Close()

}

type catalogProductEvent struct {
	idx int
	t   *testing.T
}

func (cpe catalogProductEvent) Do(_ context.Context, action string, table ddl.Table, rows [][]interface{}) error {
	sl := time.Duration(rand.Intn(100)) * time.Millisecond
	time.Sleep(sl)

	cpe.t.Logf("%d Sleep: %s => %q.%q", cpe.idx, sl, table.Schema, table.Name)
	for _, r := range rows {
		cpe.t.Logf("%#v", r)
	}
	cpe.t.Logf("\n")
	return nil
}
func (cpe catalogProductEvent) Complete(_ context.Context) error {
	return nil // errors.NewFatalf("[test] What is incomplete?")
}
func (cpe catalogProductEvent) String() string {
	return "WTF? catalogProductEvent"
}
