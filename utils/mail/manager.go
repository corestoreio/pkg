// Copyright 2015, Cyrill @ Schumacher.fm and the CoreStore contributors
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

package mail

import (
	"bytes"

	"github.com/corestoreio/csfw/config"
)

// @todo Manager

// ManagerOption can be used as an argument in NewManager to configure a manager.
type ManagerOption func(*Manager)

// Manager represents a daemon which must be created via NewManager() function
type Manager struct {
	// lastErrs a collector. While setting options, errors may occur and will
	// be accumulated here for later output in the NewManager() function.
	lastErrs []error

	config config.Reader
}

var _ error = (*Manager)(nil)

// Error implements the error interface. Returns a string where each error has
// been separated by a line break.
func (m *Manager) Error() string {
	var buf bytes.Buffer
	for _, e := range m.lastErrs {
		buf.WriteString(e.Error())
		buf.WriteString("\n")
	}
	return buf.String()
}

// Options applies optional arguments to the daemon
// struct. It returns the last set option. More info about the returned function:
// http://commandcenter.blogspot.com/2014/01/self-referential-functions-and-design.html
func (m *Manager) Option(opts ...ManagerOption) *Manager {
	for _, o := range opts {
		if o != nil {
			o(m)
		}
	}
	return m
}

// NewManager creates a new
func NewManager(opts ...ManagerOption) (*Manager, error) {
	m := &Manager{
		config: config.DefaultManager,
	}
	m.Option(opts...)

	if m.lastErrs != nil {
		return nil, m // because Manager implements error interface
	}
	return m, nil
}
