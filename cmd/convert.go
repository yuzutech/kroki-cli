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

var (

	// SupportedImageFormats maps string names back to their
	// kroki.ImageFormat constant value.
	SupportedImageFormats = []kroki.ImageFormat{
		kroki.SVG,
		kroki.PNG,
		kroki.JPEG,
		kroki.PDF,
		kroki.Base64,
	}

	// SupportedDiagramTypes maps string names back to the
	// corresponding kroki.DiagramType constant.
	SupportedDiagramTypes = []kroki.DiagramType{
		kroki.GraphViz,
		kroki.PlantUML,
		kroki.Nomnoml,
		kroki.BlockDiag,
		kroki.Mermaid,
		kroki.Svgbob,
		kroki.UMlet,
		kroki.C4PlantUML,
		kroki.SeqDiag,
		kroki.Erd,
		kroki.NwDiag,
		kroki.ActDiag,
		kroki.Ditaa,
		kroki.RackDiag,
		kroki.Vega,
		kroki.VegaLite,
		kroki.WaveDrom,
	}

	// The mappings are pre-populated with some Additional
	// forms that are matched/accepted for certain types.
	// The init() function will flesh them out with a set of
	// standard names based on the two Supported* slices.

	// Valid --format arguments for output file type
	ImageFormatNames = map[string]kroki.ImageFormat{
		"jpg": kroki.JPEG,
	}

	// File extensions matched to derive output file type
	ImageFormatExtensions = map[string]kroki.ImageFormat{
		".jpg": kroki.JPEG,
	}

	// Valid --type arguments for input file type
	DiagramTypeNames = map[string]kroki.DiagramType{
		"dot":    kroki.GraphViz,
		"er":     kroki.Erd,
		"nwdiag": kroki.NwDiag,
	}

	// Filename matching for input file type
	DiagramTypeExtensions = map[string]kroki.DiagramType{
		".dot":    kroki.GraphViz,
		".gv":     kroki.GraphViz,
		".puml":   kroki.PlantUML,
		".c4puml": kroki.C4PlantUML,
		".c4":     kroki.C4PlantUML,
		".er":     kroki.Erd,
		".vg":     kroki.Vega,
		".vgl":    kroki.VegaLite,
		".vl":     kroki.VegaLite,
	}
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
	if f, ok := ImageFormatNames[value]; ok {
		return f, nil
	}
	return kroki.ImageFormat(""), errors.Errorf(
		"invalid image format %s.",
		value)
}

func ImageFormatFromFile(filePath string) (kroki.ImageFormat, error) {
	fileExtension := filepath.Ext(filePath)
	value := strings.ToLower(fileExtension)
	if f, ok := ImageFormatExtensions[value]; ok {
		return f, nil
	}
	return kroki.ImageFormat(""), errors.Errorf(
		"invalid image format %s.",
		value)
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
	if d, ok := DiagramTypeNames[value]; ok {
		return d, nil
	}
	return kroki.DiagramType(""), errors.Errorf(
		"invalid graph format %s.",
		value)
}

func GraphFormatFromFile(filePath string) (kroki.DiagramType, error) {
	fileExtension := filepath.Ext(filePath)
	value := strings.ToLower(fileExtension)
	if d, ok := DiagramTypeExtensions[fileExtension]; ok {
		return d, nil
	}
	return kroki.DiagramType(""), errors.Errorf(
		"unable to infer the graph format from the file extension %s, please specify the diagram type using --type flag.",
		value)
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

// Set up mappings
func init() {
	for _, v := range SupportedDiagramTypes {
		DiagramTypeNames[string(v)] = v
		DiagramTypeExtensions["."+string(v)] = v
	}
	for _, v := range SupportedImageFormats {
		ImageFormatNames[string(v)] = v
		ImageFormatExtensions["."+string(v)] = v
	}
}
