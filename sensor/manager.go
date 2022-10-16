package sensor

import "sync"

type Status string

const (
	Free  Status = "FREE"
	Busy  Status = "BUSY"
	Error Status = "ERROR"
)

type SensorManager struct {
	// General sensing status
	status Status

	// Current campaign ID (if running) or empty (if free)
	campaignId string

	// Flags to pass onto the es-sensor process
	flags CommandFlags

	sync.RWMutex
}

var manager *SensorManager

func (m *SensorManager) Status() Status {
	m.RLock()
	defer m.RUnlock()
	return m.status
}

func (m *SensorManager) CampaignId() string {
	m.RLock()
	defer m.RUnlock()
	return m.campaignId
}

func Manager() *SensorManager {
	return manager
}

func Init(flags CommandFlags) {
	if manager == nil {
		manager = &SensorManager{
			flags:      flags,
			status:     Free,
			campaignId: "",
		}
	}
}
