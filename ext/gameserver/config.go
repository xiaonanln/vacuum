package gameserver

import (
	"encoding/json"

	"github.com/xiaonanln/vacuum/config"
	"github.com/xiaonanln/vacuum/vlog"
)

type gameserverConfig struct {
	GatesNum       int `json:"gates_num"`
	GatesStartPort int `json:"gates_start_port"`
	GatesPortStep  int `json:"gates_port_step"`
}

func loadGameserverConfig() *gameserverConfig {
	_gameserverConfig := config.LoadExtraConfig("gameserver")
	configData, err := json.Marshal(_gameserverConfig)
	if err != nil {
		vlog.Panic(err)
	}

	vlog.Debug("loadGameServerConfig: %s\n", string(configData))

	var gameserverConfig gameserverConfig
	json.Unmarshal(configData, &gameserverConfig)
	return &gameserverConfig
}
