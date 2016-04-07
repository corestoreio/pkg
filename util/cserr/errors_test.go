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
	"fmt"
	"strings"
	"testing"

	"github.com/corestoreio/csfw/util/cserr"
	"github.com/juju/errors"
	"github.com/stretchr/testify/assert"
)

var _ error = (*cserr.MultiErr)(nil)

func TestMultiErrors(t *testing.T) {
	t.Parallel()
	assert.Equal(t,
		"[{github.com/corestoreio/csfw/util/cserr/errors_test.go:38: Err1}]\n[{github.com/corestoreio/csfw/util/cserr/errors_test.go:38: Err2}]\n[{github.com/corestoreio/csfw/util/cserr/errors_test.go:38: Err3}]",
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
		"[{github.com/corestoreio/csfw/util/cserr/errors_test.go:47: Err5}]\n[{github.com/corestoreio/csfw/util/cserr/errors_test.go:49: Err6}]\n[{github.com/corestoreio/csfw/util/cserr/errors_test.go:50: Err7}]",
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

func TestMultiAppendRecursive(t *testing.T) {
	t.Parallel()

	me := cserr.NewMultiErr(goerr.New("Err1"))
	me.AppendErrors(cserr.NewMultiErr(goerr.New("Err2"), cserr.NewMultiErr(goerr.New("Err3"))))
	assert.Exactly(t, "Err1\nErr2\nErr3", me.Error())
	fmtd := fmt.Sprintf("%#v", me)
	// "&cserr.MultiErr{errs:[]error{(*errors.errorString)(0xc82000f590), (*errors.errorString)(0xc82000f5b0), (*errors.errorString)(0xc82000f5c0)}, details:false}" (actual)
	assert.Exactly(t, 1, strings.Count(fmtd, "MultiErr"))
	assert.Exactly(t, 3, strings.Count(fmtd, "*errors.errorString"))
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
	assert.True(t, me.Contains(fmt.Errorf("Err4")))
	assert.False(t, me.Contains(fmt.Errorf("Err5")))
	assert.False(t, me.Contains(nil))
}

func TestContains(t *testing.T) {
	t.Parallel()

	e1 := errors.New("Err1")
	e2 := errors.New("Err2")
	e3 := errors.New("Err3")
	e4 := goerr.New("Err4")

	var me *cserr.MultiErr
	me = me.AppendErrors(e2, e1, errors.Mask(e4))
	assert.NotNil(t, me)

	assert.True(t, cserr.Contains(e1, errors.New("Err1")))
	assert.False(t, cserr.Contains(e1, errors.New("Err5")))
	assert.True(t, cserr.Contains(e1, e1))
	assert.False(t, cserr.Contains(e1, e2))
	assert.False(t, cserr.Contains(e4, e2))
	assert.True(t, cserr.Contains(e4, errors.Mask(e4)))
	assert.True(t, cserr.Contains(e4, e4))
	assert.False(t, cserr.Contains(nil, e4))
	assert.False(t, cserr.Contains(e4, nil))
	assert.False(t, cserr.Contains(nil, nil))

	assert.False(t, cserr.Contains(me, nil))
	assert.True(t, cserr.Contains(me, e2))
	assert.True(t, cserr.Contains(me, e1))
	assert.True(t, cserr.Contains(e1, me))
	assert.True(t, cserr.Contains(errors.Mask(me), e4))
	assert.True(t, cserr.Contains(e4, errors.Mask(me)))
	assert.False(t, cserr.Contains(me, e3))
	assert.False(t, cserr.Contains(e3, me))

	assert.True(t, cserr.Contains(cserr.NewMultiErr(e3, e1), me))
	assert.True(t, cserr.Contains(me, cserr.NewMultiErr(e3, e1)))
	assert.False(t, cserr.Contains(me, cserr.NewMultiErr(e3)))
}

var _ error = (*cserr.Error)(nil)

func TestError(t *testing.T) {
	const err cserr.Error = "I'm a constant Error"
	assert.EqualError(t, err, "I'm a constant Error")
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
