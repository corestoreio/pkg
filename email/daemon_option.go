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
	"crypto/tls"
	"errors"

	"github.com/go-gomail/gomail"
)

// DaemonOption can be used as an argument in NewDaemon to configure a daemon.
type DaemonOption func(*Daemon) DaemonOption

// SetMessageChannel sets your custom channel to listen to.
func SetMessageChannel(mailChan chan *gomail.Message) DaemonOption {
	return func(da *Daemon) DaemonOption {
		if mailChan == nil {
			da.lastErrs = append(da.lastErrs, errors.New("*gomail.Message channel cannot be nil"))
		}
		previous := da.msgChan
		da.msgChan = mailChan
		da.closed = false
		return SetMessageChannel(previous)
	}
}

// SetDialer sets a custom dialer, e.g. for a different smtp.Auth use.
// Usually a *gomail.Dialer. If not provided falls back the plain auth dialer
// of gomail. Applying the SetDialer with set the sendFunc to nil.
func SetDialer(di Dialer) DaemonOption {
	return func(da *Daemon) DaemonOption {
		previous := da.dialer
		if di == nil {
			da.lastErrs = append(da.lastErrs, errors.New("gomail.Dialer cannot be nil"))
		}
		da.dialer = di
		da.dialerIsCustom = true
		da.sendFunc = nil
		return SetDialer(previous)
	}
}

// SetSendFunc lets you implements your email-sending function for e.g.
// to use any other third party API provider. Setting this option
// will remove the dialer. Your implementation must handle timeouts, etc.
func SetSendFunc(sf gomail.SendFunc) DaemonOption {
	return func(da *Daemon) DaemonOption {
		previous := da.sendFunc
		if sf == nil {
			da.lastErrs = append(da.lastErrs, errors.New("gomail.SendFunc cannot be nil"))
		}
		da.sendFunc = sf
		da.dialer = nil
		return SetSendFunc(previous)
	}
}

// SetTLSConfig sets the TLS configuration for a default plain dialer used for TLS
// (when the STARTTLS extension is used) or SSL connections.
func SetTLSConfig(c *tls.Config) DaemonOption {
	return func(da *Daemon) DaemonOption {
		previous := da.tlsConfig
		if nil == c {
			da.lastErrs = append(da.lastErrs, errors.New("*tls.Config cannot be nil"))
		}
		da.tlsConfig = c
		return SetTLSConfig(previous)
	}
}
