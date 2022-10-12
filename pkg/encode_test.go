package pkg

import (
	"bytes"
	"testing"
)

func TestEncodeFromReader(t *testing.T) {
	buf := bytes.NewBuffer([]byte(""))
	buf.Write([]byte("digraph G {Hello->World}"))
	result := CaptureOutput(func() {
		EncodeFromReader(buf)
	})
	expected := "eNpKyUwvSizIUHBXqPZIzcnJ17ULzy_KSakFBAAA__9sQAjG\n"
	if result != expected {
		t.Errorf("EncodeFromReader error\nexpected: %s\nactual:   %s", expected, result)
	}
}
