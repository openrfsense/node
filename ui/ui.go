package ui

import (
	"embed"
	"log"
	"net"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
	"github.com/gofiber/template/html"
	"github.com/mdlayher/ethtool"
	"github.com/mdlayher/wifi"
)

//go:embed views/*
var viewsFs embed.FS

//go:embed static/*
var staticFs embed.FS

// Initializes Fiber view engine with embedded HTML templates from views/
func NewEngine() *html.Engine {
	engine := html.NewFileSystem(http.FS(viewsFs), ".html")
	engine.Reload(true)

	return engine
}

// Configure a router for the UI. Initializes routes and view models.
func Use(router *fiber.App) {
	router.Use("/static", filesystem.New(filesystem.Config{
		Root:       http.FS(staticFs),
		PathPrefix: "static",
		Browse:     true,
	}))

	router.Get("/", renderIndex)
}

// Renders the main webpage for the UI
func renderIndex(c *fiber.Ctx) error {
	wifiMap := fiber.Map{
		"connected": false,
	}
	ethMap := fiber.Map{
		"connected": false,
	}

	wc, err := wifi.New()
	if err != nil {
		return err
	}

	wIfaces, err := wc.Interfaces()
	if err != nil {
		return err
	}

	for _, iface := range wIfaces {
		if iface.Type != wifi.InterfaceTypeStation {
			continue
		}

		bss, _ := wc.BSS(iface)
		netIface, err := net.InterfaceByIndex(iface.Index)
		if err != nil {
			log.Print(err)
			break
		}

		addrs, err := netIface.Addrs()
		if err != nil {
			log.Print(err)
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

	ethtool, err := ethtool.New()
	defer ethtool.Close()
	if err != nil {
		return err
	}

	infos, err := ethtool.LinkInfos()
	if err != nil {
		return err
	}

	for _, i := range infos {
		iface, err := net.InterfaceByIndex(i.Interface.Index)
		if err != nil {
			log.Print(err)
			break
		}

		addrs, err := iface.Addrs()
		if err != nil {
			log.Print(err)
			break
		}

		for _, addr := range addrs {
			// NOTE: IsPrivate generally returns true for virtual network devices
			if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() && !ipnet.IP.IsPrivate() {
				if ipnet.IP.To4() != nil {
					ethMap["connected"] = true
					ethMap["ip"] = ipnet.IP.String()
					ethMap["interface"] = i.Interface.Name
					break
				}
			}
		}
	}

	return c.Render("views/index", fiber.Map{
		"wifi": wifiMap,
		"eth":  ethMap,
	})
}
