package nats

import (
	nats "github.com/nats-io/nats.go"
	"github.com/openrfsense/node/sensor"
	"github.com/openrfsense/node/system"
)

type Error struct {
	SensorID string `json:"sensorId"`
	Error    error  `json:"error"`
}

type Output struct {
	SensorID string `json:"sensorId"`
	Output   string `json:"output"`
}

// Waits for sensor manager errors or command output and sends a simple
// identifiable message on the proper channel.
func sendManagerData(conn *nats.EncodedConn, errChan chan<- error) {
	for {
		select {
		case sensorErr := <-sensor.Err():
			pubErr := conn.Publish("node.all.error", Error{
				SensorID: system.ID(),
				Error:    sensorErr,
			})
			if pubErr != nil {
				errChan <- pubErr
			}
		case output := <-sensor.Output():
			pubErr := conn.Publish("node.all.output", Output{
				SensorID: system.ID(),
				Output:   output,
			})
			if pubErr != nil {
				errChan <- pubErr
			}
		}
	}
}

// Simple consumer which logs error received on a channel.
func errorLogger(errChan <-chan error) {
	for {
		err := <-errChan
		log.Error(err)
	}
}
