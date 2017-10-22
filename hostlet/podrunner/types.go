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

package podrunner

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"regexp"
	"strings"
	"sync"

	"github.com/hooto/hlog4g/hlog"

	"github.com/sysinner/incore/config"
	"github.com/sysinner/incore/inapi"
	"github.com/sysinner/incore/inutils"
)

const (
	stats_tick      int64  = 5e9
	stats_cycle_buf uint32 = 20
	stats_cycle_log uint32 = 60
	oplog_podpull          = "hostlet/pod-updater"
	oplog_ctncmd           = "hostlet/box-keeper"
)

var (
	vol_podhome_fmt      = "%s/%s.%s/home/action"
	vol_agentsys_dir_fmt = "%s/%s.%s/home/action/.sysinner"
	BoxInstanceNameReg   = regexp.MustCompile("^([0-9a-f]{16,24})-([0-9a-f]{4})-([a-z]{1}[a-z0-9]{0,19})$")
)

func vol_podhome_dir(pod_id string, rep_id uint16) string {
	return fmt.Sprintf(vol_podhome_fmt, config.Config.PodHomeDir,
		pod_id, inutils.Uint16ToHexString(rep_id))
}

func vol_agentsys_dir(pod_id string, rep_id uint16) string {
	return fmt.Sprintf(vol_agentsys_dir_fmt, config.Config.PodHomeDir,
		pod_id, inutils.Uint16ToHexString(rep_id))
}

type BoxInstance struct {
	stats_pending bool
	ID            string
	Name          string
	PodID         string
	RepId         uint16
	PodOpAction   uint32
	PodOpVersion  uint32
	Spec          inapi.PodSpecBoxBound
	Apps          inapi.AppInstances
	Ports         inapi.ServicePorts
	Retry         int
	Env           []inapi.EnvVar
	Status        inapi.PbPodBoxStatus
	Stats         *inapi.TimeStatsFeed
}

func BoxInstanceName(pod_id string, rep *inapi.PodOperateReplica, box_name string) string {
	rep_id := uint16(0)
	if rep != nil {
		rep_id = rep.Id
	}

	return fmt.Sprintf(
		"%s-%s-%s",
		pod_id, inutils.Uint16ToHexString(rep_id), box_name,
	)
}

func BoxInstanceNameParse(hostname string) (pod_id string, rep_id uint16, box_name string) {

	if ns := BoxInstanceNameReg.FindStringSubmatch(hostname); len(ns) == 4 {

		rb, _ := hex.DecodeString(ns[2])
		rep_id = binary.BigEndian.Uint16(rb)

		return ns[1], rep_id, ns[3]
	}

	return "", 0, ""
}

func (inst *BoxInstance) OpRepKey() string {
	return inapi.NsZonePodOpRepKey(inst.PodID, inst.RepId)
}

func (inst *BoxInstance) SpecDesired() bool {

	//
	if inst.Status.Name == "" {
		hlog.Printf("info", "SD")
		return true // wait init
	}

	if inst.Status.Phase == "" {
		hlog.Printf("info", "SD")
		return false
	}

	//
	if inst.Spec.Resources.CpuLimit != inst.Status.ResCpuLimit ||
		inst.Spec.Resources.MemLimit != inst.Status.ResMemLimit {
		hlog.Printf("info", "SD")
		return false
	}

	if len(inst.Ports) != len(inst.Status.Ports) {
		hlog.Printf("info", "SD")
		return false
	}

	for _, v := range inst.Ports {

		mat := false
		for _, vd := range inst.Status.Ports {

			if uint32(v.BoxPort) != vd.BoxPort {
				continue
			}

			if v.HostPort > 0 && uint32(v.HostPort) != vd.HostPort {
				hlog.Printf("info", "SD")
				return false
			}

			mat = true
			break
		}

		if !mat {
			return false
		}
	}

	//
	img2 := inapi.LabelSliceGet(inst.Status.ImageOptions, "docker/image/name")
	if img2 == nil {
		hlog.Printf("info", "SD")
		return false
	}
	img1, _ := inst.Spec.Image.Options.Get("docker/image/name")
	if img2.Value != img1.String() {
		hlog.Printf("info", "SD")
		return false
	}

	//
	if !inapi.PbVolumeMountSliceEqual(inst.Spec.Mounts, inst.Status.Mounts) {
		hlog.Printf("info", "SD")
		return false
	}

	if len(inst.Spec.Command) != len(inst.Status.Command) ||
		strings.Join(inst.Spec.Command, " ") != strings.Join(inst.Status.Command, " ") {
		hlog.Printf("info", "SD")
		return false
	}

	return true
}

