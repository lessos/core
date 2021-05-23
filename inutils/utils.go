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

package inutils

import (
	"encoding/binary"
	"encoding/hex"
	"os/exec"
	"strings"
)

var (
	hostKernel = ""
)

func ResSysHostKernel() string {

	if hostKernel == "" {
		cmd, err := exec.LookPath("uname")
		if err == nil {
			rs, err := exec.Command(cmd, "-r").Output()
			if err == nil {
				hostKernel = strings.TrimSpace(string(rs))
			}
		}
	}

	return hostKernel
}

func Uint16ToHexString(v uint16) string {
	bs := make([]byte, 2)
	binary.BigEndian.PutUint16(bs, v)
	return hex.EncodeToString(bs)
}

func Uint32ToBytes(v uint32) []byte {
	bs := make([]byte, 4)
	binary.BigEndian.PutUint32(bs, v)
	return bs
}

func Uint32ToHexString(v uint32) string {
	bs := make([]byte, 4)
	binary.BigEndian.PutUint32(bs, v)
	return hex.EncodeToString(bs)
}

func Uint64ToBytes(v uint64) []byte {
	bs := make([]byte, 8)
	binary.BigEndian.PutUint64(bs, v)
	return bs
}

func BytesToUint16(bs []byte) uint16 {
	return binary.BigEndian.Uint16(bs)
}

func BytesToUint32(bs []byte) uint32 {
	return binary.BigEndian.Uint32(bs)
}

func BytesToUint64(bs []byte) uint64 {
	return binary.BigEndian.Uint64(bs)
}

func BytesCopy(src []byte) []byte {
	dst := make([]byte, len(src))
	copy(dst, src)
	return dst
}

func ArrayUint32Has(ar []uint32, v uint32) bool {
	for _, v2 := range ar {
		if v2 == v {
			return true
		}
	}
	return false
}
