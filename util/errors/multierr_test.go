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

package errors_test

import (
	"bytes"
	goerr "errors"
	"fmt"
	"strings"
	"testing"

	"github.com/corestoreio/csfw/util/errors"
	"github.com/stretchr/testify/assert"
)

var _ error = (*errors.MultiErr)(nil)

func TestMultiErrors(t *testing.T) {
	errStr := errors.NewMultiErr(
		errors.New("Err1"),
		errors.New("Err2"),
		errors.New("Err3"),
	).Error()

	assert.Regexp(t,
		"Err1\ngithub.com/corestoreio/csfw/util/errors_test.TestMultiErrors\n\t.+/github.com/corestoreio/csfw/util/errors/multierr_test.go:32\ntesting.tRunner\n\t/usr/local/opt/go/libexec/src/testing/testing.go:610\nErr2\ngithub.com/corestoreio/csfw/util/errors_test.TestMultiErrors\n\t.+/github.com/corestoreio/csfw/util/errors/multierr_test.go:33\ntesting.tRunner\n\t/usr/local/opt/go/libexec/src/testing/testing.go:610\nErr3\ngithub.com/corestoreio/csfw/util/errors_test.TestMultiErrors\n\t.+/github.com/corestoreio/csfw/util/errors/multierr_test.go:34\ntesting.tRunner\n\t.+/testing/testing.go:610\n",
		errStr,
	)
}

func TestMultiAppend(t *testing.T) {

	e := errors.NewMultiErr().AppendErrors(
		errors.New("Err5"),
		nil,
		errors.New("Err6"),
		errors.New("Err7"),
	)
	assert.True(t, e.HasErrors())
	assert.Regexp(t,
		"Err5\ngithub.com/corestoreio/csfw/util/errors_test.TestMultiAppend\n\t.+/github.com/corestoreio/csfw/util/errors/multierr_test.go:46\ntesting.tRunner\n\t/usr/local/opt/go/libexec/src/testing/testing.go:610\nErr6\ngithub.com/corestoreio/csfw/util/errors_test.TestMultiAppend\n\t.+/github.com/corestoreio/csfw/util/errors/multierr_test.go:48\ntesting.tRunner\n\t/usr/local/opt/go/libexec/src/testing/testing.go:610\nErr7\ngithub.com/corestoreio/csfw/util/errors_test.TestMultiAppend\n\t.+/github.com/corestoreio/csfw/util/errors/multierr_test.go:49\ntesting.tRunner\n\t.+/testing/testing.go:610\n",
		e.Error(),
	)
}

func TestMultiEmpty(t *testing.T) {

	assert.False(t, errors.NewMultiErr(nil, nil).HasErrors())
	assert.Equal(t, "", errors.NewMultiErr(nil).Error())
}

func TestHasErrorsNil(t *testing.T) {

	var e *errors.MultiErr
	assert.False(t, e.HasErrors())

	e = &errors.MultiErr{}
	assert.False(t, e.HasErrors())
}

func TestMultiAppendToNil(t *testing.T) {

	var e *errors.MultiErr
	e = e.AppendErrors(errors.New("Err74"))

	assert.True(t, e.HasErrors())
	assert.Regexp(t, "Err74\ngithub.com/corestoreio/csfw/util/errors_test.TestMultiAppendToNil\n\t.+/github.com/corestoreio/csfw/util/errors/multierr_test.go:76\ntesting.tRunner\n\t.+/testing/testing.go:610\n", e.Error())
}

