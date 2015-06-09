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

// Package csMake replaces the Makefile. csMake is only used via go:generate.
package main

import (
	"os"
	"testing"

	"github.com/corestoreio/csfw/storage/csdb"
	"github.com/stretchr/testify/assert"
)

func TestCheckEnv(t *testing.T) {
	envBkp := os.Getenv(csdb.EnvDSN)
	defer os.Setenv(csdb.EnvDSN, envBkp)

	os.Setenv(csdb.EnvDSN, "testing")
	assert.NoError(t, checkEnv())
	os.Setenv(csdb.EnvDSN, "")
	assert.Error(t, checkEnv())
}
