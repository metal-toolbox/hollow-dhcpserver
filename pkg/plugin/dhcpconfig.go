package plugin

import "net"

type DHCPConfig struct {
	V4Configs []V4Config `json:"ipv4"`
	V6Configs []V6Config `json:"ipv6"`
}

type V4Config struct {
	MacAddress string   `json:"mac_address"`
	CIDR       string   `json:"cidr"`
	Gateway    net.IP   `json:"gateway"`
	Resolvers  []net.IP `json:"resolvers"`
	BootServer string   `json:"boot_server"`
	BootFile   string   `json:"boot_file"`
}

type V6Config struct {
	MacAddress string   `json:"mac_address"`
	CIDR       string   `json:"cidr"`
	Netmask    string   `json:"netmask"`
	Gateway    string   `json:"gateway"`
	Resolvers  []string `json:"resolvers"`
	BootServer string   `json:"boot_server"`
	BootFile   string   `json:"boot_file"`
}
