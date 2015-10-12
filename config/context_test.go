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
	"github.com/corestoreio/csfw/config"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
	"testing"
)

func TestContextMustReader(t *testing.T) {
	mr := config.NewMockReader()
	ctx := context.WithValue(context.Background(), config.CtxKeyReader, mr)
	assert.Exactly(t, mr, config.ContextMustReader(ctx))

	defer func() {
		if r := recover(); r != nil {
			assert.EqualError(t, r.(error), config.ErrTypeAssertionReaderFailed.Error())
		}
	}()
	ctx = context.WithValue(context.Background(), config.CtxKeyReader, "Hello")
	config.ContextMustReader(ctx)
}

func TestContextMustReaderPubSuber(t *testing.T) {
	mr := config.NewMockReader()
	ctx := context.WithValue(context.Background(), config.CtxKeyReaderPubSuber, mr)
	assert.Exactly(t, mr, config.ContextMustReaderPubSuber(ctx))

	defer func() {
		if r := recover(); r != nil {
			assert.EqualError(t, r.(error), config.ErrTypeAssertionReaderPubSuberFailed.Error())
		}
	}()
	ctx = context.WithValue(context.Background(), config.CtxKeyReaderPubSuber, "Hello")
	config.ContextMustReaderPubSuber(ctx)
}

type cWrite struct {
}

func (w cWrite) Write(_ ...config.ArgFunc) error {
	return nil
}

var _ config.Writer = (*cWrite)(nil)

func TestContextMustWriter(t *testing.T) {
	wr := cWrite{}
	ctx := context.WithValue(context.Background(), config.CtxKeyWriter, wr)
	assert.Exactly(t, wr, config.ContextMustWriter(ctx))

	defer func() {
		if r := recover(); r != nil {
			assert.EqualError(t, r.(error), config.ErrTypeAssertionWriterFailed.Error())
		}
	}()
	ctx = context.WithValue(context.Background(), config.CtxKeyWriter, "Hello")
	config.ContextMustWriter(ctx)
}
