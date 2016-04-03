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

package cserr_test

import (
	goerr "errors"
	"testing"

	"github.com/corestoreio/csfw/util/cserr"
	"github.com/juju/errors"
	"github.com/stretchr/testify/assert"
)

func TestUnwrapMasked(t *testing.T) {
	t.Parallel()
	e1 := goerr.New("PHP")
	assert.Exactly(t, e1, cserr.UnwrapMasked(errors.Mask(e1)))
	assert.Exactly(t, e1, cserr.UnwrapMasked(e1))

	e2 := errors.New("Scala")
	assert.Exactly(t, e2, cserr.UnwrapMasked(errors.Mask(e2)))
	assert.Exactly(t, nil, cserr.UnwrapMasked(errors.Mask(nil)))
	assert.Exactly(t, e2, cserr.UnwrapMasked(e2))
	assert.Exactly(t, nil, cserr.UnwrapMasked(nil))
}
