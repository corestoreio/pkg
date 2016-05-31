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

package log_test

import (
	"io"
	"sync"
	"testing"

	"github.com/corestoreio/csfw/log"
	"github.com/stretchr/testify/assert"
)

var _ io.Writer = (*log.MutexBuffer)(nil)

func TestMutexBuffer(t *testing.T) {

	mb := &log.MutexBuffer{}
	var wg sync.WaitGroup

	// detect race conditions
	wg.Add(1)
	go func(t *testing.T) {
		defer wg.Done()
		if _, err := mb.Write([]byte(`W1`)); err != nil {
			t.Fatal(err)
		}
	}(t)

	wg.Add(1)
	go func(t *testing.T) {
		defer wg.Done()
		if _, err := mb.Write([]byte(`W2`)); err != nil {
			t.Fatal(err)
		}
	}(t)
	wg.Wait()
	assert.Contains(t, mb.String(), "W1")
	assert.Contains(t, mb.String(), "W2")
}
