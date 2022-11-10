package system

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/knadh/koanf"

	"github.com/openrfsense/common/logging"
	"github.com/openrfsense/common/stats"
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

// Tries to fetch the machine vendor and model.
func GetModel() string {
	ret := ""

	if vendor, err := os.ReadFile("/sys/devices/virtual/dmi/id/board_vendor"); err == nil {
		ret = strings.TrimSpace(string(vendor))
		if name, err := os.ReadFile("/sys/devices/virtual/dmi/id/board_name"); err == nil {
			ret += " " + strings.TrimSpace(string(name))
		}
		return ret
	}

	if name, err := os.ReadFile("/sys/devices/virtual/dmi/id/product_name"); err == nil {
		ret = strings.TrimSpace(string(name))
		if version, err := os.ReadFile("/sys/devices/virtual/dmi/id/product_version"); err == nil {
			ret += " " + strings.TrimSpace(string(version))
		}
		return ret
	}

	if model, err := os.ReadFile("/sys/firmware/devicetree/base/model"); err == nil {
		return strings.TrimSpace(string(model))
	}

	return ret
}

// Returns system uptime in milliseconds.
func GetUptime() (time.Duration, error) {
	data, err := os.ReadFile("/proc/uptime")
	if err != nil {
		return 0, fmt.Errorf("%w, error reading uptime", err)
	}

	timeString := strings.Fields(string(data))[0]
	uptime, err := time.ParseDuration(timeString + "s")
	if err != nil {
		return 0, fmt.Errorf("%w, parsing /proc/uptime", err)
	}

	return uptime, nil
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

	uptime, err := GetUptime()
	if err != nil {
		return nil, err
	}

	s := &stats.Stats{
		ID:       ID(),
		Hostname: hostname,
		Model:    GetModel(),
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
