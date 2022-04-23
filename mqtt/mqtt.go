package mqtt

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	emitter "github.com/emitter-io/go/v2"

	"github.com/openrfsense/common/config"
	"github.com/openrfsense/common/keystore"
	"github.com/openrfsense/common/logging"
	"github.com/openrfsense/common/types"
	"github.com/openrfsense/node/system"
)

var (
	client     *emitter.Client
	DefaultTTL = 600 * time.Second

	log = logging.New().
		WithPrefix("mqtt").
		WithLevel(logging.DebugLevel).
		WithFlags(logging.FlagsDevelopment)
)

// Type Payload represents a payload which can be sent over MQTT.
// Since emitter rejects anything which isn't a string or byte array,
// using generics ensures the handlers always return a payload which
// can be sent. In general, message sending has to be ensured.
type Payload interface {
	~string | ~[]byte
}

// Type Handler represents a simplified emitter.MessageHandler which returns a payload.
// The payload is then managed by the Handle function.
type Handler[P Payload] func(m emitter.Message) (P, error)

// Starts the internal MQTT client (and connects to the broker) and sets up handlers for
//  messages on the right topics.
func Init() {
	brokerHost := fmt.Sprintf("%s:%d", config.Must[string]("mqtt.host"), config.Must[int]("mqtt.port"))
	brokerUrl := url.URL{
		Scheme: config.Must[string]("mqtt.protocol"),
		Host:   brokerHost,
	}

	client = emitter.NewClient(
		emitter.WithUsername(system.ID()),
		emitter.WithBrokers(brokerUrl.String()),
		emitter.WithAutoReconnect(true),
		emitter.WithConnectTimeout(10*time.Second),
		emitter.WithKeepAlive(10*time.Second),
		emitter.WithMaxReconnectInterval(2*time.Minute),
	)

	err := client.Connect()
	if err != nil {
		ticker := time.NewTicker(30 * time.Second)
		for !client.IsConnected() {
			log.Warn("Could not connect to MQTT broker, trying again")
			<-ticker.C
			client.Connect()
		}
		ticker.Stop()
	}

	client.OnConnect(func(_ *emitter.Client) {
		log.Info("Connected to MQTT broker")
	})

	Handle("/all/", "get", HandlerStatsBrief)
	Handle("stats/", "get", HandlerStats)
}

// Register a MQTT message handler for requests. The requests are actually received
// at node/method/path/ and payloads returned by the handler are sent to node/path/.
func Handle[P Payload](path string, method string, handler Handler[P]) {
	trimmedPath := strings.Trim(path, "/")
	trimmedMethod := strings.Trim(method, "/")

	// FIXME: use something better
	var pathResp, pathReq string
	if strings.HasPrefix(path, "/") {
		pathResp = fmt.Sprintf("node/%s/", trimmedPath)
		pathReq = fmt.Sprintf("node/%s/%s/", trimmedMethod, trimmedPath)
	} else {
		pathResp = fmt.Sprintf("node/%s/%s/", system.ID(), trimmedPath)
		pathReq = fmt.Sprintf("node/%s/%s/%s/", trimmedMethod, system.ID(), trimmedPath)
	}

	Subscribe(pathReq, func(_ *emitter.Client, m emitter.Message) {
		payload, err := handler(m)
		if err != nil {
			log.Error(err)
			return
		}
		log.Debugf("sending %T on %s", payload, pathResp)
		Publish(pathResp, payload)
	})
}

// Custom keystore.Retriever which fetches channel keys from the OpenRFSense backend.
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

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return "", err
		}

		return string(body), nil
	}
}
