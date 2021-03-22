package vpn

import (
	"encoding/json"
	"github.com/hashicorp/go-retryablehttp"
	"github.com/paranoco/pnc/hostid"
	"github.com/paranoco/pnc/oauth2ns"
	"golang.org/x/oauth2"
	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"
	"net/http"
)

func getWireguardSessionConfiguration(conf *VPNLoginConfiguration, host string, jwtToken string) (*WireguardSessionConfiguration, error) {
	client := retryablehttp.NewClient()
	client.RetryMax = 3

	privateKey, err := wgtypes.GeneratePrivateKey()
	if err != nil {
		return nil, err
	}

	publicKey := privateKey.PublicKey()

	ppp := &PeerRequest{
		JWT:             jwtToken,
		HostId:          host,
		WireguardPubKey: publicKey.String(),
	}

	postBody, err := json.Marshal(&ppp)
	if err != nil {
		return nil, err
	}

	req, err := retryablehttp.NewRequest(http.MethodPost, conf.VPNEndpoint+"api/pnc", postBody)
	if err != nil {
		return nil, err
	}

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var peerResponse PeerResponse

	err = json.NewDecoder(res.Body).Decode(&peerResponse)
	if err != nil {
		return nil, err
	}

	return &WireguardSessionConfiguration{
		LocalIP:                 peerResponse.YourIP,
		LocalKey:                privateKey.String(),
		RouterIP:                peerResponse.RouterIP,
		RouterKey:               peerResponse.RouterKey,
		RouterEndpoint:          peerResponse.RouterEndpoint,
		RouterNetworks:          peerResponse.RouterNetworks,
		AuthorizationExpiration: peerResponse.AuthorizationExpiration,
	}, nil
}

func getHostId() (string, error) {
	var h *hostid.HostIdentity

	h, err := hostid.LoadHostIdentity()
	if err != nil {
		h, err = hostid.GenerateIdentity()
		if err != nil {
			return "", err
		}

		err = h.Save()
		if err != nil {
			return "", err
		}
	}

	return h.HostId, nil
}

func getSSOToken(c *VPNLoginConfiguration) (string, error) {
	conf := &oauth2.Config{
		ClientID:     c.SSOClientID,
		ClientSecret: c.SSOClientSecret,
		Scopes:       []string{OPENID_SCOPE},
		Endpoint: oauth2.Endpoint{
			AuthURL:  c.SSOEndpoint + "auth",
			TokenURL: c.SSOEndpoint + "token",
		},
	}

	o := &oauth2ns.Oauth2Cli{}
	client, err := o.AuthenticateUser(conf)
	if err != nil {
		return "", err
	}

	return client.Token.AccessToken, nil
}