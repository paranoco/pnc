package main

import (
	"github.com/paranoco/pnc/vpn"
	"log"
	"os"

	"github.com/urfave/cli/v2"
)

var version string = "dev"

func main() {
	app := &cli.App{
		Name:    "pnc",
		Version: version,
		Usage:   "paranoco toolbelt",
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
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
