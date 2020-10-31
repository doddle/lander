package main

import (
	"io/ioutil"

	"github.com/withmandala/go-log"
	"gopkg.in/yaml.v2"
)

type Config struct {
	Colorschemes []Colorschemes `yaml:"colorschemes"`
}
type Colorschemes struct {
	Color string   `yaml:"color"`
	Hex   string   `yaml:"hex"`
	Verbs []string `yaml:"verbs"`
}

func initialConfig(logger *log.Logger, configFlag string) Config {
	var results Config

	// the default value for `-config` == "default"
	if configFlag == "default" {

		logger.Info("using default settings")

		var schemes []Colorschemes
		schemes = append(schemes, Colorschemes{
			Color: "light-blue accent-1",
			Hex:   "26c5e8",
			Verbs: []string{"staging", "stage", "uat", "preprod"},
		})
		schemes = append(schemes, Colorschemes{
			// Color: "red lighten-2",
			Color: "red lighten-2",
			Hex:   "d84a73",
			Verbs: []string{"prod", "live"},
		})
		schemes = append(schemes, Colorschemes{
			Color: "light-green lighten-2",
			Hex:   "4ad86f",
			Verbs: []string{"dev", "sandbox", "play"},
		})

		// TODO: always have a "default scheme"
		results := Config{
			Colorschemes: schemes,
		}
		return results
	}

	yamlFile, err := ioutil.ReadFile(configFlag)
	if err != nil {
		logger.Errorf("Error reading file: '%s', err: %v", configFlag, err)
	}
	err = yaml.Unmarshal(yamlFile, &results)
	if err != nil {
		logger.Fatalf("Invalid yaml: %v", err)
	}
	if len(results.Colorschemes) < 1 {
		logger.Fatal("'colorschemes: []' in config file was empty")
	}
	return results
}
