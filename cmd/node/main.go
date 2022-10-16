package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/gofiber/fiber/v2"
	"github.com/spf13/pflag"

	"github.com/openrfsense/common/config"
	"github.com/openrfsense/common/logging"
	"github.com/openrfsense/node/api"
	"github.com/openrfsense/node/nats"
	"github.com/openrfsense/node/sensor"
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
	err := config.Load(*configPath, DefaultConfig)
	if err != nil {
		log.Fatal(err)
	}

	log.Info("Initializing sensor manager")
	cmdFlags := sensor.CommandFlags{}
	err = config.Unmarshal("node.sensor", &cmdFlags)
	if err != nil {
		log.Fatal(err)
	}
	sensor.Init(cmdFlags)

	// Connect ot NATS only if the node is connected to the internet
	// if system.IsOnline() {
	log.Info("Connecting to NATS")
	err = nats.Init(*natsTokenPath)
	if err != nil {
		log.Fatal(err)
	}
	defer nats.Disconnect()
	// }

	log.Info("Starting internal API")

	// Start the internal backend
	router := api.Start("/api", fiber.Config{
		AppName:               "openrfsense-node",
		DisableStartupMessage: true,
		PassLocalsToViews:     true,
		Views:                 ui.NewEngine(),
	})
	// Initialize UI (templated web pages)
	ui.Init(router)
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
