// Package plugin implements a coredhcp plugin that retrieves IP lease information
// from the hollow ecosystem and uses it in the DHCP reply. The data will be
// pulled from both serverservice and instanceservice attributes.
//
// Example config:
//
// server6:
//   listen: '[::]547'
//   - example:
//   - server_id: LL aa:bb:cc:dd:ee:ff
//   - hollow: https://hollow.sh
//
// This will send requests to https://hollow.sh/api/<path>.
// OIDC environment is read from the environment.
// The follow environment variables are expected:
// 	HOLLOWDHCP_OIDC_ISSUER
// 	HOLLOWDHCP_OIDC_CLIENT_ID
// 	HOLLOWDHCP_OIDC_CLIENT_SECRET
// 	HOLLOWDHCP_OIDC_AUDIENCE
package plugin

import (
	"fmt"
	"net"
	"net/url"

	"github.com/coredhcp/coredhcp/handler"
	"github.com/coredhcp/coredhcp/logger"
	"github.com/coredhcp/coredhcp/plugins"
	"github.com/insomniacslk/dhcp/dhcpv4"
	"github.com/insomniacslk/dhcp/dhcpv6"
	serverservice "go.hollow.sh/serverservice/pkg/api/v1"
)

var log = logger.GetLogger("plugins/hollow")

// Plugin wraps plugin registration information
var Plugin = plugins.Plugin{
	Name:   "hollow",
	Setup6: setup6,
	Setup4: setup4,
}

var hollowClient *serverservice.Client

func initHollowClient(args ...string) error {
	if hollowClient != nil {
		// already initialized
		return nil
	}

	// oidcIssuer := os.Getenv("HOLLOWDHCP_OIDC_ISSUER")
	// if oidcIssuer == "" {
	// 	return fmt.Errorf("expected HOLLOWDHCP_OIDC_ISSUER to be set")
	// }

	// oidcClientID := os.Getenv("HOLLOWDHCP_OIDC_CLIENT_ID")
	// if oidcClientID == "" {
	// 	return fmt.Errorf("expected HOLLOWDHCP_OIDC_CLIENT_ID to be set")
	// }

	// oidcClientSecret := os.Getenv("HOLLOWDHCP_OIDC_CLIENT_SECRET")
	// if oidcClientSecret == "" {
	// 	return fmt.Errorf("expected HOLLOWDHCP_OIDC_CLIENT_SECRET to be set")
	// }

	// oidcAudience := os.Getenv("HOLLOWDHCP_OIDC_AUDIENCE")
	// if oidcAudience == "" {
	// 	return fmt.Errorf("expected HOLLOWDHCP_OIDC_AUDIENCE to be set")
	// }

	// ctx := context.TODO()

	// TODO: put this back once we can hit our oidc issuer and discovery works
	// provider, err := oidc.NewProvider(ctx, viper.GetString("oidc.issuer"))
	// if err != nil {
	// 	logger.Fatalw("failed to read oidc configuration", "error", err)
	// }

	// oauthConfig := clientcredentials.Config{
	// 	ClientID:     oidcClientID,
	// 	ClientSecret: oidcClientSecret,
	// 	// TokenURL:       provider.Endpoint().TokenURL,
	// 	TokenURL:       oidcIssuer + "oauth2/token",
	// 	Scopes:         []string{"read:server", "read:instance"},
	// 	EndpointParams: url.Values{"audience": []string{oidcAudience}},
	// }

	if len(args) != 1 {
		return fmt.Errorf("got %d arguments, want 1", len(args))
	}
	u, err := url.Parse(args[0])
	if err != nil {
		return fmt.Errorf("invalid URL '%s': %v", args[0], err)
	}

	hollowClient, err = serverservice.NewClientWithToken("token", u.String(), nil)

	return err
}

func setup6(args ...string) (handler.Handler6, error) {
	if err := initHollowClient(args...); err != nil {
		return nil, err
	}
	log.Info("Loaded hollow plugin for DHCPv6.")
	return hollowHandler6, nil
}

func setup4(args ...string) (handler.Handler4, error) {
	if err := initHollowClient(args...); err != nil {
		return nil, err
	}
	log.Info("Loaded hollow plugin for DHCPv4.")
	return hollowHandler4, nil
}

func hollowHandler6(req, resp dhcpv6.DHCPv6) (dhcpv6.DHCPv6, bool) {
	log.Debugf("Received DHCPv6 packet: %s", req.Summary())
	// mac, err := dhcpv6.ExtractMAC(req)
	// if err != nil {
	// 	log.Warningf("Could not find client MAC, dropping request")
	// 	return resp, false
	// }
	// // extract the IA_ID from option IA_NA
	// opt := req.GetOneOption(dhcpv6.OptionIANA)
	// if opt == nil {
	// 	log.Warningf("No option IA_NA found in request, dropping request")
	// 	return resp, false
	// }
	// iaID := opt.(*dhcpv6.OptIANA).IaId
	// log.Debugf("Retrieving IP addresses for MAC %s", mac)
	// ips, err := netBox.GetIPs(mac.String())
	// if err != nil {
	// 	log.Warningf("No IPs found for MAC %s: %v", mac.String(), err)
	// 	return resp, false
	// }
	// for _, addr := range ips {
	// 	if addr.IP.To4() == nil && addr.IP.To16() != nil {
	// 		resp.AddOption(&dhcpv6.OptIANA{
	// 			IaId: iaID,
	// 			Options: []dhcpv6.Option{
	// 				&dhcpv6.OptIAAddress{
	// 					IPv6Addr: addr.IP.To16(),
	// 					// default lifetime, can be overridden by other plugins
	// 					PreferredLifetime: 3600,
	// 					ValidLifetime:     3600,
	// 				},
	// 			},
	// 		})
	// 		break
	// 	}
	// }
	// log.Infof("Resp %s", resp.Summary())
	// return resp, true
	return resp, false
}

func hollowHandler4(req, resp *dhcpv4.DHCPv4) (*dhcpv4.DHCPv4, bool) {
	log.Debugf("Received DHCPv4 packet: %s", req.Summary())
	mac := req.ClientHWAddr.String()

	cfg, hostname, err := getV4Config(mac)
	if err != nil {
		log.Warningf("No IPs found for MAC %s: %v", mac, err)
		return resp, false
	}

	resp.Options.Update(dhcpv4.OptHostName(hostname))

	ipAddr, ipNet, err := net.ParseCIDR(cfg.CIDR)
	if err != nil {
		log.Warningf("MAC %s malformed IP %s error: %s...dropping request", req.ClientHWAddr.String(), cfg.CIDR, err)
		return resp, false
	}

	resp.YourIPAddr = ipAddr
	resp.Options.Update(dhcpv4.OptSubnetMask(ipNet.Mask))
	resp.Options.Update(dhcpv4.OptRouter(cfg.Gateway))

	if req.IsOptionRequested(dhcpv4.OptionDomainNameServer) {
		resp.Options.Update(dhcpv4.OptDNS(cfg.Resolvers...))
	}

	if cfg.BootServer != "" && cfg.BootFile != "" {
		resp.Options.Update(dhcpv4.OptTFTPServerName(cfg.BootServer))
		resp.Options.Update(dhcpv4.OptBootFileName(cfg.BootFile))
	}

	return resp, true
}
