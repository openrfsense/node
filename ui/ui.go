package ui

import (
	"embed"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
	"github.com/gofiber/template/html"
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

	router.Get("/", func(c *fiber.Ctx) error {
		return c.Render("views/index", fiber.Map{})
	})
}
