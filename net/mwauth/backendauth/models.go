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

package backendauth

import (
	"bytes"
	"net"

	"github.com/corestoreio/csfw/config"
	"github.com/corestoreio/csfw/config/cfgmodel"
	"github.com/corestoreio/csfw/util/errors"
)

type IPRange struct {
	from net.IP
	to   net.IP
}

func NewIPRange(from, to string) IPRange {
	return IPRange{
		from: net.ParseIP(from).To16(),
		to:   net.ParseIP(to).To16(),
	}
}

func (ir IPRange) In(test net.IP) bool {
	tv6 := test.To16()
	return ir.from != nil && ir.to != nil && tv6 != nil && bytes.Compare(tv6, ir.from) >= 0 && bytes.Compare(tv6, ir.to) <= 0
}

func (ir IPRange) InStr(ip string) bool {
	return ir.In(net.ParseIP(ip))
}

// ConfigIPRange defines how IP ranges are stored and handled.
// A valid IP range string looks like for example:
// 		IPv4: 74.50.153.0-74.50.153.4
// 		IPv6: ::ffff:192.0.2.128-::ffff:192.0.2.250
// 		IPv6: 2001:0db8:85a3:0000:0000:8a2e:0370:7334-2001:0db8:85a3:0000:0000:8a2e:0370:8334
// No white spaces! Multiple entries supported via line break \n
type ConfigIPRange struct {
	cfgmodel.StringCSV
}

// NewConfigIPRange ....
// A valid IP range string looks like for example:
// 		IPv4: 74.50.153.0-74.50.153.4
// 		IPv6: ::ffff:192.0.2.128-::ffff:192.0.2.250
// 		IPv6: 2001:0db8:85a3:0000:0000:8a2e:0370:7334-2001:0db8:85a3:0000:0000:8a2e:0370:8334
// No white spaces! Multiple entries supported via line break \n
func NewConfigIPRange(path string, opts ...cfgmodel.Option) ConfigIPRange {
	cip := ConfigIPRange{
		StringCSV: cfgmodel.NewStringCSV(path, opts...),
	}
	cip.Separator = '-'
	return cip
}

// Get ...
func (cc ConfigIPRange) Get(sg config.ScopedGetter) (IPRange, error) {
	raw, err := cc.StringCSV.Get(sg)
	if err != nil {
		return IPRange{}, errors.Wrap(err, "[backendauth] Str.Get")
	}
	if len(raw) != 2 {
		return IPRange{}, errors.NewNotValidf("[backendauth] IP Range %q not in expected format: IP.From-IP.To", raw)
	}
	return NewIPRange(raw[0], raw[1]), nil
}
