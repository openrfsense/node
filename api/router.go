package api

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
)

// Configure a router for the interal API. Initializes REST endpoints under the given prefix.
func Use(router *fiber.App, prefix string) {
	// TODO: is auth needed?
	router.Use(
		logger.New(),
		recover.New(),
		requestid.New(),
	)

	router.Route(prefix, func(router fiber.Router) {
		router.Post("/wifi", func(c *fiber.Ctx) error {
			log.Printf("SSID: " + c.FormValue("ssid"))
			log.Printf("Password: " + c.FormValue("password"))
			return c.SendString("")
		})
	})
}
