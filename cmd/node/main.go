package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/gofiber/fiber/v2"
	"github.com/spf13/pflag"

	"github.com/openrfsense/common/config"
	"github.com/openrfsense/common/logging"
	"github.com/openrfsense/node/api"
	"github.com/openrfsense/node/nats"
	"github.com/openrfsense/node/system"
	"github.com/openrfsense/node/ui"
)

type Node struct {
	Port int `koanf:"port"`
}

type Backend struct {
	Port  int               `koanf:"port"`
	Users map[string]string `koanf:"users"`
}

type NATS struct {
	Host string `koanf:"host"`
	Port int    `koanf:"port"`
}

type NodeConfig struct {
	Node    `koanf:"node"`
	Backend `koanf:"backend"`
	NATS    `koanf:"nats"`
}

// FIXME: move elsewhere
var DefaultConfig = NodeConfig{
	Node: Node{
		Port: 8080,
	},
	Backend: Backend{
		Port: 8080,
	},
	NATS: NATS{
		Port: 0,
	},
}

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
	natsTokenPath := pflag.StringP("token", "t", "/etc/openrfsense/token.txt", "path to yaml config file")
	pflag.Parse()

	log.Infof("Starting node %s", system.ID())

	log.Info("Loading config")
	err := config.Load(*configPath, DefaultConfig)
	if err != nil {
		log.Fatal(err)
	}

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

	router := api.Start("/api", fiber.Config{
		AppName:               "openrfsense-node",
		DisableStartupMessage: true,
		PassLocalsToViews:     true,
		Views:                 ui.NewEngine(),
	})
	ui.Init(router)
	defer router.Shutdown()

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()
	<-ctx.Done()
	log.Info("Shutting down")
}
