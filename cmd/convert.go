package cmd

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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
	outFile, err := cmd.Flags().GetString("out-file")
	if err != nil {
		exit(err)
	}
	client := GetClient(cmd)
	if filePath == "-" {
		reader := bufio.NewReader(os.Stdin)
		ConvertFromReader(client, graphFormat, imageFormat, outFile, reader)
	} else {
		ConvertFromFile(client, filePath, graphFormat, imageFormat, outFile)
	}
}

func ConvertFromReader(client kroki.Client, diagramTypeRaw string, imageFormatRaw string, outFile string, reader io.Reader) {
	if diagramTypeRaw == "" {
		exit("diagram type must be specify using --type flag")
	}
	diagramType, err := GraphFormatFromValue(diagramTypeRaw)
	if err != nil {
		exit(err)
	}
	imageFormat, err := ResolveImageFormat(imageFormatRaw, outFile)
	if err != nil {
		exit(err)
	}
	text, err := GetTextFromReader(reader)
	if err != nil {
		exit(err)
	}
	result, err := client.FromString(text, diagramType, imageFormat)
	if err != nil {
		exit(err)
	}
	if outFile == "" || outFile == "-" {
		fmt.Println(result)
	} else {
		err = client.WriteToFile(outFile, result)
		if err != nil {
			exit(err)
		}
	}
}

func GetTextFromReader(reader io.Reader) (result string, err error) {
	input, err := ioutil.ReadAll(reader)
	return string(input), err
}

func ConvertFromFile(client kroki.Client, filePath string, graphFormatRaw string, imageFormatRaw string, outFile string) {
	graphFormat, err := ResolveGraphFormat(graphFormatRaw, filePath)
	if err != nil {
		exit(err)
	}
	imageFormat, err := ResolveImageFormat(imageFormatRaw, outFile)
	if err != nil {
		exit(err)
	}
	result, err := client.FromFile(filePath, graphFormat, imageFormat)
	if err != nil {
		exit(err)
	}
	if outFile == "-" {
		fmt.Println(result)
	} else {
		err = client.WriteToFile(ResolveOutputFilePath(outFile, filePath, imageFormat), result)
		if err != nil {
			exit(err)
		}
	}
}

func ResolveOutputFilePath(outFile string, filePath string, imageFormat kroki.ImageFormat) string {
	if outFile != "" {
		return outFile
	}
	fileExtension := path.Ext(filePath)
	return filePath[0:len(filePath)-len(fileExtension)] + "." + string(imageFormat)
}

func ResolveImageFormat(imageFormatRaw string, outFile string) (kroki.ImageFormat, error) {
	if imageFormatRaw == "" {
		if outFile == "" || outFile == "-" {
			return kroki.SVG, nil
		}
		return ImageFormatFromFile(outFile)
	}
	return ImageFormatFromValue(imageFormatRaw)
}

func ImageFormatFromValue(imageFormatRaw string) (kroki.ImageFormat, error) {
	value := strings.ToLower(imageFormatRaw)
	switch value {
	case "svg":
		return kroki.SVG, nil
	case "png":
		return kroki.PNG, nil
	case "jpeg":
		return kroki.JPEG, nil
	case "pdf":
		return kroki.PDF, nil
	case "base64":
		return kroki.Base64, nil
	default:
		return kroki.ImageFormat(""), errors.Errorf(
			"invalid image format %s.",
			value)
	}
}

func ImageFormatFromFile(filePath string) (kroki.ImageFormat, error) {
	fileExtension := filepath.Ext(filePath)
	value := strings.ToLower(fileExtension)
	switch value {
	case ".svg":
		return kroki.SVG, nil
	case ".png":
		return kroki.PNG, nil
	case ".jpeg", ".jpg":
		return kroki.JPEG, nil
	case ".pdf":
		return kroki.PDF, nil
	default:
		return kroki.ImageFormat(""), errors.Errorf(
			"invalid image format %s.",
			value)
	}
}

func ResolveGraphFormat(graphFormatRaw string, filePath string) (kroki.DiagramType, error) {
	if graphFormatRaw == "" {
		return GraphFormatFromFile(filePath)
	} else {
		return GraphFormatFromValue(graphFormatRaw)
	}
}

func GraphFormatFromValue(value string) (kroki.DiagramType, error) {
	value = strings.ToLower(value)
	switch value {
	case "dot", "graphviz":
		return kroki.GraphViz, nil
	case "plantuml":
		return kroki.PlantUML, nil
	case "nomnoml":
		return kroki.Nomnoml, nil
	case "blockdiag":
		return kroki.BlockDiag, nil
	case "mermaid":
		return kroki.Mermaid, nil
	case "svgbob":
		return kroki.Svgbob, nil
	case "umlet":
		return kroki.UMlet, nil
	case "c4plantuml":
		return kroki.C4PlantUML, nil
	case "seqdiag":
		return kroki.SeqDiag, nil
	case "erd", "er":
		return kroki.Erd, nil
	case "nwdiag":
		return kroki.NwDiag, nil
	case "actdiag":
		return kroki.ActDiag, nil
	case "ditaa":
		return kroki.Ditaa, nil
	default:
		return kroki.DiagramType(""), errors.Errorf(
			"invalid graph format %s.",
			value)
	}
}

func GraphFormatFromFile(filePath string) (kroki.DiagramType, error) {
	fileExtension := filepath.Ext(filePath)
	value := strings.ToLower(fileExtension)
	switch value {
	case ".dot", ".gv", ".graphviz":
		return kroki.GraphViz, nil
	case ".puml", ".plantuml":
		return kroki.PlantUML, nil
	case ".nomnoml":
		return kroki.Nomnoml, nil
	case ".blockdiag":
		return kroki.BlockDiag, nil
	case ".mermaid":
		return kroki.Mermaid, nil
	case ".svgbob":
		return kroki.Svgbob, nil
	case ".umlet":
		return kroki.UMlet, nil
	case ".c4puml", ".c4", ".c4plantuml":
		return kroki.C4PlantUML, nil
	case ".seqdiag":
		return kroki.SeqDiag, nil
	case ".erd", ".er":
		return kroki.Erd, nil
	case ".nwdiag":
		return kroki.NwDiag, nil
	case ".actdiag":
		return kroki.ActDiag, nil
	case ".ditaa":
		return kroki.Ditaa, nil
	default:
		return kroki.DiagramType(""), errors.Errorf(
			"unable to infer the graph format from the file extension %s, please specify the diagram type using --type flag.",
			value)
	}
}

func GetClient(cmd *cobra.Command) kroki.Client {
	configFilePath, err := cmd.Flags().GetString("config")
	if err != nil {
		exit(err)
	}
	if configFilePath != "" {
		file, err := os.Open(configFilePath)
		if err != nil {
			exit(err)
		}
		err = viper.ReadConfig(file)
		if err != nil {
			exit(err)
		}
	}
	return kroki.New(kroki.Configuration{
		URL:     viper.GetString("endpoint"),
		Timeout: viper.GetDuration("timeout"),
	})
}
