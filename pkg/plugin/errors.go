package plugin

import (
	"errors"
	"fmt"
)

var (
	// ErrNoLeaseFound is the error returned when a dhcp lease for the mac address can not be found
	ErrNoLeaseFound = errors.New("no dhcp lease found for MAC address")
	// ErrDuplicateLeaseFound is the error returned when multiple dhcp leases for the sane mac addresses are found
	ErrDuplicateLeaseFound = errors.New("multiple dhcp leases found for MAC address")
	// ErrArgumentMissing is the error returned when arguments are missing

)

// ErrMissingEnvVariable is returned when an environment variable isn't provided
type ErrMissingEnvVariable struct {
	EnvVar string
}

// Error returns the ErrMissingEnvVariable in string format
func (e *ErrMissingEnvVariable) Error() string {
	return fmt.Sprintf("expected %s to be set", e.EnvVar)
}
