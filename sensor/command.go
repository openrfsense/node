package sensor

import (
	"strings"

	"github.com/fatih/structs"
)

type CommandFlags struct {
	AvgFactor       string `yaml:"averagingFactor" flag:"-a"`
	FftBatchLen     string `yaml:"fftBatchLength" flag:"-b"`
	ClkOffset       string `yaml:"clockOffset" flag:"-c"`
	DevIndex        string `yaml:"devIndex" flag:"-d"`
	Log2FFTsize     string `yaml:"log2FFTsize" flag:"-f"`
	Gain            string `yaml:"gain" flag:"-g"`
	ClkCorrPeriod   string `yaml:"clockCorrectionPeriod" flag:"-k"`
	SchemaFile      string `yaml:"schemaFile" flag:"-m"`
	SslCollector    string `yaml:"sslCollector" flag:"-n"`
	SOverlap        string `yaml:"segmentOverlap" flag:"-o"`
	FreqOverlap     string `yaml:"frequencyOverlap" flag:"-q"`
	MinTimeRes      string `yaml:"minTimeResolution" flag:"-r"`
	SampRate        string `yaml:"samplingRate" flag:"-s"`
	MonitorTime     string `yaml:"monitorTime" flag:"-t"`
	Window          string `yaml:"windowingFunction" flag:"-w"`
	AbsoluteTime    string `yaml:"absoluteTime" flag:"-x"`
	HoppingStrategy string `yaml:"hoppingStrategy" flag:"-y"`
	MeasurementType string `yaml:"measurementType" flag:"-z"`

	// These should be ignored by koanf even if found in the configuration
	SensorId   string
	CampaignId string
	MinFreq    string
	MaxFreq    string
	Command    string
}

var DefaultFlags = CommandFlags{
	Command: "orfs_sensor",
}

func generateFlags(sip CommandFlags) []string {
	ret := []string{}

	// Flagged arguments
	for _, f := range structs.Fields(sip) {
		flag := f.Tag("flag")
		// Skip field if it doesn't have a flag (prefix and suffix)
		if flag == "" {
			continue
		}

		value := f.Value().(string)
		// Only consider non-empty values
		if strings.TrimSpace(value) != "" {
			ret = append(ret, flag, value)
		}
	}

	// Suffixed arguments
	ret = append(ret, sip.SensorId, sip.CampaignId, sip.MinFreq, sip.MaxFreq)

	return ret
}
