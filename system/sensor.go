package system

import (
	"github.com/openrfsense/node/sensor"

	"github.com/openrfsense/common/stats"
)

// Type StatsSensor reports sensor status information (busy/free, campaign ID, etc.).
type StatsSensor struct {
	// General sensing status
	Status sensor.Status `json:"status"`

	// Current campaign ID (if running) or null (if free)
	CampaignId string `json:"campaignId,omitempty"`
}

// providerSensor implements stats.Provider.
var _ stats.Provider = providerSensor{}

// Stats provider for sensor information.
type providerSensor struct{}

func (providerSensor) Name() string {
	return "sensor"
}

func (providerSensor) Stats() (interface{}, error) {
	sm := sensor.Manager()

	return StatsSensor{
		Status:     sm.Status(),
		CampaignId: sm.CampaignId(),
	}, nil
}
