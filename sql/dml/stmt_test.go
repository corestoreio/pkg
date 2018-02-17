// Copyright 2015-present, Cyrill @ Schumacher.fm and the CoreStore contributors
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

package dml

import (
	"testing"

	"github.com/corestoreio/errors"
	"github.com/stretchr/testify/assert"
)

func TestStmtWrapper(t *testing.T) {
	t.Parallel()

	t.Run("PrepareContext", func(t *testing.T) {
		defer func() {
			if r := recover(); r != nil {
				if err, ok := r.(error); ok {
					assert.True(t, errors.Is(err, errors.NotAllowed), "Error should have kind errors.NotAllowed")
				} else {
					t.Errorf("Panic should contain an error but got:\n%+v", r)
				}
			} else {
				t.Error("Expecting a panic but got nothing")
			}
		}()

		sw := stmtWrapper{}
		sw.PrepareContext(nil, "SELECT X")
	})

	t.Run("ExecContext", func(t *testing.T) {
		defer func() {
			if r := recover(); r != nil {
				if err, ok := r.(error); ok {
					assert.True(t, errors.Is(err, errors.NotAllowed), "Error should have kind errors.NotAllowed")
				} else {
					t.Errorf("Panic should contain an error but got:\n%+v", r)
				}
			} else {
				t.Error("Expecting a panic but got nothing")
			}
		}()

		sw := stmtWrapper{}
		sw.ExecContext(nil, "SELECT X")
	})

	t.Run("QueryContext", func(t *testing.T) {
		defer func() {
			if r := recover(); r != nil {
				if err, ok := r.(error); ok {
					assert.True(t, errors.Is(err, errors.NotAllowed), "Error should have kind errors.NotAllowed")
				} else {
					t.Errorf("Panic should contain an error but got:\n%+v", r)
				}
			} else {
				t.Error("Expecting a panic but got nothing")
			}
		}()

		sw := stmtWrapper{}
		sw.QueryContext(nil, "SELECT X")
	})

	t.Run("QueryRowContext", func(t *testing.T) {
		defer func() {
			if r := recover(); r != nil {
				if err, ok := r.(error); ok {
					assert.True(t, errors.Is(err, errors.NotAllowed), "Error should have kind errors.NotAllowed")
				} else {
					t.Errorf("Panic should contain an error but got:\n%+v", r)
				}
			} else {
				t.Error("Expecting a panic but got nothing")
			}
		}()

		sw := stmtWrapper{}
		sw.QueryRowContext(nil, "SELECT X")
	})
}
