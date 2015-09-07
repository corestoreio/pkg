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
	"crypto/tls"
	"errors"
	"io"
	"strconv"
	"time"

	"github.com/corestoreio/csfw/config"
	"github.com/corestoreio/csfw/utils/log"
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

// @todo implement config paths and options for TLS certificates and its configuration.

const (
	defaultHost = "localhost"
	defaultPort = 25
)

// OfflineLogger represents a special email logger if mail sending has
// been deactivated for a scope ID. The underlying logger is a NullLogger.
var OfflineLogger log.Logger = new(log.NullLogger)

// OfflineSend defines a function which uses the OfflineLogger.Info function to
// log emails when SMTP has been disabled.
var OfflineSend gomail.SendFunc = offlineSend

func offlineSend(from string, to []string, msg io.WriterTo) error {
	if OfflineLogger.IsInfo() {
		var buf bytes.Buffer
		if _, err := msg.WriteTo(&buf); err != nil {
			return log.Error("mail.daemon.OfflineSend", "err", err, "buf", buf.String(), "message", msg)
		}
		OfflineLogger.Info("mail.Send", "from", from, "to", to, "msg", buf.String())
	}
	return nil
}

// ErrMailChannelClosed will be returned when the channel is closed.
var ErrMailChannelClosed = errors.New("The mail channel has been closed.")

// Dialer mocked out *gomail.Dialer for testing. Sorry for the confusion but
// *gomail.Dialer is the wrong name because ending on "er" means interface
// and not a struct.
type Dialer interface {
	Dial() (gomail.SendCloser, error)
}

// Daemon represents a daemon which must be created via NewDaemon() function
type Daemon struct {
	// lastErrs a collector. While setting options, errors may occur and will
	// be accumulated here for later output in the NewDaemon() function.
	lastErrs    []error
	msgChan     chan *gomail.Message
	dialer      Dialer // usually *gomail.Dialer
	sendFunc    gomail.SendFunc
	closed      bool
	config      config.Reader
	scopeID     config.ScopeIDer
	smtpTimeout time.Duration
	// tlsConfig caches the call to SetTLSConfig because TLS setting can only
	// be applied when the dialer has already been set.
	tlsConfig *tls.Config
}

var _ error = (*Daemon)(nil)

// Error implements the error interface. Returns a string where each error has
// been separated by a line break.
func (dm *Daemon) Error() string {
	var buf bytes.Buffer
	for _, e := range dm.lastErrs {
		buf.WriteString(e.Error())
		buf.WriteString("\n")
	}
	return buf.String()
}

// Start listens to a channel and sends all incoming messages to a SMTP server.
// Errors will be logged. Use code snippet:
//		d := NewDaemon(...)
// 		go func(){
//			if err := d.Worker(); err != nil {
// 				panic(err) // for example
// 			}
// 		}()
//		d.Send(*gomail.Message)
//		d.Stop()
func (dm *Daemon) Worker() error {
	if dm.closed {
		return ErrMailChannelClosed
	}
	if dm.sendFunc != nil {
		return dm.workerSendFunc()
	}
	return dm.workerDial()
}

func (dm *Daemon) workerSendFunc() error {
	for {
		select {
		case m, ok := <-dm.msgChan:
			if !ok {
				dm.closed = true
				return nil
			}

			if err := gomail.Send(dm.sendFunc, m); err != nil {
				log.Error("mail.daemon.Start.Send", "err", err, "message", m)
			}
		}
	}
}

func (dm *Daemon) workerDial() error {
	var s gomail.SendCloser
	var err error
	open := false
	for {
		select {
		case m, ok := <-dm.msgChan:
			if !ok {
				dm.closed = true
				return nil
			}
			if !open {
				if s, err = dm.dialer.Dial(); err != nil {
					return log.Error("mail.daemon.Start.Dial", "err", err, "message", m)
				}
				open = true
			}
			if err := gomail.Send(s, m); err != nil {
				log.Error("mail.daemon.Start.Send", "err", err, "message", m)
			}
		// Close the connection to the SMTP server if no email was sent in
		// the last n seconds.
		case <-time.After(dm.smtpTimeout):
			if open {
				if err := s.Close(); err != nil {
					return log.Error("mail.daemon.Start.Close", "err", err)
				}
				open = false
			}
		}
	}
}

