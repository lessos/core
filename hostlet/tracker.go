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

package hostlet

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
	"time"

	"github.com/hooto/hlog4g/hlog"
	"github.com/lessos/lessgo/encoding/json"
	"github.com/lessos/lessgo/types"
	"golang.org/x/net/context"

	ps_cpu "github.com/shirou/gopsutil/cpu"
	ps_disk "github.com/shirou/gopsutil/disk"
	ps_mem "github.com/shirou/gopsutil/mem"
	ps_net "github.com/shirou/gopsutil/net"

	"github.com/sysinner/incore/config"
	// "github.com/sysinner/incore/data"
	"github.com/sysinner/incore/hostlet/napi"
	"github.com/sysinner/incore/hostlet/nstatus"
	"github.com/sysinner/incore/inapi"
	"github.com/sysinner/incore/inutils"
	"github.com/sysinner/incore/rpcsrv"
	"github.com/sysinner/incore/status"
)

var (
	statsPodRepNames = []string{
		"ram/us", "ram/cc",
		"net/rs", "net/ws",
		"cpu/us",
		"fs/rn", "fs/rs", "fs/wn", "fs/ws",
	}
	statsHostNames = []string{
		"ram/us", "ram/cc",
		"net/rs", "net/ws",
		"cpu/sys", "cpu/user",
		"fs/sp/rn", "fs/sp/rs", "fs/sp/wn", "fs/sp/ws",
	}
)

func statusTracker() {

	//
	// if len(status.LocalZoneMasterList.Items) == 0 {

	// 	var zms inapi.ResZoneMasterList
	// 	if rs := data.LocalDB.PvGet(inapi.NsLocalZoneMasterList()); rs.OK() {

	// 		if err := rs.Decode(&zms); err == nil {

	// 			if synced := status.LocalZoneMasterList.SyncList(zms); synced {

	// 				cms := []inapi.HostNodeAddress{}
	// 				for _, v := range status.LocalZoneMasterList.Items {
	// 					cms = append(cms, inapi.HostNodeAddress(v.Addr))
	// 				}

	// 				if len(cms) > 0 {
	// 					config.Config.Masters = cms
	// 					config.Config.Sync()
	// 				}
	// 			}
	// 		}
	// 	}
	// }

	//
	if len(status.LocalZoneMasterList.Items) == 0 && len(config.Config.Masters) == 0 {
		hlog.Printf("warn", "No MasterList.Items Found")
		return
	}

	zms, err := msgZoneMasterHostStatusSync()
	if err != nil {
		hlog.Printf("warn", "ZoneMaster HostStatusSync: %s", err.Error())
		return
	}

	if zms.ExpPods != nil {
		for _, v := range zms.ExpPods {
			var pod inapi.Pod
			if err := json.Decode([]byte(v), &pod); err == nil {
				// inapi.ObjPrint("name", pod)
				nstatus.PodQueue.Set(&pod)
			}
		}
	}

	// fmt.Println(zms)
	if zms.Masters != nil {
		if chg := status.LocalZoneMasterList.SyncList(*zms.Masters); chg {
			cms := []inapi.HostNodeAddress{}
			for _, v := range status.LocalZoneMasterList.Items {
				cms = append(cms, inapi.HostNodeAddress(v.Addr))
			}
			if len(cms) > 0 {
				config.Config.Masters = cms
				config.Config.Sync()
				hlog.Printf("warn", "zone-master/list refreshed")
			}
		}
	}

	if zms.ExpPsmaps != nil {
		sync_nsz(zms.ExpPsmaps)
	}

	if len(zms.ZoneInpackServiceUrl) > 10 && zms.ZoneInpackServiceUrl != config.Config.InpackServiceUrl {
		config.Config.InpackServiceUrl = zms.ZoneInpackServiceUrl
		config.Config.Sync()
	}

	if err := podVolQuotaRefresh(); err != nil {
		hlog.Printf("error", "failed to refresh quota projects: %s", err.Error())
	}
}

