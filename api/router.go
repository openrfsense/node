package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"

	"github.com/openrfsense/common/logging"
)

var log = logging.New().
	WithPrefix("api").
	WithLevel(logging.DebugLevel).
	WithFlags(logging.FlagsDevelopment)

// Configure a router and use logger for the interal API. Initializes REST endpoints under the given prefix.
func Init(router *fiber.App, prefix string) {
	// TODO: is auth needed?
	router.Use(
		logger.New(),
		recover.New(),
		requestid.New(),
	)

	router.Route(prefix, func(router fiber.Router) {
		router.Post("/network/wifi", func(c *fiber.Ctx) error {
			// FIXME: implement this
			log.Debug("router", "SSID: "+c.FormValue("ssid"))
			log.Debug("router", "Password: "+c.FormValue("password"))

			return c.SendString("")
		})
	})
}
