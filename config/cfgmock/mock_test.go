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

package cfgmock_test

import (
	"fmt"
	"reflect"
	"sync"
	"testing"
	"time"

	"github.com/corestoreio/pkg/config"
	"github.com/corestoreio/pkg/config/cfgmock"
	"github.com/corestoreio/pkg/config/cfgpath"
	"github.com/corestoreio/pkg/store/scope"
	"github.com/corestoreio/errors"
	"github.com/stretchr/testify/assert"
)

var _ config.Getter = (*cfgmock.Service)(nil)
var _ config.Writer = (*cfgmock.Write)(nil)
var _ config.GetterPubSuber = (*cfgmock.Service)(nil)
var _ fmt.GoStringer = (*cfgmock.PathValue)(nil)

func TestPathValueGoStringer(t *testing.T) {
	pv := cfgmock.PathValue{
		"bb/cc/dd": true,
		"rr/ss/tt": 3.141592,
		"aa/bb/cc": 1,
	}
	const want = `cfgmock.PathValue{
"aa/bb/cc": 1,
"bb/cc/dd": true,
"rr/ss/tt": 3.141592,
}`
	assert.Exactly(t, want, pv.GoString())
}

func TestService_FnInvokes(t *testing.T) {
	called := 0
	s := cfgmock.Service{
		ByteFn: func(path string) ([]byte, error) {
			called++
			return nil, nil
		},
		StringFn: func(path string) (string, error) {
			called++
			return "", nil
		},
		BoolFn: func(path string) (bool, error) {
			called++
			return false, nil
		},
		Float64Fn: func(path string) (float64, error) {
			called++
			return 0, nil
		},
		IntFn: func(path string) (int, error) {
			called++
			return 0, nil
		},
		TimeFn: func(path string) (time.Time, error) {
			called++
			return time.Time{}, nil
		},
		DurationFn: func(path string) (time.Duration, error) {
			called++
			return 0, nil
		},
	}
	_, _ = s.Byte(cfgpath.MustNewByParts("test/service/invokes"))
	_, _ = s.String(cfgpath.MustNewByParts("test/service/invokes"))
	_, _ = s.Bool(cfgpath.MustNewByParts("test/service/invokes"))
	_, _ = s.Float64(cfgpath.MustNewByParts("test/service/invokes"))
	_, _ = s.Int(cfgpath.MustNewByParts("test/service/invokes"))
	_, _ = s.Time(cfgpath.MustNewByParts("test/service/invokes"))
	_, _ = s.Duration(cfgpath.MustNewByParts("test/service/invokes"))
	assert.Exactly(t, 7, called)
}

func TestService_FnInvokes_Map(t *testing.T) {
	s := cfgmock.Service{
		ByteFn: func(path string) ([]byte, error) {
			return nil, nil
		},
		StringFn: func(path string) (string, error) {
			return "", nil
		},
		BoolFn: func(path string) (bool, error) {
			return false, nil
		},
		Float64Fn: func(path string) (float64, error) {
			return 0, nil
		},
		IntFn: func(path string) (int, error) {
			return 0, nil
		},
		TimeFn: func(path string) (time.Time, error) {
			return time.Time{}, nil
		},
		DurationFn: func(path string) (time.Duration, error) {
			return 0, nil
		},
	}

	const iterations = 10
	var wg sync.WaitGroup
	wg.Add(iterations)
	for i := 0; i < iterations; i++ {
		// food for the race detector
		go func(wg *sync.WaitGroup) {
			defer wg.Done()
			_, _ = s.Byte(cfgpath.MustNewByParts("test/service/invokes"))
			_, _ = s.String(cfgpath.MustNewByParts("test/service/invokes"))
			_, _ = s.Bool(cfgpath.MustNewByParts("test/service/invokes"))
			_, _ = s.Float64(cfgpath.MustNewByParts("test/service/invokes"))
			_, _ = s.Int(cfgpath.MustNewByParts("test/service/invokes"))
			_, _ = s.Time(cfgpath.MustNewByParts("test/service/invokes"))
			_, _ = s.Duration(cfgpath.MustNewByParts("test/service/invokes"))
		}(&wg)
	}
	wg.Wait()

	assert.Exactly(t, iterations, s.ByteInvokes().Sum())
	assert.Exactly(t, iterations, s.StringInvokes().Sum())
	assert.Exactly(t, iterations, s.BoolInvokes().Sum())
	assert.Exactly(t, iterations, s.Float64Invokes().Sum())
	assert.Exactly(t, iterations, s.IntInvokes().Sum())
	assert.Exactly(t, iterations, s.TimeInvokes().Sum())
	assert.Exactly(t, iterations, s.DurationInvokes().Sum())

	assert.Exactly(t, 1, s.ByteInvokes().PathCount())
	assert.Exactly(t, 1, s.StringInvokes().PathCount())
	assert.Exactly(t, 1, s.BoolInvokes().PathCount())
	assert.Exactly(t, 1, s.Float64Invokes().PathCount())
	assert.Exactly(t, 1, s.IntInvokes().PathCount())
	assert.Exactly(t, 1, s.TimeInvokes().PathCount())
	assert.Exactly(t, 1, s.DurationInvokes().PathCount())

	assert.Exactly(t, []string{`default/0/test/service/invokes`}, s.ByteInvokes().Paths())
	assert.Exactly(t, []string{`default/0/test/service/invokes`}, s.StringInvokes().Paths())
	assert.Exactly(t, []string{`default/0/test/service/invokes`}, s.BoolInvokes().Paths())
	assert.Exactly(t, []string{`default/0/test/service/invokes`}, s.Float64Invokes().Paths())
	assert.Exactly(t, []string{`default/0/test/service/invokes`}, s.IntInvokes().Paths())
	assert.Exactly(t, []string{`default/0/test/service/invokes`}, s.TimeInvokes().Paths())
	assert.Exactly(t, []string{`default/0/test/service/invokes`}, s.DurationInvokes().Paths())

	assert.Exactly(t, 70, s.AllInvocations().Sum())
	assert.Exactly(t, 1, s.AllInvocations().PathCount())
}

