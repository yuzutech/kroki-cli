package cmd

import (
	"bufio"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/yuzutech/kroki-go"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
)

var gConfigFilePath string
var gDiagramType string
var gImageFormat string

var gVersion string
var gCommit string

var client kroki.Client

var RootCmd = &cobra.Command{
	Use:           "kroki convert",
	Short:         `Convert text diagram to image.
By default, the output is written to a file with the basename of the source file and the appropriate extension.
Example: kroki convert hello.dot`,
}

var convertCmd = &cobra.Command{
	Use:           "convert file",
	Short:         "Convert text diagram to image",
	Args:          cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		fileName := args[0]
		if fileName == "-" {
			reader := bufio.NewReader(os.Stdin)
			text , err := ioutil.ReadAll(reader)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			if gDiagramType == "" {
				fmt.Println("You must specify the diagram type using --type flag")
				os.Exit(1)
			}
			result, err := client.FromString(string(text[:]), kroki.Graphviz, kroki.ImageFormat(gImageFormat))
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			fmt.Println(result)
		} else {
			var diagramType kroki.GraphFormat
			var imageFormat = kroki.ImageFormat(gImageFormat)
			if gDiagramType == "" {
				diagramType = inferDiagramType(fileName)
				if diagramType == "" {
					fmt.Println("Unable to infer the diagram type, please specify the diagram type using --type flag")
					os.Exit(1)
				}
			} else {
				diagramType = kroki.GraphFormat(gDiagramType)
			}
			result, err := client.FromFile(fileName, diagramType, imageFormat)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			err = client.WriteToFile(outputFilePath(fileName, imageFormat), result)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		}
	},
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version of kroki",
	Run: func(cmd *cobra.Command, _ []string) {
		fmt.Printf("%s %s %s [https://kroki.io]\n", cmd.Parent().Name(), gVersion, gCommit)
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the RootCmd.
func Execute(version, commit string) {
	gVersion = version
	gCommit = commit
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	convertCmd.PersistentFlags().StringVarP(&gConfigFilePath, "config", "c", "", "specify an alternate config file [env KROKI_CONFIG]")
	convertCmd.PersistentFlags().StringVarP(&gDiagramType, "type", "t", "", "specify the diagram type [actdiag, blockdiag, c4plantuml, ditaa, dot, erd, graphviz, nomnoml, nwdiag, plantuml, seqdiag, svgbob, umlet] (default: infer from file extension)")
	convertCmd.PersistentFlags().StringVarP(&gImageFormat, "format", "f", string(kroki.Svg), "specify the output format (default: svg)")
	// -o, --out-file FILE              output file (default: based on path of input file); use - to output to STDOUT

	RootCmd.AddCommand(versionCmd)
	RootCmd.AddCommand(convertCmd)

	cobra.OnInitialize(buildClient)
}

func outputFilePath(fileName string, imageFormat kroki.ImageFormat) string {
	fileExtension := filepath.Ext(fileName)
	return strings.Replace(fileName, fileExtension, "." + string(imageFormat), 1)
}

func inferDiagramType(fileName string) kroki.GraphFormat {
	switch fileExtension := filepath.Ext(fileName); fileExtension {
	case ".dot", ".gv", ".graphviz":
		return kroki.Graphviz
	case ".puml", ".plantuml":
		return kroki.Plantuml
	case ".nomnoml":
		return kroki.Nomnoml
	case ".blockdiag":
		return kroki.BlockDiag
	case ".mermaid":
		return kroki.Mermaid
	case ".svgbob":
		return kroki.Svgbob
	case ".umlet":
		return kroki.Umlet
	case ".c4puml", ".c4":
		return kroki.C4plantuml
	case ".seqdiag":
		return kroki.SeqDiag
	case ".erd":
		return kroki.GraphFormat("erd")
	case ".nwdiag":
		return kroki.GraphFormat("nwdiag")
	case ".actdiag":
		return kroki.GraphFormat("actdiag")
	case ".ditaa":
		return kroki.GraphFormat("ditaa")
	default:
		return kroki.GraphFormat("")
	}
}

func buildClient() {
	client = kroki.New(kroki.Configuration{
		URL:     "https://demo.kroki.io",
		Timeout: time.Second * 20,
	})
}
