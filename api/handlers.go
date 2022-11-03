package api

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/openrfsense/node/config"
	"gopkg.in/yaml.v3"
)

func HandleConfigPost(ctx *fiber.Ctx) error {
	text := ctx.FormValue("configText")
	log.Debug("Config text: " + text)
	if len(text) == 0 {
		return fiber.ErrBadRequest
	}
	conf := config.NodeConfig{}
	err := yaml.Unmarshal([]byte(text), &conf)
	if err != nil {
		return err
	}

	log.Debugf("Got config: %#v", conf)

	return ctx.SendStatus(http.StatusOK)
}

func HandleWifiPost(ctx *fiber.Ctx) error {
	// FIXME: implement this
	log.Debug("SSID: " + ctx.FormValue("ssid"))
	log.Debug("Password: " + ctx.FormValue("password"))

	return ctx.SendStatus(http.StatusOK)
}
