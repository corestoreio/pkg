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

package config_test

import (
	"fmt"
	"sync"
	"testing"

	"github.com/corestoreio/pkg/config"
	"github.com/stretchr/testify/assert"
)

var _ config.Getter = (*config.Mock)(nil)
var _ config.GetterPubSuber = (*config.Mock)(nil)
var _ config.Putter = (*config.MockWrite)(nil)
var _ fmt.GoStringer = (*config.MockPathValue)(nil)

func TestPathValueGoStringer(t *testing.T) {
	pv := config.MockPathValue{
		"rr/ss/tt": "3.141592",
	}
	assert.Exactly(t, "config.MockPathValue{\n\"rr/ss/tt\": \"3.141592\",\n}", pv.GoString())
}

func TestService_FnInvokes(t *testing.T) {
	called := 0
	s := config.Mock{
		GetFn: func(p *config.Path) (v *config.Value) {
			called++
			return
		},
	}
	_ = s.Get(config.MustNewPath("test/service/invokes"))
	_ = s.Get(config.MustNewPath("test/service/invokes"))
	assert.Exactly(t, 2, called)
}

func TestService_FnInvokes_Map(t *testing.T) {
	s := config.Mock{
		GetFn: func(p *config.Path) (v *config.Value) {
			return
		},
	}

	const iterations = 10
	var wg sync.WaitGroup
	wg.Add(iterations)
	for i := 0; i < iterations; i++ {
		// food for the race detector
		go func(wg *sync.WaitGroup) {
			defer wg.Done()
			_ = s.Get(config.MustNewPath("test/service/invokes"))
		}(&wg)
	}
	wg.Wait()

	assert.Exactly(t, iterations, s.Invokes().Sum())

	assert.Exactly(t, 1, s.Invokes().PathCount())

	assert.Exactly(t, []string{`default/0/test/service/invokes`}, s.Invokes().Paths())

	assert.Exactly(t, 10, s.AllInvocations().Sum())
	assert.Exactly(t, 1, s.AllInvocations().PathCount())
}

// func TestInvocations_ScopeIDs(t *testing.T) {
// 	iv := config.Invocations{"websites/5/web/cors/allow_credentials": 1, "default/0/web/cors/allow_credentials": 1}
// 	assert.Exactly(t, scope.TypeIDs{scope.DefaultTypeID, scope.Website.Pack(5)}, iv.ScopeIDs())
// }
//
// func TestInvocations_Paths(t *testing.T) {
// 	iv := config.Invocations{"websites/5/web/cors/allow_credentials": 1, "default/0/web/cors/allow_credentials": 1}
// 	assert.Exactly(t, []string{"default/0/web/cors/allow_credentials", "websites/5/web/cors/allow_credentials"}, iv.Paths())
// }
//
// func TestNewServiceAllTypes(t *testing.T) {
//
// 	types := []interface{}{time.Hour, "a", int(3141), float64(2.7182) * 3.141, true, time.Now(), []byte(`H∑llo goph€r`)}
// 	p := cfgpath.MustMakeByString("aa/bb/cc")
//
// 	for iFaceIDX, wantVal := range types {
// 		mg := config.NewService(config.MockPathValue{
// 			p.String(): wantVal,
// 		})
//
// 		var haveVal interface{}
// 		var haveErr error
// 		switch wantVal.(type) {
// 		case []byte:
// 			haveVal, haveErr = mg.Byte(p)
// 		case string:
// 			haveVal, haveErr = mg.String(p)
// 		case bool:
// 			haveVal, haveErr = mg.Bool(p)
// 		case float64:
// 			haveVal, haveErr = mg.Float64(p)
// 		case int:
// 			haveVal, haveErr = mg.Int(p)
// 		case time.Time:
// 			haveVal, haveErr = mg.Time(p)
// 		case time.Duration:
// 			haveVal, haveErr = mg.Duration(p)
// 		default:
// 			t.Fatalf("Unsupported type: %#v in Index Value %d", wantVal, iFaceIDX)
// 		}
//
// 		if haveErr != nil {
// 			t.Fatal(haveErr)
// 		}
// 		if !reflect.DeepEqual(wantVal, haveVal) {
// 			t.Fatalf("Want %v Have %v", wantVal, haveVal)
// 		}
// 		assert.Exactly(t, 1, mg.AllInvocations().Sum())
// 		assert.Exactly(t, 1, mg.AllInvocations().PathCount())
// 		assert.Exactly(t, []string{`default/0/aa/bb/cc`}, mg.AllInvocations().Paths())
// 	}
// }
//
// func TestNewServiceAllTypes_NotFound(t *testing.T) {
//
// 	types := []interface{}{time.Hour, "a", int(3141), float64(2.7182) * 3.141, true, time.Now(), []byte(`H∑llo goph€r`)}
// 	p := cfgpath.MustMakeByString("xx/yy/zz")
//
// 	for iFaceIDX, wantVal := range types {
// 		mg := config.NewService()
//
// 		var haveErr error
// 		switch wantVal.(type) {
// 		case []byte:
// 			_, haveErr = mg.Byte(p)
// 		case string:
// 			_, haveErr = mg.String(p)
// 		case bool:
// 			_, haveErr = mg.Bool(p)
// 		case float64:
// 			_, haveErr = mg.Float64(p)
// 		case int:
// 			_, haveErr = mg.Int(p)
// 		case time.Time:
// 			_, haveErr = mg.Time(p)
// 		case time.Duration:
// 			_, haveErr = mg.Duration(p)
// 		default:
// 			t.Fatalf("Unsupported type: %#v in Index Value %d", wantVal, iFaceIDX)
// 		}
//
// 		assert.True(t, errors.NotFound.Match(haveErr), "%+v", haveErr)
//
// 		assert.Exactly(t, 1, mg.AllInvocations().Sum())
// 		assert.Exactly(t, 1, mg.AllInvocations().PathCount())
// 		assert.Exactly(t, []string{`default/0/xx/yy/zz`}, mg.AllInvocations().Paths())
// 	}
// }
