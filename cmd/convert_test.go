package cmd

import (
	"testing"

	"github.com/yuzutech/kroki-go"
)

func TestOutputFilePath(t *testing.T) {

	cases := []struct {
		filePath    string
		imageFormat kroki.ImageFormat
		expected    string
	}{
		{
			filePath:    "/path/to/hello.dot",
			imageFormat: kroki.Svg,
			expected:    "/path/to/hello.svg",
		},
		{
			filePath:    "/path/dot/hello.dot",
			imageFormat: kroki.Svg,
			expected:    "/path/dot/hello.svg",
		},
		{
			filePath:    "hello.dot.puml",
			imageFormat: kroki.ImageFormat("png"),
			expected:    "hello.dot.png",
		},
	}
	for _, c := range cases {

		result := OutputFilePath(c.filePath, c.imageFormat)
		if result != c.expected {
			t.Errorf("OutputFilePath error\nexpected: %s\nactual:   %s", c.expected, result)
		}
	}
}
