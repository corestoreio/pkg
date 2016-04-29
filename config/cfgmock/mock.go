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

package cfgmock

import (
	"fmt"
	"reflect"
	"sort"
	"time"

	"github.com/corestoreio/csfw/config"
	"github.com/corestoreio/csfw/config/cfgpath"
	"github.com/corestoreio/csfw/config/internal/cfgctx"
	"github.com/corestoreio/csfw/config/storage"
	"github.com/corestoreio/csfw/util/bufferpool"
	"github.com/corestoreio/csfw/util/conv"
	"github.com/corestoreio/csfw/util/errors"
	"golang.org/x/net/context"
)

type keyNotFound struct{}

func (a keyNotFound) Error() string  { return "[cfgmock] Get() Path" }
func (a keyNotFound) NotFound() bool { return true }

// Write used for testing when writing configuration values.
type Write struct {
	// WriteError gets always returned by Write
	WriteError error
	// ArgPath will be set after calling write to export the config path.
	// Values you enter here will be overwritten when calling Write
	ArgPath string
	// ArgValue contains the written data
	ArgValue interface{}
}

// Write writes to a black hole, may return an error
func (w *Write) Write(p cfgpath.Path, v interface{}) error {
	w.ArgPath = p.String()
	w.ArgValue = v
	return w.WriteError
}

// OptionFunc to initialize the NewService
type OptionFunc func(*Service)

// Service used for testing. Contains functions which will be called in the
// appropriate methods of interface config.Getter.
// Using WithPV() has precedence over the applied functions.
type Service struct {
	db              storage.Storager
	FByte           func(path string) ([]byte, error)
	FString         func(path string) (string, error)
	FBool           func(path string) (bool, error)
	FFloat64        func(path string) (float64, error)
	FInt            func(path string) (int, error)
	FTime           func(path string) (time.Time, error)
	SubscriptionID  int
	SubscriptionErr error
}

// PathValue is a required type for an option function. PV = path => value.
// This map[string]interface{} is protected by a mutex.
type PathValue map[string]interface{}

func (pv PathValue) set(db storage.Storager) {
	for fq, v := range pv {
		p, err := cfgpath.SplitFQ(fq)
		if err != nil {
			panic(err)
		}
		if err := db.Set(p, v); err != nil {
			panic(err)
		}
	}
}

// GoString creates a sorted Go syntax valid map representation.
// This function panics if it fails to write to the internal buffer.
// Panicing permitted here because this function is only used in testing.
func (pv PathValue) GoString() string {
	keys := make(sort.StringSlice, len(pv))
	i := 0
	for k := range pv {
		keys[i] = k
		i++
	}
	keys.Sort()

	buf := bufferpool.Get()
	defer bufferpool.Put(buf)
	if _, err := buf.WriteString("cfgmock.PathValue{\n"); err != nil {
		panic(err)
	}

	for _, p := range keys {
		if _, err := fmt.Fprintf(buf, "%q: %#v,\n", p, pv[p]); err != nil {
			panic(err)
		}
	}
	if _, err := buf.WriteRune('}'); err != nil {
		panic(err)
	}
	return buf.String()
}

// WithByte returns a function which can be used in the NewService().
// Your function returns a string value from a given cfgpath.
// Call priority 2.
func WithByte(f func(path string) ([]byte, error)) OptionFunc {
	return func(mr *Service) { mr.FByte = f }
}

// WithString returns a function which can be used in the NewService().
// Your function returns a string value from a given cfgpath.
// Call priority 2.
func WithString(f func(path string) (string, error)) OptionFunc {
	return func(mr *Service) { mr.FString = f }
}

// WithBool returns a function which can be used in the NewService().
// Your function returns a bool value from a given path.
// Call priority 2.
func WithBool(f func(path string) (bool, error)) OptionFunc {
	return func(mr *Service) { mr.FBool = f }
}

// WithFloat64 returns a function which can be used in the NewService().
// Your function returns a float64 value from a given path.
// Call priority 2.
func WithFloat64(f func(path string) (float64, error)) OptionFunc {
	return func(mr *Service) { mr.FFloat64 = f }
}

// WithInt returns a function which can be used in the NewService().
// Your function returns an int value from a given path.
// Call priority 2.
func WithInt(f func(path string) (int, error)) OptionFunc {
	return func(mr *Service) { mr.FInt = f }
}

// WithTime returns a function which can be used in the NewService().
// Your function returns a Time value from a given path.
// Call priority 2.
func WithTime(f func(path string) (time.Time, error)) OptionFunc {
	return func(mr *Service) {
		mr.FTime = f
	}
}

// WithPV lets you define a map of path and its values.
// Key is the fully qualified configuration path and value is the value.
// Value must be of the same type as returned by the functions.
// Panics on error.
// Call priority 1.
func WithPV(pv PathValue) OptionFunc {
	return func(mr *Service) {
		pv.set(mr.db)
	}
}

