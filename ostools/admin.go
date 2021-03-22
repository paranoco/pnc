package ostools

import (
	"github.com/pkg/errors"
	"os"
)

func IsAdmin() bool {
	return os.Geteuid() == 0
}

func EnsureAdministratorRights() error {
	if !IsAdmin() {
		return errors.New("must run as root")
	} else {
		return nil
	}
}