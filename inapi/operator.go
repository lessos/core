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
	"fmt"
	"sync"
	"time"
)

var (
	OpActionStart     uint32 = 1 << 1
	OpActionRunning   uint32 = 1 << 2
	OpActionStop      uint32 = 1 << 3
	OpActionStopped   uint32 = 1 << 4
	OpActionDestroy   uint32 = 1 << 5
	OpActionDestroyed uint32 = 1 << 6
	OpActionPending   uint32 = 1 << 11
	OpActionWarning   uint32 = 1 << 12
	OpActionResFree   uint32 = 1 << 24
	OpActionHang      uint32 = 1 << 25
	oplogListMu       sync.RWMutex
	oplogSetsMu       sync.RWMutex

	OpActionDesires = []uint32{
		OpActionStart, OpActionRunning,
		OpActionStop, OpActionStopped,
		OpActionDestroy, OpActionDestroyed,
	}
)

func OpActionValid(op uint32) bool {
	return OpActionAllow(
		OpActionStart|OpActionRunning|
			OpActionStop|OpActionStopped|
			OpActionDestroy|OpActionDestroyed|
			OpActionPending|OpActionWarning,
		op,
	)
}

func OpActionAllow(opbase, op uint32) bool {
	return (op & opbase) == op
}

func OpActionRemove(opbase, op uint32) uint32 {
	return (opbase | op) - (op)
}

func OpActionAppend(opbase, op uint32) uint32 {
	return (opbase | op)
}

func OpActionControlFilter(opbase uint32) uint32 {

	if OpActionAllow(opbase, OpActionDestroy) {
		opbase = OpActionRemove(opbase, OpActionStart|OpActionStop)
	} else if OpActionAllow(opbase, OpActionStop) {
		opbase = OpActionRemove(opbase, OpActionStart)
	} else if OpActionAllow(opbase, OpActionStart) {
		opbase = OpActionRemove(opbase, OpActionStop)
	}

	return opbase
}

func OpActionStrings(action uint32) []string {
	s := []string{}

	if OpActionAllow(action, OpActionStart) {
		s = append(s, "start")
	}

	if OpActionAllow(action, OpActionRunning) {
		s = append(s, "running")
	}

	if OpActionAllow(action, OpActionStop) {
		s = append(s, "stop")
	}

	if OpActionAllow(action, OpActionStopped) {
		s = append(s, "stopped")
	}

	if OpActionAllow(action, OpActionDestroy) {
		s = append(s, "destroy")
	}

	if OpActionAllow(action, OpActionDestroyed) {
		s = append(s, "destroyed")
	}

	if OpActionAllow(action, OpActionPending) {
		s = append(s, "pending")
	}

	if OpActionAllow(action, OpActionWarning) {
		s = append(s, "warning")
	}

	if OpActionAllow(action, OpActionResFree) {
		s = append(s, "resfree")
	}

	if OpActionAllow(action, OpActionHang) {
		s = append(s, "hang")
	}

	return s
}

//
const (
	PbOpLogOK    = "ok"
	PbOpLogInfo  = "info"
	PbOpLogWarn  = "warn"
	PbOpLogError = "error"
	PbOpLogFatal = "fatal"
)

const (
	OpLogNsZoneMasterPodScheduleCharge = "zm/ps/charge"
	OpLogNsZoneMasterPodScheduleAlloc  = "zm/ps/alloc"
)

var (
	OpLogNsZoneMasterPodScheduleRep = func(repId uint32) string {
		if repId > 65535 {
			repId = 65535
		}
		return fmt.Sprintf("zm/ps/rep/%d", repId)
	}
)

type OpLogList []*PbOpLogSets

func (ls *OpLogList) Get(sets_name string) *PbOpLogSets {
	oplogListMu.RLock()
	defer oplogListMu.RUnlock()
	return PbOpLogSetsSliceGet(*ls, sets_name)
}

func (ls *OpLogList) LogSet(sets_name string, version uint32, name, status, msg string) {

	oplogListMu.Lock()
	defer oplogListMu.Unlock()

	sets := PbOpLogSetsSliceGet(*ls, sets_name)
	if sets == nil {
		sets = &PbOpLogSets{
			Name:    sets_name,
			Version: version,
		}
		*ls, _ = PbOpLogSetsSliceSync(*ls, sets)
	}

	if version < sets.Version {
		return
	}

	sets.LogSet(version, name, status, msg)
}

func NewPbOpLogSets(sets_name string, version uint32) *PbOpLogSets {

	return &PbOpLogSets{
		Name:    sets_name,
		Version: version,
	}
}

func (rs *PbOpLogSets) LogSet(version uint32, name, status, message string) {

	oplogSetsMu.Lock()
	defer oplogSetsMu.Unlock()

	if version > 0 && version > rs.Version {
		rs.Version = version
		rs.Items = []*PbOpLogEntry{}
	}

	tn := uint64(time.Now().UnixNano() / 1e6)

	rs.Items, _ = PbOpLogEntrySliceSync(rs.Items, &PbOpLogEntry{
		Name:    name,
		Status:  status,
		Message: message,
		Updated: tn,
	})
}

func (rs *PbOpLogSets) LogSetEntry(entry *PbOpLogEntry) {
	rs.Items, _ = PbOpLogEntrySliceSync(rs.Items, entry)
}

func NewPbOpLogEntry(name, status, message string) *PbOpLogEntry {
	return &PbOpLogEntry{
		Name:    name,
		Status:  status,
		Message: message,
		Updated: uint64(time.Now().UnixNano() / 1e6),
	}
}
