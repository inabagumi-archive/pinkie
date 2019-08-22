package main

import (
	"log"
	"os"

	"github.com/inabagumi/ytc/cli"
)

var version = "dev"

func main() {
	c, err := cli.NewClient(os.Args[0], version)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	c.Run(os.Args[1:])
}
