package hostid

import (
	"encoding/json"

	"github.com/99designs/keyring"
	"github.com/mitchellh/go-homedir"
	"github.com/pkg/errors"
)

const KEYCHAIN_SERVICE_NAME = "com.paranoco"
const HOST_KEY_ID = "com.paranoco.host"
const CONFIG_KEY = ""

func keyringConfig() keyring.Config {
	pncConfigDir, err := homedir.Expand("~/.pnc/")
	if err != nil {
		panic(err)
	}

	return keyring.Config{
		ServiceName:              KEYCHAIN_SERVICE_NAME,
		KeychainTrustApplication: true,
		KeychainSynchronizable:   false,
		FileDir:                  pncConfigDir,
		FilePasswordFunc: func(prompt string) (string, error) {
			return CONFIG_KEY, nil
		},
	}
}

func (h *HostIdentity) Save() error {
	ring, err := keyring.Open(keyringConfig())
	if err != nil {
		return errors.Wrap(err, "can't open Keyring")
	}

	serializedKey, err := json.Marshal(h)
	if err != nil {
		return errors.Wrap(err, "can't json.Marshal HostIdentity")
	}

	err = ring.Set(keyring.Item{
		Key:                         HOST_KEY_ID,
		Label:                       HOST_KEY_ID,
		Description:                 "Device Private Key",
		Data:                        serializedKey,
		KeychainNotTrustApplication: false,
		KeychainNotSynchronizable:   true,
	})
	if err != nil {
		return errors.Wrap(err, "can't load Paranoco Device Private Key")
	}

	return nil
}

func LoadHostIdentity() (*HostIdentity, error) {
	ring, err := keyring.Open(keyringConfig())
	if err != nil {
		return nil, errors.Wrap(err, "can't open Keyring")
	}

	serializedKey, err := ring.Get(HOST_KEY_ID)
	if err != nil {
		return nil, err
	}

	var result HostIdentity
	err = json.Unmarshal(serializedKey.Data, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}
