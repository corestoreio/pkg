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

package cfgfile_test

import (
	"testing"

	"github.com/corestoreio/pkg/config"
	"github.com/corestoreio/pkg/config/storage"
	"github.com/corestoreio/pkg/config/storage/cfgfile"
	"github.com/corestoreio/pkg/store/scope"
	"github.com/stretchr/testify/assert"
)

func TestWithLoadData_Success(t *testing.T) {

	inMem := storage.NewMap()
	cfgSrv, err := config.NewService(
		inMem, config.Options{},
		cfgfile.WithLoadJSON("testdata/example.json"),
	)
	if err != nil {
		t.Fatalf("%+v", err)
	}
	p := config.MustNewPath("payment/stripe/user_name")

	assert.Exactly(t, `"AUserName"`, cfgSrv.Get(p).String())
	assert.Exactly(t, `"WS0Username"`, cfgSrv.Get(p.BindWebsite(0)).String())
	assert.Exactly(t, `"WS1Username"`, cfgSrv.Get(p.BindWebsite(1)).String())
	assert.Exactly(t, `"WS2Username"`, cfgSrv.Get(p.BindWebsite(2)).String())

	assert.Exactly(t, `"SO5Username"`, cfgSrv.Get(p.BindStore(5)).String())
	assert.Exactly(t, `"SO11Username"`, cfgSrv.Get(p.BindStore(11)).String())

	assert.Exactly(t, `"1234"`, cfgSrv.Get(config.MustNewPath("payment/stripe/port")).String())
	assert.Exactly(t, `"true"`, cfgSrv.Get(config.MustNewPathWithScope(scope.Website.WithID(0), "payment/stripe/enable")).String())

}
