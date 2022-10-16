package nats

import (
	"encoding/json"

	nats "github.com/nats-io/nats.go"

	"github.com/openrfsense/common/types"
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

func HandlerAggregatedMeasurement(conn *nats.EncodedConn, msg *nats.Msg) error {
	amr := types.AggregatedMeasurementRequest{}
	err := json.Unmarshal(msg.Data, &amr)
	if err != nil {
		return err
	}

	for _, id := range amr.Sensors {
		if id == system.ID() {
			log.Debugf("%#v\n", amr)
		}
	}

	return HandlerStatsBrief(conn, msg)
}
