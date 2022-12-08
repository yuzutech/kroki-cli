package pkg

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/yuzutech/kroki-go"
)

// getImageFormatExtensions returns a map of file extensions (including '.') with their corresponding image format
func getImageFormatExtensions() map[string]kroki.ImageFormat {
	imageFormatExtensions := map[string]kroki.ImageFormat{
		".jpg": kroki.JPEG,
	}
	supportedImageFormats := kroki.GetSupportedImageFormats()
	for _, v := range supportedImageFormats {
		imageFormatExtensions["."+string(v)] = v
	}
	return imageFormatExtensions
}

// getDiagramTypeNames returns a map of diagram names with their corresponding diagram type
func getDiagramTypeNames() map[string]kroki.DiagramType {
	diagramTypeNames := map[string]kroki.DiagramType{
		"dot": kroki.GraphViz,
	}
	supportedDiagramTypes := kroki.GetSupportedDiagramTypes()
	for _, v := range supportedDiagramTypes {
		diagramTypeNames[string(v)] = v
	}
	return diagramTypeNames
}

// getDiagramTypeExtensions returns a map of diagram file extensions (including '.') with their corresponding diagram type
func getDiagramTypeExtensions() map[string]kroki.DiagramType {
	diagramTypeExtensions := map[string]kroki.DiagramType{
		".d2":     kroki.D2,
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
	supportedDiagramTypes := kroki.GetSupportedDiagramTypes()
	for _, v := range supportedDiagramTypes {
		diagramTypeExtensions["."+string(v)] = v
	}
	return diagramTypeExtensions
}

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
	input, err := io.ReadAll(reader)
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
	if f, ok := getImageFormatExtensions()["."+value]; ok {
		return f, nil
	}
	return "", errors.Errorf(
		"invalid image format %s.",
		value)
}

func ImageFormatFromFile(filePath string) (kroki.ImageFormat, error) {
	fileExtension := filepath.Ext(filePath)
	value := strings.ToLower(fileExtension)
	if f, ok := getImageFormatExtensions()[value]; ok {
		return f, nil
	}
	return "", errors.Errorf(
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
	if d, ok := getDiagramTypeNames()[value]; ok {
		return d, nil
	}
	// support unrecognized type
	return kroki.DiagramType(value), nil
}

func GraphFormatFromFile(filePath string) (kroki.DiagramType, error) {
	fileExtension := filepath.Ext(filePath)
	value := strings.ToLower(fileExtension)
	if d, ok := getDiagramTypeExtensions()[fileExtension]; ok {
		return d, nil
	}
	return "", errors.Errorf(
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
