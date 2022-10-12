package pkg

import (
	"bufio"
	"bytes"
	"compress/zlib"
	"encoding/base64"
	"fmt"
	"io"
	"net/url"
	"os"
	"path"
	"strings"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

func Decode(_ *cobra.Command, args []string) {
	input := args[0]
	if input == "-" {
		reader := bufio.NewReader(os.Stdin)
		DecodeFromReader(reader)
	} else {
		DecodeFromInput(input)
	}
}

func DecodeFromReader(reader io.Reader) {
	text, err := GetTextFromReader(reader)
	if err != nil {
		exit(err)
	}
	result, err := DecodeInput(text)
	if err != nil {
		exit(err)
	}
	fmt.Println(result)
}

func DecodeFromInput(input string) {
	result, err := DecodeInput(input)
	if err != nil {
		exit(err)
	}
	fmt.Println(result)
}

// takes a string encoded using deflate + base64 format and returns a decoded string
func DecodeInput(input string) (string, error) {
	// special case to extract the encoded diagram from a URL (GET request)
	// expected format is: https://kroki.io/diagram/format/encoded
	if strings.HasPrefix(input, "https://") || strings.HasPrefix(input, "http://") {
		inputUrl, _ := getUrl(input)
		if inputUrl != nil {
			// get the last part of the URL
			input = path.Base(inputUrl.Path)
		}
	}
	result, err := base64.URLEncoding.DecodeString(input)
	if err != nil {
		return "", errors.Wrap(err, "fail to decode the input")
	}
	reader, err := zlib.NewReader(bytes.NewReader(result))
	if err != nil {
		return "", errors.Wrap(err, "fail to create the reader")
	}
	out := new(strings.Builder)
	_, _ = io.Copy(out, reader)
	return out.String(), nil
}

func getUrl(input string) (*url.URL, error) {
	_, err := url.ParseRequestURI(input)
	if err != nil {
		return nil, err
	}

	u, err := url.Parse(input)
	if err != nil || u.Scheme == "" || u.Host == "" {
		return u, errors.New("invalid URL")
	}

	return u, nil
}
