package plugin

import (
	"github.com/urfave/cli/v2"
	"golang.org/x/xerrors"

	"github.com/aquasecurity/trivy/internal/config"
	"github.com/aquasecurity/trivy/pkg/log"
	"github.com/aquasecurity/trivy/pkg/plugin"
)

// Install installs a plugin
func Install(c *cli.Context) error {
	if c.NArg() != 1 {
		cli.ShowSubcommandHelpAndExit(c, 1)
	}

	if err := initLogger(c); err != nil {
		return xerrors.Errorf("initialize error: %w", err)
	}

	url := c.Args().First()
	if _, err := plugin.Install(c.Context, url); err != nil {
		return xerrors.Errorf("plugin install error: %w", err)
	}

	return nil
}

// Run runs a plugin
func Run(c *cli.Context) error {
	if c.NArg() < 1 {
		cli.ShowSubcommandHelpAndExit(c, 1)
	}

	if err := initLogger(c); err != nil {
		return xerrors.Errorf("initialize error: %w", err)
	}

	url := c.Args().First()
	pl, err := plugin.Install(c.Context, url)
	if err != nil {
		return xerrors.Errorf("plugin install error: %w", err)
	}

	if err = pl.Run(c.Context, c.Args().Tail()); err != nil {
		return xerrors.Errorf("unable to run %s plugin: %w", pl.Name, err)
	}
	return nil
}

// LoadCommands loads plugins
func LoadCommands() cli.Commands {
	var commands cli.Commands
	plugins, _ := plugin.LoadAll()
	for _, p := range plugins {
		cmd := &cli.Command{
			Name:  p.Name,
			Usage: p.Usage,
			Action: func(c *cli.Context) error {
				if err := initLogger(c); err != nil {
					return xerrors.Errorf("initialize error: %w", err)
				}

				if err := p.Run(c.Context, c.Args().Slice()); err != nil {
					return xerrors.Errorf("plugin error: %w", err)
				}
				return nil
			},
			SkipFlagParsing: true,
		}
		commands = append(commands, cmd)
	}
	return commands
}

func initLogger(ctx *cli.Context) error {
	conf, err := config.NewGlobalConfig(ctx)
	if err != nil {
		return xerrors.Errorf("config error: %w", err)
	}

	if err = log.InitLogger(conf.Debug, conf.Quiet); err != nil {
		return xerrors.Errorf("failed to initialize a logger: %w", err)
	}
	return nil
}
