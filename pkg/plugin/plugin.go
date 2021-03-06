package plugin

import (
	"context"
	"net"
	"net/url"
	"os"
	"time"

	"github.com/coredhcp/coredhcp/handler"
	"github.com/coredhcp/coredhcp/logger"
	"github.com/coredhcp/coredhcp/plugins"
	"github.com/friendsofgo/errors"
	"github.com/insomniacslk/dhcp/dhcpv4"
	"github.com/insomniacslk/dhcp/dhcpv6"
	serverservice "go.hollow.sh/serverservice/pkg/api/v1"
	"golang.org/x/oauth2/clientcredentials"
)

const defaultLeaseDuration time.Duration = 3600 * time.Second

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

	if len(args) != 1 {
		return &ErrInvalidArgCount{Expected: 1, Provided: len(args)}
	}

	uri, err := url.Parse(args[0])
	if err != nil {
		return errors.Wrapf(err, "invalid URL '%s'", args[0])
	}

	if os.Getenv("HOLLOWDHCP_AUTH_TOKEN") != "" {
		hollowClient, err = serverservice.NewClientWithToken(os.Getenv("HOLLOWDHCP_AUTH_TOKEN"), uri.String(), nil)

		return err
	}

	oidcIssuer := os.Getenv("HOLLOWDHCP_OIDC_ISSUER")
	if oidcIssuer == "" {
		return &ErrMissingEnvVariable{EnvVar: "HOLLOWDHCP_OIDC_ISSUER"}
	}

	oidcClientID := os.Getenv("HOLLOWDHCP_OIDC_CLIENT_ID")
	if oidcClientID == "" {
		return &ErrMissingEnvVariable{EnvVar: "HOLLOWDHCP_OIDC_CLIENT_ID"}
	}

	oidcClientSecret := os.Getenv("HOLLOWDHCP_OIDC_CLIENT_SECRET")
	if oidcClientSecret == "" {
		return &ErrMissingEnvVariable{EnvVar: "HOLLOWDHCP_OIDC_CLIENT_SECRET"}
	}

	oidcAudience := os.Getenv("HOLLOWDHCP_OIDC_AUDIENCE")
	if oidcAudience == "" {
		return &ErrMissingEnvVariable{EnvVar: "HOLLOWDHCP_OIDC_AUDIENCE"}
	}

	ctx := context.TODO()

	// TODO: put this back once we can hit our oidc issuer and discovery works
	// provider, err := oidc.NewProvider(ctx, viper.GetString("oidc.issuer"))
	// if err != nil {
	// 	logger.Fatalw("failed to read oidc configuration", "error", err)
	// }

	oauthConfig := clientcredentials.Config{
		ClientID:     oidcClientID,
		ClientSecret: oidcClientSecret,
		// TokenURL:       provider.Endpoint().TokenURL,
		TokenURL:       oidcIssuer + "oauth2/token",
		Scopes:         []string{"read:server", "read:instance"},
		EndpointParams: url.Values{"audience": []string{oidcAudience}},
	}

	hollowClient, err = serverservice.NewClient(uri.String(), oauthConfig.Client(ctx))

	return err
}

func setup6(args ...string) (handler.Handler6, error) {
	if err := initHollowClient(args...); err != nil {
		return nil, err
	}

	log.Info("loaded hollow plugin for DHCPv6")

	return hollowHandler6, nil
}

func setup4(args ...string) (handler.Handler4, error) {
	if err := initHollowClient(args...); err != nil {
		return nil, err
	}

	log.Info("loaded hollow plugin for DHCPv4")

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

	cfg, hostname, err := getV4Lease(mac)
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

	// default lifetime, can be overridden by other plugins
	resp.Options.Update(dhcpv4.OptIPAddressLeaseTime(defaultLeaseDuration))

	if req.IsOptionRequested(dhcpv4.OptionDomainNameServer) && len(cfg.Resolvers) != 0 {
		resp.Options.Update(dhcpv4.OptDNS(cfg.Resolvers...))
	}

	if cfg.BootServer != "" && cfg.BootFile != "" {
		resp.Options.Update(dhcpv4.OptTFTPServerName(cfg.BootServer))
		resp.Options.Update(dhcpv4.OptBootFileName(cfg.BootFile))
	}

	return resp, true
}
