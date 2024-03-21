package configurations

import (
	"encoding/json"
	"io"
	"log"
	"os"
)

type Configuration struct {
	MaxMRUDisplay        int  `json:"MaxMRUDisplay"`
	FileFallbackBehavior bool `json:"FileFallbackBehavior"`
}

func DefaultConfigurations() Configuration {
	return Configuration{
		MaxMRUDisplay:        -1,
		FileFallbackBehavior: true,
	}
}

func (c Configuration) GetConfigurations() Configuration {

	configurations := DefaultConfigurations()

	homeDir, _ := os.UserHomeDir()
	configPath := homeDir + "/.config/ucd/"
	configFileName := "ucd.conf"

	os.MkdirAll(configPath, 0700)
	configFile, _ := os.Open(configPath + configFileName)

	configFileBytes, _ := io.ReadAll(configFile)

	if len(configFileBytes) > 0 {
		err := json.Unmarshal(configFileBytes, &configurations)

		if err != nil {
			log.Fatalln("Error decoding configurations:", err)
		}
	}

	configFileBytes, _ = json.MarshalIndent(configurations, "", "\t")
	os.WriteFile(configPath+configFileName, configFileBytes, 0644)

	return configurations
}
