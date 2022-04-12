package ui

import (
	"embed"
	"net"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
	"github.com/gofiber/template/html"
	"github.com/mdlayher/ethtool"
	"github.com/mdlayher/wifi"

	"github.com/openrfsense/common/logging"
)

var log = logging.New(
	logging.WithPrefix("ui"),
	logging.WithFlags(logging.FlagsDevelopment),
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

// Configure a router and use a logger for the UI. Initializes routes and view models.
func Init(router *fiber.App) {
	router.Use(
		"/static",
		filesystem.New(filesystem.Config{
			Root:       http.FS(staticFs),
			PathPrefix: "static",
			Browse:     true,
		}),
		compress.New(),
	)

	router.Get("/", renderIndex)
}

// Renders the main webpage for the UI.
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

	ethtool, err := ethtool.New()
	defer ethtool.Close()
	if err != nil {
		return err
	}

	// FIXME: ethtool is slow at detecting the link state, use a watcher + context?
	states, err := ethtool.LinkStates()
	if err != nil {
		log.Error("failed to get eth link infos")
		return err
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

	return c.Render("views/index", fiber.Map{
		"wifi": wifiMap,
		"eth":  ethMap,
	})
}
