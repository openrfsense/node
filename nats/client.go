package nats

import (
	"fmt"
	"strings"

	"github.com/knadh/koanf"
	nats "github.com/nats-io/nats.go"

	"github.com/openrfsense/common/logging"
	"github.com/openrfsense/node/system"
)

// Type Handler is a custom handler for NATS messages which also receives
// a reference to the current NATS (encoded) connection. Any error returned
// is logged but does not halt the application.
type Handler func(*nats.EncodedConn, *nats.Msg) error

var (
	conn *nats.EncodedConn

	log = logging.New().
		WithLevel(logging.DebugLevel).
		WithFlags(logging.FlagsDevelopment)
)

// Initializes the internal NATS connection and sets up handlers for various subjects.
// Uses the token found in tokenFile but also looks for the token in the config, under
// nats.token (ORFS_NATS_TOKEN in env variables).
func Init(config *koanf.Koanf, tokenFile string) error {
	addr := fmt.Sprintf("nats://%s:%d", config.MustString("nats.host"), config.MustInt("nats.port"))

	var token string
	var err error
	if config.Exists("nats.token") {
		token = config.MustString("nats.token")
	} else {
		token, err = GetToken(tokenFile)
		if err != nil {
			return err
		}
	}

	conn, err = connect(addr, system.ID(), token)
	if err != nil {
		return err
	}

	// Handle requests on node.all
	err = handle(conn, system.ID(), ".all", HandlerStatsBrief)
	if err != nil {
		return err
	}

	// Handle requests on node.$id.stats
	err = handle(conn, system.ID(), "stats", HandlerStats)
	if err != nil {
		return err
	}

	// Handle requests on node.all.aggregated
	err = handle(conn, system.ID(), ".all.aggregated", HandlerAggregatedMeasurement)
	if err != nil {
		return err
	}

	// Handle requests on node.all.raw
	err = handle(conn, system.ID(), ".all.raw", HandlerRawMeasurement)
	if err != nil {
		return err
	}

	return nil
}

// Drain and close the internal NATS connection.
func Disconnect() {
	if conn != nil {
		err := conn.Drain()
		if err != nil {
			log.Error(err)
		}
		conn.Close()
	}
}

// Creates an encoded connection to the specified NATS address with a client ID.
func connect(addr string, clientId string, token string) (*nats.EncodedConn, error) {
	c, err := nats.Connect(
		addr,
		nats.RetryOnFailedConnect(true),
		nats.Name(clientId),
		nats.Token(token),
	)
	if err != nil {
		return nil, err
	}

	// TODO: log error (warning) if not connected after some time (goroutine)

	conn, err := nats.NewEncodedConn(c, nats.JSON_ENCODER)
	if err != nil {
		return nil, err
	}

	return conn, nil
}

// Registers a custom message handler (see type Handler) with automatic path formatting.
// Paths beginning with '.' (the separator) are absolute and formatted to 'node.$path',
// while paths like '$path' are prepended with the client ID and become 'node.$id.$path'.
func handle(conn *nats.EncodedConn, clientId string, path string, handler Handler) error {
	trimmed := strings.Trim(path, ".")
	requestPath := fmt.Sprintf("node.%s.%s", clientId, trimmed)
	if strings.HasPrefix(path, ".") {
		requestPath = fmt.Sprintf("node.%s", trimmed)
	}

	_, err := conn.Subscribe(requestPath, func(msg *nats.Msg) {
		log.Debugf("received message on %s", requestPath)
		err := handler(conn, msg)
		if err != nil {
			log.Error(err)
		}
	})

	if err == nil {
		log.Debugf("Subscribed to %s", requestPath)
	}

	return err
}
