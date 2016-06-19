// Copyright (c) 2015, Dave Cheney <dave@cheney.net>
// All rights reserved.
//
// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions are met:
//
// * Redistributions of source code must retain the above copyright notice, this
// list of conditions and the following disclaimer.
//
// * Redistributions in binary form must reproduce the above copyright notice,
// this list of conditions and the following disclaimer in the documentation
// and/or other materials provided with the distribution.
//
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
// AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
// IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
// DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE
// FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL
// DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR
// SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER
// CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY,
// OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
// OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.

package errors

import (
	"fmt"
	"io"
	"regexp"
	"strings"
	"testing"
)

func TestFormatNew(t *testing.T) {
	tests := []struct {
		error
		format string
		want   string
	}{{
		New("error"),
		"%s",
		"error",
	}, {
		New("error"),
		"%v",
		"error",
	}, {
		New("error"),
		"%+v",
		"error\n" +
			"github.com/corestoreio/csfw/util/errors.TestFormatNew\n" +
			"\t.+/github.com/corestoreio/csfw/util/errors/format_test.go:49",
	}}

	for _, tt := range tests {
		testFormatRegexp(t, tt.error, tt.format, tt.want)
	}
}

func TestFormatErrorf(t *testing.T) {
	tests := []struct {
		error
		format string
		want   string
	}{{
		Errorf("%s", "error"),
		"%s",
		"error",
	}, {
		Errorf("%s", "error"),
		"%v",
		"error",
	}, {
		Errorf("%s", "error"),
		"%+v",
		"error\n" +
			"github.com/corestoreio/csfw/util/errors.TestFormatErrorf\n" +
			"\t.+/github.com/corestoreio/csfw/util/errors/format_test.go:75",
	}}

	for _, tt := range tests {
		testFormatRegexp(t, tt.error, tt.format, tt.want)
	}
}

func TestFormatWrap(t *testing.T) {
	tests := []struct {
		error
		format string
		want   string
	}{{
		Wrap(New("error"), "error2"),
		"%s",
		"error2: error",
	}, {
		Wrap(New("error"), "error2"),
		"%v",
		"error2: error",
	}, {
		Wrap(New("error"), "error2"),
		"%+v",
		"error\n" +
			"github.com/corestoreio/csfw/util/errors.TestFormatWrap\n" +
			"\t.+/github.com/corestoreio/csfw/util/errors/format_test.go:101",
	}, {
		Wrap(io.EOF, "error"),
		"%s",
		"error: EOF",
	}}

	for _, tt := range tests {
		testFormatRegexp(t, tt.error, tt.format, tt.want)
	}
}

func TestFormatWrapf(t *testing.T) {
	tests := []struct {
		error
		format string
		want   string
	}{{
		Wrapf(New("error"), "error%d", 2),
		"%s",
		"error2: error",
	}, {
		Wrap(io.EOF, "error"),
		"%v",
		"error: EOF",
	}, {
		Wrap(io.EOF, "error"),
		"%+v",
		"EOF\n" +
			"github.com/corestoreio/csfw/util/errors.TestFormatWrapf\n" +
			"\t.+/github.com/corestoreio/csfw/util/errors/format_test.go:131: error",
	}, {
		Wrapf(New("error"), "error%d", 2),
		"%v",
		"error2: error",
	}, {
		Wrapf(New("error"), "error%d", 2),
		"%+v",
		"error\n" +
			"github.com/corestoreio/csfw/util/errors.TestFormatWrapf\n" +
			"\t.+/github.com/corestoreio/csfw/util/errors/format_test.go:141",
	}}

	for _, tt := range tests {
		testFormatRegexp(t, tt.error, tt.format, tt.want)
	}
}

func testFormatRegexp(t *testing.T, arg interface{}, format, want string) {
	got := fmt.Sprintf(format, arg)
	lines := strings.SplitN(got, "\n", -1)
	for i, w := range strings.SplitN(want, "\n", -1) {
		match, err := regexp.MatchString(w, lines[i])
		if err != nil {
			t.Fatal(err)
		}
		if !match {
			t.Errorf("fmt.Sprintf(%q, err): got: %q, want: %q", format, got, want)
		}
	}
}
