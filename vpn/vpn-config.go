package vpn

import (
	"encoding/json"
	"fmt"
	"github.com/hashicorp/go-retryablehttp"
	"github.com/shibukawa/configdir"
	"github.com/urfave/cli/v2"
	"net/http"
)

const VENDOR = "com.paranoco"
const APPLICATION_NAME = "pnc"

func GetVpnConfig() (*VPNLoginConfiguration, error) {
	configDirs := configdir.New(VENDOR, APPLICATION_NAME)
	folders := configDirs.QueryFolders(configdir.Global)
	data, err := folders[0].ReadFile("vpn.json")
	if err != nil {
		return &VPNLoginConfiguration{}, err
	}

	var result VPNLoginConfiguration
	err = json.Unmarshal(data, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (c *VPNLoginConfiguration) Save() error {
	configDirs := configdir.New(VENDOR, APPLICATION_NAME)

	data, err := json.Marshal(c)
	if err != nil {
		return err
	}

	folders := configDirs.QueryFolders(configdir.Global)
	err = folders[0].WriteFile("vpn.json", data)
	if err != nil {
		return err
	}

	return nil
}

func FetchRemoteConfig(url string) (*VPNLoginConfiguration, error) {
	client := retryablehttp.NewClient()
	client.RetryMax = 3

	req, err := retryablehttp.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	var response VPNLoginConfiguration

	err = json.NewDecoder(res.Body).Decode(&response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

func VpnConfigCommand(c *cli.Context) error {
	config, _ := GetVpnConfig()

	if c.String("set-vpn") != "" {
		setVPNURL := "https://" + c.String("set-vpn") + "/api/pnc"
		conf, err := FetchRemoteConfig(setVPNURL)
		if err != nil {
			return err
		}
		config = conf
	}

	err := config.Save()
	if err != nil {
		return err
	}

	fmt.Printf("%#v\n", config)
	return nil
}
