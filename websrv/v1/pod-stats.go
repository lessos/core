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
	"encoding/base64"
	"strings"

	"github.com/hooto/httpsrv"
	"github.com/hooto/iam/iamapi"
	"github.com/hooto/iam/iamclient"
	"github.com/lessos/lessgo/encoding/json"
	"github.com/lessos/lessgo/types"

	"github.com/lessos/loscore/data"
	"github.com/lessos/loscore/losapi"
)

type PodStats struct {
	*httpsrv.Controller
	us iamapi.UserSession
}

func (c *PodStats) Init() int {

	//
	c.us, _ = iamclient.SessionInstance(c.Session)

	if !c.us.IsLogin() {
		c.Response.Out.WriteHeader(401)
		c.RenderJson(types.NewTypeErrorMeta(iamapi.ErrCodeUnauthorized, "Unauthorized"))
		return 1
	}

	return 0
}

func (c PodStats) FeedAction() {

	var (
		pod_id         = c.Params.Get("id")
		qry            = c.Params.Get("qry")
		qry_names      = c.Params.Get("qry_names")
		qry_time_past  = uint32(c.Params.Uint64("qry_time_past"))
		qry_time_cycle = uint32(c.Params.Uint64("qry_time_cycle"))
		pod            losapi.Pod
		fq             losapi.TimeStatsFeedQuerySet
	)

	if len(qry) > 2 {
		bs, err := base64.StdEncoding.DecodeString(qry)
		if err != nil || len(bs) < 10 {
			c.RenderJson(types.NewTypeErrorMeta("400", "Bad Request"))
			return
		}

		if err := json.Decode(bs, &fq); err != nil {
			c.RenderJson(types.NewTypeErrorMeta("400", "Bad Request"))
			return
		}
	} else {
		fq.TimePast = qry_time_past
		fq.TimeCycle = qry_time_cycle
		as := strings.Split(qry_names, ",")
		for _, v := range as {

			if fq.Get(v) != nil {
				continue
			}

			fq.Items = append(fq.Items, &losapi.TimeStatsEntryQuerySet{
				Name: v,
			})
		}
	}

	if len(fq.Items) < 1 {
		c.RenderJson(types.NewTypeErrorMeta("400", "Bad Request"))
		return
	}

	fq.Fix()

	if fq.TimeStart >= fq.TimeCutset {
		c.RenderJson(types.NewTypeErrorMeta("400", "Bad Request"))
		return
	}

	feed := losapi.NewTimeStatsFeed(fq.TimeCycle)
	defer c.RenderJson(feed)

	if obj := data.ZoneMaster.PvGet(losapi.NsGlobalPodInstance(pod_id)); obj.OK() {
		obj.Decode(&pod)
		if pod.Meta.ID == "" || pod.Meta.User != c.us.UserName {
			feed.Error = types.NewErrorMeta("404", "Pod Not Found")
			return
		}
	} else {
		feed.Error = types.NewErrorMeta("404", "Pod Not Found")
		return
	}

	// feed.Debugs.Set("time_start", fmt.Sprintf("%d", fq.TimeStart))
	// feed.Debugs.Set("time_starts", time.Unix(int64(fq.TimeStart), 0))
	// feed.Debugs.Set("time_cut", fmt.Sprintf("%d", fq.TimeCutset))
	// feed.Debugs.Set("time_cuts", time.Unix(int64(fq.TimeCutset), 0))

	if rs := data.HiMaster.ProgScan(
		losapi.NsZonePodRepStats(pod.Spec.Zone,
			pod.Meta.ID,
			0,
			"sys",
			fq.TimeStart,
		),
		losapi.NsZonePodRepStats(pod.Spec.Zone,
			pod.Meta.ID,
			0,
			"sys",
			fq.TimeCutset,
		),
		50000,
	); rs.OK() {

		ls := rs.KvList()
		var vf losapi.TimeStatsFeed
		for _, v := range ls {

			if err := v.Decode(&vf); err != nil {
				continue
			}

			for _, v2 := range vf.Items {
				if fq.Get(v2.Name) == nil {
					continue
				}
				for _, v3 := range v2.Items {
					feed.Sync(v2.Name, v3.Time, v3.Value, "ow")
				}
			}
		}
	}

	for _, v := range feed.Items {

		for i := fq.TimeStart; i <= fq.TimeCutset; i += fq.TimeCycle {
			v.Sync(i, 0, "ex", true)
		}
		v.Sort()

		if len(v.Items) < 2 {
			continue
		}

		fqi := fq.Get(v.Name)
		if fqi == nil {
			continue
		}
		if fqi.Delta {
			for i := len(v.Items) - 1; i > 0; i-- {
				if v.Items[i].Value >= v.Items[i-1].Value && v.Items[i-1].Value > 0 {
					v.Items[i].Value = v.Items[i].Value - v.Items[i-1].Value
				} else {
					v.Items[i].Value = 0
				}
			}
		}

		v.Items = v.Items[1:]
	}

	feed.Kind = "TimeStatsFeed"
}