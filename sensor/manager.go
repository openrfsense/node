package sensor

import (
	"fmt"
	"os/exec"
	"strconv"
	"sync"
	"time"

	"github.com/openrfsense/common/types"

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

	// Flags to pass onto the es-sensor process
	flags CommandFlags

	// Last command output, if any
	output string

	// Last command error, if any
	err error

	sync.RWMutex
}

var manager *sensorManager

// Starts the actual process.
func (m *sensorManager) Run() {
	flagsSlice := generateFlags(m.flags)

	go func(m *sensorManager) {
		time.Sleep(time.Until(m.begin))

		cmd := exec.Command("orfs_sensor", flagsSlice...)
		m.Lock()
		m.status = Busy
		m.Unlock()

		// TODO: send output to backend via NATS
		output, err := cmd.CombinedOutput()

		m.Lock()
		defer m.Unlock()
		m.output = string(output)
		m.err = err
		if err == nil {
			m.status = Free
		} else {
			m.status = Error
		}
	}(m)
}

// Returns the current campaign ID (assigned by the backend).
func CampaignId() string {
	manager.RLock()
	defer manager.RUnlock()
	return manager.campaignId
}

// Returns the CommandFlags object used to start the campaign.
func Flags() CommandFlags {
	return manager.flags
}

// Returns the command output for the current campaign.
func Output() string {
	manager.RLock()
	defer manager.RUnlock()
	return manager.output
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
		}
	}

	// Initialize TCP collector to the one described in the configuration
	manager.flags.TcpCollector = fmt.Sprintf(
		"%s:%d",
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
	manager.begin = amr.Begin

	if len(flags) > 0 {
		manager.flags = flags[0]
	}

	monitorTime := amr.End.Unix() - amr.Begin.Unix()
	manager.flags.MonitorTime = strconv.FormatInt(monitorTime, 10)

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
	manager.begin = rmr.Begin

	if len(flags) > 0 {
		manager.flags = flags[0]
	}

	monitorTime := rmr.End.Unix() - rmr.Begin.Unix()
	manager.flags.MonitorTime = strconv.FormatInt(monitorTime, 10)

	return manager
}
