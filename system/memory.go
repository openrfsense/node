package system

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/openrfsense/common/stats"
)

// Type StatsMemory contains memory status information parsed from /proc/meminfo.
type StatsMemory struct {
	// The amount of physical RAM, in kibibytes, left unused by the system
	Free uint64 `json:"free"`

	// Total amount of usable RAM, in kibibytes, which is physical RAM minus a number of reserved bits and the kernel binary code
	Total uint64 `json:"total"`

	// Estimate of how much memory is available for starting new applications, without swapping
	Available uint64 `json:"available"`

	// The amount, in kibibytes, of temporary storage for raw disk blocks
	Buffers uint64 `json:"buffers"`

	// The amount of physical RAM, in kibibytes, used as cache memory
	Cached uint64 `json:"cached"`

	// The total amount of swap free, in kibibytes
	SwapFree uint64 `json:"swap_free"`

	// The total amount of swap available, in kibibytes
	SwapTotal uint64 `json:"swap_total"`
}

// providerMemory implements stats.Provider.
var _ stats.Provider = providerMemory{}

// Stats provider for memory information.
type providerMemory struct{}

func (providerMemory) Name() string {
	return "memory"
}

func (providerMemory) Stats() (interface{}, error) {
	data, err := os.ReadFile("/proc/meminfo")
	if err != nil {
		return nil, fmt.Errorf("%w: error reading memory info", err)
	}
	memInfos := strings.Split(string(data), "\n")

	res := &StatsMemory{}
	for _, l := range memInfos {
		fields := strings.Fields(l)
		if len(fields) < 2 {
			continue
		}
		tag := fields[0]
		val, err := strconv.ParseUint(fields[1], 10, 64)
		if err != nil {
			continue
		}
		switch tag {
		case "MemTotal:":
			res.Total = val
		case "MemFree:":
			res.Free = val
		case "MemAvailable:":
			res.Available = val
		case "Buffers:":
			res.Buffers = val
		case "Cached:":
			res.Cached = val
		case "SwapTotal:":
			res.SwapTotal = val
		case "SwapFree:":
			res.SwapFree = val
		}
	}
	return res, nil
}
