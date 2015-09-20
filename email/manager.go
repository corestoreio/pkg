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

package email

import (
	"bytes"

	"sync"

	"github.com/corestoreio/csfw/config"
	"github.com/go-gomail/gomail"
)

// PathSmtp* defines the configuration settings for a SMTP daemon.
const (
	PathSmtpDisable         = "system/smtp/disable"           // Scope: Default, Website, Store
	PathSmtpHost            = "system/smtp/host"              // Scope: Default, Website, Store
	PathSmtpPort            = "system/smtp/port"              // Scope: Default, Website, Store
	PathSmtpUsername        = "system/smtp/username"          // Scope: Default, Website, Store
	PathSmtpPassword        = "system/smtp/password"          // Scope: Default, Website, Store
	PathSmtpSetReturnPath   = "system/smtp/set_return_path"   // Scope: Default; 0 = no, 1 = yes, 2 = specified in PathSmtpReturnPathEmail
	PathSmtpReturnPathEmail = "system/smtp/return_path_email" // Scope: Default; email address
	PathSmtpMandrillAPIKey  = "system/smtp/mandrill_api_key"  // Scope: Default, Website, Store @todo
)

const (
	defaultHost = "localhost"
	defaultPort = 25
)

// TOOD(cs) implement config paths and options for TLS certificates and its configuration.
// TOOD(cs) implement ideas from https://github.com/nathan-osman/go-cannon

// ManagerOption can be used as an argument in NewManager to configure a manager.
type ManagerOption func(*Manager)

// Manager represents a daemon which must be created via NewManager() function.
// A manager starts and stops a daemon. Restarts happens on config changes.
type Manager struct {
	// lastErrs a collector. While setting options, errors may occur and will
	// be accumulated here for later output in the NewManager() function.
	lastErrs []error

	Config config.Reader

	mu     sync.RWMutex
	dialer map[uint64]Dialer
}

var _ error = (*Manager)(nil)
var _ config.MessageReceiver = (*Manager)(nil)

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

func (m *Manager) Send(si config.ScopeIDer, m *gomail.Message) error {

	return nil
}

// SubscribeToConfigChanges subcribes the function MessageConfig to the
// config.Subscriber
func (m *Manager) SubscribeToConfigChanges(s config.Subscriber) config.MessageReceiver {

	return m
}

// MessageConfig allows subscription to the publish/subscribe message system of
// config.Manager. MessageConfig will be added via SubscribeToConfigChanges to the
// config.Subscriber.
// IF a configuration change
func (m *Manager) MessageConfig(path string, sg config.ScopeGroup, si config.ScopeIDer) {

}

func (m *Manager) getHost(s config.ScopeIDer) string {
	h := m.Config.GetString(config.Path(PathSmtpHost), config.ScopeStore(s))
	if h == "" {
		h = defaultHost
	}
	return h
}

func (m *Manager) getPort(s config.ScopeIDer) int {
	p := m.Config.GetInt(config.Path(PathSmtpPort), config.ScopeStore(s))
	if p < 1 {
		p = defaultPort
	}
	return p
}

func (m *Manager) getUsername(s config.ScopeIDer) string {
	return m.Config.GetString(config.Path(PathSmtpUsername), config.ScopeStore(s))
}

func (m *Manager) getPassword(s config.ScopeIDer) string {
	return m.Config.GetString(config.Path(PathSmtpPassword), config.ScopeStore(s))
}

func (m *Manager) allocate(dm *Daemon) Dialer {
	m.mu.Lock()
	defer m.mu.Unlock()

	id, changed := dm.UniqueID.Get()

	if dm.lastIDchanged {
		if _, ok := m.dialer[dm.lastID]; ok {
			m.dialer[dm.lastID] = nil // current dialer will be GCed
			delete(m.dialer, dm.lastID)
		}
		dm.lastIDchanged = false
	}

	if _, ok := m.dialer[id]; !ok {

		m.dialer[id] = dm.dialer

		if nil == m.dialer[id] {
			nd := &gomailPlainDialer{
				Dialer: newPlainDialer(m.getHost(), m.getPort(), m.getUsername(), m.getPassword()),
			}

			if nil != dm.tlsConfig {
				nd.TLSConfig = dm.tlsConfig
			}
			m.dialer[id] = nd
		}
	}
	return m.dialer[id]
}

func NewManager(opts ...ManagerOption) (*Manager, error) {
	m := &Manager{
		Config: config.DefaultManager,
		dialer: make(map[uint64]Dialer),
	}
	m.Option(opts...)

	if m.lastErrs != nil {
		return nil, m // because Manager implements error interface
	}
	return m, nil
}
