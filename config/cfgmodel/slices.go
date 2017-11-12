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

package cfgmodel

import (
	"bytes"
	"encoding/csv"
	"io"
	"strconv"
	"strings"

	"github.com/corestoreio/pkg/config"
	"github.com/corestoreio/pkg/store/scope"
	"github.com/corestoreio/pkg/util/bufferpool"
	"github.com/corestoreio/pkg/util/slices"
	"github.com/corestoreio/errors"
)

// CSVComma separates CSV values. Default value.
const CSVComma = ','

// WithCSVComma applies a custom CSV separator to the types StringCSV or IntCSV
func WithCSVComma(sep rune) Option {
	return func(b *optionBox) error {
		switch {
		case b.StringCSV != nil:
			b.StringCSV.Comma = sep
		case b.CSV != nil:
			b.CSV.Comma = sep
		case b.IntCSV != nil:
			b.IntCSV.Separator = sep
		}
		return nil
	}
}

// StringCSV represents a path in config.Getter which will be saved as a CSV
// string and returned as a string slice. Separator is a comma. It represents
// not a multi line string!
type StringCSV struct {
	Str
	// Comma is your custom separator, default is constant CSVComma
	Comma rune
}

// NewStringCSV creates a new CSV string type. Acts as a multiselect. Default
// separator: constant CSVComma. It represents not a multi line string! An error
// occurred in the options gets added to the field OptionError which you can
// check.
func NewStringCSV(path string, opts ...Option) StringCSV {
	ret := StringCSV{
		Comma: CSVComma,
		Str:   NewStr(path),
	}
	ret.LastError = (&ret).Option(opts...)
	return ret
}

// Option sets the options and returns the last set previous option
func (str *StringCSV) Option(opts ...Option) error {
	ob := &optionBox{
		baseValue: &str.baseValue,
		StringCSV: str,
	}
	for _, o := range opts {
		if err := o(ob); err != nil {
			return errors.Wrap(err, "[cfgmodel] StringCSV.Option")
		}
	}
	str = ob.StringCSV
	str.baseValue = *ob.baseValue
	return nil
}

// Get returns a string slice. Splits the stored string by comma. Can return
// nil,nil. Empty values will be discarded. Returns a slice containing unique
// entries. No validation will be made.
func (str StringCSV) Get(sg config.Scoped) ([]string, error) {
	s, err := str.Str.Get(sg)
	if err != nil {
		return nil, errors.Wrap(err, "[cfgmodel] Str.Get")
	}
	if s == "" {
		return nil, nil
	}
	var ret slices.String = strings.Split(s, string(str.Comma))
	return ret.Unique(), nil
}

// Write writes a slice with its scope and ID to the writer. Validates the input
// string slice for correct values if set in cfgsource.Slice.
func (str StringCSV) Write(w config.Writer, sl []string, h scope.TypeID) error {
	for _, v := range sl {
		if err := str.ValidateString(v); err != nil {
			return err
		}
	}
	return str.baseValue.Write(w, strings.Join(sl, string(str.Comma)), h)
}

// IntCSV represents a path in config.Getter which will be saved as a CSV string
// and returned as an int64 slice. Separator is a comma. It represents not a
// multi line string!
type IntCSV struct {
	Str
	// Lenient ignores errors in parsing integers
	Lenient bool
	// Separator is your custom separator, default is constant CSVComma
	Separator rune
}

// NewIntCSV creates a new int CSV type. Acts as a multiselect. It represents
// not a multi line string! An error occurred in the options gets added to the
// field OptionError which you can check.
func NewIntCSV(path string, opts ...Option) IntCSV {
	ret := IntCSV{
		Str:       NewStr(path),
		Separator: CSVComma,
	}
	ret.LastError = (&ret).Option(opts...)
	return ret
}

// Option sets the options and returns the last set previous option
func (ic *IntCSV) Option(opts ...Option) error {
	ob := &optionBox{
		baseValue: &ic.baseValue,
		IntCSV:    ic,
	}
	for _, o := range opts {
		if err := o(ob); err != nil {
			return errors.Wrap(err, "[cfgmodel] StringCSV.Option")
		}
	}
	ic = ob.IntCSV
	ic.baseValue = *ob.baseValue
	return nil
}

