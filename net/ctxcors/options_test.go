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
	"github.com/corestoreio/csfw/util/log"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestWithOptionsPassthrough(t *testing.T) {
	t.Parallel()

	c := MustNew()
	assert.False(t, c.OptionsPassthrough)
	if _, err := c.Options(WithOptionsPassthrough()); err != nil {
		t.Fatal(err)
	}
	assert.True(t, c.OptionsPassthrough)
}

func TestWithAllowCredentials(t *testing.T) {
	t.Parallel()

	c := MustNew()
	assert.False(t, c.AllowCredentials)
	if _, err := c.Options(WithAllowCredentials()); err != nil {
		t.Fatal(err)
	}
	assert.True(t, c.AllowCredentials)
}

func TestWithMaxAge(t *testing.T) {
	t.Parallel()

	c := MustNew()
	_, err := c.Options(WithMaxAge(-1 * time.Second))
	assert.EqualError(t, err, "Invalid Duration seconds: -1")

	c = MustNew()
	_, err = c.Options(WithMaxAge(2 * time.Second))
	assert.NoError(t, err)
	assert.Exactly(t, "2", c.maxAge)
}

func TestWithLogger(t *testing.T) {
	t.Parallel()

	c := MustNew()
	assert.Exactly(t, log.BlackHole{}, c.Log)

	logga := log.NewBlackHole()
	_, err := c.Options(WithLogger(&logga))
	assert.NoError(t, err)
	assert.Exactly(t, &logga, c.Log)
}

func TestWithBackend(t *testing.T) {
	t.Parallel()

	c := MustNew()
	assert.Nil(t, c.Backend)

	cfgStruct, err := NewConfigStructure()
	if err != nil {
		t.Fatal(err)
	}

	be := NewBackend(cfgStruct)
	_, err = c.Options(WithBackend(be))
	assert.NoError(t, err)
	assert.Exactly(t, be, c.Backend)
}

func TestWithBackendApplied(t *testing.T) {
	t.Parallel()
}