var (
	sync_nsz_lasts          = map[string]uint64{}
	hostSyncVolLasted int64 = 0
	sync_nszs               = []*inapi.NsPodServiceMap{}
	sync_nsz_path           = "/dev/shm/sysinner/nsz"
)

func sync_nsz(ls []*inapi.NsPodServiceMap) {

	if len(sync_nszs) == 0 {
		os.MkdirAll(sync_nsz_path, 0755)

		filepath.Walk(sync_nsz_path, func(path string, info os.FileInfo, err error) error {
			if err != nil || info.IsDir() {
				return nil
			}
			if !inapi.PodIdReg.MatchString(info.Name()) {
				return nil
			}
			var nsz inapi.NsPodServiceMap
			if err := json.DecodeFile(path, &nsz); err != nil {
				os.Remove(path)
			} else {
				sync_nszs = append(sync_nszs, &nsz)
			}
			return nil
		})
	}

	for _, v := range sync_nszs {
		if p := inapi.NsPodServiceMapSliceGet(ls, v.Id); p == nil {
			os.Remove(sync_nsz_path + "/" + v.Id)
			if _, ok := sync_nsz_lasts[v.Id]; ok {
				delete(sync_nsz_lasts, v.Id)
			}
		}
	}

	for _, v := range ls {

		if len(v.Id) < 8 {
			continue
		}

		last, ok := sync_nsz_lasts[v.Id]
		if !ok || v.Updated > last {
			json.EncodeToFile(v, sync_nsz_path+"/"+v.Id, "")
			sync_nsz_lasts[v.Id] = v.Updated
		}
	}

	if !inapi.NsPodServiceMapSliceEqual(sync_nszs, ls) {
		sync_nszs = ls
	}
}

