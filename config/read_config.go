package config

import (
	"encoding/json"
	"io/ioutil"

	log "github.com/Sirupsen/logrus"
)

const (
	CONFIG_FILENAME = "vacuum.conf"
)

type VacuumConfig struct {
	Dispatcher struct {
		Host        string `json:"host"`
		ConsoleHost string `json:"console_host"`
		PublicIP    string `json:"public_ip"`
	} `json:"dispatcher"`
	Vacuums []struct {
	} `json:"vacuums"`
}

var (
	config VacuumConfig
)

func checkError(err error) {
	if err != nil {
		log.Panic(err)
	}
}

func LoadConfig(configFile string) {
	if configFile == "" {
		configFile = CONFIG_FILENAME
	}

	data, err := ioutil.ReadFile(configFile)
	checkError(err)

	err = json.Unmarshal(data, &config)
	checkError(err)

	log.WithField("config", config).Infof("Load config: %s", configFile)
}

func GetConfig() *VacuumConfig {
	return &config
}
