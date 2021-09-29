package plugin

import "errors"

var (
	// ErrNoLeaseFound is the error returned when a dhcp lease for the mac address can not be found
	ErrNoLeaseFound = errors.New("no dhcp lease found for MAC address")
	// ErrDuplicateLeaseFound is the error returned when multiple dhcp leases for the sane mac addresses are found
	ErrDuplicateLeaseFound = errors.New("multiple dhcp leases found for MAC address")
)
