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
func (gomailPlainDialer) SetConfig(config.ScopedGetter) {
	// noop
}

type emailConfig struct {
	Config config.ScopedGetter
}

func (c *emailConfig) getHost() string {
	h := c.Config.String(PathSmtpHost)
	if h == "" {
		h = defaultHost
	}
	return h
}

func (c *emailConfig) getPort() int {
	p := c.Config.Int(PathSmtpPort)
	if p < 1 {
		p = defaultPort
	}
	return p
}

func (c *emailConfig) getUsername() string {
	return c.Config.String(PathSmtpUsername)
}

func (c *emailConfig) getPassword() string {
	return c.Config.String(PathSmtpPassword)
}

func newEmailConfig(c config.ScopedGetter) *emailConfig {
	return &emailConfig{
		Config: c,
	}
}
