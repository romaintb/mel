package main

import (
	"fmt"
	"os"

	"github.com/romaintb/mel/internal/app"
)

// Version will be set at build time via ldflags
var version = "dev"

func main() {
	if err := app.Run(version); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
