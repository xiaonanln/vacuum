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
	} `json:"dispatcher"`
}

var (
	config VacuumConfig
)

func checkError(err error) {
	if err != nil {
		log.Panic(err)
	}
}

func LoadConfig() {
	data, err := ioutil.ReadFile(CONFIG_FILENAME)
	checkError(err)

	err = json.Unmarshal(data, &config)
	checkError(err)

	log.WithField("config", config).Infof("Load config: %s", CONFIG_FILENAME)
}
