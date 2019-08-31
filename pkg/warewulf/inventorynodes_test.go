package warewulf

import (
	"testing"

	"github.com/PolarGeospatialCenter/inventory/pkg/inventory/types"
)

func TestGetGatewayWeight(t *testing.T) {
	testCases := []struct {
		name           string
		nic            *types.NICInstance
		expectedWeight int
	}{
		{
			name:           "no gateway",
			nic:            &types.NICInstance{Config: types.NicConfig{Gateway: []string{}}},
			expectedWeight: 0,
		},
		{
			name:           "single stack rfc1918",
			nic:            &types.NICInstance{Config: types.NicConfig{Gateway: []string{"192.168.0.1"}}},
			expectedWeight: 2,
		},
		{
			name:           "single stack ipv4 global unicast",
			nic:            &types.NICInstance{Config: types.NicConfig{Gateway: []string{"198.51.100.1"}}},
			expectedWeight: 4,
		},
		{
			name:           "single stack ipv6 global unicast",
			nic:            &types.NICInstance{Config: types.NicConfig{Gateway: []string{"2001:db8::1"}}},
			expectedWeight: 4,
		},
		{
			name:           "ipv6 link local",
			nic:            &types.NICInstance{Config: types.NicConfig{Gateway: []string{"fe80::1"}}},
			expectedWeight: 1,
		},
		{
			name:           "dual stack rfc1918 and link local",
			nic:            &types.NICInstance{Config: types.NicConfig{Gateway: []string{"192.168.0.1", "fe80::1"}}},
			expectedWeight: 3,
		},
		{
			name:           "dual stack ipv4 global unicast and link local",
			nic:            &types.NICInstance{Config: types.NicConfig{Gateway: []string{"198.51.100.1", "fe80::1"}}},
			expectedWeight: 5,
		},
		{
			name:           "dual stack ipv4 global unicast and ula",
			nic:            &types.NICInstance{Config: types.NicConfig{Gateway: []string{"198.51.100.1", "fd00:9bab:aa1b:735b::1"}}},
			expectedWeight: 6,
		},
		{
			name:           "dual stack rfc1918 and ula",
			nic:            &types.NICInstance{Config: types.NicConfig{Gateway: []string{"192.168.0.1", "fd00:9bab:aa1b:735b::1"}}},
			expectedWeight: 4,
		},
		{
			name:           "dual stack global unicast",
			nic:            &types.NICInstance{Config: types.NicConfig{Gateway: []string{"198.51.100.1", "2001:db8::1"}}},
			expectedWeight: 8,
		},
	}

	for _, c := range testCases {
		t.Run(c.name, func(st *testing.T) {
			if actualWeight := getGatewayWeight(c.nic); actualWeight != c.expectedWeight {
				st.Errorf("Expected weight of %d.  Got %d.", c.expectedWeight, actualWeight)
			}
		})
	}
}

func TestChooseDefaultGatewayNetwork(t *testing.T) {
	testCases := []struct {
		name            string
		networks        map[string]*types.NICInstance
		expectedNetwork string
	}{
		{
			name: "single nic",
			networks: map[string]*types.NICInstance{
				"only": &types.NICInstance{Config: types.NicConfig{Gateway: []string{"198.51.100.1"}}},
			},
			expectedNetwork: "only",
		},
		{
			name: "dual nic, one public",
			networks: map[string]*types.NICInstance{
				"public":  &types.NICInstance{Config: types.NicConfig{Gateway: []string{"198.51.100.1"}}},
				"private": &types.NICInstance{Config: types.NicConfig{Gateway: []string{"192.168.0.1"}}},
			},
			expectedNetwork: "public",
		},
		{
			name: "dual nic, one dualstack",
			networks: map[string]*types.NICInstance{
				"singlestack": &types.NICInstance{Config: types.NicConfig{Gateway: []string{"198.51.100.1"}}},
				"dualstack":   &types.NICInstance{Config: types.NicConfig{Gateway: []string{"198.51.100.1", "2001:db8::1"}}},
			},
			expectedNetwork: "dualstack",
		},
	}
	for _, c := range testCases {
		t.Run(c.name, func(st *testing.T) {
			if actualNetwork := chooseDefaultGatewayNetwork(c.networks); actualNetwork != c.expectedNetwork {
				st.Errorf("Expected network '%s' to be chosen, got '%s'", c.expectedNetwork, actualNetwork)
			}
		})
	}
}
