package ui

import (
	"embed"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
	"github.com/gofiber/template/html"

	"github.com/openrfsense/common/logging"
)

var log = logging.New().
	WithPrefix("ui").
	WithFlags(logging.FlagsDevelopment)

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
	wifiMap, err := newWifiMap()
	if err != nil {
		return err
	}

	ethMap, err := newEthMap()
	if err != nil {
		return err
	}

	configMap, err := newConfMap()
	if err != nil {
		return err
	}

	return c.Render("views/index", fiber.Map{
		"wifi":   wifiMap,
		"eth":    ethMap,
		"config": configMap,
	})
}
