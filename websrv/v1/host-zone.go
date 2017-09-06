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

	"github.com/lessos/loscore/data"
	"github.com/lessos/loscore/losapi"
)

func (c Host) ZoneListAction() {

	ls := losapi.GeneralObjectList{}
	defer c.RenderJson(&ls)

	//
	rs := data.ZoneMaster.PvScan(losapi.NsGlobalSysZone(""), "", "", 100).KvList()

	for _, v := range rs {

		var zone losapi.ResZone

		if err := v.Decode(&zone); err != nil || zone.Meta.Id == "" {
			continue
		}

		if c.Params.Get("fields") == "cells" {

			rs2 := data.ZoneMaster.PvScan(losapi.NsGlobalSysCell(zone.Meta.Id, ""), "", "", 100).KvList()

			for _, v2 := range rs2 {

				var cell losapi.ResCell

				if err := v2.Decode(&cell); err == nil {
					zone.Cells = append(zone.Cells, &cell)
				}
			}
		}

		ls.Items = append(ls.Items, zone)
	}

	ls.Kind = "HostZoneList"
}

func (c Host) ZoneEntryAction() {

	var set struct {
		losapi.GeneralObject
		losapi.ResZone
	}

	defer c.RenderJson(&set)

	if obj := data.ZoneMaster.PvGet(losapi.NsGlobalSysZone(c.Params.Get("id"))); obj.OK() {

		if err := obj.Decode(&set.ResZone); err != nil {
			set.Error = types.NewErrorMeta("400", err.Error())
		} else {

			if c.Params.Get("fields") == "cells" {

				rs2 := data.ZoneMaster.PvScan(losapi.NsGlobalSysCell(set.Meta.Id, ""), "", "", 100).KvList()

				for _, v2 := range rs2 {

					var cell losapi.ResCell

					if err := v2.Decode(&cell); err == nil {
						set.Cells = append(set.Cells, &cell)
					}
				}
			}

		}
	}

	if set.Meta.Id != "" {
		set.Kind = "HostZone"
	} else {
		set.Error = types.NewErrorMeta("404", "Item Not Found")
	}
}
