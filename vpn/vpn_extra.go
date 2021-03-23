package vpn

import (
	"fmt"
	"golang.zx2c4.com/wireguard/wgctrl"
	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"
	"log"
	"net"
	"os/exec"
	"strings"
	"time"
)

const IFCONFIG = "/sbin/ifconfig"
const ROUTE = "/sbin/route"

func stats(device string) {
	c, err := wgctrl.New()
	if err != nil {
		log.Fatalf("failed to open wgctrl: %v", err)
	}
	defer c.Close()

	d, err := c.Device(device)
	if err != nil {
		log.Fatalf("failed to get device %q: %v", device, err)
	}

	for _, p := range d.Peers {
		printPeer(p)
	}
}

func configureDevice2(device string, serverKeyStr string, myKeyStr string, endpoint string, networks []string) error {
	c, err := wgctrl.New()
	if err != nil {
		return err
	}
	defer c.Close()

	d, err := c.Device(device)
	if err != nil {
		return err
	}

	serverKey, err := wgtypes.ParseKey(serverKeyStr)
	if err != nil {
		return err
	}

	myKey, err := wgtypes.ParseKey(myKeyStr)
	if err != nil {
		return err
	}

	var allowedIps []net.IPNet
	for _, network := range networks {
		_, allowedIp, err := net.ParseCIDR(network)
		if err != nil {
			return err
		}
		allowedIps = append(allowedIps, *allowedIp)
	}

	endpointAddr, err := net.ResolveUDPAddr("udp", endpoint)
	if err != nil {
		return err
	}

	cfg := wgtypes.Config{
		PrivateKey: &myKey,
		// address
		ReplacePeers: true,
		Peers: []wgtypes.PeerConfig{
			wgtypes.PeerConfig{
				PublicKey: serverKey,
				Endpoint:  endpointAddr,
				// PersistentKeepaliveInterval: time.Second * 60,
				ReplaceAllowedIPs: true,
				AllowedIPs:        allowedIps,
			},
		},
	}

	err = c.ConfigureDevice(d.Name, cfg)
	if err != nil {
		return err
	}

	go func() {
		for {
			stats(d.Name)
			time.Sleep(10 * time.Second)
		}
	}()

	return nil
}

func printPeer(p wgtypes.Peer) {
	const f = `rcvd: %d B ## sent: %d B ## last handshake: %s
`

	fmt.Printf(
		f,
		p.ReceiveBytes,
		p.TransmitBytes,
		p.LastHandshakeTime.String(),
	)
}

func ipsString(ipns []net.IPNet) string {
	ss := make([]string, 0, len(ipns))
	for _, ipn := range ipns {
		ss = append(ss, ipn.String())
	}

	return strings.Join(ss, ", ")
}

func setInterfaceIp(iface string, ipaddr string, router string) error {
	cmd := exec.Command(IFCONFIG, iface, "inet", ipaddr, router)

	_, err := cmd.Output()
	if err != nil {
		return err
	}

	return nil
}

func addRouteTo(net string, iface string) error {
	cmd := exec.Command(ROUTE, "add", "-net", net, "-interface", iface)

	_, err := cmd.Output()
	if err != nil {
		return err
	}

	return nil
}
