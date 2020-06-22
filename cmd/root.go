package cmd

import (
	"fmt"
	"sort"

	"github.com/spf13/cobra"
)

var gVersion string
var gCommit string

var RootCmd = &cobra.Command{
	Use: "kroki convert",
	Short: `Convert text diagram to image.
By default, the output is written to a file with the basename of the source file and the appropriate extension.
Example: kroki convert hello.dot`,
}

var convertCmd = &cobra.Command{
	Use:   "convert file",
	Short: "Convert text diagram to image",
	Args:  cobra.ExactArgs(1),
	Run:   Convert,
}

var encodeCmd = &cobra.Command{
	Use:   "encode file",
	Short: "Encode text diagram in deflate + base64 format",
	Args:  cobra.ExactArgs(1),
	Run:   Encode,
}

var decodeCmd = &cobra.Command{
	Use:   "decode input",
	Short: "Decode an encoded (deflate + base64) diagram",
	Args:  cobra.ExactArgs(1),
	Run:   Decode,
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
		exit(err)
	}
}

func init() {
	diagramTypeNames := make([]string, len(SupportedDiagramTypes))
	imageFormatNames := make([]string, len(SupportedImageFormats))
	for i, v := range SupportedDiagramTypes {
		diagramTypeNames[i] = string(v)
	}
	sort.Strings(diagramTypeNames)
	for i, v := range SupportedImageFormats {
		imageFormatNames[i] = string(v)
	}
	sort.Strings(imageFormatNames)

	typeHelp := fmt.Sprintf("diagram type %s (default: infer from file extension)", diagramTypeNames)
	formatHelp := fmt.Sprintf("output format %s (default: infer from output file extension otherwise svg)", imageFormatNames)

	convertCmd.PersistentFlags().StringP("config", "c", "", "alternate config file [env KROKI_CONFIG]")
	convertCmd.PersistentFlags().StringP("type", "t", "", typeHelp)
	convertCmd.PersistentFlags().StringP("format", "f", "", formatHelp)
	convertCmd.PersistentFlags().StringP("out-file", "o", "", "output file (default: based on path of input file); use - to output to STDOUT")
	RootCmd.AddCommand(versionCmd)
	RootCmd.AddCommand(convertCmd)
	RootCmd.AddCommand(encodeCmd)
	RootCmd.AddCommand(decodeCmd)

	SetupConfig()

	cobra.OnInitialize(InitDefaultConfig)
}
