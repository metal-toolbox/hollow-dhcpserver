package plugin

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	serverservice "go.hollow.sh/serverservice/pkg/api/v1"
)

// DHCPAttributeNamespace is the namespace that is searched to find lease information
const DHCPAttributeNamespace string = "sh.hollow.dhcpserver.lease"

func getV4Config(mac string) (*V4Config, string, error) {
	ctx := context.TODO()

	mac = strings.ToLower(mac)

	params := &serverservice.ServerListParams{
		AttributeListParams: []serverservice.AttributeListParams{
			{
				Namespace: DHCPAttributeNamespace,
				Keys:      []string{"ipv4"},
				Operator:  serverservice.OperatorLike,
				Value:     fmt.Sprintf(`%%%s%%`, mac),
			},
		},
	}

	r, _, err := hollowClient.List(ctx, params)
	if err != nil {
		return nil, "", err
	}

	if len(r) == 0 {
		return nil, "", ErrNoLeaseFound
	}

	if len(r) != 1 {
		return nil, "", ErrDuplicateLeaseFound
	}

	srv := r[0]
	hostname := srv.Name

	var cfg *DHCPConfig

	for _, attr := range srv.Attributes {
		if attr.Namespace != DHCPAttributeNamespace {
			continue
		}

		jsonData, err := attr.Data.MarshalJSON()
		if err != nil {
			return nil, "", err
		}

		if err := json.Unmarshal(jsonData, &cfg); err != nil {
			return nil, "", err
		}

		break
	}

	for _, v4Cfg := range cfg.V4Configs {
		if v4Cfg.MacAddress == mac {
			return &v4Cfg, hostname, nil
		}
	}

	// if we made it here we didn't find a lease that matches the MAC address
	return nil, "", ErrNoLeaseFound
}
