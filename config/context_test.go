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

package config_test

import (
	"testing"

	"github.com/corestoreio/csfw/config"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
)

func TestContextMustReader(t *testing.T) {
	mr := config.NewMockReader()
	ctx := config.NewContextReader(context.Background(), mr)
	mrHave, ok := config.FromContextReader(ctx)
	assert.Exactly(t, mr, mrHave)
	assert.True(t, ok)

	ctx = config.NewContextReader(context.Background(), nil)
	mrHave, ok = config.FromContextReader(ctx)
	assert.Nil(t, mrHave)
	assert.False(t, ok)
}

func TestContextMustReaderPubSuber(t *testing.T) {
	mr := config.NewMockReader()
	ctx := config.NewContextReaderPubSuber(context.Background(), mr)
	mrHave, ok := config.FromContextReaderPubSuber(ctx)
	assert.Exactly(t, mr, mrHave)
	assert.True(t, ok)

	ctx = config.NewContextReaderPubSuber(context.Background(), nil)
	mrHave, ok = config.FromContextReaderPubSuber(ctx)
	assert.Nil(t, mrHave)
	assert.False(t, ok)
}

type cWrite struct {
}

func (w cWrite) Write(_ ...config.ArgFunc) error {
	return nil
}

var _ config.Writer = (*cWrite)(nil)

func TestContextMustWriter(t *testing.T) {
	wr := cWrite{}
	ctx := config.NewContextWriter(context.Background(), wr)
	wrHave, ok := config.FromContextWriter(ctx)
	assert.Exactly(t, wr, wrHave)
	assert.True(t, ok)

	ctx = config.NewContextWriter(context.Background(), nil)
	wrHave, ok = config.FromContextWriter(ctx)
	assert.Nil(t, wrHave)
	assert.False(t, ok)
}
