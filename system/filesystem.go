package system

import (
	"fmt"
	"os"
	"strings"
	"syscall"

	"github.com/openrfsense/common/stats"
)

// Type StatsFS contains a single disk/mount information parsed from /proc/self/mountinfo.
type StatsFS struct {
	// Disk device name
	Device string `json:"device"`

	// Mount type
	Type string `json:"type"`

	// Mount path
	Mount string `json:"mount"`

	// Free space
	Free uint64 `json:"free"`

	// Available space
	Available uint64 `json:"available"`

	// Total disk size
	Size uint64 `json:"size"`
}

// providerFs implements stats.Provider.
var _ stats.Provider = providerFs{}

// Stats provider for network information.
type providerFs struct{}

func (pn providerFs) Name() string {
	return "fs"
}

func (pn providerFs) Stats() (interface{}, error) {
	data, err := os.ReadFile("/proc/self/mountinfo")
	if err != nil {
		return nil, fmt.Errorf("%w: error detecting mounted filesystems", err)
	}
	mountInfos := strings.Split(string(data), "\n")

	res := []*StatsFS{}
	for _, mountInfo := range mountInfos {
		fields := strings.Fields(mountInfo)
		if len(fields) < 5 {
			continue
		}

		mount := fields[4]

		statfs := &syscall.Statfs_t{}
		if err := syscall.Statfs(mount, statfs); err != nil {
			// ignore error
			continue
		}

		fs := &StatsFS{
			Mount:     mount,
			Type:      fields[8],
			Device:    fields[9],
			Free:      statfs.Bfree * uint64(statfs.Bsize),
			Available: statfs.Bavail * uint64(statfs.Bsize),
			Size:      statfs.Blocks * uint64(statfs.Bsize),
		}
		res = append(res, fs)
	}

	return res, nil
}
