package main

import (
	_ "embed"
	"log"
	"os"

	"github.com/mitchellh/colorstring"
	"github.com/paranoco/pnc/vpn"

	"github.com/urfave/cli/v2"
)

var version string = "dev"

//go:embed help_app.template
var appHelpTemplate string

//go:embed help_command.template
var commandHelpTemplate string

//go:embed help_subcommand.template
var subcommandHelpTemplate string

func main() {
	cli.AppHelpTemplate = colorstring.Color(appHelpTemplate)
	cli.CommandHelpTemplate = colorstring.Color(commandHelpTemplate)
	cli.SubcommandHelpTemplate = colorstring.Color(subcommandHelpTemplate)

	app := &cli.App{
		Name:        "pnc",
		Version:     version,
		HideVersion: false,
		Usage:       "paranoco toolbelt",
		Commands: []*cli.Command{
			{
				Name:   "vpn",
				Usage:  "connect to a configured VPN",
				Action: vpn.VpnCommand,
			},
			{
				Name:   "vpn-config",
				Usage:  "configure a VPN to connect to",
				Action: vpn.VpnConfigCommand,
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name: "set-vpn",
					},
				},
			},
			{
				Name:    "public-ip",
				Aliases: []string{"ip"},
				Usage:   "prints your public IPv4",
				Action:  PublicIpCommand,
			},
			{
				Name:   "status",
				Usage:  "prints configuration status",
				Action: vpn.StatusCommand,
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
