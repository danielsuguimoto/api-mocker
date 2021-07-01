package main

import (
	"fmt"
	"os"

	"github.com/danielsuguimoto/api-mocker/server"
)

func main() {
	if len(os.Args) <= 1 {
		fmt.Fprint(os.Stderr, "error: resources json file is required.")
		os.Exit(2)
	}

	resourcesFilePath := os.Args[1];

	server := server.Create(3000)

	if err := server.LoadResources(resourcesFilePath); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v", err)
		os.Exit(1)
	}

	server.Listen()

	server.WatchResource(resourcesFilePath)
}
