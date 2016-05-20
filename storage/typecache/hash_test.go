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

package typecache

import (
	"hash/fnv"
	"testing"
)

func TestFnv64a_Sum64(t *testing.T) {

	var data = []byte(`// Copyright 2015-2016, Cyrill @ Schumacher.fm and the CoreStore contributors`)
	var f fnv64a
	have := f.Sum64(data)

	gof := fnv.New64a()
	if _, err := gof.Write(data); err != nil {
		t.Fatal(err)
	}
	want := gof.Sum64()
	if have != want {
		t.Errorf("Have %d Want %d", have, want)
	}
}
