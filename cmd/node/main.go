package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	emitter "github.com/emitter-io/go/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/spf13/pflag"

	"github.com/openrfsense/common/config"
	"github.com/openrfsense/common/keystore"
	"github.com/openrfsense/node/api"
	"github.com/openrfsense/node/mqtt"
	"github.com/openrfsense/node/ui"
)

func main() {
	configPath := pflag.StringP("config", "c", "", "path to yaml config file")
	pflag.Parse()

	id := uuid.New().String()
	log.Println("Starting node " + id)

	log.Println("Loading config")
	err := config.Load(*configPath)
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

	log.Printf("Node id is: %s", mqtt.Client.ID())

	// TODO: move elsewhere
	mqtt.Subscribe("sensors/all/", nil)
	mqtt.Subscribe("sensors/"+mqtt.Client.ID()+"/cmd/", func(_ *emitter.Client, m emitter.Message) {
		log.Printf("received remote command: %s", string(m.Payload()))
	})
	mqtt.Publish("sensors/"+mqtt.Client.ID()+"/output/", "hello")

	log.Println("Starting internal API")
	router := fiber.New(fiber.Config{
		Views: ui.NewEngine(),
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

	// addr := fmt.Sprintf(":%d", config.Get[int]("api.port"))
	if err := router.Listen(":9090"); err != nil {
		log.Fatal(err)
	}

	<-shutdown
}
