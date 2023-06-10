package main

import (
	"log"

	"github.com/nikhilsbhat/gocd-prometheus-exporter/cmd"
	docgen "github.com/nikhilsbhat/urfavecli-docgen"
)

//go:generate go run github.com/nikhilsbhat/gocd-prometheus-exporter/docs
func main() {
	if err := docgen.GenerateDocs(cmd.App(), "gocd_prometheus_exporter"); err != nil {
		log.Fatalln(err)
	}
}
