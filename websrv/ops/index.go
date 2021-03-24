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

package ops

import (
	"fmt"
	"html"

	"github.com/hooto/httpsrv"
	"github.com/hooto/iam/iamclient"

	incfg "github.com/sysinner/incore/config"
	status "github.com/sysinner/incore/status"
	"github.com/sysinner/inpanel"
)

type Index struct {
	*httpsrv.Controller
}

func (c Index) IndexAction() {

	c.AutoRender = false
	c.Response.Out.Header().Set("Cache-Control", "no-cache")

	login := "false"
	if !iamclient.SessionIsLogin(c.Session) {
		login = "true"
	}

	cfgLogo, ok := status.ZoneSysConfigGroupList.Value("innerstack/sys/webui",
		"cp_navbar_logo")
	if !ok || len(cfgLogo) < 1 {
		cfgLogo = "/in/~/in/cp/img/logo-10x7-light-h128.png"
	}
	cfgLogo += "?v=" + inpanel.VersionHash

	cfgTitle, ok := status.ZoneSysConfigGroupList.Value("innerstack/sys/webui",
		"html_head_title")
	if !ok {
		cfgTitle = "InnerStack"
	} else {
		cfgTitle = html.EscapeString(cfgTitle)
	}

	siteUrl := "https://www.sysinner.com"
	if incfg.Config.ZoneMain != nil && incfg.Config.ZoneMain.LocaleLang == "zh-CN" {
		siteUrl = "https://www.sysinner.cn"
	}

	c.RenderString(`<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="utf-8">
  <title>` + cfgTitle + `</title>
  <script src="/in/~/valueui/main.js?v=` + inpanel.VersionHash + `"></script>
  <link rel="stylesheet" href="/in/~/in/cp/css/base.css?v=` + inpanel.VersionHash + `" type="text/css">
  <link rel="shortcut icon" type="image/x-icon" href="/in/~/in/cp/img/logo-1x1-light.ico?v=` + inpanel.VersionHash + `">
  <script type="text/javascript">
    valueui.app_version = "` + inpanel.VersionHash + `";
    valueui.basepath = "/in/~/";
    window.onload = valueui.use("in/ops/js/main.js", function() {
      inOps.version = "` + inpanel.VersionHash + `";
      inOps.zone_id = "` + inpanel.ZoneId + `";
      inOps.nav_cluster_zone = ` + fmt.Sprintf("%v", inpanel.OpsClusterZone) + `;
      inOps.nav_cluster_cell = ` + fmt.Sprintf("%v", inpanel.OpsClusterCell) + `;
      inOps.nav_cluster_host = ` + fmt.Sprintf("%v", inpanel.OpsClusterHost) + `;
      inOps.Boot(` + login + `);
    });
  </script>
</head>
<body id="valueui-body">
<div class="incp-well" id="incp-well">
<div class="incp-well-box">
  <div class="incp-well-panel">
    <div class="body2c">
      <div class="body2c1">
        <img src="` + cfgLogo + `" width="64px">
      </div>
      <div class="body2c2">
        <div>InnerStack<br/>Operations Management</div>
      </div>
    </div>
    <div class="status status_dark" id="incp-well-status">loading</div>
  </div>
  <div class="footer">
    <span class="url-info">Powered by <a href="` + siteUrl + `" target="_blank">InnerStack</a></span>
  </div>
</div>
</div>
</body>
</html>
`)
}
