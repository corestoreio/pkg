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
	"testing"
	"time"

	"github.com/corestoreio/csfw/storage/binlogsync"
	"github.com/corestoreio/csfw/storage/csdb"
	"github.com/siddontang/go-mysql/schema"
)

func TestIntegrationNewCanal(t *testing.T) {
	dsn, err := csdb.GetParsedDSN()
	if err != nil {
		t.Fatalf("Failed to get DSN from env %q with %+v", csdb.EnvDSN, err)
	}
	c, err := binlogsync.NewCanal(dsn)
	if err != nil {
		t.Fatalf("%+v", err)
	}

	c.RegRowsEventHandler(&catalogProductEvent{t: t})

	if err := c.Start(); err != nil {
		t.Fatalf("%+v", err)
	}
	time.Sleep(time.Second * 10)
	c.Close()

}

type catalogProductEvent struct {
	t *testing.T
}

func (cpe *catalogProductEvent) Do(action string, table *schema.Table, rows [][]interface{}) error {

	cpe.t.Logf("%s.%s", table.Schema, table.Name)
	for _, r := range rows {
		cpe.t.Logf("%#v", r)
	}
	cpe.t.Logf("\n")
	return nil
}
func (cpe *catalogProductEvent) Complete() error {
	return nil // errors.NewFatalf("[test] What is incomplete?")
}
func (cpe *catalogProductEvent) String() string {
	return "WTF? catalogProductEvent"
}
