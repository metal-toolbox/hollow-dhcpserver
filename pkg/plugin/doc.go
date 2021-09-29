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