func TestInvocations_ScopeIDs(t *testing.T) {
	iv := cfgmock.Invocations{"websites/5/web/cors/allow_credentials": 1, "default/0/web/cors/allow_credentials": 1}
	assert.Exactly(t, scope.TypeIDs{scope.DefaultTypeID, scope.Website.Pack(5)}, iv.ScopeIDs())
}

func TestInvocations_Paths(t *testing.T) {
	iv := cfgmock.Invocations{"websites/5/web/cors/allow_credentials": 1, "default/0/web/cors/allow_credentials": 1}
	assert.Exactly(t, []string{"default/0/web/cors/allow_credentials", "websites/5/web/cors/allow_credentials"}, iv.Paths())
}

func TestNewServiceAllTypes(t *testing.T) {

	types := []interface{}{time.Hour, "a", int(3141), float64(2.7182) * 3.141, true, time.Now(), []byte(`H∑llo goph€r`)}
	p := cfgpath.MustNewByParts("aa/bb/cc")

	for iFaceIDX, wantVal := range types {
		mg := cfgmock.NewService(cfgmock.PathValue{
			p.String(): wantVal,
		})

		var haveVal interface{}
		var haveErr error
		switch wantVal.(type) {
		case []byte:
			haveVal, haveErr = mg.Byte(p)
		case string:
			haveVal, haveErr = mg.String(p)
		case bool:
			haveVal, haveErr = mg.Bool(p)
		case float64:
			haveVal, haveErr = mg.Float64(p)
		case int:
			haveVal, haveErr = mg.Int(p)
		case time.Time:
			haveVal, haveErr = mg.Time(p)
		case time.Duration:
			haveVal, haveErr = mg.Duration(p)
		default:
			t.Fatalf("Unsupported type: %#v in Index Value %d", wantVal, iFaceIDX)
		}

		if haveErr != nil {
			t.Fatal(haveErr)
		}
		if !reflect.DeepEqual(wantVal, haveVal) {
			t.Fatalf("Want %v Have %v", wantVal, haveVal)
		}
		assert.Exactly(t, 1, mg.AllInvocations().Sum())
		assert.Exactly(t, 1, mg.AllInvocations().PathCount())
		assert.Exactly(t, []string{`default/0/aa/bb/cc`}, mg.AllInvocations().Paths())
	}
}

func TestNewServiceAllTypes_NotFound(t *testing.T) {

	types := []interface{}{time.Hour, "a", int(3141), float64(2.7182) * 3.141, true, time.Now(), []byte(`H∑llo goph€r`)}
	p := cfgpath.MustNewByParts("xx/yy/zz")

	for iFaceIDX, wantVal := range types {
		mg := cfgmock.NewService()

		var haveErr error
		switch wantVal.(type) {
		case []byte:
			_, haveErr = mg.Byte(p)
		case string:
			_, haveErr = mg.String(p)
		case bool:
			_, haveErr = mg.Bool(p)
		case float64:
			_, haveErr = mg.Float64(p)
		case int:
			_, haveErr = mg.Int(p)
		case time.Time:
			_, haveErr = mg.Time(p)
		case time.Duration:
			_, haveErr = mg.Duration(p)
		default:
			t.Fatalf("Unsupported type: %#v in Index Value %d", wantVal, iFaceIDX)
		}

		assert.True(t, errors.IsNotFound(haveErr), "%+v", haveErr)

		assert.Exactly(t, 1, mg.AllInvocations().Sum())
		assert.Exactly(t, 1, mg.AllInvocations().PathCount())
		assert.Exactly(t, []string{`default/0/xx/yy/zz`}, mg.AllInvocations().Paths())
	}
}
