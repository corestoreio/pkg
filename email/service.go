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
	"sync"

	"github.com/corestoreio/csfw/config"
	"github.com/corestoreio/csfw/store/scope"
	"github.com/corestoreio/csfw/util"
	"github.com/go-gomail/gomail"
)

// PathSmtp* defines the configuration settings for a SMTP daemon.
const (
	PathSmtp                = "system/smtp"                   // Used for pubsub
	PathSmtpDisable         = PathSmtp + "/disable"           // Scope: Default, Website, Store
	PathSmtpHost            = PathSmtp + "/host"              // Scope: Default, Website, Store
	PathSmtpPort            = PathSmtp + "/port"              // Scope: Default, Website, Store
	PathSmtpUsername        = PathSmtp + "/username"          // Scope: Default, Website, Store
	PathSmtpPassword        = PathSmtp + "/password"          // Scope: Default, Website, Store
	PathSmtpSetReturnPath   = PathSmtp + "/set_return_path"   // Scope: Default; 0 = no, 1 = yes, 2 = specified in PathSmtpReturnPathEmail
	PathSmtpReturnPathEmail = PathSmtp + "/return_path_email" // Scope: Default; email address
	PathSmtpMandrillAPIKey  = PathSmtp + "/mandrill_api_key"  // Scope: Default, Website, Store @todo
)

// TODO(cs) implement config paths and options for TLS certificates and its configuration.
// TODO(cs) implement ideas from https://github.com/nathan-osman/go-cannon

// ManagerOption can be used as an argument in NewManager to configure a manager.
type ServiceOption func(*Service)

// Manager represents a daemon which must be created via NewManager() function.
// A manager starts and stops a daemon. Restarts happens on config changes.
type Service struct {
	// lastErrs a collector. While setting options, errors may occur and will
	// be accumulated here for later output in the NewManager() function.
	lastErrs []error

	*emailConfig

	mu     sync.RWMutex
	dialer map[uint64]Dialer
}

var _ error = (*Service)(nil)
var _ config.MessageReceiver = (*Service)(nil)

// Error implements the error interface. Returns a string where each error has
// been separated by a line break.
func (s *Service) Error() string {
	return util.Errors(s.lastErrs...)
}

// Options applies optional arguments to the daemon
// struct. It returns the last set option. More info about the returned function:
// http://commandcenter.blogspot.com/2014/01/self-referential-functions-and-design.html
func (s *Service) Option(opts ...ServiceOption) *Service {
	for _, o := range opts {
		if o != nil {
			o(s)
		}
	}
	return s
}

func (s *Service) Send(sc scope.Scope, id int64, m *gomail.Message) error {

	return nil
}

// SubscribeToConfigChanges subscribes the function MessageConfig to the
// config.Subscriber
func (s *Service) SubscribeToConfigChanges(sub config.Subscriber) (subscriptionID int, err error) {
	return sub.Subscribe(PathSmtp, s)
}

// MessageConfig allows subscription to the publish/subscribe message system of
// config.Service. MessageConfig will be added via SubscribeToConfigChanges to the
// config.Subscriber.
// IF a configuration change
func (s *Service) MessageConfig(path string, sc scope.Scope, id int64) error {
	switch path {
	case PathSmtpHost, PathSmtpPort, PathSmtpUsername:
		// start and stop the daemon for the corresponding scope group and id
	case PathSmtpDisable:
		// stop daemon and replace dialer
	}
	return nil
}

func (s *Service) allocate(dm *Daemon) Dialer {
	s.mu.Lock()
	defer s.mu.Unlock()

	//	id, changed := dm.UniqueID.Get()
	//
	//	if dm.lastIDchanged {
	//		if _, ok := m.dialer[dm.lastID]; ok {
	//			m.dialer[dm.lastID] = nil // current dialer will be GCed
	//			delete(m.dialer, dm.lastID)
	//		}
	//		dm.lastIDchanged = false
	//	}
	//
	//	if _, ok := m.dialer[id]; !ok {
	//
	//		m.dialer[id] = dm.dialer
	//
	//		if nil == m.dialer[id] {
	//			nd := &gomailPlainDialer{
	//				Dialer: newPlainDialer(m.getHost(), m.getPort(), m.getUsername(), m.getPassword()),
	//			}
	//
	//			if nil != dm.tlsConfig {
	//				nd.TLSConfig = dm.tlsConfig
	//			}
	//			m.dialer[id] = nd
	//		}
	//	}
	//	return m.dialer[id]
	return nil
}

func NewService(opts ...ServiceOption) (*Service, error) {
	// load available configs from DB ... implement auto refresh
	s := &Service{
		//emailConfig: newEmailConfig(config.DefaultService),
		dialer: make(map[uint64]Dialer),
	}
	s.Option(opts...)

	if s.lastErrs != nil {
		return nil, s // because Manager implements error interface
	}
	return s, nil
}