func msgZoneMasterHostStatusSync() (*inapi.ResHostBound, error) {

	//
	addr := status.LocalZoneMasterList.LeaderAddr()
	if addr == "" {
		if len(config.Config.Masters) > 0 {
			addr = string(config.Config.Masters[0])
		}
		if addr == "" {
			return nil, errors.New("No MasterList.LeaderAddr Found")
		}
	}
	// hlog.Printf("info", "No MasterList.LeaderAddr Addr: %s", addr)

	//
	conn, err := rpcsrv.ClientConn(addr)
	if err != nil {
		hlog.Printf("error", "No MasterList.LeaderAddr Found: %s", addr)
		return nil, err
	}

	// host/meta
	if status.Host.Meta == nil {
		status.Host.Meta = &inapi.ObjectMeta{
			Id: status.Host.Meta.Id,
		}
	}

	// host/spec
	if status.Host.Spec == nil {

		status.Host.Spec = &inapi.ResHostSpec{
			PeerLanAddr: string(config.Config.Host.LanAddr),
			PeerWanAddr: string(config.Config.Host.WanAddr),
		}
	}

	//
	if status.Host.Spec.Platform == nil {

		os, arch, _ := inutils.ResSysHostEnvDistArch()

		status.Host.Spec.Platform = &inapi.ResPlatform{
			Os:     os,
			Arch:   arch,
			Kernel: inutils.ResSysHostKernel(),
		}
	}

	//
	if status.Host.Spec.Capacity == nil {

		vm, _ := ps_mem.VirtualMemory()

		status.Host.Spec.Capacity = &inapi.ResHostResource{
			Cpu: uint64(runtime.NumCPU()) * 1000,
			Mem: vm.Total,
		}
	}

	tn := time.Now().Unix()

	if len(status.Host.Status.Volumes) == 0 ||
		tn-hostSyncVolLasted > 600 {

		var (
			devs, _ = ps_disk.Partitions(false)
			vols    = []*inapi.ResHostVolume{}
		)

		sort.Slice(devs, func(i, j int) bool {
			if strings.Compare(devs[i].Device+devs[i].Mountpoint, devs[j].Device+devs[j].Device) < 0 {
				return true
			}
			return false
		})

		ars := types.ArrayString{}
		for _, dev := range devs {

			if ars.Has(dev.Device) {
				continue
			}
			ars.Set(dev.Device)

			if !strings.HasPrefix(dev.Device, "/dev/") ||
				strings.HasPrefix(dev.Mountpoint, "/boot") ||
				strings.Contains(dev.Mountpoint, "/devicemapper/mnt/") {
				continue
			}

			if st, err := ps_disk.Usage(dev.Mountpoint); err == nil {
				vols = append(vols, &inapi.ResHostVolume{
					Name:  dev.Mountpoint,
					Total: st.Total,
					Used:  st.Used,
				})
			}
		}

		if len(vols) > 0 {

			sort.Slice(vols, func(i, j int) bool {
				if strings.Compare(vols[i].Name, vols[j].Name) < 0 {
					return true
				}
				return false
			})

			status.Host.Status.Volumes = vols
		}

		hostSyncVolLasted = tn
	}

	// fmt.Println("pod", len(podrunner.PodRepActives))
	// fmt.Println("box", len(podrunner.BoxActives))

	if stats := hostStatsFeed(); stats != nil {
		status.Host.Status.Stats = stats
		// js, _ := json.Encode(stats, "  ")
		// fmt.Println("send", string(js))
	}

	// Pod Rep Status
	status.Host.Prs = []*inapi.PbPodRepStatus{}
	nstatus.PodRepActives.Each(func(pod *inapi.Pod) {

		if pod.Operate.Replica == nil {
			return
		}

		if inapi.OpActionAllow(pod.Operate.Replica.Action, inapi.OpActionStop|inapi.OpActionStopped) ||
			inapi.OpActionAllow(pod.Operate.Replica.Action, inapi.OpActionDestroy|inapi.OpActionDestroyed) {
			return
		}

		repStatus := &inapi.PbPodRepStatus{
			PodId:  pod.Meta.ID,
			RepId:  pod.Operate.Replica.RepId,
			OpLog:  nstatus.PodRepOpLogs.Get(pod.OpRepKey()),
			Action: 0,
		}

		if repStatus.OpLog == nil {
			hlog.Printf("warn", "No OpLog Found %s", pod.OpRepKey())
			return
		}

		//

		instName := napi.BoxInstanceName(pod.Meta.ID, pod.Operate.Replica)

		boxInst := nstatus.BoxActives.Get(instName)
		if boxInst == nil {
			repStatus.Action = inapi.OpActionPending
		} else {

			repStatus.Action = boxInst.Status.Action
			repStatus.Started = boxInst.Status.Started
			repStatus.Updated = boxInst.Status.Updated
			repStatus.Ports = boxInst.Status.Ports

			if boxInst.Stats != nil {

				//
				var (
					feed            = inapi.NewPbStatsSampleFeed(napi.BoxStatsLogCycle)
					ec_time  uint32 = 0
					ec_value int64  = 0
				)
				for _, name := range statsPodRepNames {
					if ec_time, ec_value = boxInst.Stats.Extract(name, napi.BoxStatsLogCycle, ec_time); ec_value >= 0 {
						feed.SampleSync(name, ec_time, ec_value)
					}
				}

				if len(feed.Items) > 0 {
					repStatus.Stats = feed
				}
			}
		}

		if inapi.OpActionAllow(pod.Operate.Replica.Action, inapi.OpActionDestroy) &&
			inapi.OpActionAllow(repStatus.Action, inapi.OpActionDestroyed) {

			dir := napi.PodVolSysDir(pod.Meta.ID, pod.Operate.Replica.RepId)
			if _, err := os.Stat(dir); err == nil {
				dir2 := napi.PodVolSysDirArch(pod.Meta.ID, pod.Operate.Replica.RepId)
				if err = os.Rename(dir, dir2); err != nil {
					hlog.Printf("error", "pod %s, rep %d, archive vol-sys err %s",
						pod.Meta.ID, pod.Operate.Replica.RepId, err.Error())
					return
				}
			}

		}

		// hlog.Printf("debug", "PodRep %s Phase %s", pod.OpRepKey(), repStatus.Phase)

		// js, _ := json.Encode(repStatus, "  ")
		// fmt.Println(string(js))

		if proj := quotaConfig.Fetch(pod.OpRepKey()); proj != nil {
			repStatus.Volumes = []*inapi.PbVolumeStatus{
				{
					MountPath: "/home/action",
					Used:      proj.Used,
				},
			}
		}

		status.Host.Prs = append(status.Host.Prs, repStatus)

		var box_oplog inapi.PbOpLogSets
		fpath := fmt.Sprintf(napi.AgentBoxStatus, config.Config.PodHomeDir, pod.OpRepKey())
		if err := json.DecodeFile(fpath, &box_oplog); err == nil {
			if box_oplog.Version >= repStatus.OpLog.Version {
				for _, vlog := range box_oplog.Items {
					repStatus.OpLog.LogSetEntry(vlog)
				}
			}
		}

	})

	return inapi.NewApiZoneMasterClient(conn).HostStatusSync(
		context.Background(), &status.Host,
	)
}

