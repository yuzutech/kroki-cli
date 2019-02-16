package cmd

import (
	"bufio"
	"fmt"
	"github.com/pkg/errors"
	"io"
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
	outFile, err := cmd.Flags().GetString("out-file")
	if filePath == "-" {
		writer, err := GetWriter(filePath, outFile)
		if err != nil {
			exit(err)
		}
		ConvertFromStdin(graphFormat, imageFormat, outFile, writer)
	} else {
		ConvertFromFile(filePath, graphFormat, imageFormat, outFile)
	}
}

func GetWriter(filePath string, outFile string) (io.Writer, error) {
	if outFile == "" || outFile == "-" {
		return os.Stdout, nil
	}
	file, err := os.Create(filePath)
	if err != nil {
		return nil, errors.Wrapf(err, "fail to create file '%s'", filePath)
	}
	defer file.Close()
	return file, nil
}

func ConvertFromStdin(diagramTypeRaw string, imageFormatRaw string, outFile string, writer io.Writer) {
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
	text, err := GetTextFromStdin()
	if err != nil {
		exit(err)
	}
	result, err := client.FromString(text, diagramType, imageFormat)
	if err != nil {
		exit(err)
	}
	_, err = writer.Write([]byte(result))
	if err != nil {
		exit(err)
	}
}

func GetTextFromStdin() (result string, err error) {
	reader := bufio.NewReader(os.Stdin)
	input, err := ioutil.ReadAll(reader)
	return string(input), err
}

func ConvertFromFile(filePath string, graphFormatRaw string, imageFormatRaw string, outFile string) {
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
		if  outFile == "" || outFile == "-" {
			return kroki.Svg, nil
		}
		return ImageFormatFromFile(outFile)
	}
	return ImageFormatFromValue(imageFormatRaw)
}

func ImageFormatFromValue(imageFormatRaw string) (kroki.ImageFormat, error) {
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

func ImageFormatFromFile(filePath string) (kroki.ImageFormat, error) {
	fileExtension := filepath.Ext(filePath)
	value := strings.ToLower(fileExtension)
	switch value {
	case ".svg":
		return kroki.Svg, nil
	case ".png":
		return kroki.ImageFormat("png"), nil
	case ".jpeg", ".jpg":
		return kroki.ImageFormat("jpeg"), nil
	case ".pdf":
		return kroki.ImageFormat("pdf"), nil
	default:
		return kroki.ImageFormat(""), errors.Errorf(
			"invalid image format %s.",
			value)
	}
}

func ResolveGraphFormat(graphFormatRaw string, filePath string) (kroki.GraphFormat, error) {
	if graphFormatRaw == "" {
		return GraphFormatFromFile(filePath)
	} else {
		return GraphFormatFromValue(graphFormatRaw)
	}
}

func GraphFormatFromValue(value string) (kroki.GraphFormat, error) {
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

func GraphFormatFromFile(filePath string) (kroki.GraphFormat, error) {
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
