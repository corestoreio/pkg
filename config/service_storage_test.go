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
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSimpleStorage(t *testing.T) {
	sp := newSimpleStorage()

	sp.Set("k2", 19.99)
	assert.Exactly(t, 19.99, sp.Get("k2").(float64))

	sp.Set("k1", 4711)
	assert.Exactly(t, 4711, sp.Get("k1").(int))

	assert.Nil(t, sp.Get("k1a"))

	assert.Exactly(t, []string{"k1", "k2"}, sp.AllKeys())
}