// Stop closes the channel stops the daemon
func (dm *Daemon) Stop() error {
	if dm.closed {
		return ErrMailChannelClosed
	}
	close(dm.msgChan)
	dm.closed = true
	return nil
}

// Send sends a mail
func (dm *Daemon) Send(m *gomail.Message) error {
	if dm.closed {
		return ErrMailChannelClosed
	}
	dm.msgChan <- m
	return nil
}

// SendPlain sends a simple plain text email
func (dm *Daemon) SendPlain(from, to, subject, body string) error {
	return dm.sendMsg(from, to, subject, body, false)
}

// SendPlain sends a simple html email
func (dm *Daemon) SendHtml(from, to, subject, body string) error {
	return dm.sendMsg(from, to, subject, body, true)
}

func (dm *Daemon) sendMsg(from, to, subject, body string, isHtml bool) error {
	if dm.closed {
		return ErrMailChannelClosed
	}
	m := gomail.NewMessage()
	m.SetHeader("From", from)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	contentType := "text/plain"
	if isHtml {
		contentType = "text/html"
	}
	m.SetBody(contentType, body)
	dm.Send(m)
	return nil
}

// Options applies optional arguments to the daemon
// struct. It returns the last set option. More info about the returned function:
// http://commandcenter.blogspot.com/2014/01/self-referential-functions-and-design.html
func (dm *Daemon) Option(opts ...DaemonOption) (previous DaemonOption) {
	for _, o := range opts {
		if o != nil {
			previous = o(dm)
		}
	}
	return previous
}

// IsOffline checks if SMTP sending for the current scope ID has been deactivated.
// If disabled the output will be logged.
func (dm *Daemon) IsOffline() bool {
	if nil == dm.config {
		return true
	}
	return dm.config.GetBool(config.Path(PathSmtpDisable), config.ScopeStore(dm.scopeID))
}

func (dm *Daemon) getHost() string {
	h := dm.config.GetString(config.Path(PathSmtpHost), config.ScopeStore(dm.scopeID))
	if h == "" {
		h = defaultHost
	}
	return h
}

func (dm *Daemon) getPort() int {
	p := dm.config.GetInt(config.Path(PathSmtpPort), config.ScopeStore(dm.scopeID))
	if p < 1 {
		p = defaultPort
	}
	return p
}

func (dm *Daemon) getUsername() string {
	return dm.config.GetString(config.Path(PathSmtpUsername), config.ScopeStore(dm.scopeID))
}

func (dm *Daemon) getPassword() string {
	// @todo decryption of the encrpted stored password? or extend config.Reader API?
	return dm.config.GetString(config.Path(PathSmtpPassword), config.ScopeStore(dm.scopeID))
}

// internalID with which you can identify a daemon connection to the same SMTP server
// independent of the scope ID. Only used so far for the dialer pool.
func (dm *Daemon) internalID() []byte {
	return []byte(dm.getHost() + strconv.Itoa(dm.getPort()) + dm.getUsername())
}

// NewDaemon creates a new mail sending daemon to send to a SMTP server.
// Per default it uses localhost:25, creates an unbuffered channel, uses the
// config.DefaultManager, applies the admin scope (0) and sets the SMTP
// timeout to 30s.
func NewDaemon(opts ...DaemonOption) (*Daemon, error) {
	d := &Daemon{
		config:      config.DefaultManager,
		scopeID:     config.ScopeID(0), // Default Scope aka Admin Scope
		smtpTimeout: time.Second * 30,
	}
	d.Option(opts...)

	if d.IsOffline() {
		SetSendFunc(OfflineSend)(d)
	}

	if d.msgChan == nil {
		d.msgChan = make(chan *gomail.Message)
	}

	if nil == d.dialer && nil == d.sendFunc {
		d.dialer = dialerPool.allocatePlain(d)
	}

	if d.lastErrs != nil {
		return nil, d // because Daemon implements error interface
	}
	return d, nil
}
