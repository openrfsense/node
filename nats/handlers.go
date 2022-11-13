package nats

import (
	"encoding/json"

	nats "github.com/nats-io/nats.go"

	"github.com/openrfsense/common/types"
	"github.com/openrfsense/node/sensor"
	"github.com/openrfsense/node/stats"
	"github.com/openrfsense/node/system"
)

// Responds with full system stats (system.GetStats).
func HandlerStats(conn *nats.EncodedConn, msg *nats.Msg) error {
	stat, err := stats.GetStats()
	if err != nil {
		return err
	}

	return conn.Publish(msg.Reply, *stat)
}

// Responds with brief system stats (system.GetStatsBrief).
func HandlerStatsBrief(conn *nats.EncodedConn, msg *nats.Msg) error {
	stat, err := stats.GetStatsBrief()
	if err != nil {
		return err
	}

	return conn.Publish(msg.Reply, *stat)
}

// Starts an aggregated measurement and sends back brief stats.
func HandlerAggregatedMeasurement(conn *nats.EncodedConn, msg *nats.Msg) error {
	amr := types.AggregatedMeasurementRequest{}
	err := json.Unmarshal(msg.Data, &amr)
	if err != nil {
		return err
	}

	for _, id := range amr.Sensors {
		if id == system.ID() {
			log.Debugf("got measurement request: %#v\n", amr)
			go sensor.WithAggregated(amr).Run()
			return HandlerStatsBrief(conn, msg)
		}
	}

	return nil
}

// Starts a raw measurement and sends back brief stats.
func HandlerRawMeasurement(conn *nats.EncodedConn, msg *nats.Msg) error {
	rmr := types.RawMeasurementRequest{}
	err := json.Unmarshal(msg.Data, &rmr)
	if err != nil {
		return err
	}

	for _, id := range rmr.Sensors {
		if id == system.ID() {
			log.Debugf("got measurement request: %#v\n", rmr)
			go sensor.WithRaw(rmr).Run()
			return HandlerStatsBrief(conn, msg)
		}
	}

	return nil
}
