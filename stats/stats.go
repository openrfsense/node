package stats

import (
	"os"

	"github.com/knadh/koanf"

	"github.com/openrfsense/common/logging"
	"github.com/openrfsense/common/stats"
	"github.com/openrfsense/node/system"
)

var staticLocation providerLocation

var log = logging.New().
	WithPrefix("system").
	WithLevel(logging.DebugLevel).
	WithFlags(logging.FlagsDevelopment)

// Initializes the reporting system, mainly fetches and copies in memory static information
// needed by providers.
func Init(config *koanf.Koanf) {
	staticLocation = providerLocation{
		LocationName: config.String("location.name"),
		Elevation:    config.MustFloat64("location.elevation"),
		Latitude:     config.MustFloat64("location.latitude"),
		Longitude:    config.MustFloat64("location.longitude"),
	}
}

// Returns full system stats.
func GetStats() (*stats.Stats, error) {
	s, err := GetStatsBrief()
	if err != nil {
		return nil, err
	}

	// Provider errors are logged but not propagated since they
	// are just extra information
	err = s.Provide(
		providerMemory{},
		providerFs{},
		providerNetwork{},
	)
	if err != nil {
		log.Error(err)
	}

	return s, nil
}

// Returns brief system stats, enough to identify the machine. For more in-depth metrics, use GetStats.
func GetStatsBrief() (*stats.Stats, error) {
	hostname, err := os.Hostname()
	if err != nil {
		return nil, err
	}

	uptime, err := system.GetUptime()
	if err != nil {
		return nil, err
	}

	s := &stats.Stats{
		ID:       system.ID(),
		Hostname: hostname,
		Model:    system.GetModel(),
		Uptime:   uptime,
	}

	// Provider errors are logged but not propagated since they
	// are just extra information
	err = s.Provide(
		staticLocation,
		providerSensor{},
	)
	if err != nil {
		log.Error(err)
	}

	return s, nil
}
