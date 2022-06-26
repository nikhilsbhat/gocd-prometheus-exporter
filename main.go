package main

import (
	"log"
	"os"

	"github.com/nikhilsbhat/gocd-prometheus-exporter/cmd"
)

func main() {
	app := cmd.App()
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
