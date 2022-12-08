package pkg

import (
	"bufio"
	"fmt"
	"io"
	"os"

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
	content, err := os.ReadFile(filePath)
	if err != nil {
		exit(fmt.Errorf("fail to read file %s: %w", filePath, err))
	}
	input := string(content)
	result, err := kroki.CreatePayload(input)
	if err != nil {
		exit(err)
	}
	fmt.Println(result)
}
