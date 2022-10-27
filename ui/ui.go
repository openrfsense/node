package ui

import (
	"embed"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
	"github.com/gofiber/template/html"
	"github.com/knadh/koanf"

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
func Init(config *koanf.Koanf, router *fiber.App) {
	router.Use(
		"/static",
		compress.New(compress.Config{
			Level: compress.LevelBestSpeed,
		}),
		func(c *fiber.Ctx) error {
			c.Set("Cache-Control", "public, max-age=31536000")
			return c.Next()
		},
		filesystem.New(filesystem.Config{
			Root:       http.FS(staticFs),
			PathPrefix: "static",
			Browse:     true,
		}),
	)

	router.Get("/", renderIndex(config))
}

// Renders the main webpage for the UI.
func renderIndex(config *koanf.Koanf) fiber.Handler {
	return func(c *fiber.Ctx) error {
		wifiMap, err := newWifiMap()
		if err != nil {
			return err
		}

		ethMap, err := newEthMap()
		if err != nil {
			return err
		}

		configMap, err := newConfMap(config)
		if err != nil {
			return err
		}

		return c.Render("views/index", fiber.Map{
			"wifi":   wifiMap,
			"eth":    ethMap,
			"config": configMap,
		})
	}
}
