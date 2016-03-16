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
	"reflect"
	"testing"
	"time"

	"github.com/corestoreio/csfw/config"
	"github.com/corestoreio/csfw/config/cfgmock"
	"github.com/corestoreio/csfw/config/cfgpath"
)

var _ config.Getter = (*cfgmock.Service)(nil)
var _ config.Writer = (*cfgmock.Write)(nil)
var _ config.GetterPubSuber = (*cfgmock.Service)(nil)

func TestNewMockGetterAllTypes(t *testing.T) {
	t.Parallel()

	types := []interface{}{"a", int(3141), float64(2.7182) * 3.141, true, time.Now()}
	p := cfgpath.MustNewByParts("aa/bb/cc")

	for iFaceIDX, wantVal := range types {
		mg := cfgmock.NewService(cfgmock.WithPV(
			cfgmock.PathValue{
				p.String(): wantVal,
			},
		))

		var haveVal interface{}
		var haveErr error
		switch wantVal.(type) {
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
			t.Fatalf("Want %v Have %V", wantVal, haveVal)
		}
	}

}