func (inst *BoxInstance) OpActionDesired() bool {

	if (inapi.OpActionAllow(inst.PodOpAction, inapi.OpActionStart) && inst.Status.Phase == inapi.OpStatusRunning) ||
		(inapi.OpActionAllow(inst.PodOpAction, inapi.OpActionStop) && inst.Status.Phase == inapi.OpStatusStopped) ||
		(inapi.OpActionAllow(inst.PodOpAction, inapi.OpActionDestroy) && inst.Status.Phase == inapi.OpStatusDestroyed) {
		return true
	}

	return false
}

func (inst *BoxInstance) volume_mounts_refresh() {

	ls := []*inapi.PbVolumeMount{
		{
			Name:      "home",
			MountPath: "/home/action",
			HostDir:   vol_podhome_dir(inst.PodID, inst.RepId),
			ReadOnly:  false,
		},
		{
			Name:      "sysinner/nsz",
			MountPath: "/dev/shm/sysinner/nsz",
			HostDir:   "/dev/shm/sysinner/nsz",
			ReadOnly:  true,
		},
	}

	for _, app := range inst.Apps {

		for _, pkg := range app.Spec.Packages {

			ls, _ = inapi.PbVolumeMountSliceSync(ls, &inapi.PbVolumeMount{
				Name:      "ipm-" + pkg.Name,
				MountPath: ipm_mountpath(pkg.Name, pkg.Version),
				HostDir:   ipm_hostdir(pkg.Name, pkg.Version, pkg.Release, pkg.Dist, pkg.Arch),
				ReadOnly:  true,
			})
		}
	}

	if !inapi.PbVolumeMountSliceEqual(inst.Spec.Mounts, ls) {
		inst.Spec.Mounts = ls
	}
}

func (inst *BoxInstance) volume_mounts_export() []string {

	bindVolumes := []string{}

	for _, v := range inst.Spec.Mounts {

		bindVolume := v.HostDir + ":" + v.MountPath
		if v.ReadOnly {
			bindVolume += ":ro"
		}

		bindVolumes = append(bindVolumes, bindVolume)
	}

	return bindVolumes
}

var box_sets_mu sync.RWMutex

type BoxInstanceSets []*BoxInstance

func (ls *BoxInstanceSets) Get(name string) *BoxInstance {

	box_sets_mu.RLock()
	defer box_sets_mu.RUnlock()

	for _, v := range *ls {
		if v.Name == name {
			return v
		}
	}
	return nil
}

func (ls *BoxInstanceSets) Set(item *BoxInstance) {

	box_sets_mu.Lock()
	defer box_sets_mu.Unlock()

	for i, v := range *ls {
		if v.Name == item.Name {
			(*ls)[i] = item
			return
		}
	}

	*ls = append(*ls, item)
}

func (ls *BoxInstanceSets) Size() int {

	box_sets_mu.RLock()
	defer box_sets_mu.RUnlock()

	return len(*ls)
}

func (ls *BoxInstanceSets) Each(fn func(item *BoxInstance)) {

	box_sets_mu.RLock()
	defer box_sets_mu.RUnlock()

	for _, v := range *ls {
		fn(v)
	}
}
