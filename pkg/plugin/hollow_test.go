package plugin

import (
	"encoding/json"
	"fmt"
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	serverservice "go.hollow.sh/serverservice/pkg/api/v1"
)

func TestGetV4Lease(t *testing.T) {
	srvCfg := LeaseData{
		V4Leases: []V4Lease{
			{
				MacAddress: "happy:mac",
				CIDR:       "10.1.2.10/24",
				Gateway:    net.ParseIP("10.1.2.1"),
				Resolvers:  []net.IP{net.ParseIP("1.1.1.1")},
			},
		},
	}

	jsonSrvCfg, err := json.Marshal(srvCfg)
	require.NoError(t, err)

	srv := &serverservice.Server{
		Name: "testServer",
		Attributes: []serverservice.Attributes{
			{
				Namespace: DHCPAttributeNamespace,
				Data:      json.RawMessage([]byte(jsonSrvCfg)),
			},
		},
	}

	jsonSrv, err := json.Marshal(srv)
	require.NoError(t, err)

	var testCases = []struct {
		testName         string
		macAddress       string
		responseCode     int
		responseBody     string
		expectedError    error
		expectedHostname string
	}{
		{
			"request fails",
			"MAC",
			401,
			`{"message": "error"}`,
			serverservice.ServerError{Message: "error", StatusCode: 401},
			"",
		},
		{
			"err: no leases found",
			"UNKNOWN:MAC",
			200,
			`{"records": []}`,
			ErrNoLeaseFound,
			"",
		},
		{
			"err: duplicate leases found",
			"DUPLICATE:MAC",
			200,
			`{"records": [{}, {}]}`,
			ErrDuplicateLeaseFound,
			"",
		},
		{
			"err: server found but somehow mac not found",
			"not:happy:MAC",
			200,
			fmt.Sprintf(`{"records": [%s]}`, jsonSrv),
			ErrNoLeaseFound,
			"",
		},
		{
			"happy path: server found",
			"happy:MAC",
			200,
			fmt.Sprintf(`{"records": [%s]}`, jsonSrv),
			nil,
			"testServer",
		},
	}

	for _, tt := range testCases {
		t.Run(tt.testName, func(t *testing.T) {
			jsonResponse := json.RawMessage([]byte(tt.responseBody))
			fmt.Printf("Response Body:\n%s\n", tt.responseBody)
			hollowClient = mockServerServiceClient(string(jsonResponse), tt.responseCode)

			cfg, hostname, err := getV4Lease(tt.macAddress)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.ErrorIs(t, err, tt.expectedError)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, cfg)
				assert.NotEmpty(t, hostname)
				assert.Equal(t, tt.expectedHostname, hostname)
			}
		})
	}
}
