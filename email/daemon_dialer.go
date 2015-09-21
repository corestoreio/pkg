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
	"github.com/corestoreio/csfw/config"
	"github.com/go-gomail/gomail"
)

const (
	defaultHost = "localhost"
	defaultPort = 25
)

// newPlainDialer stubbed out for tests
var newPlainDialer func(host string, port int, username, password string) *gomail.Dialer = gomail.NewPlainDialer

var _ Dialer = (*gomailPlainDialer)(nil)

// gomailPlainDialer is a wrapper for the interface Dialer.
type gomailPlainDialer struct {
	*gomail.Dialer
}

// SetConfigReader noop method to comply with the interface Dialer.
func (gomailPlainDialer) SetConfigReader(config.Reader) {
	// noop
}

type emailConfig struct {
	Config config.Reader
}

func (c *emailConfig) getHost(s config.ScopeIDer) string {
	h := c.Config.GetString(config.Path(PathSmtpHost), config.ScopeStore(s))
	if h == "" {
		h = defaultHost
	}
	return h
}

func (c *emailConfig) getPort(s config.ScopeIDer) int {
	p := c.Config.GetInt(config.Path(PathSmtpPort), config.ScopeStore(s))
	if p < 1 {
		p = defaultPort
	}
	return p
}

func (c *emailConfig) getUsername(s config.ScopeIDer) string {
	return c.Config.GetString(config.Path(PathSmtpUsername), config.ScopeStore(s))
}

func (c *emailConfig) getPassword(s config.ScopeIDer) string {
	return c.Config.GetString(config.Path(PathSmtpPassword), config.ScopeStore(s))
}

func newEmailConfig(c config.Reader) *emailConfig {
	return &emailConfig{
		Config: c,
	}
}
