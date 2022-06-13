package main

import (
	"github.com/nikhilsbhat/gocd-prometheus-exporter/cmd"
	"log"
	"os"
)

func main() {
	app := cmd.App()
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
