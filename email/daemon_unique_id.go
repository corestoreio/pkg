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
	"hash"
	"hash/fnv"
	"strconv"

	"github.com/corestoreio/csfw/config"
)

type uniqueID struct {
	config config.ScopedGetter
	lastID uint64
}

func (u *uniqueID) SetConfig(r config.ScopedGetter) {
	u.config = r
}

func (u *uniqueID) getHost() string {
	h := u.config.String(PathSmtpHost)
	if h == "" {
		h = defaultHost
	}
	return h
}

func (u *uniqueID) getPort() int {
	p := u.config.Int(PathSmtpPort)
	if p < 1 {
		p = defaultPort
	}
	return p
}

func (u *uniqueID) getUsername() string {
	return u.config.String(PathSmtpUsername)
}

// ID with which you can identify a daemon connection to the same SMTP server
// independent of the scope ID.
func (u *uniqueID) Get() (id uint64, hasChanged bool) {
	var h hash.Hash64
	h = fnv.New64()
	data := []byte(u.getHost() + strconv.Itoa(u.getPort()) + u.getUsername())
	if _, err := h.Write(data); err != nil {
		PkgLog.Info("mail.daemon.ID", "err", err, "hashWrite", string(data))
		return
	}
	if u.lastID != h.Sum64() {
		u.lastID = h.Sum64()
		return u.lastID, true // ID has changed, means some one updated the configuration.
	}
	return h.Sum64(), false // has not changed
}
