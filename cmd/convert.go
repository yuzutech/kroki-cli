package cmd

import (
	"bufio"
	"fmt"
	"github.com/pkg/errors"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/yuzutech/kroki-go"
)

func Convert(cmd *cobra.Command, args []string) {
	filePath := args[0]
	graphFormat, err := cmd.Flags().GetString("type")
	if err != nil {
		exit(err)
	}
	imageFormat, err := cmd.Flags().GetString("format")
	if err != nil {
		exit(err)
	}
	if filePath == "-" {
		ConvertFromStdin(graphFormat, imageFormat)
	} else {
		ConvertFromFile(filePath, graphFormat, imageFormat)
	}
}

func ConvertFromStdin(diagramType string, imageFormat string) {
	if diagramType == "" {
		exit("diagram type must be specify using --type flag")
	}
	text, err := GetTextFromStdin()
	if err != nil {
		exit(err)
	}
	result, err := client.FromString(text, kroki.Graphviz, kroki.ImageFormat(imageFormat))
	if err != nil {
		exit(err)
	}
	fmt.Println(result)
}

func GetTextFromStdin() (result string, err error) {
	reader := bufio.NewReader(os.Stdin)
	input, err := ioutil.ReadAll(reader)
	return string(input), err
}

func ConvertFromFile(filePath string, graphFormatRaw string, imageFormatRaw string) {
	graphFormat, err := ResolveGraphFormat(graphFormatRaw, filePath)
	if err != nil {
		exit(err)
	}
	imageFormat, err := ResolveImageFormat(imageFormatRaw)
	if err != nil {
		exit(err)
	}
	result, err := client.FromFile(filePath, graphFormat, imageFormat)
	if err != nil {
		exit(err)
	}
	err = client.WriteToFile(OutputFilePath(filePath, imageFormat), result)
	if err != nil {
		exit(err)
	}
}

func OutputFilePath(filePath string, imageFormat kroki.ImageFormat) string {
	fileExtension := path.Ext(filePath)
	return filePath[0:len(filePath)-len(fileExtension)] + "." + string(imageFormat)
}

func ResolveImageFormat(imageFormatRaw string) (result kroki.ImageFormat, err error) {
	value := strings.ToLower(imageFormatRaw)
	switch value {
	case "svg":
		return kroki.Svg, nil
	case "png":
		return kroki.ImageFormat("png"), nil
	case "jpeg":
		return kroki.ImageFormat("jpeg"), nil
	case "pdf":
		return kroki.ImageFormat("pdf"), nil
	default:
		return kroki.ImageFormat(""), errors.Errorf(
			"invalid image format %s.",
			value)
	}
}

func ResolveGraphFormat(graphFormatRaw string, filePath string) (result kroki.GraphFormat, err error) {
	if graphFormatRaw == "" {
		return GraphFormatFromFile(filePath)
	} else {
		return GraphFormatFromValue(graphFormatRaw)
	}
}

func GraphFormatFromValue(value string) (result kroki.GraphFormat, err error) {
	value = strings.ToLower(value)
	switch value {
	case "dot", "graphviz":
		return kroki.Graphviz, nil
	case "plantuml":
		return kroki.Plantuml, nil
	case "nomnoml":
		return kroki.Nomnoml, nil
	case "blockdiag":
		return kroki.BlockDiag, nil
	case "mermaid":
		return kroki.Mermaid, nil
	case "svgbob":
		return kroki.Svgbob, nil
	case "umlet":
		return kroki.Umlet, nil
	case "c4plantuml":
		return kroki.C4plantuml, nil
	case "seqdiag":
		return kroki.SeqDiag, nil
	case "erd":
		return kroki.GraphFormat("erd"), nil
	case "nwdiag":
		return kroki.GraphFormat("nwdiag"), nil
	case "actdiag":
		return kroki.GraphFormat("actdiag"), nil
	case "ditaa":
		return kroki.GraphFormat("ditaa"), nil
	default:
		return kroki.GraphFormat(""), errors.Errorf(
			"invalid graph format %s.",
			value)
	}
}

func GraphFormatFromFile(filePath string) (result kroki.GraphFormat, err error) {
	fileExtension := filepath.Ext(filePath)
	value := strings.ToLower(fileExtension)
	switch value {
	case ".dot", ".gv", ".graphviz":
		return kroki.Graphviz, nil
	case ".puml", ".plantuml":
		return kroki.Plantuml, nil
	case ".nomnoml":
		return kroki.Nomnoml, nil
	case ".blockdiag":
		return kroki.BlockDiag, nil
	case ".mermaid":
		return kroki.Mermaid, nil
	case ".svgbob":
		return kroki.Svgbob, nil
	case ".umlet":
		return kroki.Umlet, nil
	case ".c4puml", ".c4", ".c4plantuml":
		return kroki.C4plantuml, nil
	case ".seqdiag":
		return kroki.SeqDiag, nil
	case ".erd":
		return kroki.GraphFormat("erd"), nil
	case ".nwdiag":
		return kroki.GraphFormat("nwdiag"), nil
	case ".actdiag":
		return kroki.GraphFormat("actdiag"), nil
	case ".ditaa":
		return kroki.GraphFormat("ditaa"), nil
	default:
		return kroki.GraphFormat(""), errors.Errorf(
			"unable to infer the graph format from the file extension %s, please specify the diagram type using --type flag.",
			value)
	}
}