// WithContextGetter adds a cfgmock.Service to a context.
func WithContextGetter(ctx context.Context, opts ...OptionFunc) context.Context {
	return context.WithValue(ctx, cfgctx.KeyGetter{}, NewService(opts...))
}

// WithContextScopedGetter adds a scoped cfgmock.Service to a context.
func WithContextScopedGetter(websiteID, storeID int64, ctx context.Context, opts ...OptionFunc) context.Context {
	return context.WithValue(ctx, cfgctx.KeyScopedGetter{}, NewService(opts...).NewScoped(websiteID, storeID))
}

// WithStorage sets another storage engine to the mock service. This option function
// must be applied at first or your already added values will be lost.
func WithStorage(s storage.Storager) OptionFunc {
	return func(mr *Service) {
		mr.db = s
	}
}

// NewService creates a new Service used in testing.
// Allows you to set different options or you can set the struct fields afterwards.
// WithPV() option has priority over With<T>() functions.
// The simple KV acts as the default storage engine.
func NewService(opts ...OptionFunc) *Service {
	mr := &Service{
		db: storage.NewKV(),
	}
	for _, opt := range opts {
		opt(mr)
	}
	return mr
}

// UpdateValues adds or overwrites the internal path => value map.
func (mr *Service) UpdateValues(pathValues PathValue) {
	pathValues.set(mr.db)
}

func (mr *Service) hasVal(p cfgpath.Path) bool {
	v, err := mr.db.Get(p)
	if err != nil && !errors.IsNotFound(err) {
		println("Mock.Service.hasVal error:", err.Error(), "path", p.String())
	}
	return v != nil && err == nil
}

func (mr *Service) getVal(p cfgpath.Path) interface{} {
	v, err := mr.db.Get(p)
	if err != nil && !errors.IsNotFound(err) {
		println("Mock.Service.getVal error:", err.Error(), "path", p.String())
		return nil
	}
	v = indirect(v)
	return v
}

// Byte returns a byte slice value
func (mr *Service) Byte(p cfgpath.Path) ([]byte, error) {
	switch {
	case mr.hasVal(p):
		return conv.ToByteE(mr.getVal(p))
	case mr.FByte != nil:
		return mr.FByte(p.String())
	default:
		return nil, keyNotFound{}
	}
}

// String returns a string value
func (mr *Service) String(p cfgpath.Path) (string, error) {
	switch {
	case mr.hasVal(p):
		return conv.ToStringE(mr.getVal(p))
	case mr.FString != nil:
		return mr.FString(p.String())
	default:
		return "", keyNotFound{}
	}
}

// Bool returns a bool value
func (mr *Service) Bool(p cfgpath.Path) (bool, error) {
	switch {
	case mr.hasVal(p):
		return conv.ToBoolE(mr.getVal(p))
	case mr.FBool != nil:
		return mr.FBool(p.String())
	default:
		return false, keyNotFound{}
	}
}

// Float64 returns a float64 value
func (mr *Service) Float64(p cfgpath.Path) (float64, error) {
	switch {
	case mr.hasVal(p):
		return conv.ToFloat64E(mr.getVal(p))
	case mr.FFloat64 != nil:
		return mr.FFloat64(p.String())
	default:
		return 0.0, keyNotFound{}
	}
}

// Int returns an integer value
func (mr *Service) Int(p cfgpath.Path) (int, error) {
	switch {
	case mr.hasVal(p):
		return conv.ToIntE(mr.getVal(p))
	case mr.FInt != nil:
		return mr.FInt(p.String())
	default:
		return 0, keyNotFound{}
	}
}

// Time returns a time value
func (mr *Service) Time(p cfgpath.Path) (time.Time, error) {
	switch {
	case mr.hasVal(p):
		return conv.ToTimeE(mr.getVal(p))
	case mr.FTime != nil:
		return mr.FTime(p.String())
	default:
		return time.Time{}, keyNotFound{}
	}
}

// Subscribe returns the before applied SubscriptionID and SubscriptionErr
// Does not start any underlying Goroutines.
func (mr *Service) Subscribe(_ cfgpath.Route, s config.MessageReceiver) (subscriptionID int, err error) {
	return mr.SubscriptionID, mr.SubscriptionErr
}

// NewScoped creates a new config.ScopedReader which uses the underlying
// mocked paths and values.
func (mr *Service) NewScoped(websiteID, storeID int64) config.ScopedGetter {
	return config.NewScopedService(mr, websiteID, storeID)
}

// From html/template/content.go
// Copyright 2011 The Go Authors. All rights reserved.
// indirect returns the value, after dereferencing as many times
// as necessary to reach the base type (or nil).
func indirect(a interface{}) interface{} {
	if a == nil {
		return nil
	}
	if t := reflect.TypeOf(a); t.Kind() != reflect.Ptr {
		// Avoid creating a reflect.Value if it's not a pointer.
		return a
	}
	v := reflect.ValueOf(a)
	for v.Kind() == reflect.Ptr && !v.IsNil() {
		v = v.Elem()
	}
	return v.Interface()
}
