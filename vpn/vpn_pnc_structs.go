package vpn

import "time"

const OPENID_SCOPE = "email"

type VPNLoginConfiguration struct {
	SSOClientID     string
	SSOClientSecret string
	SSOEndpoint     string
	VPNEndpoint     string
}

type PeerResponse struct {
	YourIP                  string
	RouterIP                string
	RouterKey               string
	RouterEndpoint          string
	RouterNetworks          []string
	AuthorizationExpiration time.Time
}

type WireguardSessionConfiguration struct {
	LocalIP                 string
	LocalKey                string
	RouterIP                string
	RouterKey               string
	RouterEndpoint          string
	RouterNetworks          []string
	AuthorizationExpiration time.Time
}

type PeerRequest struct {
	JWT             string
	HostId          string
	WireguardPubKey string
}