func xxxTestMultiErr_CustomFormatter(t *testing.T) {

	m1 := errors.NewMultiErr(errors.New("Hello1"))
	m1.AppendErrors(
		errors.NewMultiErr(errors.NewAlreadyClosedf("Brain"),
			errors.NewNotFoundf("Mind"),
		),
		errors.New("Hello2"),
	)

	assert.Regexp(t,
		"Hello1\ngithub.com/corestoreio/csfw/util/errors_test.TestMultiErr_CustomFormatter\n\t.+/github.com/corestoreio/csfw/util/errors/multierr_test.go:84\ntesting.tRunner\n\t/usr/local/opt/go/libexec/src/testing/testing.go:610\nAlready closed\ngithub.com/corestoreio/csfw/util/errors.NewAlreadyClosedf\n\t.+/github.com/corestoreio/csfw/util/errors/behaviour.go:414: Brain\nNot found\ngithub.com/corestoreio/csfw/util/errors.NewNotFoundf\n\t.+/github.com/corestoreio/csfw/util/errors/behaviour.go:234: Mind\nHello2\ngithub.com/corestoreio/csfw/util/errors_test.TestMultiErr_CustomFormatter\n\t.+/github.com/corestoreio/csfw/util/errors/multierr_test.go:89\ntesting.tRunner\n\t.+/testing/testing.go:610\n",
		m1.Error())

	m1.Formatter = func(errs []error) string {
		var buf bytes.Buffer
		for _, err := range errs {
			buf.WriteString(`* `)
			buf.WriteString(err.Error())
			buf.WriteRune('\n')
		}
		return buf.String()
	}
	assert.Exactly(t,
		"* Hello1\n* Brain: Already closed\n* Mind: Not found\n* Hello2\n",
		m1.Error())
}

func TestMultiAppendNilToNil1(t *testing.T) {

	var e *errors.MultiErr
	e = e.AppendErrors()
	assert.False(t, e.HasErrors())
	assert.Nil(t, e)
}

func TestMultiAppendNilToNil2(t *testing.T) {

	var e *errors.MultiErr
	e = e.AppendErrors(nil, nil)
	assert.False(t, e.HasErrors())
	assert.Nil(t, e)
}

func TestMultiAppendRecursive(t *testing.T) {

	me := errors.NewMultiErr(goerr.New("Err1")).
		AppendErrors(errors.NewMultiErr(goerr.New("Err2"), errors.NewMultiErr(goerr.New("Err3"))))
	assert.Exactly(t, "Err1\nErr2\nErr3\n", me.Error())
	fmtd := fmt.Sprintf("%#v", me)
	// "&errors.MultiErr{errs:[]error{(*errors.errorString)(0xc82000f590), (*errors.errorString)(0xc82000f5b0), (*errors.errorString)(0xc82000f5c0)}, details:false}" (actual)
	assert.Exactly(t, 1, strings.Count(fmtd, "MultiErr"))
	assert.Exactly(t, 3, strings.Count(fmtd, "*errors.errorString"))
}

var _ error = (*errors.Error)(nil)

func TestError(t *testing.T) {
	const err errors.Error = "I'm a constant Error"
	assert.EqualError(t, err, "I'm a constant Error")
}

func TestMultiErrContainsAll(t *testing.T) {
	tests := []struct {
		me   error
		vf   []errors.BehaviourFunc
		want bool
	}{
		{errors.NewMultiErr(nil, errors.NewNotValidf("r1")), []errors.BehaviourFunc{errors.IsNotValid}, true},
		{errors.NewMultiErr(nil, errors.NewNotValidf("r1")), []errors.BehaviourFunc{errors.IsNotFound}, false},
		{errors.NewMultiErr(), []errors.BehaviourFunc{errors.IsNotFound}, false},
		{errors.New("random"), []errors.BehaviourFunc{errors.IsNotFound}, false},
		{nil, []errors.BehaviourFunc{errors.IsNotFound}, false},
		{errors.NewMultiErr(nil, errors.NewNotValidf("r1"), errors.NewNotFoundf("r2")), []errors.BehaviourFunc{errors.IsNotFound}, true}, // 5
		{errors.NewMultiErr(nil, errors.NewNotValidf("r1"), errors.NewNotFoundf("r2")), []errors.BehaviourFunc{errors.IsNotFound, errors.IsTemporary}, false},
		{errors.NewMultiErr(nil, errors.NewNotValidf("r1"), errors.NewNotFoundf("r2")), []errors.BehaviourFunc{errors.IsNotFound, errors.IsNotValid}, true},
		{errors.NewMultiErr(nil, errors.NewNotValidf("r1"), errors.NewNotFoundf("r2")), []errors.BehaviourFunc{errors.IsNotFound, errors.IsNotValid, errors.IsAlreadyExists}, true},
		{errors.NewMultiErr(nil, errors.NewNotValidf("r1"), errors.NewNotFoundf("r2")), []errors.BehaviourFunc{errors.IsAlreadyClosed, errors.IsNotValid, errors.IsAlreadyExists}, false},
		{errors.NewMultiErr(nil, errors.NewNotValidf("r1"), errors.NewNotFoundf("r2")), nil, false},
		{errors.NewMultiErr(nil), nil, false},
		{nil, nil, false},
		{errors.NewMultiErr(nil, errors.NewNotValidf("r1"), errors.NewNotFoundf("r2"), errors.NewMultiErr(errors.Error("r3"), errors.NewMultiErr(errors.Error("r4"), errors.NewNotImplementedf("r5")))),
			[]errors.BehaviourFunc{errors.IsNotImplemented},
			true},
	}
	for i, test := range tests {
		if have, want := errors.MultiErrContainsAll(test.me, test.vf...), test.want; have != want {
			t.Errorf("Index %d: Have %t Want %t", i, have, want)
		}
	}
}

