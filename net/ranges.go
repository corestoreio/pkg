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

package net

import (
	"bytes"
	gonet "net"
)

// IPRange contains a from and to IP address. v4 or v6 doesn't matter.
type IPRange struct {
	from gonet.IP
	to   gonet.IP
}

// IPRanges contains multiple IPRange entries. v4 or v6 doesn't matter.
type IPRanges []IPRange

// NewIPRange creates a new instance.
func NewIPRange(from, to string) IPRange {
	return IPRange{
		from: gonet.ParseIP(from).To16(),
		to:   gonet.ParseIP(to).To16(),
	}
}

// In checks if test IP lies within the range.
func (ir IPRange) In(test gonet.IP) bool {
	tv6 := test.To16()
	return ir.from != nil && ir.to != nil && tv6 != nil && bytes.Compare(tv6, ir.from) >= 0 && bytes.Compare(tv6, ir.to) <= 0
}

// InStr checks if the test IP address string lies within the range.
func (ir IPRange) InStr(ip string) bool {
	return ir.In(gonet.ParseIP(ip))
}

// In checks if test IP lies within the ranges.
func (s IPRanges) In(test gonet.IP) bool {
	for _, ir := range s {
		if ir.In(test) {
			return true
		}
	}
	return false
}

// InStr checks if the test IP address string lies within the ranges.
func (s IPRanges) InStr(ip string) bool {
	return s.In(gonet.ParseIP(ip))
}

// PrivateIPRanges defines a list of private IP subnets.
// Can be modified by yourself.
var PrivateIPRanges = IPRanges{
	NewIPRange("10.0.0.0", "10.255.255.255"),
	NewIPRange("100.64.0.0", "100.127.255.255"),
	NewIPRange("172.16.0.0", "172.31.255.255"),
	NewIPRange("192.0.0.0", "192.0.0.255"),
	NewIPRange("192.168.0.0", "192.168.255.255"),
	NewIPRange("198.18.0.0", "198.19.255.255"),
	// NewIPRange("fc00::/7", "fc00::/7"),
}
