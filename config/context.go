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

package config

import (
	"errors"
	"golang.org/x/net/context"
)

// ctxKey type is unexported to prevent collisions with context keys defined in
// other packages.
type ctxKey uint

// Key* defines the keys to access a value in a context.Context
const (
	CtxKeyReader ctxKey = iota
	CtxKeyReaderPubSuber
	CtxKeyWriter
)

var (
	ErrTypeAssertionReaderFailed         = errors.New("Type assertion to config.Reader failed. Maybe missing?")
	ErrTypeAssertionReaderPubSuberFailed = errors.New("Type assertion to config.ReaderPubSuber failed. Maybe missing?")
	ErrTypeAssertionWriterFailed         = errors.New("Type assertion to config.Writer failed. Maybe missing?")
)

// ContextMustReader returns a config.Reader from a context.
func ContextMustReader(ctx context.Context) Reader {
	r, ok := ctx.Value(CtxKeyReader).(Reader)
	if !ok {
		panic(ErrTypeAssertionReaderFailed)
	}
	return r
}

// ContextMustReaderPubSuber returns a config.ReaderPubSuber from a context.
func ContextMustReaderPubSuber(ctx context.Context) ReaderPubSuber {
	r, ok := ctx.Value(CtxKeyReaderPubSuber).(ReaderPubSuber)
	if !ok {
		panic(ErrTypeAssertionReaderPubSuberFailed)
	}
	return r
}

// ContextMustWriter returns a config.Writer from a context.
func ContextMustWriter(ctx context.Context) Writer {
	r, ok := ctx.Value(CtxKeyWriter).(Writer)
	if !ok {
		panic(ErrTypeAssertionWriterFailed)
	}
	return r
}
