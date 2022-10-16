package sensor

import (
	"strings"
	"testing"
)

func TestGenerateFlags(t *testing.T) {
	in := CommandFlags{
		Command:  "scan",
		MinFreq:  "160000000",
		MaxFreq:  "180000000",
		DevIndex: "1",
	}

	result := strings.Join(generateFlags(in), " ")
	expected := "scan -d 1 160000000 180000000"

	if result != expected {
		t.Fail()
	}
}
