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

func BuildEsSensorCommand(sip CommandFlags) {
	// flags := generateFlags(sip)
	// cmd := exec.Command("es_sensor", flags...)
}
