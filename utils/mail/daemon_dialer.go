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
	"hash"
	"hash/fnv"
	"sync"

	"github.com/go-gomail/gomail"
)

// DefaultDialer connects to localhost on port 25.
var DefaultDialer = gomail.NewPlainDialer("localhost", 25, "", "")

var uniqueDialerCheck = newUniqueDialer()

type uniqueDialer struct {
	mu     sync.RWMutex
	dialer map[uint64]*gomail.Dialer
	hash   hash.Hash64
}

func newUniqueDialer() *uniqueDialer {
	return &uniqueDialer{
		dialer: make(map[uint64]*gomail.Dialer),
		hash:   fnv.New64(),
	}
}

func (ud *uniqueDialer) register(host, port, username string, d *gomail.Dialer) {
	h := fnv.New64()
	h.Write([]byte(host + port + username))
	ud.mu.Lock()
	ud.dialer[h.Sum64()] = d
	ud.mu.Unlock()
}
