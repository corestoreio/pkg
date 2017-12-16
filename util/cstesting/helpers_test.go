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

package cstesting_test

import (
	"os"
	"strings"
	"testing"

	"github.com/corestoreio/pkg/util/cstesting"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestChangeEnv(t *testing.T) {
	// cannot run parallel

	key := "X_CORESTORE_TESTING"
	val := "X_CORESTORE_TESTING_VAL1"
	require.NoError(t, os.Setenv(key, val))

	f := cstesting.ChangeEnv(t, key, val+"a")
	assert.Exactly(t, val+"a", os.Getenv(key))
	f()
	assert.Exactly(t, val, os.Getenv(key))
}

func TestChangeDir(t *testing.T) {
	wdOld, err := os.Getwd()
	require.NoError(t, err)

	f := cstesting.ChangeDir(t, os.TempDir())
	wdNew, err := os.Getwd()
	require.NoError(t, err)
	wdNew = strings.Replace(wdNew, "/private", "", 1)
	assert.Exactly(t, os.TempDir(), wdNew+string(os.PathSeparator))
	f()

	wdCurrent, err := os.Getwd()
	require.NoError(t, err)
	assert.Exactly(t, wdOld, wdCurrent)
}
