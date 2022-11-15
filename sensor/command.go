package sensor

import (
	"strings"

	"github.com/fatih/structs"
)

type CommandFlags struct {
	AvgFactor       string `yaml:"avgFactor" flag:"-a"`
	FftBatchLen     string `yaml:"fftBatchLen" flag:"-b"`
	ClkOffset       string `yaml:"clkOffset" flag:"-c"`
	DevIndex        string `yaml:"devIndex" flag:"-d"`
	Log2FFTsize     string `yaml:"log2FFTsize" flag:"-f"`
	Gain            string `yaml:"gain" flag:"-g"`
	ClkCorrPeriod   string `yaml:"clkCorrPeriod" flag:"-k"`
	SchemaFile      string `yaml:"schemaFile" flag:"-m"`
	SslCollector    string `yaml:"sslCollector" flag:"-n"`
	SOverlap        string `yaml:"soverlap" flag:"-o"`
	FreqOverlap     string `yaml:"freqOverlap" flag:"-q"`
	MinTimeRes      string `yaml:"minTimeRes" flag:"-r"`
	SampRate        string `yaml:"sampRate" flag:"-s"`
	MonitorTime     string `yaml:"monitorTime" flag:"-t"`
	Reserved        string `yaml:"reserved" flag:"-u"`
	Window          string `yaml:"window" flag:"-w"`
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
	HoppingStrategy: "sequential",
	MinTimeRes:      "0",
	ClkCorrPeriod:   "3600",
	SampRate:        "2400000",
	MonitorTime:     "0",
	Gain:            "32.8",
	MinFreq:         "24000000",
	MaxFreq:         "1766000000",
	Command:         "orfs_sensor",
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
