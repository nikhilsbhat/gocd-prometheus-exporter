package main

import (
	"log"
	"os"

	"github.com/nikhilsbhat/gocd-prometheus-exporter/cmd"
)

func main() {
	app := cmd.App()
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
