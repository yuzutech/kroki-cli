package main

import (
	"log"

	"github.com/yuzutech/kroki-cli/pkg"
)

var (
	// version comes from the tag (during the build)
	version = "dev"
	// commit represents the HEAD commit (during the build)
	commit = "n/a"
)

func main() {
	log.SetFlags(0)
	pkg.Execute(version, commit)
}
