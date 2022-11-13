package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/gofiber/fiber/v2"
	"github.com/spf13/pflag"

	"github.com/openrfsense/common/logging"
	"github.com/openrfsense/node/api"
	"github.com/openrfsense/node/config"
	"github.com/openrfsense/node/nats"
	"github.com/openrfsense/node/sensor"
	"github.com/openrfsense/node/stats"
	"github.com/openrfsense/node/system"
	"github.com/openrfsense/node/ui"
)

var (
	version = ""
	commit  = ""
	date    = ""

	log = logging.New().
		WithPrefix("main").
		WithLevel(logging.DebugLevel).
		WithFlags(logging.FlagsDevelopment)
)

func main() {
	configPath := pflag.StringP("config", "c", "/etc/openrfsense/config.yml", "path to yaml config file")
	natsTokenPath := pflag.StringP("token", "t", "/etc/openrfsense/token.txt", "path to token file")
	showVersion := pflag.BoolP("version", "v", false, "shows program version and build info")
	pflag.Parse()

	if *showVersion {
		fmt.Printf("openrfsense-node v%s (%s) built on %s\n", version, commit, date)
		return
	}

	log.Infof("Starting node %s", system.ID())

	log.Info("Loading config")
	konfig, err := config.Load(*configPath)
	if err != nil {
		log.Fatal(err)
	}

	stats.Init(konfig)

	log.Info("Initializing sensor manager")
	err = sensor.Init(konfig)
	if err != nil {
		log.Fatal(err)
	}

	// Connect to NATS
	log.Info("Connecting to NATS")
	err = nats.Init(konfig, *natsTokenPath)
	if err != nil {
		log.Fatal(err)
	}
	defer nats.Disconnect()

	log.Info("Starting internal API")

	// Start the internal backend
	router := api.Start(konfig, "/api", fiber.Config{
		AppName:               "openrfsense-node",
		DisableStartupMessage: true,
		PassLocalsToViews:     true,
		Views:                 ui.NewEngine(),
	})
	// Initialize UI (templated web pages)
	ui.Init(konfig, router)
	defer func() {
		err = router.Shutdown()
		if err != nil {
			log.Fatal(err)
		}
	}()

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()
	<-ctx.Done()
	log.Info("Shutting down")
}
