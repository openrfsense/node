package mqtt

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	emitter "github.com/emitter-io/go/v2"

	"github.com/openrfsense/common/config"
	"github.com/openrfsense/common/keystore"
	"github.com/openrfsense/common/logging"
	"github.com/openrfsense/common/types"
)

var (
	Client     *emitter.Client
	DefaultTTL = 600 * time.Second

	log = logging.New(
		logging.WithPrefix("mqtt"),
		logging.WithLevel(logging.DebugLevel),
		logging.WithFlags(logging.FlagsDevelopment),
	)
)

// TODO: make a better init procedure and/or move to openrfsense-common
func Init() {
	brokerHost := fmt.Sprintf("%s:%d", config.Must[string]("mqtt.host"), config.Must[int]("mqtt.port"))
	brokerUrl := url.URL{
		Scheme: config.Must[string]("mqtt.protocol"),
		Host:   brokerHost,
	}

	Client = emitter.NewClient(
		emitter.WithUsername(ID()),
		emitter.WithBrokers(brokerUrl.String()),
		emitter.WithAutoReconnect(true),
		emitter.WithConnectTimeout(10*time.Second),
		emitter.WithKeepAlive(10*time.Second),
		emitter.WithMaxReconnectInterval(2*time.Minute),
	)

	// FIXME: remove this
	Client.OnMessage(func(_ *emitter.Client, msg emitter.Message) {
		log.Debugf("mqtt", "[emitter] -> [B] received: '%s' topic: '%s'\n", msg.Payload(), msg.Topic())
	})

	err := Client.Connect()
	if err != nil {
		ticker := time.NewTicker(30 * time.Second)
		for !Client.IsConnected() {
			logging.Warn("mqtt", "Could not connect to MQTT broker, trying again")
			<-ticker.C
			Client.Connect()
		}
		ticker.Stop()
	}

	Client.OnConnect(func(_ *emitter.Client) {
		logging.Info("mqtt", "Connected to MQTT broker")
	})
}

// Custom keystore.Retriever which fetches channel keys from the OpenRFSense backend
func NewBackendRetriever() keystore.Retriever {
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	return func(channel string, access string) (string, error) {
		apiUrl := fmt.Sprintf(
			"http://%s:%d/api/v1/key",
			config.Must[string]("backend.host"),
			config.Must[int]("backend.port"),
		)

		keyReq := types.KeyRequest{
			Channel: channel,
			Access:  access,
		}

		data, err := json.Marshal(keyReq)
		if err != nil {
			return "", err
		}

		// Using Fiber Agent would be cool but its API is still unstable
		req, err := http.NewRequest(http.MethodPost, apiUrl, bytes.NewBuffer(data))
		if err != nil {
			return "", err
		}
		req.Header.Set("Content-Type", "application/json")
		req.SetBasicAuth(
			config.Must[string]("backend.username"),
			config.Must[string]("backend.password"),
		)
		resp, err := client.Do(req)
		if err != nil {
			return "", err
		}
		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return "", err
		}

		return string(body), nil
	}
}

// Disconnect will end the connection with the server, but not before waiting the specified
// time for existing work to be completed.
func Disconnect(waitTime time.Duration) {
	Client.Disconnect(waitTime)
}

// Wrapper around emitter.Presence with automatic key management.
// Presence sends a presence request to the broker.
func Presence(channel string, status, changes bool) error {
	key, err := keystore.Must(channel, "p")
	if err != nil {
		return err
	}
	return Client.Presence(key, channel, status, changes)
}

// Wrapper around emitter.Publish with automatic key management.
// Publish will publish a message with the specified QoS and content to the specified topic.
// Returns a token to track delivery of the message to the broker
func Publish(channel string, payload interface{}, options ...emitter.Option) error {
	key, err := keystore.Must(channel, "w")
	if err != nil {
		return err
	}
	return Client.Publish(key, channel, payload, options...)
}

// Wrapper around emitter.Subscribe with automatic key management.
// Subscribe starts a new subscription. Provide a MessageHandler to be executed when a
// message is published on the topic provided.
func Subscribe(channel string, optionalHandler emitter.MessageHandler, options ...emitter.Option) error {
	key, err := keystore.Must(channel, "r")
	if err != nil {
		return err
	}
	return Client.Subscribe(key, channel, optionalHandler, options...)
}
