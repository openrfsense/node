package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	emitter "github.com/emitter-io/go/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/spf13/pflag"

	"github.com/openrfsense/common/config"
	"github.com/openrfsense/common/keystore"
	"github.com/openrfsense/node/api"
	"github.com/openrfsense/node/mqtt"
	"github.com/openrfsense/node/ui"
)

type Node struct {
	Port int
}

type Backend struct {
	Port  int
	Users map[string]string
}

type MQTT struct {
	Protocol string
	Host     string
	Port     int
	Secret   string
}

type BackendConfig struct {
	Node
	Backend
	MQTT
}

var DefaultConfig = BackendConfig{
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

func main() {
	configPath := pflag.StringP("config", "c", "", "path to yaml config file")
	pflag.Parse()

	log.Println("Starting node " + mqtt.ID())

	log.Println("Loading config")
	err := config.Load(*configPath, DefaultConfig)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Starting keystore")
	keystore.Init(mqtt.NewBackendRetriever(), mqtt.DefaultTTL)

	log.Println("Connecting to MQTT")
	err = mqtt.InitClient()
	if err != nil {
		log.Fatalf("could not connect to MQTT: %v", err)
	}

	log.Printf("Remote node id is: %s", mqtt.Client.ID())

	// FIXME: move elsewhere
	mqtt.Subscribe("sensors/all/", nil)
	mqtt.Subscribe("sensors/"+mqtt.Client.ID()+"/cmd/", func(_ *emitter.Client, m emitter.Message) {
		log.Printf("received remote command: %s", string(m.Payload()))
	})
	mqtt.Publish("sensors/"+mqtt.Client.ID()+"/output/", "hello")

	log.Println("Starting internal API")
	router := fiber.New(fiber.Config{
		PassLocalsToViews: true,
		Views:             ui.NewEngine(),
	})

	api.Use(router, "/api")
	ui.Use(router)

	// Graceful shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	shutdown := make(chan struct{})

	go func() {
		<-c
		router.Shutdown()
		log.Println("Shutting down")
		shutdown <- struct{}{}
	}()

	addr := fmt.Sprintf(":%d", config.GetWeakInt("node.port"))
	if err := router.Listen(addr); err != nil {
		log.Fatal(err)
	}

	<-shutdown
}
