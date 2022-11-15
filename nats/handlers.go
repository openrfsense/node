package nats

import (
	"github.com/openrfsense/common/types"
	"github.com/openrfsense/node/sensor"
	"github.com/openrfsense/node/stats"
	"github.com/openrfsense/node/system"
)

// Responds with full system stats (system.GetStats).
func HandlerStats(subject string, reply string, _ interface{}) {
	stat, err := stats.GetStats()
	if err != nil {
		errors <- err
		return
	}

	_ = conn.Publish(reply, *stat)
}

// Responds with brief system stats (system.GetStatsBrief).
func HandlerStatsBrief(subject string, reply string, _ interface{}) {
	stat, err := stats.GetStatsBrief()
	if err != nil {
		errors <- err
		return
	}

	_ = conn.Publish(reply, *stat)
}

// Starts an aggregated measurement and sends back brief stats.
func HandlerAggregatedMeasurement(subject string, reply string, amr *types.AggregatedMeasurementRequest) {
	for _, id := range amr.Sensors {
		if id == system.ID() {
			log.Debugf("got measurement request: %#v\n", amr)
			go sensor.WithAggregated(*amr).Run()
			HandlerStatsBrief("", reply, nil)
			return
		}
	}
}

// Starts a raw measurement and sends back brief stats.
func HandlerRawMeasurement(subject string, reply string, rmr *types.RawMeasurementRequest) {
	for _, id := range rmr.Sensors {
		if id == system.ID() {
			log.Debugf("got measurement request: %#v\n", rmr)
			go sensor.WithRaw(*rmr).Run()
			HandlerStatsBrief("", reply, nil)
			return
		}
	}
}
