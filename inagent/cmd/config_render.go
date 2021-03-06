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

package cmd

import (
	"errors"

	"github.com/sysinner/incore/inapi"
	"github.com/sysinner/incore/inconf"
	"github.com/sysinner/incore/inutils/filerender"
)

type configRenderCommand struct {
	cmd  *inapi.BaseCommand
	args struct {
		AppSpec string
		Input   string
		Output  string
	}
}

func NewConfigRenderCommand() *inapi.BaseCommand {

	c := &configRenderCommand{
		cmd: &inapi.BaseCommand{
			Use:     "config-render",
			Aliases: []string{"confrender"},
			Short:   "read input file and render with config data, then write to output file",
		},
	}

	c.cmd.FParseErrWhitelist.UnknownFlags = true

	c.cmd.Flags().StringVar(&c.args.AppSpec, "app-spec",
		"",
		`app-spec id`,
	)

	c.cmd.Flags().StringVar(&c.args.Input, "in",
		"",
		`input file path (template of text, json, toml, yaml, ini)`,
	)

	c.cmd.Flags().StringVar(&c.args.Output, "out",
		"",
		`output file path`,
	)

	c.cmd.RunE = c.run

	return c.cmd
}

func (it *configRenderCommand) run(cmd *inapi.BaseCommand, args []string) error {

	if err := podSetup(); err != nil {
		return err
	}

	appCfr, err := appSetup(it.args.AppSpec)
	if err != nil {
		return err
	}

	if it.args.Input == "" {
		return errors.New("--in input file not setup")
	}

	if it.args.Output == "" {
		return errors.New("--out output file not setup")
	}

	if err := cfgRender(appCfr, it.args.Input, it.args.Output); err != nil {
		return err
	}

	return nil
}

func cfgRender(appCfr *inconf.AppConfigurator, src, dst string) error {

	sets := varParams(appCfr)

	return filerender.Render(src, dst, 0644, sets)
}
