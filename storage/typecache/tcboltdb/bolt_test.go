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

package tcboltdb

import (
	"io/ioutil"
	"math"
	"testing"

	"os"
	"path/filepath"

	"github.com/corestoreio/csfw/storage/typecache"
	"github.com/corestoreio/csfw/util/errors"
	"github.com/stretchr/testify/assert"
)

var _ typecache.Cacher = (*wrapper)(nil)

func getTempFile(t *testing.T) string {
	f, err := ioutil.TempFile("", "tcboltdb_")
	if err != nil {
		t.Fatal(err)
	}
	return f.Name()
}

func TestWithBolt_Success(t *testing.T) {
	fn := getTempFile(t)
	defer os.Remove(fn)

	p, err := typecache.NewProcessor(WithFile(fn, 0600))
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := p.Cache.Close(); err != nil {
			t.Fatal(err)
		}
	}()
	var key = []byte(`key1`)
	if err := p.Set(key, math.Pi); err != nil {
		t.Fatal(err)
	}

	var newVal float64
	if err := p.Get(key, &newVal); err != nil {
		t.Fatal(err)
	}
	assert.Exactly(t, math.Pi, newVal)

}

func TestWithBolt_Error(t *testing.T) {
	p, err := typecache.NewProcessor(WithFile(filepath.Join("non", "existent"), 0400))
	assert.Nil(t, p)
	assert.True(t, errors.IsFatal(err), "Error: %s", err)
}
