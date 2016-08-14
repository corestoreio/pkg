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

package storemock_test

import (
	"testing"

	"github.com/corestoreio/csfw/store/storemock"
	"github.com/corestoreio/csfw/util/errors"
	"github.com/stretchr/testify/assert"
)

func TestFind_IsAllowedStoreID(t *testing.T) {
	f := storemock.Find{
		Allowed:      true,
		AllowedCode:  "a",
		AllowedError: errors.NewFatalf("Upps.."),
	}
	isA, ac, err := f.IsAllowedStoreID(0, 0)
	assert.True(t, isA)
	assert.Exactly(t, "a", ac)
	assert.True(t, errors.IsFatal(err))
}

func TestFind_DefaultStoreID(t *testing.T) {
	f := storemock.Find{
		StoreIDDefault:   -1,
		WebsiteIDDefault: -2,
		StoreIDError:     errors.NewFatalf("Upps.."),
	}
	sID, wID, err := f.DefaultStoreID(0)
	assert.Exactly(t, int64(-1), sID)
	assert.Exactly(t, int64(-2), wID)
	assert.True(t, errors.IsFatal(err))
}

func TestFind_StoreIDbyCode(t *testing.T) {
	f := storemock.Find{
		IDByCodeStoreID:   -1,
		IDByCodeWebsiteID: -2,
		IDByCodeError:     errors.NewFatalf("Upps.."),
	}
	sID, wID, err := f.StoreIDbyCode(0, "x")
	assert.Exactly(t, int64(-1), sID)
	assert.Exactly(t, int64(-2), wID)
	assert.True(t, errors.IsFatal(err))
}