// Get returns an int slice. Int string gets splited by comma. Can return
// nil,nil. If multiple values cannot be casted to int then the last known error
// gets returned.
func (ic IntCSV) Get(sg config.Scoped) ([]int, error) {
	s, err := ic.Str.Get(sg)
	if err != nil {
		return nil, errors.Wrap(err, "[cfgmodel] Str.Get")
	}
	if s == "" {
		return nil, nil
	}

	csv := strings.Split(s, string(ic.Separator))

	ret := make([]int, 0, len(csv))

	for _, line := range csv {
		line = strings.TrimSpace(line)
		if line != "" {
			v, err := strconv.Atoi(line)
			if err != nil && false == ic.Lenient {
				return ret, errors.NewNotValidf(errIntCSVFailedToConvertToInt, line, err)
			}
			if err == nil {
				ret = append(ret, v)
			}
		}
	}
	return ret, nil
}

// Write writes int values as a CSV string
func (ic IntCSV) Write(w config.Writer, sl []int, h scope.TypeID) error {

	val := bufferpool.Get()
	defer bufferpool.Put(val)

	for i, v := range sl {

		if err := ic.ValidateInt(v); err != nil {
			return errors.Wrap(err, "[cfgmodel] ValidateInt")
		}

		if _, err := val.WriteString(strconv.Itoa(v)); err != nil {
			return errors.Wrapf(err, "[cfgmodel] Value %v", v)
		}
		if i < len(sl)-1 {
			if _, err := val.WriteRune(ic.Separator); err != nil {
				return errors.Wrap(err, "[cfgmodel] WriteRune")
			}
		}
	}
	return ic.baseValue.Write(w, val.String(), h)
}

// CSV represents a path in config.Getter which will be saved as a CSV multi
// line string and returned as a string slice slice. Separator is a comma. New
// line separator: \r and/or \n.
type CSV struct {
	Str
	Comma     rune // field delimiter (set to ',' by NewReader)
	Comment   rune // comment character for start of line
	NewReader func(r io.Reader) *csv.Reader
	NewWriter func(w io.Writer) *csv.Writer
}

// NewCSV creates a new CSV string type which can parse multi line strings with
// a separator. Default separator: constant CSVComma. An error occurred in the
// options gets added to the field OptionError which you can check.
func NewCSV(path string, opts ...Option) CSV {
	ret := CSV{
		Str:       NewStr(path),
		Comma:     CSVComma,
		Comment:   '#',
		NewReader: csv.NewReader,
		NewWriter: csv.NewWriter,
	}
	ret.LastError = (&ret).Option(opts...)
	return ret
}

// Option sets the options and returns the last set previous option
func (c *CSV) Option(opts ...Option) error {
	ob := &optionBox{
		baseValue: &c.baseValue,
		CSV:       c,
	}
	for _, o := range opts {
		if err := o(ob); err != nil {
			return errors.Wrap(err, "[cfgmodel] StringCSV.Option")
		}
	}
	c.Comma = ob.CSV.Comma
	c.Comment = ob.CSV.Comment
	return nil
}

// Get returns a string slice. Splits the stored string by comma and new lines
// by \r and/or \n. Can return nil,nil. Error behaviour: NotValid
func (c CSV) Get(sg config.Scoped) ([][]string, error) {
	s, err := c.Str.Get(sg)
	if err != nil {
		return nil, errors.Wrap(err, "[cfgmodel] Str.Get")
	}
	if s == "" {
		return nil, nil
	}
	r := c.NewReader(bytes.NewBufferString(s))
	r.Comma = c.Comma
	r.Comment = c.Comment // not possible to set currently the comment
	res, err := r.ReadAll()
	if err != nil {
		return nil, errors.NewNotValidf("[cfgmodel] CSV.NewReader.ReadAll: %v", err)
	}
	return res, nil
}

// Write writes a slice with its scope and ID to the writer. Validates the input
// string slice for correct values if set in cfgsource.Slice.
func (c CSV) Write(w config.Writer, csv [][]string, h scope.TypeID) error {
	buf := bufferpool.Get()
	defer bufferpool.Put(buf)

	cw := c.NewWriter(buf)
	cw.Comma = c.Comma
	if err := cw.WriteAll(csv); err != nil {
		return errors.NewNotValidf("[cfgmodel] CSV.NewWriter.WriteAll: %v", err)
	}

	return c.baseValue.Write(w, buf.String(), h)
}
