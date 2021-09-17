package config

import (
	"github.com/taosdata/go-utils/json"
	"github.com/taosdata/go-utils/mqtt"
	"github.com/taosdata/go-utils/util"
	"os"
)

type Config struct {
	TDengine *TDengine    `json:"TDengine"`
	MQTT     *mqtt.Config `json:"mqtt"`
	ShowSql  bool         `json:"showSql"`
}
type TDengine struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	DB       string `json:"db"`
}

func Init(path string) *Config {
	if !util.PathExist(path) {
		panic("config path not exist")
	}
	f, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	var c Config
	err = json.NewDecoder(f).Decode(&c)
	if err != nil {
		panic(err)
	}
	return &c
}
