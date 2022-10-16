package system

import (
	"github.com/openrfsense/common/config"
	"github.com/openrfsense/common/stats"
)

// Type StatsLocation provides static location information as defined in the
// node local configuration. The type is GeoJSON compatible.
type StatsLocation struct {
	// Name of the location. Not required
	Name string `json:"name,omitempty"`

	// For the purpose of this project, this should always be "Point"
	Type string `json:"type"`

	// Coordinate array for the GeoJSON object. Longitude goes first as per spec
	Coordinates []float64 `json:"coordinates"`

	// Sensor elevation/altitude
	Elevation float64 `json:"elevation"`
}

// providerLocation implements stats.Provider
var _ stats.Provider = providerLocation{}

type providerLocation struct{}

func (providerLocation) Name() string {
	return "location"
}

// TODO: make location required/sanitize configuration on start
// Web UI pre check?
func (providerLocation) Stats() (interface{}, error) {
	lat := config.GetOrDefault("node.location.latitude", 0.0)
	long := config.GetOrDefault("node.location.longitude", 0.0)
	alt := config.GetOrDefault("node.location.elevation", 0.0)

	sl := &StatsLocation{
		Name:        config.Get[string]("node.location.name"),
		Type:        "Point",
		Coordinates: []float64{long, lat},
		Elevation:   alt,
	}

	return sl, nil
}
