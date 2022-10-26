package sensor

import (
	"strings"

	"github.com/fatih/structs"
)

type CommandFlags struct {
	DevIndex        string `koanf:"devIndex" flag:"-d"`
	ClkOffset       string `koanf:"clkOffset" flag:"-c"`
	ClkCorrPeriod   string `koanf:"clkCorrPeriod" flag:"-k"`
	Gain            string `koanf:"gain" flag:"-g"`
	HoppingStrategy string `koanf:"hoppingStrategy" flag:"-y"`
	SampRate        string `koanf:"sampRate" flag:"-s"`
	Log2FFTsize     string `koanf:"log2FFTsize" flag:"-f"`
	FftBatchLen     string `koanf:"fftBatchLen" flag:"-b"`
	AvgFactor       string `koanf:"avgFactor" flag:"-a"`
	SOverlap        string `koanf:"soverlap" flag:"-o"`
	FreqOverlap     string `koanf:"freqOverlap" flag:"-q"`
	MonitorTime     string `koanf:"monitorTime" flag:"-t"`
	MinTimeRes      string `koanf:"minTimeRes" flag:"-r"`
	Window          string `koanf:"window" flag:"-w"`
	TcpCollector    string `koanf:"tcpCollector" flag:"-m"`
	SslCollector    string `koanf:"sslCollector" flag:"-n"`
	AbsoluteTime    string `koanf:"absoluteTime" flag:"-x"`
	MeasurementType string `koanf:"measurementType" flag:"-z"`
	Reserved        string `koanf:"reserved" flag:"-u"`

	// These should be ignored by koanf even if found in the configuration
	MinFreq string
	MaxFreq string
	Command string
}

var DefaultFlags = CommandFlags{
	ClkOffset:       "0",
	DevIndex:        "0",
	Log2FFTsize:     "8",
	FreqOverlap:     "0.167",
	HoppingStrategy: "sequential",
	AvgFactor:       "5",
	MinTimeRes:      "0",
	Window:          "hanning",
	ClkCorrPeriod:   "3600",
	FftBatchLen:     "10",
	SOverlap:        "128",
	SampRate:        "2400000",
	MonitorTime:     "0",
	Gain:            "32.8",
	MinFreq:         "24000000",
	MaxFreq:         "1766000000",
	Command:         "es_sensor",
}

func generateFlags(sip CommandFlags) []string {
	ret := []string{}

	// Prefixed arguments
	ret = append(ret, sip.Command)

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
	ret = append(ret, sip.MinFreq, sip.MaxFreq)

	return ret
}
