package main

import (
	"fmt"
	"os"

	"git.arcline.it/ArclineIT/nexus/cmd/server"
)

func main() {
	if err := server.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "nexus: %v\n", err)
		os.Exit(1)
	}
}

