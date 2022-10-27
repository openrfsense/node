package api

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/knadh/koanf"

	"github.com/openrfsense/common/logging"
)

var log = logging.New().
	WithPrefix("api").
	WithLevel(logging.DebugLevel).
	WithFlags(logging.FlagsDevelopment)

// Create a router and use logger for the internal API. Initializes REST endpoints under the given prefix.
func Start(config *koanf.Koanf, prefix string, routerConfig ...fiber.Config) *fiber.App {
	router := fiber.New(routerConfig...)

	// TODO: is auth needed?
	router.Use(
		logger.New(),
		recover.New(),
		requestid.New(),
	)

	router.Route(prefix, func(router fiber.Router) {
		router.Post("/network/wifi", HandleWifiPost)
		router.Post("/config", HandleConfigPost)
	})

	addr := fmt.Sprintf(":%d", config.MustInt("node.port"))

	go func() {
		if err := router.Listen(addr); err != nil {
			log.Fatal(err)
		}
	}()

	return router
}
