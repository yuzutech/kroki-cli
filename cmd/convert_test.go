package cmd

import (
	"testing"

	"github.com/yuzutech/kroki-go"
)

func TestResolveOutputFilePath(t *testing.T) {
	cases := []struct {
		filePath    string
		outFile     string
		imageFormat kroki.ImageFormat
		expected    string
	}{
		{
			filePath:    "/path/to/hello.dot",
			outFile:     "",
			imageFormat: kroki.Svg,
			expected:    "/path/to/hello.svg",
		},
		{
			filePath:    "/path/dot/hello.dot",
			outFile:     "",
			imageFormat: kroki.Svg,
			expected:    "/path/dot/hello.svg",
		},
		{
			filePath:    "hello.dot.puml",
			outFile:     "",
			imageFormat: kroki.ImageFormat("png"),
			expected:    "hello.dot.png",
		},
		{
			filePath:    "hello.dot.puml",
			outFile:     "out.png",
			imageFormat: kroki.ImageFormat("png"),
			expected:    "out.png",
		},
	}
	for _, c := range cases {

		result := ResolveOutputFilePath(c.outFile, c.filePath, c.imageFormat)
		if result != c.expected {
			t.Errorf("ResolveOutputFilePath error\nexpected: %s\nactual:   %s", c.expected, result)
		}
	}
}

func TestResolveImageFormat(t *testing.T) {
	cases := []struct {
		imageFormatRaw string
		outFile        string
		expected       kroki.ImageFormat
	}{
		{
			imageFormatRaw: "",
			outFile:        "",
			expected:       kroki.Svg, // default value
		},
		{
			imageFormatRaw: "SVG",
			outFile:        "",
			expected:       kroki.Svg,
		},
		{
			imageFormatRaw: "SVG",
			outFile:        "out.png",
			expected:       kroki.Svg, // --format flag has priority over output file extension
		},
		{
			imageFormatRaw: "PNG",
			outFile:        "out.png",
			expected:       kroki.ImageFormat("png"),
		},
		{
			imageFormatRaw: "",
			outFile:        "out.png",
			expected:       kroki.ImageFormat("png"),
		},
		{
			imageFormatRaw: "",
			outFile:        "out.dot.jpg",
			expected:       kroki.ImageFormat("jpeg"),
		},
		{
			imageFormatRaw: "txt",
			outFile:        "",
			expected:       "",
		},
		{
			imageFormatRaw: "",
			outFile:        "out.txt",
			expected:       "",
		},
		{
			imageFormatRaw: "jpeg",
			outFile:        "out.txt",
			expected:       kroki.ImageFormat("jpeg"), // output file extension is ignored, use format flag (jpeg)
		},
	}
	for _, c := range cases {
		result, _ := ResolveImageFormat(c.imageFormatRaw, c.outFile)
		if result != c.expected {
			t.Errorf("ResolveImageFormat error\nexpected: %s\nactual:   %s", c.expected, result)
		}
	}
}
