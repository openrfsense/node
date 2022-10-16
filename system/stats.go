package system

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/openrfsense/common/stats"
)

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

	return uptime / time.Millisecond, nil
}

// Returns full system stats.
func GetStats() (*stats.Stats, error) {
	s, err := GetStatsBrief()
	if err != nil {
		return nil, err
	}

	err = s.Provide(
		providerLocation{},
		providerMemory{},
		providerFs{},
		providerNetwork{},
	)
	if err != nil {
		return nil, err
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

	err = s.Provide(
		providerSensor{},
		providerLocation{},
	)
	if err != nil {
		return nil, err
	}

	return s, nil
}
