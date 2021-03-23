package main

import (
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
