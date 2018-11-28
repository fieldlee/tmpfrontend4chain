/**
 *	Author Fanxu(746439274@qq.com)
 */

package config

import (
	"encoding/json"
	"fmt"
	"frontend4chain/constant"
	"io/ioutil"
	"log"
)

var __config *Config

func All() *Config {
	return __config
}

type ServerConfig struct {
	Ip   string `json:"ip"`
	Port string `json:"port"`
}
type Couchdb struct {
	Ip       string `json:"ip"`
	Port     int    `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
}
type Config struct {
	Couchdb Couchdb      `json:"couchdb"`
	Listen  ServerConfig `json:"listen"`
}

func InitConf(args []string) {
	var confFile = fmt.Sprint(constant.CONFIGPATH, "/online.json")
	for i := 1; i < len(args); i++ {
		switch args[i] {
		case "-f":
			if i == len(args) {
				log.Fatalln("invalid config file")
			}
			confFile = args[i+1]
		}
	}
	data, err := ioutil.ReadFile(confFile)
	if err != nil {
		log.Fatalln(err.Error(), "cannot find the file: online.json")
	}
	err = json.Unmarshal(data, &__config)
	if err != nil {
		log.Fatalln(err.Error(), "cannot parse the file: online.json")
	}
}
