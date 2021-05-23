// Copyright 2015 Eryx <evorui аt gmаil dοt cοm>, All rights reserved.
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

package inapi

import (
	"errors"
	"fmt"
	"net"
	"strconv"
	"strings"
)

// default: 127.0.0.1:9529
type HostNodeAddress string

type HostNodeAddresses []HostNodeAddress

func (addr HostNodeAddress) String() string {
	return string(addr)
}

func (addr HostNodeAddress) Valid() bool {

	if addr.Port() < 1 {
		return false
	}

	if v := net.ParseIP(addr.IP()); v == nil || len(v) != 16 {
		return false
	}

	return true
}

func (addr HostNodeAddress) IP() string {

	if pos := strings.LastIndex(string(addr), ":"); pos > 0 {
		return string(addr)[:pos]
	}

	return string(addr)
}

func (addr HostNodeAddress) Port() uint16 {

	if pos := strings.LastIndex(string(addr), ":"); pos > 0 {
		port, _ := strconv.Atoi(string(addr)[pos+1:])
		return uint16(port)
	}

	return 0
}

func (addr *HostNodeAddress) SetIP(ip string) error {

	if addr.Port() > 0 {
		*addr = HostNodeAddress(fmt.Sprintf("%s:%d", ip, addr.Port()))
	} else {
		*addr = HostNodeAddress(ip)
	}

	return nil
}

func (addr *HostNodeAddress) SetPort(port uint16) error {

	if (*addr).IP() != "" {
		*addr = HostNodeAddress(fmt.Sprintf("%s:%d", addr.IP(), port))
	} else {
		*addr = HostNodeAddress(fmt.Sprintf(":%d", port))
	}

	return nil
}

func (it HostNodeAddresses) Equal(it2 HostNodeAddresses) bool {
	if len(it) == len(it2) {
		for _, v := range it {
			found := false
			for _, v2 := range it2 {
				if v == v2 {
					found = true
					break
				}
			}
			if !found {
				return false
			}
		}
		return true
	}
	return false
}

func PrivateIPValid(ipAddr string) error {

	// Private IPv4
	// 10.0.0.0 ~ 10.255.255.255
	// 172.16.0.0 ~ 172.31.255.255
	// 192.168.0.0 ~ 192.168.255.255

	ip := net.ParseIP(ipAddr)
	if ip == nil {
		return errors.New("invalid ip address")
	}

	ip = ip.To4()

	ipa := int(ip[0])
	ipb := int(ip[1])

	if ipa == 10 ||
		(ipa == 172 && ipb >= 16 && ipb <= 31) ||
		(ipa == 192 && ipb == 168) {
		return nil
	}

	return errors.New("invalid private ip address")
}
