package configurations

import (
	"log"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"errors"
)

type Configuration struct {
	MaxMRU string `json:"MaxMRU"`
	SomethingElse string `json:"SomethingElse"`
}

func DefaultConfigurations() Configuration {
	return Configuration {
		MaxMRU: "10",
		SomethingElse: "3",
	}
}

func (c Configuration) GetConfigurations() Configuration {
	
	configurations := DefaultConfigurations()

	homeDir, _ := os.UserHomeDir()
	configPath := homeDir + "/.config/ucd.conf"
	configFile, err := os.Open(configPath)

	if errors.Is(err, os.ErrNotExist) {
		log.Fatalln("File does not exist")
		// marshall default configurations into configPath
	} else {
		configFileBytes, _ := ioutil.ReadAll(configFile)
		err := json.Unmarshal(configFileBytes, &configurations)

		if err != nil {
			log.Fatalln("Error decoding configurations:", err)
		}
	}

	fmt.Printf("contents of decoded json is: %#v\r\n", configurations)

	return configurations
}