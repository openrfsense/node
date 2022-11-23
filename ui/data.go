package ui

import (
	"sort"

	gonm "github.com/Wifx/gonetworkmanager"
	"github.com/gofiber/fiber/v2"
	"github.com/openrfsense/node/system"
)

func newEthMap() (fiber.Map, error) {
	nm, err := gonm.NewNetworkManager()
	if err != nil {
		return nil, err
	}

	ret := fiber.Map{
		"connected": false,
		"ip":        "",
		"interface": "",
	}

	devices, _ := nm.GetDevices()
	for _, d := range devices {
		devType, _ := d.GetPropertyDeviceType()
		if devType != gonm.NmDeviceTypeEthernet {
			continue
		}

		ethDev, err := gonm.NewDeviceWired(d.GetPath())
		if err != nil {
			continue
		}

		plugged, _ := ethDev.GetPropertyCarrier()
		ret["connected"] = plugged
		// Try to find at least one wired connection
		if !plugged {
			continue
		}

		conn, _ := ethDev.GetPropertyActiveConnection()
		connState, _ := conn.GetPropertyState()
		if connState == gonm.NmActiveConnectionStateActivated {
			ret["connected"] = true
			ip4, _ := conn.GetPropertyIP4Config()
			addrs, _ := ip4.GetPropertyAddresses()
			ret["ip"] = addrs[0].Address
		}

		ifName, _ := d.GetPropertyInterface()
		ret["interface"] = ifName
	}

	return ret, nil
}

func newWifiMap() (fiber.Map, error) {
	nm, err := gonm.NewNetworkManager()
	if err != nil {
		return nil, err
	}

	ret := fiber.Map{
		"connected": false,
		"ip":        "",
		"interface": "",
		"ssid":      "",
		"available": []string{},
		"saved":     []string{},
	}

	wirelessDev, err := system.GetPrimaryWirelessDevice(nm)
	if err != nil {
		return nil, err
	}

	conn, _ := wirelessDev.GetPropertyActiveConnection()
	connState, _ := conn.GetPropertyState()
	if connState == gonm.NmActiveConnectionStateActivated {
		ret["connected"] = true
		ip4, _ := conn.GetPropertyIP4Config()
		addrs, _ := ip4.GetPropertyAddresses()
		ret["ip"] = addrs[0].Address
	}

	ifName, _ := wirelessDev.GetPropertyInterface()
	ret["interface"] = ifName

	activeAp, _ := wirelessDev.GetPropertyActiveAccessPoint()
	activeSsid, _ := activeAp.GetPropertySSID()
	ret["ssid"] = activeSsid

	allAps, _ := wirelessDev.GetAccessPoints()
	allSsids := []string{}
	for _, ap := range allAps {
		ssid, _ := ap.GetPropertySSID()
		allSsids = append(allSsids, ssid)
	}
	sort.Strings(allSsids)
	ret["available"] = allSsids

	savedSsids := []string{}
	settings, _ := gonm.NewSettings()
	allConns, _ := settings.ListConnections()
	for _, conn := range allConns {
		connSettings, _ := conn.GetSettings()
		if connSettings["connection"]["interface-name"] == ifName {
			ssidBytes := connSettings["802-11-wireless"]["ssid"].([]byte)
			savedSsids = append(savedSsids, string(ssidBytes))
		}
	}
	sort.Strings(savedSsids)
	ret["saved"] = savedSsids

	return ret, nil
}
