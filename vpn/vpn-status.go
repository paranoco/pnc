package vpn

import (
	"fmt"

	"github.com/99designs/keyring"
	"github.com/paranoco/alexandria/ostools"
	"github.com/urfave/cli/v2"
)

func keyringCapabilities() string {
	var r string
	for _, b := range keyring.AvailableBackends() {
		r = r + string(b) + ","
	}
	return r[0 : len(r)-1]
}

func StatusCommand(c *cli.Context) error {
	ostools.EnsureAdministratorRights()

	config, _ := GetVpnConfig()
	host, _ := getHostId()

	fmt.Printf("keyring_capabilities = %#v\n", keyringCapabilities())
	fmt.Printf("vpn_config = %#v\n", config)
	fmt.Printf("host = %#v\n", host)

	return nil
}
