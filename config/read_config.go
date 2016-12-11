package config

import (
	"encoding/json"
	"io/ioutil"

	"github.com/xiaonanln/vacuum/vlog"
)

const (
	DEFAULT_CONFIG_FILENAME = "vacuum.conf"
)

var (
	usingConfigFilename string
)

type StorageConfig struct {
	Type       string `json:"type"`
	Directory  string `json:"directory"`
	Url        string `json:"url"`
	DB         string `json:"DB"`
	Concurrent int    `json:"concurrent"`
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
	usingConfigFilename = configFile
	if usingConfigFilename == "" {
		usingConfigFilename = DEFAULT_CONFIG_FILENAME
	}

	data, err := ioutil.ReadFile(usingConfigFilename)
	checkError(err)

	err = json.Unmarshal(data, &config)
	checkError(err)

	vlog.Info("Load config: %s, config=%v", usingConfigFilename, config)
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

func LoadExtraConfig(fieldName string) map[string]interface{} {
	if usingConfigFilename == "" {
		vlog.Panicf("config filename is unknown")
	}

	data, err := ioutil.ReadFile(usingConfigFilename)
	checkError(err)

	var totalConfig map[string]interface{}
	err = json.Unmarshal(data, &totalConfig)
	checkError(err)

	config := totalConfig[fieldName]
	return config.(map[string]interface{})
}
