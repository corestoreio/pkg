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
	"sync"

	"github.com/go-gomail/gomail"
)

// dialerPool avoids using duplicated instances for one and the same connection setting.
// If we have e.g. 50 stores views and each with a different mail setting then you
// have 50 different dailers. This dialerPool uses a hash of host, port and username
// to return an already created dialer for the same settings.
var dialerPool = newPlainDialerPool()

// newPlainDialer stubbed out for tests
var newPlainDialer func(host string, port int, username, password string) *gomail.Dialer = gomail.NewPlainDialer

type dialerContainer struct {
	mu     sync.RWMutex
	dialer map[uint64]Dialer
}

func newPlainDialerPool() *dialerContainer {
	return &dialerContainer{
		dialer: make(map[uint64]Dialer),
	}
}

func (dc *dialerContainer) allocatePlain(dm *Daemon) Dialer {
	dc.mu.Lock()
	defer dc.mu.Unlock()

	id := dm.ID()
	if _, ok := dc.dialer[id]; !ok {
		gmd := newPlainDialer(dm.getHost(), dm.getPort(), dm.getUsername(), dm.getPassword())

		if nil != dm.tlsConfig {
			gmd.TLSConfig = dm.tlsConfig
		}
		dc.dialer[id] = gmd
	}
	return dc.dialer[id]
}
