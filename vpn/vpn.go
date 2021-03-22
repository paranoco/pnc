package vpn

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/paranoco/pnc/ostools"
	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"
	"golang.zx2c4.com/wireguard/device"
)

const (
	ExitSetupSuccess = 0
	ExitSetupFailed  = 1
)

const (
	ENV_WG_TUN_FD  = "WG_TUN_FD"
	ENV_WG_UAPI_FD = "WG_UAPI_FD"
)

type ConfigurationFunc func(interfaceName string) error

// Creates a wireguard interface, calls the configurationFunc, and then
// forever services that interface. This function only returns upon error
// or upon Ctrl-C.
func runWireguard(configurationFunc ConfigurationFunc) error {
	var interfaceName string = "utun"

	// get log level (default: info)

	logLevel := func() int {
		switch os.Getenv("LOG_LEVEL") {
		case "verbose", "debug":
			return device.LogLevelDebug
		case "error":
			return device.LogLevelError
		case "silent":
			return device.LogLevelSilent
		}
		return device.LogLevelInfo
	}()

	// open TUN device (or use supplied fd)

	tun, err := setupInterfaceTun(interfaceName)

	if err == nil {
		realInterfaceName, err2 := tun.Name()
		if err2 == nil {
			interfaceName = realInterfaceName
		}
	}

	logger := device.NewLogger(
		logLevel,
		fmt.Sprintf("(%s) ", interfaceName),
	)

	if err != nil {
		logger.Debug.Printf("Failed to create TUN device: %v", err)
		os.Exit(ExitSetupFailed)
	}

	device := device.NewDevice(tun, logger)

	logger.Debug.Printf("Device %s started", interfaceName)

	errs := make(chan error)
	term := make(chan os.Signal, 1)

	uapi, err := uapiListen(interfaceName, logger)
	if err != nil {
		logger.Debug.Printf("Failed to listen on uapi socket: %v", err)
		os.Exit(ExitSetupFailed)
	}

	go func() {
		for {
			conn, err := uapi.Accept()
			if err != nil {
				errs <- err
				return
			}
			go device.IpcHandle(conn)
		}
	}()

	logger.Debug.Printf("UAPI listener started")

	configurationFunc(interfaceName)

	logger.Debug.Printf("Interface configured")

	// wait for program to terminate
	signal.Notify(term, syscall.SIGTERM)
	signal.Notify(term, os.Interrupt)

	select {
	case <-term:
	case <-errs:
	case <-device.Wait():
	}

	// clean up

	uapi.Close()
	device.Close()

	logger.Debug.Printf("Shutting down")
	return nil
}

func VpnCommand(c *cli.Context) error {
	conf, err := GetVpnConfig()
	if err != nil {
		return errors.Wrap(err, "you must set a VPN configuration with pnc vpn-config")
	}

	err = ostools.EnsureAdministratorRights()
	if err != nil {
		return err
	}

	// get host id
	host, err := getHostId()
	if err != nil {
		return err
	}

	// get sso id
	token, err := getSSOToken(conf)
	if err != nil {
		return err
	}

	sessionConf, err := getWireguardSessionConfiguration(conf, host, token)
	if err != nil {
		return err
	}

	return runWireguard(func(interfaceName string) error {
		err := setInterfaceIp(interfaceName, sessionConf.LocalIP, sessionConf.RouterIP)
		if err != nil {
			return err
		}

		err = configureDevice2(interfaceName, sessionConf.RouterKey, sessionConf.LocalKey, sessionConf.RouterEndpoint, sessionConf.RouterNetworks)
		if err != nil {
			return err
		}

		for _, network := range sessionConf.RouterNetworks {
			addRouteTo(network, interfaceName)
		}

		return nil
	})
}
