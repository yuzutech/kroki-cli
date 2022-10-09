package pkg

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/yuzutech/kroki-go"
)

func Encode(_ *cobra.Command, args []string) {
	filePath := args[0]
	if filePath == "-" {
		reader := bufio.NewReader(os.Stdin)
		EncodeFromReader(reader)
	} else {
		EncodeFromFile(filePath)
	}
}

func EncodeFromReader(reader io.Reader) {
	text, err := GetTextFromReader(reader)
	if err != nil {
		exit(err)
	}
	result, err := kroki.CreatePayload(text)
	if err != nil {
		exit(err)
	}
	fmt.Println(result)
}

func EncodeFromFile(filePath string) {
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		exit(errors.Wrapf(err, "fail to read file '%s'", filePath))
	}
	input := string(content)
	result, err := kroki.CreatePayload(input)
	if err != nil {
		exit(err)
	}
	fmt.Println(result)
}
