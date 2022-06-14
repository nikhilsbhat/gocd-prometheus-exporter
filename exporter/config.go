package exporter

import (
	"errors"
	"io/ioutil"
	"log"
	"os"

	"github.com/imdario/mergo"
	"gopkg.in/yaml.v2"
)

func GetConfig(conf Config, path string) (*Config, error) {
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		log.Printf("config file %s not found, dropping configurations from file", path)
		return &conf, err
	}

	fileOUT, err := ioutil.ReadFile(path)
	if err != nil {
		log.Println("failed to read the config file, dropping configurations from file")
		return &conf, err
	}

	var newConfig Config
	if err = yaml.Unmarshal(fileOUT, &newConfig); err != nil {
		log.Println("failed to unmarshall configurations, dropping configurations from file")
		return &conf, err
	}
	if err = mergo.Merge(&newConfig, &conf, mergo.WithOverride); err != nil {
		log.Println("failed to merge configurations, dropping configurations from file")
		return &conf, err
	}
	return &newConfig, nil
}
