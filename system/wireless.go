package system

import (
	"fmt"
	"strings"
	"time"

	gonm "github.com/Wifx/gonetworkmanager"
	"github.com/google/uuid"
	"github.com/openrfsense/common/logging"
)

var log = logging.New().
	WithPrefix("system").
	WithFlags(logging.FlagsDevelopment).
	WithLevel(logging.DebugLevel)

const (
	defaultHotspotConnName = "Hotspot"
	defaultHotspotTimeout  = 5 * time.Minute
)

var connectionBase = gonm.ConnectionSettings{
	"connection": {
		"type":           "wifi",
		"interface-name": "wlan0",
		"autoconnect":    "true",
	},
	"wifi": {
		"mode": "infrastructure",
	},
	"wifi-security": {},
	"ipv4": {
		"method": "auto",
	},
	"ipv6": {
		"method": "auto",
	},
}

// Looks for a connection with a specific name (default is "Hotspot") and
// activates it on the primary wireless device, creating an access point.
func EnableHotspot() {
	nm, err := gonm.NewNetworkManager()
	if err != nil {
		return
	}

	wirelessDev, err := GetPrimaryWirelessDevice(nm)
	if err != nil {
		return
	}

	settings, _ := gonm.NewSettings()
	conns, _ := settings.ListConnections()
	for _, conn := range conns {
		connSettings, _ := conn.GetSettings()
		if connSettings["connection"]["id"] == defaultHotspotConnName {
			_, _ = nm.ActivateConnection(conn, wirelessDev, nil)
		}
	}
}

// Looks for an active wireless connection with a specific name (default is "Hotspot")
// and disables it after some time. Blocking.
func StartHotspotDisabler() {
	<-time.After(defaultHotspotTimeout)

	nm, err := gonm.NewNetworkManager()
	if err != nil {
		return
	}

	wirelessDev, err := GetPrimaryWirelessDevice(nm)
	if err != nil {
		return
	}

	activeConn, _ := wirelessDev.GetPropertyActiveConnection()
	id, _ := activeConn.GetPropertyID()
	if id != defaultHotspotConnName {
		log.Info("hotspot is not active")
	}

	_ = nm.DeactivateConnection(activeConn)
}

// Connects to an arbitrary wireless network with the first realized interface. Adds a new connection
// to NetworkManager if required.
func WirelessConnect(ssid string, password string, security string) (gonm.ActiveConnection, error) {
	nm, err := gonm.NewNetworkManager()
	if err != nil {
		return nil, err
	}

	wirelessDev, err := GetPrimaryWirelessDevice(nm)
	if err != nil {
		return nil, err
	}

	var ap gonm.AccessPoint
	allAps, _ := wirelessDev.GetAllAccessPoints()
	for _, currentAp := range allAps {
		currentSsid, _ := currentAp.GetPropertySSID()
		if currentSsid == ssid {
			ap = currentAp
		}
	}

	connections, _ := wirelessDev.GetPropertyAvailableConnections()
	for _, conn := range connections {
		connSettings, _ := conn.GetSettings()
		// The connection already exists
		if connSettings["connection"]["id"] == ssid {
			return nm.ActivateWirelessConnection(conn, wirelessDev, ap)
		}
	}

	if conn, exists := WirelessConnectionExists(ssid); exists {
		return nm.ActivateWirelessConnection(conn, wirelessDev, ap)
	}

	connection, err := generateConnection(ssid, password, security, connectionBase)
	if err != nil {
		return nil, err
	}

	return nm.AddAndActivateWirelessConnection(connection, wirelessDev, ap)
}

// If a wireless connection to the network with the given SSID is already present
// in NetworkManager, returns the connection object and "true".
func WirelessConnectionExists(ssid string) (gonm.Connection, bool) {
	settings, _ := gonm.NewSettings()
	allConns, _ := settings.ListConnections()
	for _, conn := range allConns {
		connSettings, _ := conn.GetSettings()

		ssidBytes := connSettings["802-11-wireless"]["ssid"].([]byte)
		if string(ssidBytes) == ssid {
			return conn, true
		}
	}

	return nil, false
}

func GetPrimaryWirelessDevice(nm gonm.NetworkManager) (gonm.DeviceWireless, error) {
	var wirelessDev gonm.DeviceWireless
	var err error

	// Wait for NetworkManager to be ready
	startup := true
	log.Info("Waiting fot NetworkManager")
	for startup {
		<-time.After(time.Second)
		startup, _ = nm.GetPropertyStartup()
	}

	devices, _ := nm.GetDevices()
	for _, d := range devices {
		devType, _ := d.GetPropertyDeviceType()
		if devType != gonm.NmDeviceTypeWifi {
			continue
		}

		wirelessDev, err = gonm.NewDeviceWireless(d.GetPath())
		if err != nil {
			continue
		}
	}

	if wirelessDev == nil {
		return nil, fmt.Errorf("could not find a connected wireless device")
	}

	return wirelessDev, nil
}

// Fill in the connection map template with the given parameters and a fresh UUID.
func generateConnection(ssid string, password string, security string, connection gonm.ConnectionSettings) (gonm.ConnectionSettings, error) {
	if strings.TrimSpace(ssid) == "" {
		return nil, fmt.Errorf("ssid must not be empty")
	}

	if strings.TrimSpace(security) == "" {
		return nil, fmt.Errorf("security must not be empty")
	}

	if strings.TrimSpace(password) == "" && security != "none" {
		return nil, fmt.Errorf("password can be empty only for unsecured access points")
	}

	connection["connection"]["id"] = ssid
	connection["connection"]["uuid"] = uuid.New()
	connection["wifi"]["ssid"] = ssid
	connection["wifi-security"]["key-mgmt"] = security
	connection["wifi-security"]["psk"] = password

	return connection, nil
}
