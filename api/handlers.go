package api

import (
	"net/http"

	"github.com/openrfsense/node/config"
	"github.com/openrfsense/node/system"

	"github.com/gofiber/fiber/v2"
	"gopkg.in/yaml.v3"
)

func HandleConfigPost(ctx *fiber.Ctx) error {
	text := ctx.FormValue("configText")
	if len(text) == 0 {
		return fiber.ErrBadRequest
	}
	conf := config.NodeConfig{}
	err := yaml.Unmarshal([]byte(text), &conf)
	if err != nil {
		return err
	}

	err = config.Save(text)
	if err != nil {
		return err
	}

	return ctx.SendStatus(http.StatusOK)
}

func HandleWifiPost(ctx *fiber.Ctx) error {
	_, err := system.WirelessConnect(ctx.FormValue("ssid"), ctx.FormValue("password"), ctx.FormValue("security"))
	if err != nil {
		return err
	}

	return ctx.SendStatus(http.StatusOK)
}
