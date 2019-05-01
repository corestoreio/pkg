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

package mock_test

import (
	"testing"

	"github.com/alecthomas/repr"
	"github.com/corestoreio/pkg/store/mock"
	"github.com/stretchr/testify/assert"
)

func TestNewEurozzyService_Euro(t *testing.T) {

	ns := mock.NewServiceEuroOZ()
	assert.NotNil(t, ns)

	s, err := ns.Store(4)
	if err != nil {
		t.Fatal(err)
	}
	assert.Exactly(t, "uk", s.Code.String)

	s, err = ns.Store(3)
	if err != nil {
		t.Fatal(err)
	}
	assert.Exactly(t, "ch", s.Code.String)
}

func TestNewEurozzyService_ANZ(t *testing.T) {

	ns := mock.NewServiceEuroOZ()
	assert.NotNil(t, ns)

	s, err := ns.Store(4)
	if err != nil {
		t.Fatal(err)
	}
	assert.Exactly(t, "uk", s.Code.String)

	s, err = ns.Store(3)
	if err != nil {
		t.Fatal(err)
	}
	assert.Exactly(t, "ch", s.Code.String)
	assert.Exactly(t, int64(1), s.WebsiteID)

	s, err = ns.DefaultStoreView()
	if err != nil {
		t.Fatal(err)
	}
	assert.Exactly(t, "at", s.Code.String)
}

func TestNewServiceEuroW11G11S19(t *testing.T) {
	srv := mock.NewServiceEuroW11G11S19()

	repr.Println(srv)

}
