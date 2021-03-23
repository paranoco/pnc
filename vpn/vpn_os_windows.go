package vpn

import (
	"golang.zx2c4.com/wireguard/device"
	"golang.zx2c4.com/wireguard/tun"
	"github.com/pkg/errors"
	"net"
)

func setupInterfaceTun(interfaceName string) (tun.Device, error) {
	return func() (tun.Device, error) {
		return tun.CreateTUN(interfaceName, device.DefaultMTU)
	}()
}

func uapiListen(interfaceName string, logger *device.Logger) (net.Listener, error) {
	return nil, errors.New("not implemented")
}