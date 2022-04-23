package mqtt

import (
	"encoding/json"

	emitter "github.com/emitter-io/go/v2"

	"github.com/openrfsense/node/system"
)

// Returns full system stats for this node.
func HandlerStats(m emitter.Message) ([]byte, error) {
	stats, err := system.GetStats()
	if err != nil {
		return nil, err
	}

	return json.Marshal(stats)
}

// Returns brief system stats (without any providers) for this node.
func HandlerStatsBrief(m emitter.Message) ([]byte, error) {
	stats, err := system.GetStatsBrief()
	if err != nil {
		return nil, err
	}

	return json.Marshal(stats)
}
