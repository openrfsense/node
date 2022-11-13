package sensor

import (
	"context"
	"fmt"
	"os/exec"
	"strconv"
	"sync"
	"time"

	"github.com/openrfsense/common/logging"
	"github.com/openrfsense/common/types"
	"github.com/openrfsense/node/system"

	"github.com/knadh/koanf"
)

// Type StatusEnum describes the current status of the sensor
type StatusEnum string

const (
	Free  StatusEnum = "FREE"
	Busy  StatusEnum = "BUSY"
	Error StatusEnum = "ERROR"
)

// Type sensorManager holds the necessary information to manage an es_sensor
// process and run a campaign, reporting eventual errors and command output
type sensorManager struct {
	// General sensing status
	status StatusEnum

	// Current campaign ID (if running) or empty (if free)
	campaignId string

	// Campaign start datetime
	begin time.Time

	// Campaign end datetime
	end time.Time

	// Flags to pass onto the es-sensor process
	flags CommandFlags

	// Last command output, if any
	output chan string

	// Last command error, if any
	err chan error

	sync.RWMutex
}

var manager *sensorManager

var log = logging.New().
	WithPrefix("sensor").
	WithLevel(logging.DebugLevel).
	WithFlags(logging.FlagsDevelopment)

// Starts the actual process.
func (m *sensorManager) Run() {
	m.flags.CampaignId = m.campaignId
	m.flags.SensorId = system.ID()
	flagsSlice := generateFlags(m.flags)

	log.Debugf("starting manager: %#v", m)

	// time.Sleep(time.Until(m.begin))
	log.Debugf("starting campaign %s", m.campaignId)

	ctx, cancel := context.WithDeadline(context.Background(), m.end)
	defer cancel()
	cmd := exec.CommandContext(ctx, m.flags.Command, flagsSlice...)
	log.Debug(cmd.String())

	m.status = Busy
	// TODO: send output to backend via NATS
	output, err := cmd.CombinedOutput()
	log.Debug(string(output))
	m.output <- string(output)

	m.Lock()
	m.campaignId = ""
	if err != nil {
		log.Error(err)
		m.err <- err
		m.status = Error
	} else {
		m.status = Free
	}
	m.Unlock()
}

// Open channel where command output is sent after completion.
func Output() <-chan string {
	return manager.output
}

// Open channel where command errors are sent after completion.
func Err() <-chan error {
	return manager.err
}

// Returns the current campaign ID (assigned by the backend).
func CampaignId() string {
	manager.RLock()
	defer manager.RUnlock()
	return manager.campaignId
}

// Returns the current status of the sensor manager.
func Status() StatusEnum {
	manager.RLock()
	defer manager.RUnlock()
	return manager.status
}

// Initializes a SensorManager singleton. Also loads default command line flags
// from the configuration.
func Init(config *koanf.Koanf) error {
	err := config.Unmarshal("node.sensor", &DefaultFlags)
	if err != nil {
		return err
	}

	if manager == nil {
		manager = &sensorManager{
			flags:      DefaultFlags,
			status:     Free,
			campaignId: "",
			output:     make(chan string, 1),
			err:        make(chan error, 1),
		}
	}

	// Initialize TCP collector to the one described in the configuration
	manager.flags.SslCollector = fmt.Sprintf(
		"%s:%d#",
		config.String("collector.host"),
		config.MustInt("collector.port"),
	)

	return nil
}

// Starts an aggregated measurement campaign by running orfs_sensor with the given flags.
func WithAggregated(amr types.AggregatedMeasurementRequest, flags ...CommandFlags) *sensorManager {
	manager.Lock()
	defer manager.Unlock()

	manager.campaignId = amr.CampaignId
	manager.begin = time.Unix(amr.Begin, 0)
	manager.end = time.Unix(amr.End, 0)

	if len(flags) > 0 {
		manager.flags = flags[0]
	}

	monitorTime := amr.End - amr.Begin
	manager.flags.MonitorTime = strconv.FormatInt(monitorTime, 10)

	// Set type-specific command parameters
	manager.flags.MeasurementType = "PSD"
	manager.flags.MinFreq = strconv.FormatInt(amr.FreqMin, 10)
	manager.flags.MaxFreq = strconv.FormatInt(amr.FreqMax, 10)
	manager.flags.MinTimeRes = strconv.FormatInt(amr.TimeRes, 10)

	return manager
}

// Starts a raw measurement campaign by running orfs_sensor with the given flags
func WithRaw(rmr types.RawMeasurementRequest, flags ...CommandFlags) *sensorManager {
	manager.Lock()
	defer manager.Unlock()

	manager.campaignId = rmr.CampaignId
	manager.begin = time.Unix(rmr.Begin, 0)
	manager.end = time.Unix(rmr.End, 0)

	if len(flags) > 0 {
		manager.flags = flags[0]
	}

	monitorTime := rmr.End - rmr.Begin
	manager.flags.MonitorTime = strconv.FormatInt(monitorTime, 10)

	// Set type-specific command parameters
	manager.flags.MeasurementType = "IQ"
	manager.flags.MinFreq = fmt.Sprint(rmr.FreqCenter)
	manager.flags.MaxFreq = fmt.Sprint(rmr.FreqCenter)

	return manager
}
