package cmd

import (
	"fmt"
	"os"
)

func exit(a ...interface{}) {
	fmt.Fprintln(os.Stderr, a...)
	os.Exit(1)
}
