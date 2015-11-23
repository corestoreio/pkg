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

package email

import (
	"testing"

	"github.com/corestoreio/csfw/config"
	"github.com/stretchr/testify/assert"
)

var configMock = config.NewMockGetter(
	config.WithMockInt(func(path string) int {
		//		println("int", path)
		switch path {
		case "stores/5015/system/smtp/port":
			return 0
		case "stores/6023/system/smtp/port":
			return 4040
		default:
			return 0
		}
	}),
	config.WithMockString(func(path string) string {
		//		println("string", path)
		switch path {
		case "stores/5015/system/smtp/host":
			return ""
		case "stores/5015/system/smtp/username":
			return ""
		case "stores/6023/system/smtp/host":
			return "smtp.fastmail.com"
		case "stores/6023/system/smtp/username":
			return "2522e71a49e"
		case "stores/6023/system/smtp/password":
			return "9512e71a49f"
		default:
			return ""
		}

	}),
	config.WithMockBool(func(path string) bool {
		return false
	}),
)

func TestDialerPoolDefaultConfig(t *testing.T) {
	dm, err := NewDaemon(
		SetConfig(configMock),
		SetScope(config.ScopeID(5015)),
	)
	assert.NoError(t, err)
	assert.NotNil(t, dm)
	assert.Equal(t, uint64(0xcc72e0b18f4a60fb), dm.ID()) // "localhost25"
}

func TestDialerPoolSingleton(t *testing.T) {
	dm1, err := NewDaemon(
		SetConfig(configMock),
		SetScope(config.ScopeID(6023)),
	)
	assert.NoError(t, err)
	assert.NotNil(t, dm1)
	assert.Equal(t, uint64(0x96b8eb270abcef94), dm1.ID()) // "smtp.fastmail.com40402522e71a49e"

	dm2, err := NewDaemon(
		SetConfig(configMock),
		SetScope(config.ScopeID(6023)),
	)
	assert.NoError(t, err)
	assert.NotNil(t, dm2)

	//	t.Logf("%p == %p", dm1.dialer, dm2.dialer)

	dp1 := dm1.dialer
	dp2 := dm2.dialer
	assert.True(t, dp1 == dp2, "dm1.dialer != dm2.dialer but must be equal")

	dm3, err := NewDaemon(
		SetConfig(configMock),
		SetScope(config.ScopeID(7077)),
	)
	assert.NoError(t, err)
	assert.NotNil(t, dm3)

	dp3 := dm3.dialer
	assert.True(t, dp1 == dp2 && dp1 != dp3 && dp2 != dp3, "dm1.dialer == dm2.dialer && dm1.dialer != dm3.dialer")

}
