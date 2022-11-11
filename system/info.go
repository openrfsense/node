package system

import (
	"bytes"
	"fmt"
	"os"
	"strings"
	"time"
)

// Tries to fetch the machine vendor and model.
func GetModel() string {
	if vendor, err := os.ReadFile("/sys/devices/virtual/dmi/id/board_vendor"); err == nil {
		ret := bytes.TrimSpace(vendor)
		if board, err := os.ReadFile("/sys/devices/virtual/dmi/id/board_name"); err == nil {
			ret = append(ret, byte(' '))
			ret = append(ret, bytes.TrimSpace(board)...)
		}
		return string(ret)
	}

	if name, err := os.ReadFile("/sys/devices/virtual/dmi/id/product_name"); err == nil {
		ret := bytes.TrimSpace(name)
		if version, err := os.ReadFile("/sys/devices/virtual/dmi/id/product_version"); err == nil {
			ret = append(ret, byte(' '))
			ret = append(ret, bytes.TrimSpace(version)...)
		}
		return string(ret)
	}

	if model, err := os.ReadFile("/sys/firmware/devicetree/base/model"); err == nil {
		// Raspberry Pis tend to have trailing zeroes in the model name
		model = bytes.Trim(model, "\u0000")
		return string(bytes.TrimSpace(model))
	}

	return "OpenRFSense Node"
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
