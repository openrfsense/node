package system

import (
	"github.com/Wifx/gonetworkmanager"
	"github.com/openrfsense/common/stats"
)

// Type StatsNetwork contains network device information by using NetworkManager's DBus interface.
type StatsNetwork struct {
	// List of realized network devices as returned by gonetworkmanager.NetworkManager.GetDevices()
	Devices []gonetworkmanager.Device `json:"devices"`

	// General connection status as reported by gonetworkmanager.NetworkManager.State()
	Status string `json:"status"`
}

// providerNetwork implements stats.Provider.
var _ stats.Provider = providerNetwork{}

// Stats provider for network information.
type providerNetwork struct{}

func (pn providerNetwork) Name() string {
	return "network"
}

func (pn providerNetwork) Stats() (interface{}, error) {
	nm, err := gonetworkmanager.NewNetworkManager()
	if err != nil {
		return nil, err
	}

	devices, err := nm.GetDevices()
	if err != nil {
		return nil, err
	}

	state, err := nm.State()
	if err != nil {
		return nil, err
	}

	return StatsNetwork{
		Devices: devices,
		Status:  state.String(),
	}, nil
}
