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
			log.Debugf("Got measurement request: %#v\n", amr)
			return HandlerStatsBrief(conn, msg)
		}
	}

	return nil
}

func HandlerRawMeasurement(conn *nats.EncodedConn, msg *nats.Msg) error {
	rmr := types.RawMeasurementRequest{}
	err := json.Unmarshal(msg.Data, &rmr)
	if err != nil {
		return err
	}

	for _, id := range rmr.Sensors {
		if id == system.ID() {
			log.Debugf("Got measurement request: %#v\n", rmr)
			return HandlerStatsBrief(conn, msg)
		}
	}

	return nil
}
