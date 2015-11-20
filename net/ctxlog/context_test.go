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

package ctxlog_test

import (
	"testing"

	"github.com/corestoreio/csfw/net/ctxlog"
	"github.com/corestoreio/csfw/utils/log"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
)

func TestContext(t *testing.T) {
	l := log.NewStdLogger()
	ctx := context.Background()
	ctx = ctxlog.NewContext(ctx, l)

	haveL := ctxlog.FromContext(ctx)
	assert.Exactly(t, l, haveL)

	haveL2 := ctxlog.FromContext(context.TODO())
	assert.Exactly(t, log.BlackHole{}, haveL2)
}
