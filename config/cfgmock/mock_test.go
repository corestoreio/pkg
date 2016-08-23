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

	"github.com/corestoreio/csfw/config"
	"github.com/corestoreio/csfw/config/cfgmock"
	"github.com/corestoreio/csfw/config/cfgpath"
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
		}(&wg)
	}
	wg.Wait()

	assert.Exactly(t, iterations, s.ByteFnInvokes())
	assert.Exactly(t, iterations, s.StringFnInvokes())
	assert.Exactly(t, iterations, s.BoolFnInvokes())
	assert.Exactly(t, iterations, s.Float64FnInvokes())
	assert.Exactly(t, iterations, s.IntFnInvokes())
	assert.Exactly(t, iterations, s.TimeFnInvokes())
}

func TestNewServiceAllTypes(t *testing.T) {

	types := []interface{}{"a", int(3141), float64(2.7182) * 3.141, true, time.Now(), []byte(`H∑llo goph€r`)}
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
		default:
			t.Fatalf("Unsupported type: %#v in Index Value %d", wantVal, iFaceIDX)
		}

		if haveErr != nil {
			t.Fatal(haveErr)
		}
		if false == reflect.DeepEqual(wantVal, haveVal) {
			t.Fatalf("Want %v Have %v", wantVal, haveVal)
		}
	}

}
