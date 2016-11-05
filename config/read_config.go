package config

import (
	"encoding/json"
	"io/ioutil"

	"github.com/xiaonanln/vacuum/vlog"
)

const (
	CONFIG_FILENAME = "vacuum.conf"
)

type StorageConfig struct {
	Type      string `json:"type"`
	Directory string `json:"directory"`
}

type VacuumConfig struct {
	Storage *StorageConfig `json:"storage"`
}

type DispatcherConfig struct {
	Host        string `json:"host"`
	ConsoleHost string `json:"console_host"`
	PublicIP    string `json:"public_ip"`
}

type Config struct {
	Dispatcher   DispatcherConfig `json:"dispatcher"`
	VacuumCommon VacuumConfig     `json:"vacuum_common"`
	Vacuums      []VacuumConfig   `json:"vacuums"`
}

func (c *Config) GetVacuumConfig(serverID int) (ret VacuumConfig) {
	ret = c.VacuumCommon
	moreConfig := &c.Vacuums[serverID-1]
	if moreConfig.Storage != nil {
		ret.Storage = moreConfig.Storage
	}

	return
}

var (
	config Config
)

func checkError(err error) {
	if err != nil {
		vlog.Panic(err)
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

	vlog.Infof("Load config: %s, config=%v", configFile, config)
}

func GetConfig() *Config {
	return &config
}

func FormatConfig(c interface{}) string {
	data, err := json.MarshalIndent(c, "", "\t")
	if err != nil {
		panic(err)
	}
	return string(data)
}
