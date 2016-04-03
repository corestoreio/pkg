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

var _ error = (*cserr.MultiErr)(nil)

func TestMultiErrors(t *testing.T) {
	t.Parallel()
	assert.Equal(t,
		"[{github.com/corestoreio/csfw/util/cserr/errors_test.go:36: Err1}]\n[{github.com/corestoreio/csfw/util/cserr/errors_test.go:36: Err2}]\n[{github.com/corestoreio/csfw/util/cserr/errors_test.go:36: Err3}]",
		cserr.NewMultiErr(
			errors.New("Err1"),
			errors.New("Err2"),
			errors.New("Err3"),
		).VerboseErrors().Error(),
	)
}

func TestMultiAppend(t *testing.T) {
	t.Parallel()

	e := cserr.NewMultiErr()
	e.AppendErrors(
		errors.New("Err5"),
		nil,
		errors.New("Err6"),
		errors.New("Err7"),
	)
	assert.True(t, e.HasErrors())
	assert.Equal(t,
		"[{github.com/corestoreio/csfw/util/cserr/errors_test.go:45: Err5}]\n[{github.com/corestoreio/csfw/util/cserr/errors_test.go:47: Err6}]\n[{github.com/corestoreio/csfw/util/cserr/errors_test.go:48: Err7}]",
		e.VerboseErrors().Error(),
	)
}

func TestMultiEmpty(t *testing.T) {
	t.Parallel()
	assert.False(t, cserr.NewMultiErr(nil, nil).HasErrors())
	assert.Equal(t, "", cserr.NewMultiErr(nil).Error())
}

func TestHasErrorsNil(t *testing.T) {
	t.Parallel()
	var e *cserr.MultiErr
	assert.False(t, e.HasErrors())

	e = &cserr.MultiErr{}
	assert.False(t, e.HasErrors())
}

func TestMultiAppendToNil(t *testing.T) {
	t.Parallel()
	var e *cserr.MultiErr
	e = e.AppendErrors(errors.New("Err74"))

	assert.True(t, e.HasErrors())
	assert.Equal(t, "Err74", e.Error())
}

func TestMultiAppendNilToNil1(t *testing.T) {
	t.Parallel()
	var e *cserr.MultiErr
	e = e.AppendErrors()
	assert.False(t, e.HasErrors())
	assert.Nil(t, e)
}

func TestMultiAppendNilToNil2(t *testing.T) {
	t.Parallel()
	var e *cserr.MultiErr
	e = e.AppendErrors(nil, nil)
	assert.False(t, e.HasErrors())
	assert.Nil(t, e)
}

func TestMultiErrContains(t *testing.T) {
	t.Parallel()
	var me *cserr.MultiErr

	e1 := errors.New("Err1")
	e2 := errors.New("Err2")
	e3 := errors.New("Err3")
	e4 := goerr.New("Err4")

	me = me.AppendErrors(e2, e1, errors.Mask(e4))
	assert.NotNil(t, me)
	assert.False(t, me.Contains(e3))
	assert.True(t, me.Contains(e2))
	assert.True(t, me.Contains(errors.Mask(e2)))
	assert.True(t, me.Contains(e1))
	assert.True(t, me.Contains(e4))
	assert.False(t, me.Contains(nil))
}

var benchmarkError string

// BenchmarkError-4	  500000	      3063 ns/op	    1312 B/op	      22 allocs/op
// BenchmarkError-4	  500000	      3763 ns/op	    1936 B/op	      26 allocs/op
func BenchmarkError(b *testing.B) {
	// errors.Details(e) produces those high allocs
	e := cserr.NewMultiErr().VerboseErrors()
	e.AppendErrors(
		errors.New("Err5"),
		nil,
		errors.New("Err6"),
		errors.New("Err7"),
	)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		benchmarkError = e.Error()
	}
}
