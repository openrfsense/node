package nats

import (
	"fmt"
	"math"
	"strings"

	"github.com/knadh/koanf"
	nats "github.com/nats-io/nats.go"

	"github.com/openrfsense/common/logging"
	"github.com/openrfsense/node/system"
)

// Type Route represents a channel to subscribe to and the relative handler.
type Route struct {
	Subject string
	Handler interface{}
}

// The subjects used by the node and relative handlers
var routes = []Route{
	{".all", HandlerStatsBrief},
	{"stats", HandlerStats},
	{".all.aggregated", HandlerAggregatedMeasurement},
	{".all.raw", HandlerRawMeasurement},
}

var (
	conn   *nats.EncodedConn
	errors chan error

	log = logging.New().
		WithPrefix("nats").
		WithLevel(logging.DebugLevel).
		WithFlags(logging.FlagsDevelopment)
)

// Initializes the internal NATS connection and sets up handlers for various subjects.
// Uses the token found in tokenFile but also looks for the token in the config, under
// nats.token (ORFS_NATS_TOKEN in env variables).
func Init(config *koanf.Koanf, tokenFile string) error {
	addr := fmt.Sprintf("nats://%s:%d", config.String("nats.host"), config.MustInt("nats.port"))
	token := config.MustString("nats.token")
	errors = make(chan error, 1)

	// Connect and encode the connection
	var err error
	conn, err = connect(addr, system.ID(), token)
	if err != nil {
		return err
	}

	// Register the routes
	for _, route := range routes {
		err = handle(conn, system.ID(), route.Subject, route.Handler)
		if err != nil {
			log.Error(err)
		}
	}

	// Start async error logger
	go errorLogger(errors)
	// Start manager data sender
	go sendManagerData(conn, errors)

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
		nats.Name(clientId),
		nats.Token(token),
		nats.RetryOnFailedConnect(true),
		nats.MaxReconnects(math.MaxInt),
		nats.ReconnectHandler(func(c *nats.Conn) {
			log.Info("Connection estabilished")
		}),
		nats.DisconnectErrHandler(func(_ *nats.Conn, err error) {
			log.Warnf("Connection lost: %v", err)
		}),
	)
	if err != nil {
		return nil, err
	}

	conn, err := nats.NewEncodedConn(c, nats.JSON_ENCODER)
	if err != nil {
		return nil, err
	}

	return conn, nil
}

// Registers a custom message handler (see type Handler) with automatic path formatting.
// Paths beginning with '.' (the separator) are absolute and formatted to 'node.$path',
// while paths like '$path' are prefixed with the client ID and become 'node.$id.$path'.
func handle(conn *nats.EncodedConn, clientId string, path string, handler interface{}) error {
	subject := []string{"node"}
	if !strings.HasPrefix(path, ".") {
		subject = append(subject, clientId)
	}
	subject = append(subject, strings.Trim(path, "."))
	requestPath := strings.Join(subject, ".")

	_, err := conn.Subscribe(requestPath, handler)
	if err == nil {
		log.Debugf("registered subject %s", requestPath)
	}

	return err
}
