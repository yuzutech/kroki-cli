package cmd

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
		w.Write([]byte("<svg>Hello</svg>"))
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
		w.Write([]byte("<svg>Hello</svg>"))
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
	result, err := ioutil.ReadFile(outFilePath)
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
		io.Copy(&buf, reader)
		out <- buf.String()
	}()
	wg.Wait()
	f()
	writer.Close()
	return <-out
}