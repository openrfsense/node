package stats

import (
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

type providerLocation struct {
	LocationName string
	Latitude     float64
	Longitude    float64
	Elevation    float64
}

func (providerLocation) Name() string {
	return "location"
}

func (p providerLocation) Stats() (interface{}, error) {
	sl := &StatsLocation{
		Name:        p.LocationName,
		Type:        "Point",
		Coordinates: []float64{p.Longitude, p.Latitude},
		Elevation:   p.Elevation,
	}

	return sl, nil
}
