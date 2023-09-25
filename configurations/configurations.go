package configurations

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"os"
)

type Configuration struct {
	MaxMRUDisplay int `json:"MaxMRUDisplay"`
}

func DefaultConfigurations() Configuration {
	return Configuration{
		MaxMRUDisplay: 10,
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
	}

	configFileBytes, _ := io.ReadAll(configFile)

	if len(configFileBytes) > 0 {
		err := json.Unmarshal(configFileBytes, &configurations)

		if err != nil {
			log.Fatalln("Error decoding configurations:", err)
		}
	} else {
		defaultConfigs, _ := json.MarshalIndent(configurations, "", "\t")
		os.WriteFile(configPath+configFileName, defaultConfigs, 0644)
	}

	return configurations
}
