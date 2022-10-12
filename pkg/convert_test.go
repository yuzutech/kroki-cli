package pkg

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

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
			imageFormat: kroki.SVG,
			expected:    "/path/to/hello.svg",
		},
		{
			filePath:    "/path/dot/hello.dot",
			outFile:     "",
			imageFormat: kroki.SVG,
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

func TestGraphFormatFromFile(t *testing.T) {
	cases := []struct {
		filePath    string
		expected    kroki.DiagramType
	}{
		{
			filePath:    "/path/to/hello.dot",
			expected:    kroki.GraphViz,
		},
		{
			filePath:    "/path/to/hello.puml",
			expected:    kroki.PlantUML,
		},
		{
			filePath:    "/path/to/hello.plantuml",
			expected:    kroki.PlantUML,
		},
		{
			filePath:    "/path/to/hello.vega",
			expected:    kroki.Vega,
		},
		{
			filePath:    "/path/to/hello.vg",
			expected:    kroki.Vega,
		},
		{
			filePath:    "hello.vl",
			expected:    kroki.VegaLite,
		},
		{
			filePath:    "hello.c4",
			expected:    kroki.C4PlantUML,
		},
		{
			filePath:    "hello.wavedrom",
			expected:    kroki.WaveDrom,
		},
		{
			filePath:    "hello.bpmn",
			expected:    kroki.BPMN,
		},
		{
			filePath:    "hello.excalidraw",
			expected:    kroki.Excalidraw,
		},
		{
			filePath:    "hello.bytefield",
			expected:    kroki.Bytefield,
		},
	}
	for _, c := range cases {

		result, _ := GraphFormatFromFile(c.filePath)
		if result != c.expected {
			t.Errorf("GraphFormatFromFile error\nexpected: %s\nactual:   %s", c.expected, result)
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
			expected:       kroki.SVG, // default value
		},
		{
			imageFormatRaw: "SVG",
			outFile:        "",
			expected:       kroki.SVG,
		},
		{
			imageFormatRaw: "SVG",
			outFile:        "out.png",
			expected:       kroki.SVG, // --format flag has priority over output file extension
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

func TestGraphFormatFromValue(t *testing.T) {
	cases := []struct {
		diagramTypeRaw string
		expected       kroki.DiagramType
	}{
		{
			diagramTypeRaw: "VegaLite",
			expected:       kroki.VegaLite,
		},
		{
			diagramTypeRaw: "diagramsnet",
			expected:       kroki.DiagramType("diagramsnet"),
		},
		{
			diagramTypeRaw: "Structurizr",
			expected:       kroki.DiagramType("structurizr"),
		},
	}
	for _, c := range cases {
		result, _ := GraphFormatFromValue(c.diagramTypeRaw)
		if result != c.expected {
			t.Errorf("GraphFormatFromValue error\nexpected: %s\nactual:   %s", c.expected, result)
		}
	}
}

func TestConvertFromReader(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		expected := "eNpKyUwvSizIUHBXqPZIzcnJ17ULzy_KSakFBAAA__9sQAjG"
		uri := strings.Split(r.RequestURI, "/")
		payload := uri[len(uri)-1]
		if payload != expected {
			t.Errorf("ConvertFromReader error\nexpected: %s\nactual:   %s", expected, payload)
		}
		imageFormat := uri[len(uri)-2]
		if imageFormat != string("svg") {
			t.Errorf("ConvertFromReader error\nexpected: %s\nactual:   %s", "svg", imageFormat)
		}
		diagramType := uri[len(uri)-3]
		if diagramType != "graphviz" {
			t.Errorf("ConvertFromReader error\nexpected: %s\nactual:   %s", "graphviz", diagramType)
		}
		_, _ = w.Write([]byte("<svg>Hello</svg>"))
	}))
	defer ts.Close()
	port, err := strconv.ParseUint(strings.Split(ts.URL, ":")[2], 10, 16)
	if err != nil {
		t.Errorf("error getting the port :\n%+v", err)
	}
	client := kroki.New(kroki.Configuration{
		URL:     fmt.Sprintf("http://localhost:%d", port),
		Timeout: time.Second * 10,
	})
	buf := bytes.NewBuffer([]byte(""))
	buf.Write([]byte("digraph G {Hello->World}"))
	result := CaptureOutput(func() {
		ConvertFromReader(client, "dot", "svg", "-", buf)
	})
	expected := "<svg>Hello</svg>\n"
	if result != expected {
		t.Errorf("ConvertFromReader error\nexpected: %s\nactual:   %s", expected, result)
	}
}

func TestConvertFromReaderOutFile(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		expected := "eNpKyUwvSizIUHBXqPZIzcnJ17ULzy_KSakFBAAA__9sQAjG"
		uri := strings.Split(r.RequestURI, "/")
		payload := uri[len(uri)-1]
		if payload != expected {
			t.Errorf("ConvertFromReaderOutFile error\nexpected: %s\nactual:   %s", expected, payload)
		}
		imageFormat := uri[len(uri)-2]
		if imageFormat != string("svg") {
			t.Errorf("ConvertFromReaderOutFile error\nexpected: %s\nactual:   %s", "svg", imageFormat)
		}
		diagramType := uri[len(uri)-3]
		if diagramType != "graphviz" {
			t.Errorf("ConvertFromReaderOutFile error\nexpected: %s\nactual:   %s", "graphviz", diagramType)
		}
		_, _ = w.Write([]byte("<svg>Hello</svg>"))
	}))
	defer ts.Close()
	port, err := strconv.ParseUint(strings.Split(ts.URL, ":")[2], 10, 16)
	if err != nil {
		t.Errorf("error getting the port :\n%+v", err)
	}
	client := kroki.New(kroki.Configuration{
		URL:     fmt.Sprintf("http://localhost:%d", port),
		Timeout: time.Second * 10,
	})
	buf := bytes.NewBuffer([]byte(""))
	buf.Write([]byte("digraph G {Hello->World}"))
	outFilePath := "../tests/out.ignore.test.svg"
	defer os.Remove(outFilePath)
	ConvertFromReader(client, "dot", "", outFilePath, buf)
	result, _ := ioutil.ReadFile(outFilePath)
	expected := "<svg>Hello</svg>"
	if string(result) != expected {
		t.Errorf("ConvertFromReaderOutFile error\nexpected: %s\nactual:   %s", expected, string(result))
	}
}

func CaptureOutput(f func()) string {
	reader, writer, err := os.Pipe()
	if err != nil {
		panic(err)
	}
	stdout := os.Stdout
	stderr := os.Stderr
	defer func() {
		os.Stdout = stdout
		os.Stderr = stderr
		log.SetOutput(os.Stderr)
	}()
	os.Stdout = writer
	os.Stderr = writer
	log.SetOutput(writer)
	out := make(chan string)
	wg := new(sync.WaitGroup)
	wg.Add(1)
	go func() {
		var buf bytes.Buffer
		wg.Done()
		_, _ = io.Copy(&buf, reader)
		out <- buf.String()
	}()
	wg.Wait()
	f()
	_ = writer.Close()
	return <-out
}
