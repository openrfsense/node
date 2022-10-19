package ui

import (
	"net"

	"github.com/gofiber/fiber/v2"
	"github.com/mdlayher/ethtool"
	"github.com/mdlayher/wifi"
	"github.com/openrfsense/common/config"
)

func newConfMap() (fiber.Map, error) {
	configMap := fiber.Map{
		"nats": fiber.Map{
			"token": config.GetOrDefault("nats.token", ""),
		},
	}
	return configMap, nil
}

func newEthMap() (fiber.Map, error) {
	ethMap := fiber.Map{
		"connected": false,
	}

	ethtool, err := ethtool.New()
	if err != nil {
		return nil, err
	}
	defer ethtool.Close()

	// FIXME: ethtool is slow/broken on Pi, use gonetworkmanager
	states, err := ethtool.LinkStates()
	if err != nil {
		log.Error("failed to get eth link infos")
		return nil, err
	}

	for _, state := range states {
		if !state.Link {
			continue
		}

		iface, err := net.InterfaceByIndex(state.Interface.Index)
		if err != nil {
			log.Error(err)
			break
		}

		addrs, err := iface.Addrs()
		if err != nil {
			log.Error(err)
			break
		}

		for _, addr := range addrs {
			// NOTE: IsPrivate generally returns true for virtual network devices
			if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() && !ipnet.IP.IsPrivate() {
				if ipnet.IP.To4() != nil {
					ethMap["connected"] = true
					ethMap["ip"] = ipnet.IP.String()
					ethMap["interface"] = state.Interface.Name
					break
				}
			}
		}
	}

	return ethMap, nil
}

func newWifiMap() (fiber.Map, error) {
	wifiMap := fiber.Map{
		"connected": false,
	}

	wc, err := wifi.New()
	if err != nil {
		return nil, err
	}

	wIfaces, err := wc.Interfaces()
	if err != nil {
		return nil, err
	}

	for _, iface := range wIfaces {
		if iface.Type != wifi.InterfaceTypeStation {
			continue
		}

		bss, _ := wc.BSS(iface)
		netIface, err := net.InterfaceByIndex(iface.Index)
		if err != nil {
			log.Error(err)
			break
		}

		addrs, err := netIface.Addrs()
		if err != nil {
			log.Error(err)
			break
		}

		for _, addr := range addrs {
			if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
				// TODO: IPv6 support?
				if ipnet.IP.To4() != nil {
					wifiMap["connected"] = true
					wifiMap["ssid"] = bss.SSID
					wifiMap["ip"] = ipnet.IP.String()
					wifiMap["interface"] = iface.Name
					break
				}
			}
		}
	}

	return wifiMap, nil
}
