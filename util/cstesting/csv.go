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

package cstesting

import (
	"database/sql/driver"
	"encoding/csv"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/corestoreio/csfw/storage/text"
	"github.com/corestoreio/csfw/util/errors"
)

// CSVOptions applies options to the CSV reader
type csvOptions func(*config)

type config struct {
	cc   *CSVConfig
	path string
	test bool
}

// WithFile sets the file name. File path prefix is always RootPath variable.
func WithFile(elem ...string) csvOptions {
	return func(c *config) { c.path = filepath.Join(elem...) }
}

// CSVConfig allows to set special options when parsing the csv file.
type CSVConfig struct {
	// Comma is the field delimiter.
	// It is set to comma (',') by NewReader.
	Comma rune
	// Comment, if not 0, is the comment character. Lines beginning with the
	// Comment character without preceding whitespace are ignored. With leading
	// whitespace the Comment character becomes part of the field, even if
	// TrimLeadingSpace is true.
	Comment rune
	// FieldsPerRecord is the number of expected fields per record. If
	// FieldsPerRecord is positive, Read requires each record to have the given
	// number of fields. If FieldsPerRecord is 0, Read sets it to the number of
	// fields in the first record, so that future records must have the same
	// field count. If FieldsPerRecord is negative, no check is made and records
	// may have a variable number of fields.
	FieldsPerRecord int
	// If LazyQuotes is true, a quote may appear in an unquoted field and a
	// non-doubled quote may appear in a quoted field.
	LazyQuotes bool
	// If TrimLeadingSpace is true, leading white space in a field is ignored.
	// This is done even if the field delimiter, Comma, is white space.
	TrimLeadingSpace bool
}

// WithReaderConfig sets CSV reader options
func WithReaderConfig(cr CSVConfig) csvOptions {
	return func(c *config) { c.cc = &cr }
}

// WithTestMode allows better testing. Converts []bytes in driver.Value to
// text.Chars
func WithTestMode() csvOptions {
	return func(c *config) { c.test = true }
}

// LoadCSV loads a csv file for mocked database testing. Like
// github.com/DATA-DOG/go-sqlmock does.
// CSV file should be comma separated.
func LoadCSV(opts ...csvOptions) (columns []string, rows [][]driver.Value, err error) {
	cfg := new(config)
	for _, opt := range opts {
		opt(cfg)
	}

	f, err := os.Open(cfg.path)
	if err != nil {
		err = errors.Wrap(err, "[cstesting] os.Open")
		return
	}

	csvReader := csv.NewReader(f)
	if cfg.cc != nil {
		if cfg.cc.Comma > 0 {
			csvReader.Comma = cfg.cc.Comma
		}
		if cfg.cc.Comment > 0 {
			csvReader.Comment = cfg.cc.Comment
		}
		if cfg.cc.FieldsPerRecord > 0 {
			csvReader.FieldsPerRecord = cfg.cc.FieldsPerRecord
		}
		csvReader.LazyQuotes = cfg.cc.LazyQuotes
		csvReader.TrimLeadingSpace = cfg.cc.TrimLeadingSpace
	}

	j := 0
	for {
		var res []string
		res, err = csvReader.Read()
		switch {
		case err == io.EOF:
			err = nil
			return
		case err != nil:
			err = errors.Wrap(err, "[cstesting] csvReader.Read")
			return
		case res == nil:
			err = errors.NewFatalf("[cstesting] Cannot read from csv %q", cfg.path)
			return
		}
		if j == 0 {
			columns = make([]string, len(res))
		}

		row := make([]driver.Value, len(res))
		for i, v := range res {
			v = strings.TrimSpace(v)
			if j == 0 {
				columns[i] = v
			} else {
				b := parseCol(cfg, v)
				row[i] = b
			}
		}
		if j > 0 {
			rows = append(rows, row)
		}
		j++
	}
}

func parseCol(c *config, s string) driver.Value {
	switch {
	case strings.ToLower(s) == "null":
		return nil
	}
	if c.test {
		return text.Chars(s)
	}
	return []byte(s)
}

// MockRows same as LoadCSV() but creates a fully functional driver.Rows
// interface from a CSV file.
func MockRows(opts ...csvOptions) (sqlmock.Rows, error) {
	csvHead, csvRows, err := LoadCSV(opts...)
	if err != nil {
		return nil, errors.Wrap(err, "[cstesting] LoadCSV")
	}
	rows := sqlmock.NewRows(csvHead)
	for _, row := range csvRows {
		rows.AddRow(row...)
	}
	return rows, nil
}

// MustMockRows same as MockRows but panics on error
func MustMockRows(opts ...csvOptions) sqlmock.Rows {
	r, err := MockRows(opts...)
	if err != nil {
		panic(err)
	}
	return r
}

var whiteSpaceRemover = regexp.MustCompile("\\s+")

func SQLMockQuoteMeta(s string) string {
	s = regexp.QuoteMeta(s)
	return whiteSpaceRemover.ReplaceAllString(s, " ")
}
