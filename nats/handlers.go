package nats

import (
	nats "github.com/nats-io/nats.go"

	"github.com/openrfsense/node/system"
)

func HandlerStats(conn *nats.EncodedConn, msg *nats.Msg) error {
	stats, err := system.GetStats()
	if err != nil {
		return err
	}

	return conn.Publish(msg.Reply, *stats)
}

func HandlerStatsBrief(conn *nats.EncodedConn, msg *nats.Msg) error {
	stats, err := system.GetStatsBrief()
	if err != nil {
		return err
	}

	return conn.Publish(msg.Reply, *stats)
}