func TestMultiErrContainsAny(t *testing.T) {
	tests := []struct {
		me   error
		vf   []errors.BehaviourFunc
		want bool
	}{
		{errors.NewMultiErr(nil, errors.NewNotValidf("r1")), []errors.BehaviourFunc{errors.IsNotValid}, true},
		{errors.NewMultiErr(nil, errors.NewNotValidf("r1")), []errors.BehaviourFunc{errors.IsNotFound}, false},
		{errors.NewMultiErr(), []errors.BehaviourFunc{errors.IsNotFound}, false},
		{errors.New("random"), []errors.BehaviourFunc{errors.IsNotFound}, false},
		{nil, []errors.BehaviourFunc{errors.IsNotFound}, false},
		{errors.NewMultiErr(nil, errors.NewNotValidf("r1"), errors.NewNotFoundf("r2")), []errors.BehaviourFunc{errors.IsNotFound}, true},                     // 5
		{errors.NewMultiErr(nil, errors.NewNotValidf("r1"), errors.NewNotFoundf("r2")), []errors.BehaviourFunc{errors.IsNotFound, errors.IsTemporary}, true}, // 6 different to All
		{errors.NewMultiErr(nil, errors.NewNotValidf("r1"), errors.NewNotFoundf("r2")), []errors.BehaviourFunc{errors.IsNotFound, errors.IsNotValid}, true},
		{errors.NewMultiErr(nil, errors.NewNotValidf("r1"), errors.NewNotFoundf("r2")), []errors.BehaviourFunc{errors.IsNotFound, errors.IsNotValid, errors.IsAlreadyExists}, true},
		{errors.NewMultiErr(nil, errors.NewNotValidf("r1"), errors.NewNotFoundf("r2")), []errors.BehaviourFunc{errors.IsAlreadyClosed, errors.IsNotValid, errors.IsAlreadyExists}, true}, // 9 different to All
		{errors.NewMultiErr(nil, errors.NewNotValidf("r1"), errors.NewNotFoundf("r2")), nil, false},
		{errors.NewMultiErr(nil), nil, false},
		{nil, nil, false},
		{errors.NewMultiErr(nil, errors.NewNotValidf("r1"), errors.NewNotFoundf("r2"), errors.NewMultiErr(errors.Error("r3"), errors.NewMultiErr(errors.Error("r4"), errors.NewNotImplementedf("r5")))),
			[]errors.BehaviourFunc{errors.IsNotImplemented},
			true},
	}
	for i, test := range tests {
		if have, want := errors.MultiErrContainsAny(test.me, test.vf...), test.want; have != want {
			t.Errorf("Index %d: Have %t Want %t", i, have, want)
		}
	}
}
