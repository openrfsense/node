package sensor

import (
	"strings"

	"github.com/fatih/structs"
)

type CommandFlags struct {
	DevIndex        string `yaml:"devIndex" flag:"-d"`
	ClkOffset       string `yaml:"clkOffset" flag:"-c"`
	ClkCorrPeriod   string `yaml:"clkCorrPeriod" flag:"-k"`
	Gain            string `yaml:"gain" flag:"-g"`
	HoppingStrategy string `yaml:"hoppingStrategy" flag:"-y"`
	SampRate        string `yaml:"sampRate" flag:"-s"`
	Log2FFTsize     string `yaml:"log2FFTsize" flag:"-f"`
	FftBatchLen     string `yaml:"fftBatchLen" flag:"-b"`
	AvgFactor       string `yaml:"avgFactor" flag:"-a"`
	SOverlap        string `yaml:"soverlap" flag:"-o"`
	FreqOverlap     string `yaml:"freqOverlap" flag:"-q"`
	MonitorTime     string `yaml:"monitorTime" flag:"-t"`
	MinTimeRes      string `yaml:"minTimeRes" flag:"-r"`
	Window          string `yaml:"window" flag:"-w"`
	TcpCollector    string `yaml:"tcpCollector" flag:"-m"`
	SslCollector    string `yaml:"sslCollector" flag:"-n"`
	AbsoluteTime    string `yaml:"absoluteTime" flag:"-x"`
	MeasurementType string `yaml:"measurementType" flag:"-z"`
	Reserved        string `yaml:"reserved" flag:"-u"`

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
