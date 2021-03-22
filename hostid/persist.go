package hostid

import (
	"errors"
	"fmt"
	"github.com/99designs/keyring"
	"github.com/paranoco/pnc/conflang"
	"log"
)

func (h *HostIdentity) Save() error {
	ring, err := keyring.Open(keyring.Config{
		ServiceName: "example",
	})
	if err != nil {
		return err
	}

	err = ring.Set(keyring.Item{
		Key: "foo",
		Data: []byte("secret-bar"),
	})
	if err != nil {
		return err
	}

	_, err = ring.Get("foo")
	if err != nil {
		return err
	}

	return nil
}

func Load() (*HostIdentity, error) {
	return nil, errors.New("not implemented")
}


func ToDo() {
	fmt.Printf("data\n\n")

	ring, _ := keyring.Open(keyring.Config{
		ServiceName:                    "com.paranoco.pnc23",
		KeychainTrustApplication:       true,
		KeychainSynchronizable:         true,
		KeychainAccessibleWhenUnlocked: true,
	})

	_ = ring.Set(keyring.Item{
		Key:         "foo",
		Data:        []byte("secret-bar"),
		Label:       "human friendly name",
		Description: "kind",
	})

	i, _ := ring.Get("foo")

	fmt.Printf("data: %s\n\n\n", i.Data)

	v, d := conflang.ParseConfig("conflang/schema.hcl", "conflang/test.hcl")
	if d.HasErrors() {
		log.Fatal(d)
	}
	fmt.Printf("%#v\n", v)

}