var (
	hostStats = inapi.NewPbStatsSampleFeed(napi.BoxStatsSampleCycle)
)

func hostStatsFeed() *inapi.PbStatsSampleFeed {

	timo := uint32(time.Now().Unix())

	// RAM
	vm, _ := ps_mem.VirtualMemory()
	hostStats.SampleSync("ram/us", timo, int64(vm.Used))
	hostStats.SampleSync("ram/cc", timo, int64(vm.Cached))

	// Networks
	nio, _ := ps_net.IOCounters(false)
	if len(nio) > 0 {
		hostStats.SampleSync("net/rs", timo, int64(nio[0].BytesRecv))
		hostStats.SampleSync("net/ws", timo, int64(nio[0].BytesSent))
	}

	// CPU
	cio, _ := ps_cpu.Times(false)
	if len(cio) > 0 {
		hostStats.SampleSync("cpu/sys", timo, int64(cio[0].User*float64(1e9)))
		hostStats.SampleSync("cpu/user", timo, int64(cio[0].System*float64(1e9)))
	}

	// Storage IO
	devs, _ := ps_disk.Partitions(false)
	if dev_name := diskDevName(devs, config.Config.PodHomeDir); dev_name != "" {
		if diom, err := ps_disk.IOCounters(dev_name); err == nil {
			if dio, ok := diom[dev_name]; ok {
				hostStats.SampleSync("fs/sp/rn", timo, int64(dio.ReadCount))
				hostStats.SampleSync("fs/sp/rs", timo, int64(dio.ReadBytes))
				hostStats.SampleSync("fs/sp/wn", timo, int64(dio.WriteCount))
				hostStats.SampleSync("fs/sp/ws", timo, int64(dio.WriteBytes))
			}
		}
	}

	//
	var (
		feed            = inapi.NewPbStatsSampleFeed(napi.BoxStatsLogCycle)
		ec_time  uint32 = 0
		ec_value int64  = 0
	)
	for _, name := range statsHostNames {
		if ec_time, ec_value = hostStats.Extract(name, napi.BoxStatsLogCycle, ec_time); ec_value >= 0 {
			feed.SampleSync(name, ec_time, ec_value)
		}
	}

	if len(feed.Items) > 0 {
		return feed
	}

	return nil
}

func diskDevName(pls []ps_disk.PartitionStat, path string) string {

	path = filepath.Clean(path)

	for {

		for _, v := range pls {
			if path == v.Mountpoint {
				if strings.HasPrefix(v.Device, "/dev/") {
					return v.Device[5:]
				}
				return ""
			}
		}

		if i := strings.LastIndex(path, "/"); i > 0 {
			path = path[:i]
		} else if len(path) > 1 && path[0] == '/' {
			path = "/"
		} else {
			break
		}
	}

	return ""
}
