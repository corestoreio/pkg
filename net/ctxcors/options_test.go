// Copyright (c) 2014 Olivier Poitrey <rs@dailymotion.com>
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is furnished
// to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package ctxcors

import (
	"testing"

	"github.com/corestoreio/csfw/store/scope"
	"github.com/corestoreio/csfw/util/cstesting"
	"github.com/stretchr/testify/assert"
)

func noopOption() Option {
	return func(s *Service) {}
}

func TestWithOptionsPassthrough(t *testing.T) {
	tests := []struct {
		scp  scope.Scope
		id   int64
		opt  Option
		want bool
	}{
		{scope.Default, 0, noopOption(), false},
		{scope.Default, 0, WithOptionsPassthrough(scope.Default, 0), true},
	}
	for _, test := range tests {

		c := MustNew(test.opt)
		scp, err := c.getConfigByScopeID(true, scope.NewHash(test.scp, test.id))
		cstesting.FatalIfError(t, err)
		assert.Exactly(t, test.want, scp.optionsPassthrough)
	}
}

//
//func TestWithAllowCredentials(t *testing.T) {
//
//	c := MustNew()
//	assert.False(t, c.AllowCredentials)
//	if _, err := c.Options(WithAllowCredentials()); err != nil {
//		t.Fatal(err)
//	}
//	assert.True(t, c.AllowCredentials)
//}
//
//func TestWithMaxAge(t *testing.T) {
//
//	c := MustNew()
//	_, err := c.Options(WithMaxAge(-1 * time.Second))
//	assert.EqualError(t, err, "MaxAge: Invalid Duration seconds: -1")
//
//	c = MustNew()
//	_, err = c.Options(WithMaxAge(2 * time.Second))
//	assert.NoError(t, err)
//	assert.Exactly(t, "2", c.maxAge)
//}
//
//func TestWithLogger(t *testing.T) {
//
//	c := MustNew()
//	assert.Exactly(t, log.BlackHole{}, c.Log)
//
//	logga := log.NewBlackHole()
//	_, err := c.Options(WithLogger(&logga))
//	assert.NoError(t, err)
//	assert.Exactly(t, &logga, c.Log)
//}
