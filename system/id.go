package system

import (
	"github.com/openrfsense/common/id"
)

var systemId string

// Returns (or generates if needed) the 23-character ID for this node using
// the board's vendor-model string (as reported by SysFS) as seed.
func ID() string {
	if systemId != "" {
		return systemId
	}

	model := GetModel()
	systemId = id.GenerateFromBytes([]byte(model), 23)
	return systemId
}
