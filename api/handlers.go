package api

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
)

func HandleConfigPost(ctx *fiber.Ctx) error {
	// TODO: save token to file
	log.Debug("NATS token: " + ctx.FormValue("nats-token"))

	return ctx.SendStatus(http.StatusOK)
}

func HandleWifiPost(ctx *fiber.Ctx) error {
	// FIXME: implement this
	log.Debug("SSID: " + ctx.FormValue("ssid"))
	log.Debug("Password: " + ctx.FormValue("password"))

	return ctx.SendStatus(http.StatusOK)
}
