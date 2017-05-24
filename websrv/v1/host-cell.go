// Copyright 2015 Authors, All rights reserved.
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

package v1

import (
	"github.com/lessos/lessgo/types"

	"code.hooto.com/lessos/loscore/data"
	"code.hooto.com/lessos/loscore/losapi"
)

func (c Host) CellListAction() {

	var sets losapi.GeneralObjectList
	defer c.RenderJson(&sets)

	zones := []string{}

	if zoneid := c.Params.Get("zoneid"); zoneid != "" {
		if rs := data.ZoneMaster.PvGet(losapi.NsGlobalSysZone(zoneid)); !rs.OK() {
			sets.Error = types.NewErrorMeta("404", "Zone Not Found")
			return
		}
		zones = append(zones, zoneid)
	} else {

		rss := data.ZoneMaster.PvScan(losapi.NsGlobalSysZone(""), "", "", 100).KvList()
		for _, v := range rss {
			var zone losapi.ResZone
			if err := v.Decode(&zone); err == nil {
				zones = append(zones, zone.Meta.Id)
			}
		}
	}

	//
	for _, z := range zones {
		rss := data.ZoneMaster.PvScan(losapi.NsGlobalSysCell(z, ""), "", "", 100).KvList()
		for _, v := range rss {
			var cell losapi.ResCell
			if err := v.Decode(&cell); err == nil {
				sets.Items = append(sets.Items, cell)
			}
		}
	}

	sets.Kind = "HostCellList"
}

func (c Host) CellEntryAction() {

	var set struct {
		losapi.GeneralObject
		losapi.ResCell
	}
	defer c.RenderJson(&set)

	if rs := data.ZoneMaster.PvGet(
		losapi.NsGlobalSysCell(c.Params.Get("zoneid"), c.Params.Get("cellid")),
	); rs.OK() {
		rs.Decode(&set.ResCell)
	}

	if set.Meta == nil || set.Meta.Id == "" {
		set.Error = types.NewErrorMeta("404", "Cell Not Found")
	} else {
		set.Kind = "HostCell"
	}
}
