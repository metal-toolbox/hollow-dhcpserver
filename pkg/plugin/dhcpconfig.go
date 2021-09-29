package plugin

import "net"

// LeaseData represents the structure of the attribute data of lease information. It can hold multiple V4 and V6 leases.
type LeaseData struct {
	V4Leases []V4Lease `json:"ipv4,omitempty"`
	V6Leases []V6Lease `json:"ipv6,omitempty"`
}

// V4Lease represents an IPv4 DHCP lease
type V4Lease struct {
	MacAddress string   `json:"mac_address"`
	CIDR       string   `json:"cidr"`
	Gateway    net.IP   `json:"gateway"`
	Resolvers  []net.IP `json:"resolvers,omitempty"`
	BootServer string   `json:"boot_server,omitempty"`
	BootFile   string   `json:"boot_file,omitempty"`
}

// V6Lease represents an IPv6 DHCP lease
type V6Lease struct {
	MacAddress string   `json:"mac_address"`
	CIDR       string   `json:"cidr"`
	Gateway    net.IP   `json:"gateway"`
	Resolvers  []string `json:"resolvers,omitempty"`
	BootServer string   `json:"boot_server,omitempty"`
	BootFile   string   `json:"boot_file,omitempty"`
}
