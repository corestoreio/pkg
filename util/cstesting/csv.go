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
	"strings"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/corestoreio/csfw/storage/text"
	"github.com/juju/errgo"
)

// CSVOptions applies options to the CSV reader
type csvOptions func(*config)

type config struct {
	r    *csv.Reader
	path string
	test bool
}

// WithFile sets the file name. File path prefix is always RootPath variable.
func WithFile(elem ...string) csvOptions {
	return func(c *config) { c.path = filepath.Join(append([]string{RootPath}, elem...)...) }
}

// WithReaderConfig sets CSV reader options
func WithReaderConfig(cr *csv.Reader) csvOptions {
	return func(c *config) { c.r = cr }
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
		err = errgo.Mask(err)
		return
	}

	csvReader := csv.NewReader(f)
	if cfg.r != nil {
		csvReader.Comma = cfg.r.Comma
		csvReader.Comment = cfg.r.Comment
		csvReader.FieldsPerRecord = cfg.r.FieldsPerRecord
		csvReader.LazyQuotes = cfg.r.LazyQuotes
		csvReader.TrailingComma = cfg.r.TrailingComma
		csvReader.TrimLeadingSpace = cfg.r.TrimLeadingSpace
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
			return
		case res == nil:
			err = errgo.New("Cannot read from csv")
			return
		}
		if j == 0 {
			columns = make([]string, len(res), len(res))
		}

		row := make([]driver.Value, len(res))
		for i, v := range res {
			v = strings.TrimSpace(v)
			if j == 0 {
				columns[i] = v
			} else {
				b := parseCol(v)
				row[i] = b
				if cfg.test {
					row[i] = text.Chars(b)
				}
			}
		}
		if j > 0 {
			rows = append(rows, row)
		}
		j++
	}
	return
}

func parseCol(s string) []byte {
	switch {
	case strings.ToLower(s) == "null":
		return nil
	}
	return []byte(s)
}

// MockRows same as LoadCSV() but creates a fully functional driver.Rows
// interface from a CSV file.
func MockRows(opts ...csvOptions) (sqlmock.Rows, error) {
	csvHead, csvRows, err := LoadCSV(opts...)
	if err != nil {
		return nil, err
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
