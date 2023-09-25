package configurations

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"os"
)

type Configuration struct {
	MaxMRUDisplay string `json:"MaxMRUDisplay"`
}

func DefaultConfigurations() Configuration {
	return Configuration{
		MaxMRUDisplay: "10",
	}
}

func (c Configuration) GetConfigurations() Configuration {

	configurations := DefaultConfigurations()

	homeDir, _ := os.UserHomeDir()
	configPath := homeDir + "/.config/ucd/"
	configFileName := "ucd.conf"
	configFile, err := os.Open(configPath + configFileName)

	if errors.Is(err, os.ErrNotExist) {
		err := os.MkdirAll(configPath, 0700)
		if err != nil {
			log.Fatalln(err)
		}

		configFile, err := os.Create(configPath + configFileName)
		if err != nil {
			log.Fatalln(err)
		}
		defer configFile.Close()

		defaultConfigs, err := json.MarshalIndent(configurations, "", "\t")
		if err != nil {
			log.Fatalln(err)
		}

		configFile.Write(defaultConfigs)
		configFile.Sync()
		configFile.Close()

	} else {
		configFileBytes, _ := ioutil.ReadAll(configFile)
		err := json.Unmarshal(configFileBytes, &configurations)

		if err != nil {
			log.Fatalln("Error decoding configurations:", err)
		}
	}

	return configurations
}
