// Copyright 2019 Eryx <evorui аt gmаil dοt cοm>, All rights reserved.
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

package cmd

import (
	"errors"
	"fmt"

	"github.com/lessos/lessgo/encoding/json"
	"github.com/lessos/lessgo/net/httpclient"
	"github.com/lessos/lessgo/types"

	"github.com/sysinner/incore/inapi"
)

var (
	healthActions = types.ArrayString([]string{"active", "setup"})
)

type healthSyncCommand struct {
	cmd  *inapi.BaseCommand
	args struct {
		Action string
	}
}

func NewHealthSyncCommand() *inapi.BaseCommand {

	c := &healthSyncCommand{
		cmd: &inapi.BaseCommand{
			Use:   "health-sync",
			Short: "synchronize heartbeat information to the central zone main nodes",
		},
	}

	c.cmd.FParseErrWhitelist.UnknownFlags = true

	c.cmd.Flags().StringVar(&c.args.Action, "action",
		"active",
		`action value (ex: active, setup)`,
	)

	c.cmd.RunE = c.run

	return c.cmd
}

func (it *healthSyncCommand) run(cmd *inapi.BaseCommand, args []string) error {

	if it.args.Action == "" {
		return errors.New("no --action=value found")
	}

	if !healthActions.Has(it.args.Action) {
		return errors.New("invalid --action=value")
	}

	hc := httpclient.Get("http://unix/in/v1/health/sync?action=" + it.args.Action)
	defer hc.Close()

	hc.SetUnixDomainSocket("/home/action/.sysinner/inagent.sock")

	if msg, err := hc.ReplyString(); err != nil {
		return err
	} else if msg != "OK" {
		return errors.New("error " + msg)
	}

	return nil
}

type healthStatusCommand struct {
	cmd *inapi.BaseCommand
}

func NewHealthStatusCommand() *inapi.BaseCommand {

	c := &healthStatusCommand{
		cmd: &inapi.BaseCommand{
			Use:   "health-status",
			Short: "output app status",
		},
	}

	c.cmd.FParseErrWhitelist.UnknownFlags = true

	c.cmd.RunE = c.run

	return c.cmd
}

func (it *healthStatusCommand) run(cmd *inapi.BaseCommand, args []string) error {

	hc := httpclient.Get("http://unix/in/v1/health/status")
	defer hc.Close()

	hc.SetUnixDomainSocket("/home/action/.sysinner/inagent.sock")

	var rsp inapi.HealthStatus
	if err := hc.ReplyJson(&rsp); err != nil {
		return err
	}

	js, _ := json.Encode(rsp, "  ")
	fmt.Printf("\nStatus:\n%s\n\n", string(js))

	return nil
}
