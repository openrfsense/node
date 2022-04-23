package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/spf13/pflag"

	"github.com/openrfsense/common/config"
	"github.com/openrfsense/common/keystore"
	"github.com/openrfsense/common/logging"
	"github.com/openrfsense/node/api"
	"github.com/openrfsense/node/mqtt"
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

type MQTT struct {
	Protocol string `koanf:"protocol"`
	Host     string `koanf:"host"`
	Port     int    `koanf:"port"`
}

type NodeConfig struct {
	Node    `koanf:"node"`
	Backend `koanf:"backend"`
	MQTT    `koanf:"mqtt"`
}

// FIXME: move elsewhere
var DefaultConfig = NodeConfig{
	Node: Node{
		Port: 8080,
	},
	Backend: Backend{
		Port: 8080,
	},
	MQTT: MQTT{
		Protocol: "tcp",
		Port:     8080,
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
	configPath := pflag.StringP("config", "c", "", "path to yaml config file")
	pflag.Parse()

	log.Infof("Starting node %s", system.ID())

	log.Info("Loading config")
	err := config.Load(*configPath, DefaultConfig)
	if err != nil {
		log.Fatal(err)
	}

	log.Info("Starting keystore")
	keystore.Init(mqtt.NewBackendRetriever(), mqtt.DefaultTTL)

	log.Info("Connecting to MQTT")
	mqtt.Init()

	log.Info("Starting internal API")
	router := fiber.New(fiber.Config{
		AppName:               "openrfsense-node",
		DisableStartupMessage: true,
		PassLocalsToViews:     true,
		Views:                 ui.NewEngine(),
	})

	api.Init(router, "/api")
	ui.Init(router)

	// Graceful shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	shutdown := make(chan struct{})

	go func() {
		<-c
		router.Shutdown()
		mqtt.Disconnect(time.Second)
		log.Info("Shutting down")
		shutdown <- struct{}{}
	}()

	addr := fmt.Sprintf(":%d", config.GetWeakInt("node.port"))
	if err := router.Listen(addr); err != nil {
		log.Fatal(err)
	}

	<-shutdown
}
