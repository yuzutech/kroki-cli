package pkg

import (
	"bytes"
	"testing"
)

func TestDecodeFromReader(t *testing.T) {
	cases := []struct {
		input    string
		expected string
	}{
		{
			input:    "eNpKyUwvSizIUHBXqPZIzcnJ17ULzy_KSakFBAAA__9sQAjG",
			expected: "digraph G {Hello->World}\n",
		},
		{
			input:    "https://kroki.io/graphviz/svg/eNpKyUwvSizIUHBXqPZIzcnJ17ULzy_KSakFBAAA__9sQAjG",
			expected: "digraph G {Hello->World}\n",
		},
		{
			input:    "http://localhost:8000/graphviz/svg/eNpKyUwvSizIUHBXqPZIzcnJ17ULzy_KSakFBAAA__9sQAjG",
			expected: "digraph G {Hello->World}\n",
		},
	}
	for _, c := range cases {
		buf := bytes.NewBuffer([]byte(""))
		buf.Write([]byte(c.input))
		result := CaptureOutput(func() {
			DecodeFromReader(buf)
		})
		if result != c.expected {
			t.Errorf("DecodeFromReader error\nexpected: %s\nactual:   %s", c.expected, result)
		}
	}
}
